package domain

import "time"

type Endereco struct {
	Completo     string
	Logradouro   string
	Numero       string
	Bairro       string
	Cidade       string
	UF           string
	Regiao       string
	DistanciaMar string
	Latitude     float64
	Longitude    float64
}

type Obra struct {
	Status         string
	Progresso      int
	Inicio         *time.Time
	Entrega        *time.Time
	AnoLancamento  int
}

type Edificacao struct {
	Pavimentos       int
	UnidadesPorAndar int
	TotalUnidades    int
	Coberturas       int
	AlturaTorre      string
	ProntoParaMorar  bool
	Mobiliado        bool
}

type FaixaInt struct {
	Min int
	Max int
}

type FaixaFloat struct {
	Min float64
	Max float64
}

type Caracteristicas struct {
	Dormitorios FaixaInt
	Suites      FaixaInt
	Vagas       FaixaInt
	Metragem    FaixaFloat
	Condominio  float64
	IPTU        float64
}

type PontoInteresse struct {
	Nome      string
	Distancia string
	Tempo     string
}

type Secao struct {
	Nome      string
	Categoria string
	Ordem     int
}

type GaleriaImagem struct {
	S3Path string
	Ordem  int
	Tipo   string
}

type Comercializacao struct {
	UnidadesDisponiveis int
	UltimaAtualizacao   *time.Time
}

type Empreendimento struct {
	ID                     string
	TenantID               string
	ConstrutoraID          string
	Nome                   string
	Subtitulo              string
	Slogan                 string
	Descricao              string
	FonteID                string
	Endereco               Endereco
	Obra                   Obra
	Incorporacao           string
	Edificacao             Edificacao
	Caracteristicas        Caracteristicas
	DiferenciaisUnidade    []string
	DiferenciaisCondominio []string
	PontosInteresse        []PontoInteresse
	Secoes                 []Secao
	Galeria                []GaleriaImagem
	CatalogoS3Path         string
	Foto                   string
	Comercializacao        Comercializacao
	ConstrutoraNome        string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type Pagination struct {
	Page     int
	PageSize int
}

type BoundingBox struct {
	SwLat float64
	SwLng float64
	NeLat float64
	NeLng float64
}

type EmpreendimentoFilter struct {
	ConstrutoraID *string
	Cidade        *string
	UF            *string
	StatusObra    *string
	Search        *string
	Bounds        *BoundingBox
	Pagination    Pagination
}
