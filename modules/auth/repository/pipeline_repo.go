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

type negocioDocument struct {
	ID                      bson.ObjectID `bson:"_id,omitempty"`
	TenantID                bson.ObjectID `bson:"tenant_id"`
	LeadID                  bson.ObjectID `bson:"lead_id"`
	LeadNome                string        `bson:"lead_nome"`
	LeadEmail               string        `bson:"lead_email"`
	LeadTelefone            string        `bson:"lead_telefone"`
	LeadPrazo               string        `bson:"lead_prazo,omitempty"`
	LeadFormaPagamento      string        `bson:"lead_forma_pagamento,omitempty"`
	Etapa                   string        `bson:"etapa"`
	Prioridade              string        `bson:"prioridade"`
	ValorNegocio            float64       `bson:"valor_negocio"`
	ProbabilidadeFechamento int           `bson:"probabilidade_fechamento"`
	DataCriacao             time.Time     `bson:"data_criacao"`
	DataUltimaInteracao     time.Time     `bson:"data_ultima_interacao"`
	DataMovimentacao        time.Time     `bson:"data_movimentacao"`
	DiasNaEtapa             int           `bson:"dias_na_etapa"`
	ProximaAcao             string        `bson:"proxima_acao,omitempty"`
	UltimaAnotacao          string        `bson:"ultima_anotacao,omitempty"`
	CorretorResponsavel     string        `bson:"corretor_responsavel,omitempty"`
	Tags                    []string      `bson:"tags,omitempty"`
	MotivoPerda             string        `bson:"motivo_perda,omitempty"`
	CreatedAt               time.Time     `bson:"created_at"`
	UpdatedAt               time.Time     `bson:"updated_at"`
}

func (d *negocioDocument) toDomain() *domain.Negocio {
	return &domain.Negocio{
		ID: d.ID.Hex(), TenantID: d.TenantID.Hex(), LeadID: d.LeadID.Hex(),
		LeadNome: d.LeadNome, LeadEmail: d.LeadEmail, LeadTelefone: d.LeadTelefone,
		LeadPrazo: d.LeadPrazo, LeadFormaPagamento: d.LeadFormaPagamento,
		Etapa: d.Etapa, Prioridade: d.Prioridade, ValorNegocio: d.ValorNegocio,
		ProbabilidadeFechamento: d.ProbabilidadeFechamento,
		DataCriacao: d.DataCriacao, DataUltimaInteracao: d.DataUltimaInteracao,
		DataMovimentacao: d.DataMovimentacao, DiasNaEtapa: d.DiasNaEtapa,
		ProximaAcao: d.ProximaAcao, UltimaAnotacao: d.UltimaAnotacao,
		CorretorResponsavel: d.CorretorResponsavel, Tags: d.Tags,
		MotivoPerda: d.MotivoPerda, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type PipelineRepo struct {
	col *mongo.Collection
}

func NewPipelineRepository(db *mongo.Database) domain.PipelineRepository {
	return &PipelineRepo{col: db.Collection("pipeline")}
}

func (r *PipelineRepo) Create(ctx context.Context, n *domain.Negocio) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(n.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}
	leadOID, err := bson.ObjectIDFromHex(n.LeadID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &negocioDocument{
		ID: bson.NewObjectID(), TenantID: tenantOID, LeadID: leadOID,
		LeadNome: n.LeadNome, LeadEmail: n.LeadEmail, LeadTelefone: n.LeadTelefone,
		LeadPrazo: n.LeadPrazo, LeadFormaPagamento: n.LeadFormaPagamento,
		Etapa: n.Etapa, Prioridade: n.Prioridade, ValorNegocio: n.ValorNegocio,
		ProbabilidadeFechamento: n.ProbabilidadeFechamento,
		DataCriacao: now, DataUltimaInteracao: now, DataMovimentacao: now,
		DiasNaEtapa: 0, ProximaAcao: n.ProximaAcao, UltimaAnotacao: n.UltimaAnotacao,
		CorretorResponsavel: n.CorretorResponsavel, Tags: n.Tags,
		CreatedAt: now, UpdatedAt: now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "pipeline")
	}

	n.ID = doc.ID.Hex()
	n.DataCriacao = now
	n.DataMovimentacao = now
	n.CreatedAt = now
	n.UpdatedAt = now
	return nil
}

func (r *PipelineRepo) FindByID(ctx context.Context, id string) (*domain.Negocio, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc negocioDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, database.MapError(err, "pipeline")
	}
	return doc.toDomain(), nil
}

func (r *PipelineRepo) List(ctx context.Context, tenantID string) ([]domain.Negocio, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	cursor, err := r.col.Find(ctx, bson.M{"tenant_id": tenantOID},
		options.Find().SetSort(bson.D{{Key: "data_criacao", Value: -1}}))
	if err != nil {
		return nil, database.MapError(err, "pipeline")
	}
	defer cursor.Close(ctx)

	var docs []negocioDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, database.MapError(err, "pipeline")
	}

	result := make([]domain.Negocio, len(docs))
	for i, doc := range docs {
		result[i] = *doc.toDomain()
	}
	return result, nil
}

func (r *PipelineRepo) Update(ctx context.Context, n *domain.Negocio) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(n.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"etapa": n.Etapa, "prioridade": n.Prioridade,
		"valor_negocio": n.ValorNegocio, "probabilidade_fechamento": n.ProbabilidadeFechamento,
		"proxima_acao": n.ProximaAcao, "ultima_anotacao": n.UltimaAnotacao,
		"corretor_responsavel": n.CorretorResponsavel, "tags": n.Tags,
		"motivo_perda": n.MotivoPerda, "data_ultima_interacao": now, "updated_at": now,
	}})
	return database.MapError(err, "pipeline")
}

func (r *PipelineRepo) UpdateEtapa(ctx context.Context, id, etapa string, probabilidade int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"etapa": etapa, "probabilidade_fechamento": probabilidade,
		"data_movimentacao": now, "dias_na_etapa": 0,
		"data_ultima_interacao": now, "updated_at": now,
	}})
	return database.MapError(err, "pipeline")
}

func (r *PipelineRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "pipeline")
}

func EnsurePipelineIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("pipeline")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "etapa", Value: 1}, {Key: "data_criacao", Value: -1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "lead_id", Value: 1}}},
	})
	return err
}
