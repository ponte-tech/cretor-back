package domain

import (
	"context"
	"time"
)

type Observacao struct {
	ID         string
	TenantID   string
	NegocioID  string
	Texto      string
	Autor      string
	CreatedAt  time.Time
}

type ObservacaoRepository interface {
	Create(ctx context.Context, obs *Observacao) error
	ListByNegocio(ctx context.Context, negocioID string) ([]Observacao, error)
}
