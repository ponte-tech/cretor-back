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

type observacaoDocument struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	TenantID  bson.ObjectID `bson:"tenant_id"`
	NegocioID bson.ObjectID `bson:"negocio_id"`
	Texto     string        `bson:"texto"`
	Autor     string        `bson:"autor,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
}

func (d *observacaoDocument) toDomain() *domain.Observacao {
	return &domain.Observacao{
		ID:        d.ID.Hex(),
		TenantID:  d.TenantID.Hex(),
		NegocioID: d.NegocioID.Hex(),
		Texto:     d.Texto,
		Autor:     d.Autor,
		CreatedAt: d.CreatedAt,
	}
}

type ObservacaoRepo struct {
	col *mongo.Collection
}

func NewObservacaoRepository(db *mongo.Database) domain.ObservacaoRepository {
	return &ObservacaoRepo{col: db.Collection("observacoes")}
}

func (r *ObservacaoRepo) Create(ctx context.Context, o *domain.Observacao) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(o.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}
	negocioOID, err := bson.ObjectIDFromHex(o.NegocioID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &observacaoDocument{
		ID:        bson.NewObjectID(),
		TenantID:  tenantOID,
		NegocioID: negocioOID,
		Texto:     o.Texto,
		Autor:     o.Autor,
		CreatedAt: now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "observacoes")
	}

	o.ID = doc.ID.Hex()
	o.CreatedAt = now
	return nil
}

func (r *ObservacaoRepo) ListByNegocio(ctx context.Context, negocioID string) ([]domain.Observacao, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	negocioOID, err := bson.ObjectIDFromHex(negocioID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	cursor, err := r.col.Find(ctx, bson.M{"negocio_id": negocioOID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, database.MapError(err, "observacoes")
	}
	defer cursor.Close(ctx)

	var docs []observacaoDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, database.MapError(err, "observacoes")
	}

	result := make([]domain.Observacao, len(docs))
	for i, doc := range docs {
		result[i] = *doc.toDomain()
	}
	return result, nil
}

func EnsureObservacaoIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("observacoes")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "negocio_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "negocio_id", Value: 1}}},
	})
	return err
}
