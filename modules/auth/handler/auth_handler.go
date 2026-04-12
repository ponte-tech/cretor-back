package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/dto"
	"github.com/ponte-tech/cretor-back/modules/auth/service"
	"github.com/ponte-tech/cretor-back/shared/middleware"
	"github.com/ponte-tech/cretor-back/shared/response"
	"github.com/ponte-tech/cretor-back/shared/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service *service.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(svc *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: svc,
		logger:  logger,
	}
}

// POST /auth/signup
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req dto.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	tenantID := middleware.GetTenantID(r.Context())

	result, err := h.service.Signup(r.Context(), service.SignupInput{
		TenantID: tenantID,
		Nome:     req.Nome,
		Email:    req.Email,
		Senha:    req.Senha,
		Telefone: req.Telefone,
		Role:     req.Role,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := dto.ToAuthResponse(result.Usuario, result.AccessToken, result.RefreshToken, result.ExpiresIn)
	response.JSON(w, resp, http.StatusCreated)
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	tenantID := middleware.GetTenantID(r.Context())

	result, err := h.service.Login(r.Context(), service.LoginInput{
		TenantID: tenantID,
		Email:    req.Email,
		Senha:    req.Senha,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := dto.ToAuthResponse(result.Usuario, result.AccessToken, result.RefreshToken, result.ExpiresIn)
	response.JSON(w, resp, http.StatusOK)
}

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	result, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := dto.ToAuthResponse(result.Usuario, result.AccessToken, result.RefreshToken, result.ExpiresIn)
	response.JSON(w, resp, http.StatusOK)
}

// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	usuarioID := middleware.GetUsuarioID(r.Context())

	if err := h.service.Logout(r.Context(), usuarioID); err != nil {
		h.handleError(w, err)
		return
	}

	response.Message(w, "logged out successfully", http.StatusOK)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailExists):
		response.Error(w, "email already exists", http.StatusConflict)
	case errors.Is(err, domain.ErrUnauthorized):
		response.Error(w, "invalid email or password", http.StatusUnauthorized)
	case errors.Is(err, domain.ErrInactiveUser):
		response.Error(w, "user is inactive", http.StatusForbidden)
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidEntity):
		response.Error(w, "invalid data", http.StatusBadRequest)
	case errors.Is(err, domain.ErrInvalidID):
		response.Error(w, "invalid id format", http.StatusBadRequest)
	default:
		h.logger.Error("unhandled error", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
