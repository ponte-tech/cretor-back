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

type ConstrutoraRepository interface {
	Create(ctx context.Context, c *Construtora) error
	FindByID(ctx context.Context, id string) (*Construtora, error)
	List(ctx context.Context, tenantID string, filter ConstrutoraFilter) ([]Construtora, error)
	Update(ctx context.Context, c *Construtora) error
	Delete(ctx context.Context, id string) error
}

type EmpreendimentoRepository interface {
	Create(ctx context.Context, e *Empreendimento) error
	FindByID(ctx context.Context, id string) (*Empreendimento, error)
	List(ctx context.Context, tenantID string, filter EmpreendimentoFilter) ([]Empreendimento, int64, error)
	Search(ctx context.Context, tenantID, query string, page, pageSize int) ([]Empreendimento, int64, error)
	Update(ctx context.Context, e *Empreendimento) error
	Delete(ctx context.Context, id string) error
}

type UnidadeRepository interface {
	Create(ctx context.Context, u *Unidade) error
	CreateMany(ctx context.Context, unidades []Unidade) error
	FindByID(ctx context.Context, id string) (*Unidade, error)
	List(ctx context.Context, tenantID string, filter UnidadeFilter) ([]Unidade, error)
	Update(ctx context.Context, u *Unidade) error
	Delete(ctx context.Context, id string) error
	DeleteByEmpreendimento(ctx context.Context, empreendimentoID string) error
}
