package domain

import "time"

type Construtora struct {
	ID         string
	TenantID   string
	Nome       string
	LogoS3Path string
	Website    string
	Descricao  string
	FonteID    string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ConstrutoraFilter struct {
	Status *string
	Search *string
}
