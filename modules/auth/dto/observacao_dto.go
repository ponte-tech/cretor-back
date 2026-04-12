package dto

import (
	"github.com/ponte-tech/cretor-back/modules/auth/domain"
)

type CreateObservacaoRequest struct {
	Texto string `json:"texto" validate:"required,min=1,max=2000"`
	Autor string `json:"autor"`
}

type ObservacaoResponse struct {
	ID        string `json:"id"`
	NegocioID string `json:"negocio_id"`
	Texto     string `json:"texto"`
	Autor     string `json:"autor,omitempty"`
	CreatedAt string `json:"created_at"`
}

func ToObservacaoResponse(o *domain.Observacao) ObservacaoResponse {
	return ObservacaoResponse{
		ID:        o.ID,
		NegocioID: o.NegocioID,
		Texto:     o.Texto,
		Autor:     o.Autor,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
