package dto

import (
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
)

type CreateNegocioRequest struct {
	LeadID              string   `json:"lead_id" validate:"required"`
	LeadNome            string   `json:"lead_nome" validate:"required"`
	LeadEmail           string   `json:"lead_email"`
	LeadTelefone        string   `json:"lead_telefone"`
	LeadPrazo           string   `json:"lead_prazo"`
	LeadFormaPagamento  string   `json:"lead_forma_pagamento"`
	Etapa               string   `json:"etapa" validate:"required,oneof=primeiro_contato qualificado visita_agendada proposta_enviada negociacao fechado perdido"`
	Prioridade          string   `json:"prioridade" validate:"required,oneof=baixa media alta urgente"`
	ValorNegocio        float64  `json:"valor_negocio"`
	ProximaAcao         string   `json:"proxima_acao"`
	UltimaAnotacao      string   `json:"ultima_anotacao"`
	CorretorResponsavel string   `json:"corretor_responsavel"`
	Tags                []string `json:"tags"`
}

type UpdateNegocioRequest struct {
	Etapa                   string   `json:"etapa" validate:"required,oneof=primeiro_contato qualificado visita_agendada proposta_enviada negociacao fechado perdido"`
	Prioridade              string   `json:"prioridade" validate:"required,oneof=baixa media alta urgente"`
	ValorNegocio            float64  `json:"valor_negocio"`
	ProbabilidadeFechamento int      `json:"probabilidade_fechamento"`
	ProximaAcao             string   `json:"proxima_acao"`
	UltimaAnotacao          string   `json:"ultima_anotacao"`
	CorretorResponsavel     string   `json:"corretor_responsavel"`
	Tags                    []string `json:"tags"`
	MotivoPerda             string   `json:"motivo_perda"`
}

type MoveNegocioRequest struct {
	Etapa string `json:"etapa" validate:"required,oneof=primeiro_contato qualificado visita_agendada proposta_enviada negociacao fechado perdido"`
}

type NegocioResponse struct {
	ID                      string   `json:"id"`
	LeadID                  string   `json:"lead_id"`
	LeadNome                string   `json:"lead_nome"`
	LeadEmail               string   `json:"lead_email"`
	LeadTelefone            string   `json:"lead_telefone"`
	LeadPrazo               string   `json:"lead_prazo,omitempty"`
	LeadFormaPagamento      string   `json:"lead_forma_pagamento,omitempty"`
	Etapa                   string   `json:"etapa"`
	Prioridade              string   `json:"prioridade"`
	ValorNegocio            float64  `json:"valor_negocio"`
	ProbabilidadeFechamento int      `json:"probabilidade_fechamento"`
	DataCriacao             string   `json:"data_criacao"`
	DataUltimaInteracao     string   `json:"data_ultima_interacao"`
	DataMovimentacao        string   `json:"data_movimentacao"`
	DiasNaEtapa             int      `json:"dias_na_etapa"`
	ProximaAcao             string   `json:"proxima_acao,omitempty"`
	UltimaAnotacao          string   `json:"ultima_anotacao,omitempty"`
	CorretorResponsavel     string   `json:"corretor_responsavel,omitempty"`
	Tags                    []string `json:"tags"`
	MotivoPerda             string   `json:"motivo_perda,omitempty"`
}

func ToNegocioResponse(n *domain.Negocio) NegocioResponse {
	tags := n.Tags
	if tags == nil {
		tags = []string{}
	}
	return NegocioResponse{
		ID: n.ID, LeadID: n.LeadID,
		LeadNome: n.LeadNome, LeadEmail: n.LeadEmail, LeadTelefone: n.LeadTelefone,
		LeadPrazo: n.LeadPrazo, LeadFormaPagamento: n.LeadFormaPagamento,
		Etapa: n.Etapa, Prioridade: n.Prioridade, ValorNegocio: n.ValorNegocio,
		ProbabilidadeFechamento: n.ProbabilidadeFechamento,
		DataCriacao:         n.DataCriacao.Format(time.RFC3339),
		DataUltimaInteracao: n.DataUltimaInteracao.Format(time.RFC3339),
		DataMovimentacao:    n.DataMovimentacao.Format(time.RFC3339),
		DiasNaEtapa: n.DiasNaEtapa, ProximaAcao: n.ProximaAcao,
		UltimaAnotacao: n.UltimaAnotacao, CorretorResponsavel: n.CorretorResponsavel,
		Tags: tags, MotivoPerda: n.MotivoPerda,
	}
}
