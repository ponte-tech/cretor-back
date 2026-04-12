package domain

import "context"

type UsuarioRepository interface {
	FindByID(ctx context.Context, id string) (*Usuario, error)
	FindByEmail(ctx context.Context, tenantID, email string) (*Usuario, error)
	Create(ctx context.Context, usuario *Usuario) error
	Update(ctx context.Context, usuario *Usuario) error
	UpdateLastLogin(ctx context.Context, id string) error
	SaveRefreshToken(ctx context.Context, id, token string) error
	ClearRefreshToken(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
