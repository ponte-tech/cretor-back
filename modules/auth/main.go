package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ponte-tech/cretor-back/modules/auth/handler"
	"github.com/ponte-tech/cretor-back/modules/auth/repository"
	"github.com/ponte-tech/cretor-back/modules/auth/service"
	"github.com/ponte-tech/cretor-back/shared/auth"
	"github.com/ponte-tech/cretor-back/shared/config"
	"github.com/ponte-tech/cretor-back/shared/database"
	"github.com/ponte-tech/cretor-back/shared/middleware"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	var logger *zap.Logger
	var err error
	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := database.NewClient(ctx, cfg.MongoDBURI)
	if err != nil {
		logger.Fatal("failed to connect to MongoDB", zap.Error(err))
	}
	logger.Info("connected to MongoDB")

	db := mongoClient.Database("cretor")

	// Ensure indexes
	if err := repository.EnsureIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure indexes", zap.Error(err))
	}
	if err := repository.EnsureLeadIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure lead indexes", zap.Error(err))
	}
	if err := repository.EnsurePipelineIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure pipeline indexes", zap.Error(err))
	}
	if err := repository.EnsureObservacaoIndexes(ctx, db); err != nil {
		logger.Fatal("failed to ensure observacao indexes", zap.Error(err))
	}
	logger.Info("indexes ensured")

	// Wire dependencies
	usuarioRepo := repository.NewUsuarioRepository(db)
	leadRepo := repository.NewLeadRepository(db)
	pipelineRepo := repository.NewPipelineRepository(db)
	observacaoRepo := repository.NewObservacaoRepository(db)
	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret, cfg.JWTExpirationHrs, cfg.RefreshTokenHrs)
	authService := service.NewAuthService(usuarioRepo, jwtProvider, logger)
	authHandler := handler.NewAuthHandler(authService, logger)
	leadHandler := handler.NewLeadHandler(leadRepo, pipelineRepo, logger)
	pipelineHandler := handler.NewPipelineHandler(pipelineRepo, logger)
	observacaoHandler := handler.NewObservacaoHandler(observacaoRepo, logger)
	emailHandler := handler.NewEmailHandler(logger)

	// Router
	r := buildRouter(authHandler, leadHandler, pipelineHandler, observacaoHandler, emailHandler, jwtProvider, logger)

	// Lambda or Local
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running as AWS Lambda (API Gateway v2 / HTTP API)
		adapter := chiadapter.NewV2(r)
		lambda.Start(adapter.ProxyWithContextV2)
	} else {
		// Running locally
		startLocalServer(r, cfg.Port, mongoClient, logger)
	}
}

func buildRouter(authHandler *handler.AuthHandler, leadHandler *handler.LeadHandler, pipelineHandler *handler.PipelineHandler, observacaoHandler *handler.ObservacaoHandler, emailHandler *handler.EmailHandler, jwtProvider *auth.JWTProvider, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Rate limiter for public endpoints (5 requests/minute per IP)
	publicLimiter := middleware.NewRateLimiter(5, time.Minute)

	// Public auth routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireTenant)

		r.Post("/auth/signup", authHandler.Signup)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.RefreshToken)
	})

	// Public lead route (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireTenant)
		r.Use(publicLimiter.Handler)

		r.Post("/leads", leadHandler.Create)
	})

	// Protected lead routes (CRUD)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireTenant)
		r.Use(middleware.RequireAuth(jwtProvider))

		r.Get("/leads", leadHandler.List)
		r.Get("/leads/{id}", leadHandler.GetByID)
		r.Put("/leads/{id}", leadHandler.Update)
		r.Delete("/leads/{id}", leadHandler.Delete)

		r.Post("/pipeline", pipelineHandler.Create)
		r.Get("/pipeline", pipelineHandler.List)
		r.Put("/pipeline/{id}", pipelineHandler.Update)
		r.Patch("/pipeline/{id}/move", pipelineHandler.Move)
		r.Delete("/pipeline/{id}", pipelineHandler.Delete)

		r.Get("/pipeline/{id}/observacoes", observacaoHandler.List)
		r.Post("/pipeline/{id}/observacoes", observacaoHandler.Create)

		r.Post("/email/send", emailHandler.Send)
	})

	// Protected auth routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireTenant)
		r.Use(middleware.RequireAuth(jwtProvider))

		r.Post("/auth/logout", authHandler.Logout)
	})

	return r
}

func startLocalServer(r *chi.Mux, port string, mongoClient interface{ Disconnect(ctx context.Context) error }, logger *zap.Logger) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("auth module starting on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		logger.Error("failed to disconnect MongoDB", zap.Error(err))
	}

	logger.Info("auth module stopped")
}
