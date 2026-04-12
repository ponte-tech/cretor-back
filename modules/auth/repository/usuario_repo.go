package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/shared/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type usuarioDocument struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	TenantID     bson.ObjectID `bson:"tenant_id"`
	Nome         string        `bson:"nome"`
	Email        string        `bson:"email"`
	SenhaHash    string        `bson:"senha_hash"`
	Telefone     string        `bson:"telefone"`
	Foto         *string       `bson:"foto,omitempty"`
	Role         string        `bson:"role"`
	Status       string        `bson:"status"`
	UltimoLogin  *time.Time    `bson:"ultimo_login,omitempty"`
	RefreshToken *string       `bson:"refresh_token,omitempty"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

func (d *usuarioDocument) toDomain() *domain.Usuario {
	return &domain.Usuario{
		ID:           d.ID.Hex(),
		TenantID:     d.TenantID.Hex(),
		Nome:         d.Nome,
		Email:        d.Email,
		SenhaHash:    d.SenhaHash,
		Telefone:     d.Telefone,
		Foto:         d.Foto,
		Role:         domain.Role(d.Role),
		Status:       domain.Status(d.Status),
		UltimoLogin:  d.UltimoLogin,
		RefreshToken: d.RefreshToken,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type UsuarioRepo struct {
	col *mongo.Collection
}

func NewUsuarioRepository(db *mongo.Database) domain.UsuarioRepository {
	return &UsuarioRepo{col: db.Collection("usuarios")}
}

func (r *UsuarioRepo) FindByID(ctx context.Context, id string) (*domain.Usuario, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc usuarioDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find usuario %s: %w", id, err)
	}
	return doc.toDomain(), nil
}

func (r *UsuarioRepo) FindByEmail(ctx context.Context, tenantID, email string) (*domain.Usuario, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc usuarioDocument
	err = r.col.FindOne(ctx, bson.M{
		"tenant_id": tenantOID,
		"email":     email,
	}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find usuario by email: %w", err)
	}
	return doc.toDomain(), nil
}

func (r *UsuarioRepo) Create(ctx context.Context, usuario *domain.Usuario) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(usuario.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &usuarioDocument{
		ID:        bson.NewObjectID(),
		TenantID:  tenantOID,
		Nome:      usuario.Nome,
		Email:     usuario.Email,
		SenhaHash: usuario.SenhaHash,
		Telefone:  usuario.Telefone,
		Foto:      usuario.Foto,
		Role:      string(usuario.Role),
		Status:    string(usuario.Status),
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "usuario")
	}

	usuario.ID = doc.ID.Hex()
	usuario.CreatedAt = doc.CreatedAt
	usuario.UpdatedAt = doc.UpdatedAt
	return nil
}

func (r *UsuarioRepo) Update(ctx context.Context, usuario *domain.Usuario) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(usuario.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"nome":       usuario.Nome,
			"email":      usuario.Email,
			"telefone":   usuario.Telefone,
			"foto":       usuario.Foto,
			"role":       string(usuario.Role),
			"status":     string(usuario.Status),
			"updated_at": now,
		}},
	)
	return database.MapError(err, "usuario")
}

func (r *UsuarioRepo) UpdateLastLogin(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"ultimo_login": now,
			"updated_at":   now,
		}},
	)
	return database.MapError(err, "usuario")
}

func (r *UsuarioRepo) SaveRefreshToken(ctx context.Context, id, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"refresh_token": token,
			"updated_at":    time.Now().UTC(),
		}},
	)
	return database.MapError(err, "usuario")
}

func (r *UsuarioRepo) ClearRefreshToken(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"refresh_token": nil,
			"updated_at":    time.Now().UTC(),
		}},
	)
	return database.MapError(err, "usuario")
}

func (r *UsuarioRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "usuario")
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("usuarios")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "role", Value: 1}, {Key: "status", Value: 1}},
		},
	})

	return err
}
