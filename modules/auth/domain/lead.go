package domain

import (
	"context"
	"time"
)

type Lead struct {
	ID              string
	TenantID        string
	Nome            string
	Whatsapp        string
	Email           string
	Prazo           string
	FormaPagamento  string
	Origem          string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type LeadFilter struct {
	Status *string
	Search *string
}

type LeadRepository interface {
	Create(ctx context.Context, lead *Lead) error
	FindByID(ctx context.Context, id string) (*Lead, error)
	List(ctx context.Context, tenantID string, filter LeadFilter) ([]Lead, error)
	Update(ctx context.Context, lead *Lead) error
	Delete(ctx context.Context, id string) error
}
