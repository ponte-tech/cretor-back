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

type leadDocument struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	TenantID       bson.ObjectID `bson:"tenant_id"`
	Nome           string        `bson:"nome"`
	Whatsapp       string        `bson:"whatsapp"`
	Email          string        `bson:"email"`
	Prazo          string        `bson:"prazo"`
	FormaPagamento string        `bson:"forma_pagamento"`
	Origem         string        `bson:"origem"`
	Status         string        `bson:"status"`
	CreatedAt      time.Time     `bson:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at"`
}

type LeadRepo struct {
	col *mongo.Collection
}

func NewLeadRepository(db *mongo.Database) domain.LeadRepository {
	return &LeadRepo{col: db.Collection("leads")}
}

func (r *LeadRepo) Create(ctx context.Context, lead *domain.Lead) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(lead.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &leadDocument{
		ID:             bson.NewObjectID(),
		TenantID:       tenantOID,
		Nome:           lead.Nome,
		Whatsapp:       lead.Whatsapp,
		Email:          lead.Email,
		Prazo:          lead.Prazo,
		FormaPagamento: lead.FormaPagamento,
		Origem:         lead.Origem,
		Status:         lead.Status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "lead")
	}

	lead.ID = doc.ID.Hex()
	lead.CreatedAt = doc.CreatedAt
	lead.UpdatedAt = doc.UpdatedAt
	return nil
}

func (r *LeadRepo) FindByID(ctx context.Context, id string) (*domain.Lead, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc leadDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, database.MapError(err, "lead")
	}

	return &domain.Lead{
		ID: doc.ID.Hex(), TenantID: doc.TenantID.Hex(),
		Nome: doc.Nome, Whatsapp: doc.Whatsapp, Email: doc.Email,
		Prazo: doc.Prazo, FormaPagamento: doc.FormaPagamento,
		Origem: doc.Origem, Status: doc.Status,
		CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt,
	}, nil
}

func (r *LeadRepo) List(ctx context.Context, tenantID string, filter domain.LeadFilter) ([]domain.Lead, error) {
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
		search := *filter.Search
		f["$or"] = bson.A{
			bson.M{"nome": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"whatsapp": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	cursor, err := r.col.Find(ctx, f, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, database.MapError(err, "lead")
	}
	defer cursor.Close(ctx)

	var docs []leadDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, database.MapError(err, "lead")
	}

	leads := make([]domain.Lead, len(docs))
	for i, doc := range docs {
		leads[i] = domain.Lead{
			ID: doc.ID.Hex(), TenantID: doc.TenantID.Hex(),
			Nome: doc.Nome, Whatsapp: doc.Whatsapp, Email: doc.Email,
			Prazo: doc.Prazo, FormaPagamento: doc.FormaPagamento,
			Origem: doc.Origem, Status: doc.Status,
			CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt,
		}
	}
	return leads, nil
}

func (r *LeadRepo) Update(ctx context.Context, lead *domain.Lead) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(lead.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"nome":             lead.Nome,
			"whatsapp":         lead.Whatsapp,
			"email":            lead.Email,
			"prazo":            lead.Prazo,
			"forma_pagamento":  lead.FormaPagamento,
			"status":           lead.Status,
			"updated_at":       now,
		}},
	)
	return database.MapError(err, "lead")
}

func (r *LeadRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "lead")
}

func EnsureLeadIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("leads")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "email", Value: 1}},
		},
	})

	return err
}
