package repository

import (
	"context"
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/shared/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type enderecoDoc struct {
	Completo     string  `bson:"completo"`
	Logradouro   string  `bson:"logradouro"`
	Numero       string  `bson:"numero"`
	Bairro       string  `bson:"bairro"`
	Cidade       string  `bson:"cidade"`
	UF           string  `bson:"uf"`
	Regiao       string  `bson:"regiao"`
	DistanciaMar string  `bson:"distancia_mar"`
	Latitude     float64 `bson:"latitude"`
	Longitude    float64 `bson:"longitude"`
}

type obraDoc struct {
	Status        string     `bson:"status"`
	Progresso     int        `bson:"progresso"`
	Inicio        *time.Time `bson:"inicio,omitempty"`
	Entrega       *time.Time `bson:"entrega,omitempty"`
	AnoLancamento int        `bson:"ano_lancamento"`
}

type edificacaoDoc struct {
	Pavimentos       int    `bson:"pavimentos"`
	UnidadesPorAndar int    `bson:"unidades_por_andar"`
	TotalUnidades    int    `bson:"total_unidades"`
	Coberturas       int    `bson:"coberturas"`
	AlturaTorre      string `bson:"altura_torre"`
	ProntoParaMorar  bool   `bson:"pronto_para_morar"`
	Mobiliado        bool   `bson:"mobiliado"`
}

type faixaIntDoc struct {
	Min int `bson:"min"`
	Max int `bson:"max"`
}

type faixaFloatDoc struct {
	Min float64 `bson:"min"`
	Max float64 `bson:"max"`
}

type caracteristicasDoc struct {
	Dormitorios faixaIntDoc   `bson:"dormitorios"`
	Suites      faixaIntDoc   `bson:"suites"`
	Vagas       faixaIntDoc   `bson:"vagas"`
	Metragem    faixaFloatDoc `bson:"metragem"`
	Condominio  float64       `bson:"condominio"`
	IPTU        float64       `bson:"iptu"`
}

type pontoInteresseDoc struct {
	Nome      string `bson:"nome"`
	Distancia string `bson:"distancia"`
	Tempo     string `bson:"tempo"`
}

type secaoDoc struct {
	Nome      string `bson:"nome"`
	Categoria string `bson:"categoria"`
	Ordem     int    `bson:"ordem"`
}

type galeriaImagemDoc struct {
	S3Path string `bson:"s3_path"`
	Ordem  int    `bson:"ordem"`
	Tipo   string `bson:"tipo"`
}

type comercializacaoDoc struct {
	UnidadesDisponiveis int        `bson:"unidades_disponiveis"`
	UltimaAtualizacao   *time.Time `bson:"ultima_atualizacao,omitempty"`
}

type empreendimentoDocument struct {
	ID                     bson.ObjectID        `bson:"_id,omitempty"`
	TenantID               bson.ObjectID        `bson:"tenant_id"`
	ConstrutoraID          bson.ObjectID        `bson:"construtora_id"`
	Nome                   string               `bson:"nome"`
	Subtitulo              string               `bson:"subtitulo"`
	Slogan                 string               `bson:"slogan"`
	Descricao              string               `bson:"descricao"`
	FonteID                string               `bson:"fonte_id"`
	Endereco               enderecoDoc          `bson:"endereco"`
	Obra                   obraDoc              `bson:"obra"`
	Incorporacao           string               `bson:"incorporacao"`
	Edificacao             edificacaoDoc        `bson:"edificacao"`
	Caracteristicas        caracteristicasDoc   `bson:"caracteristicas"`
	DiferenciaisUnidade    []string             `bson:"diferenciais_unidade"`
	DiferenciaisCondominio []string             `bson:"diferenciais_condominio"`
	PontosInteresse        []pontoInteresseDoc  `bson:"pontos_interesse"`
	Secoes                 []secaoDoc           `bson:"secoes"`
	Galeria                []galeriaImagemDoc   `bson:"galeria"`
	CatalogoS3Path         string               `bson:"catalogo_s3_path"`
	Foto                   string               `bson:"foto"`
	Comercializacao        comercializacaoDoc   `bson:"comercializacao"`
	ConstrutoraNome        string               `bson:"construtora_nome"`
	CreatedAt              time.Time            `bson:"created_at"`
	UpdatedAt              time.Time            `bson:"updated_at"`
}

func (d *empreendimentoDocument) toDomain() domain.Empreendimento {
	pontosInteresse := make([]domain.PontoInteresse, len(d.PontosInteresse))
	for i, p := range d.PontosInteresse {
		pontosInteresse[i] = domain.PontoInteresse{Nome: p.Nome, Distancia: p.Distancia, Tempo: p.Tempo}
	}
	secoes := make([]domain.Secao, len(d.Secoes))
	for i, s := range d.Secoes {
		secoes[i] = domain.Secao{Nome: s.Nome, Categoria: s.Categoria, Ordem: s.Ordem}
	}
	galeria := make([]domain.GaleriaImagem, len(d.Galeria))
	for i, g := range d.Galeria {
		galeria[i] = domain.GaleriaImagem{S3Path: g.S3Path, Ordem: g.Ordem, Tipo: g.Tipo}
	}

	return domain.Empreendimento{
		ID: d.ID.Hex(), TenantID: d.TenantID.Hex(), ConstrutoraID: d.ConstrutoraID.Hex(),
		Nome: d.Nome, Subtitulo: d.Subtitulo, Slogan: d.Slogan,
		Descricao: d.Descricao, FonteID: d.FonteID,
		Endereco: domain.Endereco{
			Completo: d.Endereco.Completo, Logradouro: d.Endereco.Logradouro,
			Numero: d.Endereco.Numero, Bairro: d.Endereco.Bairro,
			Cidade: d.Endereco.Cidade, UF: d.Endereco.UF,
			Regiao: d.Endereco.Regiao, DistanciaMar: d.Endereco.DistanciaMar,
			Latitude: d.Endereco.Latitude, Longitude: d.Endereco.Longitude,
		},
		Obra: domain.Obra{
			Status: d.Obra.Status, Progresso: d.Obra.Progresso,
			Inicio: d.Obra.Inicio, Entrega: d.Obra.Entrega,
			AnoLancamento: d.Obra.AnoLancamento,
		},
		Incorporacao: d.Incorporacao,
		Edificacao: domain.Edificacao{
			Pavimentos: d.Edificacao.Pavimentos, UnidadesPorAndar: d.Edificacao.UnidadesPorAndar,
			TotalUnidades: d.Edificacao.TotalUnidades, Coberturas: d.Edificacao.Coberturas,
			AlturaTorre: d.Edificacao.AlturaTorre, ProntoParaMorar: d.Edificacao.ProntoParaMorar,
			Mobiliado: d.Edificacao.Mobiliado,
		},
		Caracteristicas: domain.Caracteristicas{
			Dormitorios: domain.FaixaInt{Min: d.Caracteristicas.Dormitorios.Min, Max: d.Caracteristicas.Dormitorios.Max},
			Suites:      domain.FaixaInt{Min: d.Caracteristicas.Suites.Min, Max: d.Caracteristicas.Suites.Max},
			Vagas:       domain.FaixaInt{Min: d.Caracteristicas.Vagas.Min, Max: d.Caracteristicas.Vagas.Max},
			Metragem:    domain.FaixaFloat{Min: d.Caracteristicas.Metragem.Min, Max: d.Caracteristicas.Metragem.Max},
			Condominio:  d.Caracteristicas.Condominio,
			IPTU:        d.Caracteristicas.IPTU,
		},
		DiferenciaisUnidade:    d.DiferenciaisUnidade,
		DiferenciaisCondominio: d.DiferenciaisCondominio,
		PontosInteresse:        pontosInteresse,
		Secoes:                 secoes,
		Galeria:                galeria,
		CatalogoS3Path:         d.CatalogoS3Path,
		Foto:                   d.Foto,
		Comercializacao: domain.Comercializacao{
			UnidadesDisponiveis: d.Comercializacao.UnidadesDisponiveis,
			UltimaAtualizacao:   d.Comercializacao.UltimaAtualizacao,
		},
		ConstrutoraNome: d.ConstrutoraNome,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type EmpreendimentoRepo struct {
	col *mongo.Collection
}

func NewEmpreendimentoRepository(db *mongo.Database) domain.EmpreendimentoRepository {
	return &EmpreendimentoRepo{col: db.Collection("empreendimentos")}
}

func (r *EmpreendimentoRepo) Create(ctx context.Context, e *domain.Empreendimento) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(e.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}
	construtoraOID, err := bson.ObjectIDFromHex(e.ConstrutoraID)
	if err != nil {
		return domain.ErrInvalidID
	}

	pontosDoc := make([]pontoInteresseDoc, len(e.PontosInteresse))
	for i, p := range e.PontosInteresse {
		pontosDoc[i] = pontoInteresseDoc{Nome: p.Nome, Distancia: p.Distancia, Tempo: p.Tempo}
	}
	secoesDoc := make([]secaoDoc, len(e.Secoes))
	for i, s := range e.Secoes {
		secoesDoc[i] = secaoDoc{Nome: s.Nome, Categoria: s.Categoria, Ordem: s.Ordem}
	}
	galeriaDoc := make([]galeriaImagemDoc, len(e.Galeria))
	for i, g := range e.Galeria {
		galeriaDoc[i] = galeriaImagemDoc{S3Path: g.S3Path, Ordem: g.Ordem, Tipo: g.Tipo}
	}

	now := time.Now().UTC()
	doc := &empreendimentoDocument{
		ID:            bson.NewObjectID(),
		TenantID:      tenantOID,
		ConstrutoraID: construtoraOID,
		Nome:          e.Nome, Subtitulo: e.Subtitulo, Slogan: e.Slogan,
		Descricao: e.Descricao, FonteID: e.FonteID,
		Endereco: enderecoDoc{
			Completo: e.Endereco.Completo, Logradouro: e.Endereco.Logradouro,
			Numero: e.Endereco.Numero, Bairro: e.Endereco.Bairro,
			Cidade: e.Endereco.Cidade, UF: e.Endereco.UF,
			Regiao: e.Endereco.Regiao, DistanciaMar: e.Endereco.DistanciaMar,
			Latitude: e.Endereco.Latitude, Longitude: e.Endereco.Longitude,
		},
		Obra: obraDoc{
			Status: e.Obra.Status, Progresso: e.Obra.Progresso,
			Inicio: e.Obra.Inicio, Entrega: e.Obra.Entrega,
			AnoLancamento: e.Obra.AnoLancamento,
		},
		Incorporacao: e.Incorporacao,
		Edificacao: edificacaoDoc{
			Pavimentos: e.Edificacao.Pavimentos, UnidadesPorAndar: e.Edificacao.UnidadesPorAndar,
			TotalUnidades: e.Edificacao.TotalUnidades, Coberturas: e.Edificacao.Coberturas,
			AlturaTorre: e.Edificacao.AlturaTorre, ProntoParaMorar: e.Edificacao.ProntoParaMorar,
			Mobiliado: e.Edificacao.Mobiliado,
		},
		Caracteristicas: caracteristicasDoc{
			Dormitorios: faixaIntDoc{Min: e.Caracteristicas.Dormitorios.Min, Max: e.Caracteristicas.Dormitorios.Max},
			Suites:      faixaIntDoc{Min: e.Caracteristicas.Suites.Min, Max: e.Caracteristicas.Suites.Max},
			Vagas:       faixaIntDoc{Min: e.Caracteristicas.Vagas.Min, Max: e.Caracteristicas.Vagas.Max},
			Metragem:    faixaFloatDoc{Min: e.Caracteristicas.Metragem.Min, Max: e.Caracteristicas.Metragem.Max},
			Condominio:  e.Caracteristicas.Condominio,
			IPTU:        e.Caracteristicas.IPTU,
		},
		DiferenciaisUnidade:    e.DiferenciaisUnidade,
		DiferenciaisCondominio: e.DiferenciaisCondominio,
		PontosInteresse:        pontosDoc,
		Secoes:                 secoesDoc,
		Galeria:                galeriaDoc,
		CatalogoS3Path:         e.CatalogoS3Path,
		Comercializacao: comercializacaoDoc{
			UnidadesDisponiveis: e.Comercializacao.UnidadesDisponiveis,
			UltimaAtualizacao:   e.Comercializacao.UltimaAtualizacao,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "empreendimento")
	}

	e.ID = doc.ID.Hex()
	e.CreatedAt = doc.CreatedAt
	e.UpdatedAt = doc.UpdatedAt
	return nil
}

func (r *EmpreendimentoRepo) FindByID(ctx context.Context, id string) (*domain.Empreendimento, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc empreendimentoDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, database.MapError(err, "empreendimento")
	}

	result := doc.toDomain()
	return &result, nil
}

// listProjection excludes large embedded arrays not needed for card listing.
var listProjection = bson.M{
	"galeria":         0,
	"pontos_interesse": 0,
	"secoes":          0,
	"diferenciais_unidade":    0,
	"diferenciais_condominio": 0,
	"descricao":       0,
}

func (r *EmpreendimentoRepo) List(ctx context.Context, tenantID string, filter domain.EmpreendimentoFilter) ([]domain.Empreendimento, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, 0, domain.ErrInvalidID
	}

	f := bson.M{"tenant_id": tenantOID, "foto": bson.M{"$ne": ""}}
	if filter.ConstrutoraID != nil {
		cOID, err := bson.ObjectIDFromHex(*filter.ConstrutoraID)
		if err != nil {
			return nil, 0, domain.ErrInvalidID
		}
		f["construtora_id"] = cOID
	}
	if filter.Cidade != nil {
		f["endereco.cidade"] = bson.M{"$regex": *filter.Cidade, "$options": "i"}
	}
	if filter.UF != nil {
		f["endereco.uf"] = *filter.UF
	}
	if filter.StatusObra != nil {
		f["obra.status"] = *filter.StatusObra
	}
	if filter.Search != nil && *filter.Search != "" {
		f["nome"] = bson.M{"$regex": *filter.Search, "$options": "i"}
	}

	page := filter.Pagination.Page
	pageSize := filter.Pagination.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	skip := (page - 1) * pageSize

	// Single aggregation: $facet does count + data in ONE round-trip
	pipeline := bson.A{
		bson.M{"$match": f},
		bson.M{"$facet": bson.M{
			"data": bson.A{
				bson.M{"$sort": bson.M{"nome": 1}},
				bson.M{"$skip": skip},
				bson.M{"$limit": pageSize},
				bson.M{"$project": listProjection},
			},
			"count": bson.A{
				bson.M{"$count": "total"},
			},
		}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, database.MapError(err, "empreendimento")
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []empreendimentoDocument `bson:"data"`
		Count []struct {
			Total int64 `bson:"total"`
		} `bson:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, database.MapError(err, "empreendimento")
	}

	if len(results) == 0 {
		return nil, 0, nil
	}

	var total int64
	if len(results[0].Count) > 0 {
		total = results[0].Count[0].Total
	}

	docs := results[0].Data
	emps := make([]domain.Empreendimento, len(docs))
	for i, doc := range docs {
		emps[i] = doc.toDomain()
	}
	return emps, total, nil
}

func (r *EmpreendimentoRepo) Search(ctx context.Context, tenantID, query string, page, pageSize int) ([]domain.Empreendimento, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, 0, domain.ErrInvalidID
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	skip := (page - 1) * pageSize

	f := bson.M{
		"tenant_id": tenantOID,
		"$text":     bson.M{"$search": query},
	}

	// Single aggregation with $facet for count + data
	pipeline := bson.A{
		bson.M{"$match": f},
		bson.M{"$addFields": bson.M{"score": bson.M{"$meta": "textScore"}}},
		bson.M{"$facet": bson.M{
			"data": bson.A{
				bson.M{"$sort": bson.M{"score": -1}},
				bson.M{"$skip": skip},
				bson.M{"$limit": pageSize},
				bson.M{"$project": listProjection},
			},
			"count": bson.A{
				bson.M{"$count": "total"},
			},
		}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, database.MapError(err, "empreendimento")
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []empreendimentoDocument `bson:"data"`
		Count []struct {
			Total int64 `bson:"total"`
		} `bson:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, database.MapError(err, "empreendimento")
	}

	if len(results) == 0 {
		return nil, 0, nil
	}

	var total int64
	if len(results[0].Count) > 0 {
		total = results[0].Count[0].Total
	}

	docs := results[0].Data
	emps := make([]domain.Empreendimento, len(docs))
	for i, doc := range docs {
		emps[i] = doc.toDomain()
	}
	return emps, total, nil
}

func (r *EmpreendimentoRepo) Update(ctx context.Context, e *domain.Empreendimento) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(e.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"nome":       e.Nome,
			"subtitulo":  e.Subtitulo,
			"slogan":     e.Slogan,
			"descricao":  e.Descricao,
			"updated_at": now,
		}},
	)
	return database.MapError(err, "empreendimento")
}

func (r *EmpreendimentoRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "empreendimento")
}

func EnsureEmpreendimentoIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("empreendimentos")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "construtora_id", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "endereco.cidade", Value: 1}, {Key: "endereco.uf", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "obra.status", Value: 1}}},
		// Text index for full-text search (Portuguese)
		{
			Keys: bson.D{
				{Key: "nome", Value: "text"},
				{Key: "endereco.cidade", Value: "text"},
				{Key: "endereco.bairro", Value: "text"},
				{Key: "endereco.logradouro", Value: "text"},
				{Key: "endereco.uf", Value: "text"},
			},
			Options: options.Index().
				SetDefaultLanguage("portuguese").
				SetName("text_search_idx"),
		},
	})
	return err
}
