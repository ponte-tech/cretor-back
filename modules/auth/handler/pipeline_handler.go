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

var etapaProbabilidade = map[string]int{
	"primeiro_contato": 20, "qualificado": 40, "visita_agendada": 55,
	"proposta_enviada": 65, "negociacao": 80, "fechado": 100, "perdido": 0,
}

type PipelineHandler struct {
	repo   domain.PipelineRepository
	logger *zap.Logger
}

func NewPipelineHandler(repo domain.PipelineRepository, logger *zap.Logger) *PipelineHandler {
	return &PipelineHandler{repo: repo, logger: logger}
}

// POST /pipeline
func (h *PipelineHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNegocioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	tenantID := middleware.GetTenantID(r.Context())
	prob := etapaProbabilidade[req.Etapa]

	negocio := &domain.Negocio{
		TenantID: tenantID, LeadID: req.LeadID,
		LeadNome: req.LeadNome, LeadEmail: req.LeadEmail, LeadTelefone: req.LeadTelefone,
		LeadPrazo: req.LeadPrazo, LeadFormaPagamento: req.LeadFormaPagamento,
		Etapa: req.Etapa, Prioridade: req.Prioridade, ValorNegocio: req.ValorNegocio,
		ProbabilidadeFechamento: prob, ProximaAcao: req.ProximaAcao,
		UltimaAnotacao: req.UltimaAnotacao, CorretorResponsavel: req.CorretorResponsavel,
		Tags: req.Tags,
	}

	if err := h.repo.Create(r.Context(), negocio); err != nil {
		h.logger.Error("failed to create negocio", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.JSON(w, dto.ToNegocioResponse(negocio), http.StatusCreated)
}

// GET /pipeline
func (h *PipelineHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	negocios, err := h.repo.List(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("failed to list pipeline", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result := make([]dto.NegocioResponse, len(negocios))
	for i, n := range negocios {
		result[i] = dto.ToNegocioResponse(&n)
	}
	response.JSON(w, result, http.StatusOK)
}

// PUT /pipeline/{id}
func (h *PipelineHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateNegocioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	negocio := &domain.Negocio{
		ID: id, Etapa: req.Etapa, Prioridade: req.Prioridade,
		ValorNegocio: req.ValorNegocio, ProbabilidadeFechamento: req.ProbabilidadeFechamento,
		ProximaAcao: req.ProximaAcao, UltimaAnotacao: req.UltimaAnotacao,
		CorretorResponsavel: req.CorretorResponsavel, Tags: req.Tags,
		MotivoPerda: req.MotivoPerda,
	}

	if err := h.repo.Update(r.Context(), negocio); err != nil {
		h.logger.Error("failed to update negocio", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	updated, _ := h.repo.FindByID(r.Context(), id)
	if updated != nil {
		response.JSON(w, dto.ToNegocioResponse(updated), http.StatusOK)
	} else {
		response.Message(w, "updated", http.StatusOK)
	}
}

// PATCH /pipeline/{id}/move
func (h *PipelineHandler) Move(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.MoveNegocioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	prob := etapaProbabilidade[req.Etapa]

	if err := h.repo.UpdateEtapa(r.Context(), id, req.Etapa, prob); err != nil {
		h.logger.Error("failed to move negocio", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	updated, _ := h.repo.FindByID(r.Context(), id)
	if updated != nil {
		response.JSON(w, dto.ToNegocioResponse(updated), http.StatusOK)
	} else {
		response.Message(w, "moved", http.StatusOK)
	}
}

// DELETE /pipeline/{id}
func (h *PipelineHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete negocio", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.Message(w, "deleted", http.StatusOK)
}
