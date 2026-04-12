package domain

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleGerente  Role = "gerente"
	RoleCorretor Role = "corretor"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleGerente, RoleCorretor:
		return true
	}
	return false
}

type Status string

const (
	StatusAtivo   Status = "ativo"
	StatusInativo Status = "inativo"
)

type Usuario struct {
	ID           string
	TenantID     string
	Nome         string
	Email        string
	SenhaHash    string
	Telefone     string
	Foto         *string
	Role         Role
	Status       Status
	UltimoLogin  *time.Time
	RefreshToken *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
