package dto

import (
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
)

type CreateLeadRequest struct {
	Nome           string `json:"nome" validate:"required,min=2"`
	Whatsapp       string `json:"whatsapp" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Prazo          string `json:"prazo" validate:"required"`
	FormaPagamento string `json:"forma_pagamento" validate:"required"`
}

type UpdateLeadRequest struct {
	Nome           string `json:"nome" validate:"required,min=2"`
	Whatsapp       string `json:"whatsapp" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Prazo          string `json:"prazo" validate:"required"`
	FormaPagamento string `json:"forma_pagamento" validate:"required"`
	Status         string `json:"status" validate:"required,oneof=novo contatado qualificado convertido descartado"`
}

type LeadResponse struct {
	ID             string `json:"id"`
	Nome           string `json:"nome"`
	Whatsapp       string `json:"whatsapp"`
	Email          string `json:"email"`
	Prazo          string `json:"prazo"`
	FormaPagamento string `json:"forma_pagamento"`
	Origem         string `json:"origem"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

func ToLeadResponse(l *domain.Lead) LeadResponse {
	return LeadResponse{
		ID:             l.ID,
		Nome:           l.Nome,
		Whatsapp:       l.Whatsapp,
		Email:          l.Email,
		Prazo:          l.Prazo,
		FormaPagamento: l.FormaPagamento,
		Origem:         l.Origem,
		Status:         l.Status,
		CreatedAt:      l.CreatedAt.Format(time.RFC3339),
	}
}
