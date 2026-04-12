package domain

import (
	"context"
	"time"
)

type Negocio struct {
	ID                       string
	TenantID                 string
	LeadID                   string
	LeadNome                 string
	LeadEmail                string
	LeadTelefone             string
	LeadPrazo                string
	LeadFormaPagamento       string
	Etapa                    string
	Prioridade               string
	ValorNegocio             float64
	ProbabilidadeFechamento  int
	DataCriacao              time.Time
	DataUltimaInteracao      time.Time
	DataMovimentacao         time.Time
	DiasNaEtapa              int
	ProximaAcao              string
	UltimaAnotacao           string
	CorretorResponsavel      string
	Tags                     []string
	MotivoPerda              string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type PipelineRepository interface {
	Create(ctx context.Context, negocio *Negocio) error
	FindByID(ctx context.Context, id string) (*Negocio, error)
	List(ctx context.Context, tenantID string) ([]Negocio, error)
	Update(ctx context.Context, negocio *Negocio) error
	UpdateEtapa(ctx context.Context, id, etapa string, probabilidade int) error
	Delete(ctx context.Context, id string) error
}
