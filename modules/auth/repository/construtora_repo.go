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

type construtoraDocument struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	TenantID   bson.ObjectID `bson:"tenant_id"`
	Nome       string        `bson:"nome"`
	LogoS3Path string        `bson:"logo_s3_path"`
	Website    string        `bson:"website"`
	Descricao  string        `bson:"descricao"`
	FonteID    string        `bson:"fonte_id"`
	Status     string        `bson:"status"`
	CreatedAt  time.Time     `bson:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at"`
}

func (d *construtoraDocument) toDomain() domain.Construtora {
	return domain.Construtora{
		ID: d.ID.Hex(), TenantID: d.TenantID.Hex(),
		Nome: d.Nome, LogoS3Path: d.LogoS3Path,
		Website: d.Website, Descricao: d.Descricao,
		FonteID: d.FonteID, Status: d.Status,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type ConstrutoraRepo struct {
	col *mongo.Collection
}

func NewConstrutoraRepository(db *mongo.Database) domain.ConstrutoraRepository {
	return &ConstrutoraRepo{col: db.Collection("construtoras")}
}

func (r *ConstrutoraRepo) Create(ctx context.Context, c *domain.Construtora) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(c.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &construtoraDocument{
		ID:         bson.NewObjectID(),
		TenantID:   tenantOID,
		Nome:       c.Nome,
		LogoS3Path: c.LogoS3Path,
		Website:    c.Website,
		Descricao:  c.Descricao,
		FonteID:    c.FonteID,
		Status:     c.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "construtora")
	}

	c.ID = doc.ID.Hex()
	c.CreatedAt = doc.CreatedAt
	c.UpdatedAt = doc.UpdatedAt
	return nil
}

func (r *ConstrutoraRepo) FindByID(ctx context.Context, id string) (*domain.Construtora, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc construtoraDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, database.MapError(err, "construtora")
	}

	result := doc.toDomain()
	return &result, nil
}

func (r *ConstrutoraRepo) List(ctx context.Context, tenantID string, filter domain.ConstrutoraFilter) ([]domain.Construtora, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	f := bson.M{"tenant_id": tenantOID}
	if filter.Status != nil {
		f["status"] = *filter.Status
	}
	if filter.Search != nil && *filter.Search != "" {
		f["nome"] = bson.M{"$regex": *filter.Search, "$options": "i"}
	}

	cursor, err := r.col.Find(ctx, f, options.Find().SetSort(bson.D{{Key: "nome", Value: 1}}))
	if err != nil {
		return nil, database.MapError(err, "construtora")
	}
	defer cursor.Close(ctx)

	var docs []construtoraDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, database.MapError(err, "construtora")
	}

	result := make([]domain.Construtora, len(docs))
	for i, doc := range docs {
		result[i] = doc.toDomain()
	}
	return result, nil
}

func (r *ConstrutoraRepo) Update(ctx context.Context, c *domain.Construtora) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(c.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"nome":         c.Nome,
			"logo_s3_path": c.LogoS3Path,
			"website":      c.Website,
			"descricao":    c.Descricao,
			"status":       c.Status,
			"updated_at":   now,
		}},
	)
	return database.MapError(err, "construtora")
}

func (r *ConstrutoraRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "construtora")
}

func EnsureConstrutoraIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("construtoras")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: "nome", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}},
		},
	})
	return err
}
