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

type LeadHandler struct {
	repo         domain.LeadRepository
	pipelineRepo domain.PipelineRepository
	logger       *zap.Logger
}

func NewLeadHandler(repo domain.LeadRepository, pipelineRepo domain.PipelineRepository, logger *zap.Logger) *LeadHandler {
	return &LeadHandler{repo: repo, pipelineRepo: pipelineRepo, logger: logger}
}

// POST /leads
func (h *LeadHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Honeypot: if bot filled the hidden field, silently accept but don't save
	if req.Website != "" {
		h.logger.Warn("honeypot triggered", zap.String("ip", r.Header.Get("X-Forwarded-For")))
		response.JSON(w, dto.LeadResponse{ID: "ok"}, http.StatusCreated)
		return
	}

	tenantID := middleware.GetTenantID(r.Context())
	lead := &domain.Lead{
		TenantID: tenantID, Nome: req.Nome, Whatsapp: req.Whatsapp,
		Email: req.Email, Prazo: req.Prazo, FormaPagamento: req.FormaPagamento,
		Origem: "landing_page_alto_padrao", Status: "novo",
	}

	if err := h.repo.Create(r.Context(), lead); err != nil {
		h.logger.Error("failed to create lead", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Auto-create pipeline entry
	negocio := &domain.Negocio{
		TenantID:                tenantID,
		LeadID:                  lead.ID,
		LeadNome:                lead.Nome,
		LeadEmail:               lead.Email,
		LeadTelefone:            lead.Whatsapp,
		LeadPrazo:               lead.Prazo,
		LeadFormaPagamento:      lead.FormaPagamento,
		Etapa:                   "primeiro_contato",
		Prioridade:              "media",
		ProbabilidadeFechamento: 20,
		Tags:                    []string{lead.Origem},
	}
	if err := h.pipelineRepo.Create(r.Context(), negocio); err != nil {
		h.logger.Error("failed to auto-create pipeline entry", zap.Error(err))
	}

	h.logger.Info("lead created", zap.String("email", lead.Email))
	response.JSON(w, dto.ToLeadResponse(lead), http.StatusCreated)
}

// GET /leads
func (h *LeadHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	filter := domain.LeadFilter{}
	if s := r.URL.Query().Get("status"); s != "" {
		filter.Status = &s
	}
	if q := r.URL.Query().Get("q"); q != "" {
		filter.Search = &q
	}

	leads, err := h.repo.List(r.Context(), tenantID, filter)
	if err != nil {
		h.logger.Error("failed to list leads", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result := make([]dto.LeadResponse, len(leads))
	for i, l := range leads {
		result[i] = dto.ToLeadResponse(&l)
	}

	response.JSON(w, result, http.StatusOK)
}

// GET /leads/{id}
func (h *LeadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	lead, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find lead", zap.Error(err))
		response.Error(w, "not found", http.StatusNotFound)
		return
	}

	response.JSON(w, dto.ToLeadResponse(lead), http.StatusOK)
}

// PUT /leads/{id}
func (h *LeadHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	lead := &domain.Lead{
		ID: id, Nome: req.Nome, Whatsapp: req.Whatsapp,
		Email: req.Email, Prazo: req.Prazo, FormaPagamento: req.FormaPagamento,
		Status: req.Status,
	}

	if err := h.repo.Update(r.Context(), lead); err != nil {
		h.logger.Error("failed to update lead", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch updated lead to return
	updated, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		response.JSON(w, dto.ToLeadResponse(lead), http.StatusOK)
		return
	}

	response.JSON(w, dto.ToLeadResponse(updated), http.StatusOK)
}

// DELETE /leads/{id}
func (h *LeadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete lead", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Also delete associated pipeline entries
	if err := h.pipelineRepo.DeleteByLeadID(r.Context(), id); err != nil {
		h.logger.Error("failed to delete pipeline entries for lead", zap.Error(err))
	}

	h.logger.Info("lead deleted", zap.String("id", id))
	response.Message(w, "lead deleted", http.StatusOK)
}
