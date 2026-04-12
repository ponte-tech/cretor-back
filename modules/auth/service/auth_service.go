package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/shared/auth"
	"go.uber.org/zap"
)

type AuthService struct {
	repo   domain.UsuarioRepository
	jwt    *auth.JWTProvider
	logger *zap.Logger
}

func NewAuthService(repo domain.UsuarioRepository, jwtProvider *auth.JWTProvider, logger *zap.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		jwt:    jwtProvider,
		logger: logger,
	}
}

type SignupInput struct {
	TenantID string
	Nome     string
	Email    string
	Senha    string
	Telefone string
	Role     string
}

type LoginInput struct {
	TenantID string
	Email    string
	Senha    string
}

type AuthResult struct {
	Usuario      *domain.Usuario
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func (s *AuthService) Signup(ctx context.Context, input SignupInput) (*AuthResult, error) {
	role := domain.Role(input.Role)
	if !role.Valid() {
		return nil, fmt.Errorf("role %s: %w", input.Role, domain.ErrInvalidEntity)
	}

	existing, err := s.repo.FindByEmail(ctx, input.TenantID, input.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrEmailExists
	}

	hash, err := auth.HashPassword(input.Senha)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	usuario := &domain.Usuario{
		TenantID:  input.TenantID,
		Nome:      input.Nome,
		Email:     input.Email,
		SenhaHash: hash,
		Telefone:  input.Telefone,
		Role:      role,
		Status:    domain.StatusAtivo,
	}

	if err := s.repo.Create(ctx, usuario); err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return nil, domain.ErrEmailExists
		}
		return nil, fmt.Errorf("create usuario: %w", err)
	}

	accessToken, expiresIn, err := s.jwt.GenerateAccessToken(usuario.ID, usuario.TenantID, usuario.Email, string(usuario.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(usuario.ID, usuario.TenantID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.repo.SaveRefreshToken(ctx, usuario.ID, refreshToken); err != nil {
		s.logger.Error("failed to save refresh token", zap.Error(err))
	}

	s.logger.Info("user signed up", zap.String("email", usuario.Email), zap.String("tenant", usuario.TenantID))

	return &AuthResult{
		Usuario:      usuario,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	usuario, err := s.repo.FindByEmail(ctx, input.TenantID, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, fmt.Errorf("find usuario: %w", err)
	}

	if usuario.Status == domain.StatusInativo {
		return nil, domain.ErrInactiveUser
	}

	if !auth.ComparePassword(usuario.SenhaHash, input.Senha) {
		return nil, domain.ErrUnauthorized
	}

	accessToken, expiresIn, err := s.jwt.GenerateAccessToken(usuario.ID, usuario.TenantID, usuario.Email, string(usuario.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(usuario.ID, usuario.TenantID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.repo.UpdateLastLogin(ctx, usuario.ID); err != nil {
		s.logger.Error("failed to update last login", zap.Error(err))
	}
	if err := s.repo.SaveRefreshToken(ctx, usuario.ID, refreshToken); err != nil {
		s.logger.Error("failed to save refresh token", zap.Error(err))
	}

	now := time.Now().UTC()
	usuario.UltimoLogin = &now

	s.logger.Info("user logged in", zap.String("email", usuario.Email), zap.String("tenant", usuario.TenantID))

	return &AuthResult{
		Usuario:      usuario,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	claims, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	usuario, err := s.repo.FindByID(ctx, claims.UsuarioID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if usuario.Status == domain.StatusInativo {
		return nil, domain.ErrInactiveUser
	}

	if usuario.RefreshToken == nil || *usuario.RefreshToken != refreshToken {
		return nil, domain.ErrUnauthorized
	}

	newAccessToken, expiresIn, err := s.jwt.GenerateAccessToken(usuario.ID, usuario.TenantID, usuario.Email, string(usuario.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(usuario.ID, usuario.TenantID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.repo.SaveRefreshToken(ctx, usuario.ID, newRefreshToken); err != nil {
		s.logger.Error("failed to save refresh token", zap.Error(err))
	}

	return &AuthResult{
		Usuario:      usuario,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, usuarioID string) error {
	if err := s.repo.ClearRefreshToken(ctx, usuarioID); err != nil {
		return fmt.Errorf("clear refresh token: %w", err)
	}

	s.logger.Info("user logged out", zap.String("usuario_id", usuarioID))
	return nil
}
