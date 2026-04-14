package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/service"
	"github.com/ponte-tech/cretor-back/shared/middleware"
	"github.com/ponte-tech/cretor-back/shared/response"
	"go.uber.org/zap"
)

type EmpreendimentoHandler struct {
	service *service.EmpreendimentoService
	logger  *zap.Logger
}

func NewEmpreendimentoHandler(svc *service.EmpreendimentoService, logger *zap.Logger) *EmpreendimentoHandler {
	return &EmpreendimentoHandler{service: svc, logger: logger}
}

// List handles GET /empreendimentos?cidade=&uf=&search=&status_obra=&construtora_id=
func (h *EmpreendimentoHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	filter := domain.EmpreendimentoFilter{}
	if v := r.URL.Query().Get("cidade"); v != "" {
		filter.Cidade = &v
	}
	if v := r.URL.Query().Get("uf"); v != "" {
		filter.UF = &v
	}
	if v := r.URL.Query().Get("search"); v != "" {
		filter.Search = &v
	}
	if v := r.URL.Query().Get("status_obra"); v != "" {
		filter.StatusObra = &v
	}
	if v := r.URL.Query().Get("construtora_id"); v != "" {
		filter.ConstrutoraID = &v
	}
	if v := r.URL.Query().Get("dormitorios"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Dormitorios = append(filter.Dormitorios, n)
			}
		}
	}
	if v := r.URL.Query().Get("suites"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Suites = append(filter.Suites, n)
			}
		}
	}
	if v := r.URL.Query().Get("vagas"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.Vagas = append(filter.Vagas, n)
			}
		}
	}
	if v := r.URL.Query().Get("metragem_min"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MetragemMin = &f
		}
	}
	if v := r.URL.Query().Get("metragem_max"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MetragemMax = &f
		}
	}
	if v := r.URL.Query().Get("diferenciais_unidade"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if t := strings.TrimSpace(s); t != "" {
				filter.DiferenciaisUnidade = append(filter.DiferenciaisUnidade, t)
			}
		}
	}
	if v := r.URL.Query().Get("diferenciais_condominio"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if t := strings.TrimSpace(s); t != "" {
				filter.DiferenciaisCondominio = append(filter.DiferenciaisCondominio, t)
			}
		}
	}

	cards, err := h.service.ListEmpreendimentos(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list empreendimentos", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.JSON(w, cards, http.StatusOK)
}

// GetByID handles GET /empreendimentos/{id}
func (h *EmpreendimentoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := chi.URLParam(r, "id")

	detail, err := h.service.GetEmpreendimentoDetail(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("failed to get empreendimento", zap.String("id", id), zap.Error(err))
		response.Error(w, "not found", http.StatusNotFound)
		return
	}

	response.JSON(w, detail, http.StatusOK)
}

// ListUnidades handles GET /empreendimentos/{id}/unidades?secao=&status=
func (h *EmpreendimentoHandler) ListUnidades(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	empreendimentoID := chi.URLParam(r, "id")

	filter := domain.UnidadeFilter{
		EmpreendimentoID: empreendimentoID,
	}
	if v := r.URL.Query().Get("secao"); v != "" {
		filter.Secao = &v
	}
	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}

	unidades, err := h.service.ListUnidades(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list unidades", zap.String("empreendimento_id", empreendimentoID), zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.JSON(w, unidades, http.StatusOK)
}

// ListConstrutoras handles GET /construtoras?search=
func (h *EmpreendimentoHandler) ListConstrutoras(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	filter := domain.ConstrutoraFilter{}
	if v := r.URL.Query().Get("search"); v != "" {
		filter.Search = &v
	}

	construtoras, err := h.service.ListConstrutoras(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list construtoras", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.JSON(w, construtoras, http.StatusOK)
}
