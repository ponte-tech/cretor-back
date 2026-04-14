package domain

import "time"

type Unidade struct {
	ID               string
	TenantID         string
	EmpreendimentoID string
	Secao            string
	Numero           string
	Andar            int
	Metragem         float64
	Tipo             string
	Valor            *float64
	ValorTexto       string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UnidadeFilter struct {
	EmpreendimentoID string
	Secao            *string
	Status           *string
	ValorMin         *float64
	ValorMax         *float64
}
