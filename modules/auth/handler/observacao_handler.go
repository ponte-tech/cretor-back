package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/dto"
	"github.com/ponte-tech/cretor-back/shared/middleware"
	"github.com/ponte-tech/cretor-back/shared/response"
	"github.com/ponte-tech/cretor-back/shared/validator"
	"go.uber.org/zap"
)

type ObservacaoHandler struct {
	repo   domain.ObservacaoRepository
	logger *zap.Logger
}

func NewObservacaoHandler(repo domain.ObservacaoRepository, logger *zap.Logger) *ObservacaoHandler {
	return &ObservacaoHandler{repo: repo, logger: logger}
}

// POST /pipeline/{id}/observacoes
func (h *ObservacaoHandler) Create(w http.ResponseWriter, r *http.Request) {
	negocioID := chi.URLParam(r, "id")

	var req dto.CreateObservacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	tenantID := middleware.GetTenantID(r.Context())

	obs := &domain.Observacao{
		TenantID:  tenantID,
		NegocioID: negocioID,
		Texto:     req.Texto,
		Autor:     req.Autor,
	}

	if err := h.repo.Create(r.Context(), obs); err != nil {
		h.logger.Error("failed to create observacao", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.JSON(w, dto.ToObservacaoResponse(obs), http.StatusCreated)
}

// GET /pipeline/{id}/observacoes
func (h *ObservacaoHandler) List(w http.ResponseWriter, r *http.Request) {
	negocioID := chi.URLParam(r, "id")

	observacoes, err := h.repo.ListByNegocio(r.Context(), negocioID)
	if err != nil {
		h.logger.Error("failed to list observacoes", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result := make([]dto.ObservacaoResponse, len(observacoes))
	for i, o := range observacoes {
		result[i] = dto.ToObservacaoResponse(&o)
	}
	response.JSON(w, result, http.StatusOK)
}
