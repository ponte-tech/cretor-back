package dto

import (
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
)

// --- Requests ---

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required"`
}

type SignupRequest struct {
	Nome     string `json:"nome" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Senha    string `json:"senha" validate:"required,min=8,max=50"`
	Telefone string `json:"telefone" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin gerente corretor"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// --- Responses ---

type UsuarioResponse struct {
	ID        string  `json:"id"`
	Nome      string  `json:"nome"`
	Email     string  `json:"email"`
	Telefone  string  `json:"telefone"`
	Foto      *string `json:"foto,omitempty"`
	Role      string  `json:"role"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type AuthResponse struct {
	Usuario      UsuarioResponse `json:"usuario"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	ExpiresIn    int64           `json:"expires_in"`
}

// --- Mappers ---

func ToUsuarioResponse(u *domain.Usuario) UsuarioResponse {
	return UsuarioResponse{
		ID:        u.ID,
		Nome:      u.Nome,
		Email:     u.Email,
		Telefone:  u.Telefone,
		Foto:      u.Foto,
		Role:      string(u.Role),
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAuthResponse(u *domain.Usuario, accessToken, refreshToken string, expiresIn int64) AuthResponse {
	return AuthResponse{
		Usuario:      ToUsuarioResponse(u),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}
