package dto

import (
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
)

// --- Responses ---

type EnderecoResponse struct {
	Completo     string  `json:"completo"`
	Logradouro   string  `json:"logradouro"`
	Numero       string  `json:"numero"`
	Bairro       string  `json:"bairro"`
	Cidade       string  `json:"cidade"`
	UF           string  `json:"uf"`
	Regiao       string  `json:"regiao"`
	DistanciaMar string  `json:"distancia_mar"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

type ObraResponse struct {
	Status        string `json:"status"`
	Progresso     int    `json:"progresso"`
	Inicio        string `json:"inicio,omitempty"`
	Entrega       string `json:"entrega,omitempty"`
	AnoLancamento int    `json:"ano_lancamento"`
}

type EdificacaoResponse struct {
	Pavimentos       int    `json:"pavimentos"`
	UnidadesPorAndar int    `json:"unidades_por_andar"`
	TotalUnidades    int    `json:"total_unidades"`
	Coberturas       int    `json:"coberturas"`
	AlturaTorre      string `json:"altura_torre"`
	ProntoParaMorar  bool   `json:"pronto_para_morar"`
	Mobiliado        bool   `json:"mobiliado"`
}

type FaixaIntResponse struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type FaixaFloatResponse struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type CaracteristicasResponse struct {
	Dormitorios FaixaIntResponse   `json:"dormitorios"`
	Suites      FaixaIntResponse   `json:"suites"`
	Vagas       FaixaIntResponse   `json:"vagas"`
	Metragem    FaixaFloatResponse `json:"metragem"`
	Condominio  float64            `json:"condominio"`
	IPTU        float64            `json:"iptu"`
}

type PontoInteresseResponse struct {
	Nome      string `json:"nome"`
	Distancia string `json:"distancia"`
	Tempo     string `json:"tempo"`
}

type SecaoResponse struct {
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
	Ordem     int    `json:"ordem"`
}

type GaleriaImagemResponse struct {
	S3Path string `json:"s3_path"`
	Ordem  int    `json:"ordem"`
	Tipo   string `json:"tipo"`
}

type ComercializacaoResponse struct {
	UnidadesDisponiveis int    `json:"unidades_disponiveis"`
	UltimaAtualizacao   string `json:"ultima_atualizacao,omitempty"`
}

type ConstrutoraResponse struct {
	ID         string `json:"id"`
	Nome       string `json:"nome"`
	LogoS3Path string `json:"logo_s3_path"`
	Website    string `json:"website"`
	Descricao  string `json:"descricao"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type EmpreendimentoResponse struct {
	ID                     string                   `json:"id"`
	ConstrutoraID          string                   `json:"construtora_id"`
	Construtora            *ConstrutoraResponse      `json:"construtora,omitempty"`
	Nome                   string                   `json:"nome"`
	Subtitulo              string                   `json:"subtitulo"`
	Slogan                 string                   `json:"slogan"`
	Descricao              string                   `json:"descricao"`
	Endereco               EnderecoResponse         `json:"endereco"`
	Obra                   ObraResponse             `json:"obra"`
	Incorporacao           string                   `json:"incorporacao"`
	Edificacao             EdificacaoResponse       `json:"edificacao"`
	Caracteristicas        CaracteristicasResponse  `json:"caracteristicas"`
	DiferenciaisUnidade    []string                 `json:"diferenciais_unidade"`
	DiferenciaisCondominio []string                 `json:"diferenciais_condominio"`
	PontosInteresse        []PontoInteresseResponse `json:"pontos_interesse"`
	Secoes                 []SecaoResponse          `json:"secoes"`
	Galeria                []GaleriaImagemResponse  `json:"galeria"`
	CatalogoS3Path         string                   `json:"catalogo_s3_path"`
	Comercializacao        ComercializacaoResponse  `json:"comercializacao"`
	CreatedAt              string                   `json:"created_at"`
}

// PaginatedEmpreendimentosResponse wraps card responses with pagination metadata.
type PaginatedEmpreendimentosResponse struct {
	Data       []EmpreendimentoCardResponse `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	TotalPages int                          `json:"total_pages"`
}

// Card response para listagem (dados reduzidos)
type EmpreendimentoCardResponse struct {
	ID              string                  `json:"id"`
	ConstrutoraID   string                  `json:"construtora_id"`
	ConstrutoraNome string                  `json:"construtora_nome"`
	Nome            string                  `json:"nome"`
	Endereco        EnderecoResponse        `json:"endereco"`
	Obra            ObraResponse            `json:"obra"`
	Caracteristicas CaracteristicasResponse `json:"caracteristicas"`
	Comercializacao ComercializacaoResponse `json:"comercializacao"`
	Foto            string                  `json:"foto"`
}

type UnidadeResponse struct {
	ID               string   `json:"id"`
	EmpreendimentoID string   `json:"empreendimento_id"`
	Secao            string   `json:"secao"`
	Numero           string   `json:"numero"`
	Andar            int      `json:"andar"`
	Metragem         float64  `json:"metragem"`
	Tipo             string   `json:"tipo"`
	Valor            *float64 `json:"valor"`
	ValorTexto       string   `json:"valor_texto"`
	Status           string   `json:"status"`
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func ToConstrutoraResponse(c *domain.Construtora) ConstrutoraResponse {
	return ConstrutoraResponse{
		ID: c.ID, Nome: c.Nome, LogoS3Path: c.LogoS3Path,
		Website: c.Website, Descricao: c.Descricao, Status: c.Status,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}

func ToEmpreendimentoResponse(e *domain.Empreendimento) EmpreendimentoResponse {
	pontos := make([]PontoInteresseResponse, len(e.PontosInteresse))
	for i, p := range e.PontosInteresse {
		pontos[i] = PontoInteresseResponse{Nome: p.Nome, Distancia: p.Distancia, Tempo: p.Tempo}
	}
	secoes := make([]SecaoResponse, len(e.Secoes))
	for i, s := range e.Secoes {
		secoes[i] = SecaoResponse{Nome: s.Nome, Categoria: s.Categoria, Ordem: s.Ordem}
	}
	galeria := make([]GaleriaImagemResponse, len(e.Galeria))
	for i, g := range e.Galeria {
		galeria[i] = GaleriaImagemResponse{S3Path: g.S3Path, Ordem: g.Ordem, Tipo: g.Tipo}
	}

	return EmpreendimentoResponse{
		ID: e.ID, ConstrutoraID: e.ConstrutoraID,
		Nome: e.Nome, Subtitulo: e.Subtitulo, Slogan: e.Slogan, Descricao: e.Descricao,
		Endereco: EnderecoResponse{
			Completo: e.Endereco.Completo, Logradouro: e.Endereco.Logradouro,
			Numero: e.Endereco.Numero, Bairro: e.Endereco.Bairro,
			Cidade: e.Endereco.Cidade, UF: e.Endereco.UF,
			Regiao: e.Endereco.Regiao, DistanciaMar: e.Endereco.DistanciaMar,
			Latitude: e.Endereco.Latitude, Longitude: e.Endereco.Longitude,
		},
		Obra: ObraResponse{
			Status: e.Obra.Status, Progresso: e.Obra.Progresso,
			Inicio: formatTime(e.Obra.Inicio), Entrega: formatTime(e.Obra.Entrega),
			AnoLancamento: e.Obra.AnoLancamento,
		},
		Incorporacao: e.Incorporacao,
		Edificacao: EdificacaoResponse{
			Pavimentos: e.Edificacao.Pavimentos, UnidadesPorAndar: e.Edificacao.UnidadesPorAndar,
			TotalUnidades: e.Edificacao.TotalUnidades, Coberturas: e.Edificacao.Coberturas,
			AlturaTorre: e.Edificacao.AlturaTorre, ProntoParaMorar: e.Edificacao.ProntoParaMorar,
			Mobiliado: e.Edificacao.Mobiliado,
		},
		Caracteristicas: CaracteristicasResponse{
			Dormitorios: FaixaIntResponse{Min: e.Caracteristicas.Dormitorios.Min, Max: e.Caracteristicas.Dormitorios.Max},
			Suites:      FaixaIntResponse{Min: e.Caracteristicas.Suites.Min, Max: e.Caracteristicas.Suites.Max},
			Vagas:       FaixaIntResponse{Min: e.Caracteristicas.Vagas.Min, Max: e.Caracteristicas.Vagas.Max},
			Metragem:    FaixaFloatResponse{Min: e.Caracteristicas.Metragem.Min, Max: e.Caracteristicas.Metragem.Max},
			Condominio:  e.Caracteristicas.Condominio,
			IPTU:        e.Caracteristicas.IPTU,
		},
		DiferenciaisUnidade:    e.DiferenciaisUnidade,
		DiferenciaisCondominio: e.DiferenciaisCondominio,
		PontosInteresse:        pontos,
		Secoes:                 secoes,
		Galeria:                galeria,
		CatalogoS3Path:         e.CatalogoS3Path,
		Comercializacao: ComercializacaoResponse{
			UnidadesDisponiveis: e.Comercializacao.UnidadesDisponiveis,
			UltimaAtualizacao:   formatTime(e.Comercializacao.UltimaAtualizacao),
		},
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
	}
}

func ToEmpreendimentoCardResponse(e *domain.Empreendimento, construtoraNome string) EmpreendimentoCardResponse {
	foto := e.Foto
	return EmpreendimentoCardResponse{
		ID: e.ID, ConstrutoraID: e.ConstrutoraID,
		ConstrutoraNome: construtoraNome,
		Nome:            e.Nome,
		Endereco: EnderecoResponse{
			Completo: e.Endereco.Completo, Bairro: e.Endereco.Bairro,
			Cidade: e.Endereco.Cidade, UF: e.Endereco.UF,
			Latitude: e.Endereco.Latitude, Longitude: e.Endereco.Longitude,
		},
		Obra: ObraResponse{
			Status: e.Obra.Status, Progresso: e.Obra.Progresso,
			Entrega: formatTime(e.Obra.Entrega),
		},
		Caracteristicas: CaracteristicasResponse{
			Dormitorios: FaixaIntResponse{Min: e.Caracteristicas.Dormitorios.Min, Max: e.Caracteristicas.Dormitorios.Max},
			Suites:      FaixaIntResponse{Min: e.Caracteristicas.Suites.Min, Max: e.Caracteristicas.Suites.Max},
			Vagas:       FaixaIntResponse{Min: e.Caracteristicas.Vagas.Min, Max: e.Caracteristicas.Vagas.Max},
			Metragem:    FaixaFloatResponse{Min: e.Caracteristicas.Metragem.Min, Max: e.Caracteristicas.Metragem.Max},
		},
		Comercializacao: ComercializacaoResponse{
			UnidadesDisponiveis: e.Comercializacao.UnidadesDisponiveis,
		},
		Foto: foto,
	}
}

func ToUnidadeResponse(u *domain.Unidade) UnidadeResponse {
	return UnidadeResponse{
		ID: u.ID, EmpreendimentoID: u.EmpreendimentoID,
		Secao: u.Secao, Numero: u.Numero, Andar: u.Andar,
		Metragem: u.Metragem, Tipo: u.Tipo,
		Valor: u.Valor, ValorTexto: u.ValorTexto, Status: u.Status,
	}
}
