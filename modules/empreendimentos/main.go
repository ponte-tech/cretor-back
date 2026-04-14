package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/dto"
	"github.com/ponte-tech/cretor-back/modules/auth/repository"
	"github.com/ponte-tech/cretor-back/modules/auth/service"
	"github.com/ponte-tech/cretor-back/shared/config"
	"github.com/ponte-tech/cretor-back/shared/database"
	"github.com/ponte-tech/cretor-back/shared/middleware"
	"github.com/ponte-tech/cretor-back/shared/response"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoClient, err := database.NewClient(ctx, cfg.MongoDBURI)
	if err != nil {
		logger.Fatal("failed to connect to MongoDB", zap.Error(err))
	}
	logger.Info("connected to MongoDB")

	db := mongoClient.Database("cretor")

	// Ensure indexes
	if err := repository.EnsureConstrutoraIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure construtora indexes", zap.Error(err))
	}
	if err := repository.EnsureEmpreendimentoIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure empreendimento indexes", zap.Error(err))
	}
	if err := repository.EnsureUnidadeIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure unidade indexes", zap.Error(err))
	}
	logger.Info("indexes ensured")

	// Wire dependencies
	empreendimentoRepo := repository.NewEmpreendimentoRepository(db)
	construtoraRepo := repository.NewConstrutoraRepository(db)
	unidadeRepo := repository.NewUnidadeRepository(db)
	svc := service.NewEmpreendimentoService(empreendimentoRepo, construtoraRepo, unidadeRepo, logger)
	h := newHandler(svc, logger)

	r := buildRouter(h, logger)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := chiadapter.NewV2(r)
		lambda.Start(adapter.ProxyWithContextV2)
	} else {
		startLocalServer(r, cfg.Port, mongoClient, logger)
	}
}

// --- Handler ---

type empreendimentoHandler struct {
	svc    *service.EmpreendimentoService
	logger *zap.Logger
}

func newHandler(svc *service.EmpreendimentoService, logger *zap.Logger) *empreendimentoHandler {
	return &empreendimentoHandler{svc: svc, logger: logger}
}

func parseIntParam(q string, defaultVal int) int {
	if q == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(q)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}

func (h *empreendimentoHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	q := r.URL.Query()

	page := parseIntParam(q.Get("page"), 1)
	pageSize := parseIntParam(q.Get("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}

	filter := domain.EmpreendimentoFilter{
		Pagination: domain.Pagination{Page: page, PageSize: pageSize},
	}
	if v := q.Get("search"); v != "" {
		filter.Search = &v
	}
	if v := q.Get("cidade"); v != "" {
		filter.Cidade = &v
	}
	if v := q.Get("uf"); v != "" {
		filter.UF = &v
	}
	if v := q.Get("status_obra"); v != "" {
		filter.StatusObra = &v
	}
	if v := q.Get("construtora_id"); v != "" {
		filter.ConstrutoraID = &v
	}
	if v := q.Get("dormitorios"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Dormitorios = append(filter.Dormitorios, n)
			}
		}
	}
	if v := q.Get("suites"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Suites = append(filter.Suites, n)
			}
		}
	}
	if v := q.Get("vagas"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Vagas = append(filter.Vagas, n)
			}
		}
	}
	if v := q.Get("metragem_min"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MetragemMin = &f
		}
	}
	if v := q.Get("metragem_max"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MetragemMax = &f
		}
	}
	if v := q.Get("diferenciais_unidade"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if t := strings.TrimSpace(s); t != "" {
				filter.DiferenciaisUnidade = append(filter.DiferenciaisUnidade, t)
			}
		}
	}
	if v := q.Get("diferenciais_condominio"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if t := strings.TrimSpace(s); t != "" {
				filter.DiferenciaisCondominio = append(filter.DiferenciaisCondominio, t)
			}
		}
	}

	// Bounding box filter (map search)
	if q.Get("sw_lat") != "" && q.Get("sw_lng") != "" && q.Get("ne_lat") != "" && q.Get("ne_lng") != "" {
		swLat, _ := strconv.ParseFloat(q.Get("sw_lat"), 64)
		swLng, _ := strconv.ParseFloat(q.Get("sw_lng"), 64)
		neLat, _ := strconv.ParseFloat(q.Get("ne_lat"), 64)
		neLng, _ := strconv.ParseFloat(q.Get("ne_lng"), 64)
		if swLat != 0 || swLng != 0 || neLat != 0 || neLng != 0 {
			filter.Bounds = &domain.BoundingBox{SwLat: swLat, SwLng: swLng, NeLat: neLat, NeLng: neLng}
		}
	}

	// Default: sem filtro nenhum e sem bounds, foca em Balneário Camboriú
	hasFilter := filter.Search != nil || q.Get("cidade") != "" || q.Get("uf") != "" || q.Get("status_obra") != "" || q.Get("construtora_id") != "" || filter.Bounds != nil || len(filter.Dormitorios) > 0 || len(filter.Suites) > 0 || len(filter.Vagas) > 0 || filter.MetragemMin != nil || filter.MetragemMax != nil || len(filter.DiferenciaisUnidade) > 0 || len(filter.DiferenciaisCondominio) > 0
	if !hasFilter {
		defaultCidade := "Camboriú"
		filter.Cidade = &defaultCidade
	}

	result, err := h.svc.ListEmpreendimentos(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list empreendimentos", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, result, http.StatusOK)
}

func (h *empreendimentoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := chi.URLParam(r, "id")

	detail, err := h.svc.GetEmpreendimentoDetail(r.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "empreendimento: not found" {
			response.Error(w, "empreendimento not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to get empreendimento", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, detail, http.StatusOK)
}

func (h *empreendimentoHandler) ListUnidades(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	empID := chi.URLParam(r, "id")
	q := r.URL.Query()

	filter := domain.UnidadeFilter{EmpreendimentoID: empID}
	if v := q.Get("secao"); v != "" {
		filter.Secao = &v
	}
	if v := q.Get("status"); v != "" {
		filter.Status = &v
	}

	unidades, err := h.svc.ListUnidades(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list unidades", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, unidades, http.StatusOK)
}

func (h *empreendimentoHandler) ListConstrutoras(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	q := r.URL.Query()

	filter := domain.ConstrutoraFilter{}
	if v := q.Get("search"); v != "" {
		filter.Search = &v
	}

	construtoras, err := h.svc.ListConstrutoras(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list construtoras", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, construtoras, http.StatusOK)
}

func (h *empreendimentoHandler) Filters(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	filters, err := h.svc.GetFilters(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("failed to get filters", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, filters, http.StatusOK)
}

// Ensure handler uses dto package (compile check)
var _ = dto.ToEmpreendimentoResponse

// --- Router ---

func buildRouter(h *empreendimentoHandler, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "module": "empreendimentos"})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireTenant)

		r.Get("/empreendimentos", h.List)
		r.Get("/empreendimentos/filters", h.Filters)
		r.Get("/empreendimentos/{id}", h.GetByID)
		r.Get("/empreendimentos/{id}/unidades", h.ListUnidades)
		r.Get("/construtoras", h.ListConstrutoras)
	})

	return r
}

// --- Local Server ---

func startLocalServer(r *chi.Mux, port string, mongoClient interface{ Disconnect(ctx context.Context) error }, logger *zap.Logger) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("empreendimentos module starting", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient.Disconnect(ctx)
	server.Shutdown(ctx)
}
