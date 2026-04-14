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

type unidadeDocument struct {
	ID               bson.ObjectID `bson:"_id,omitempty"`
	TenantID         bson.ObjectID `bson:"tenant_id"`
	EmpreendimentoID bson.ObjectID `bson:"empreendimento_id"`
	Secao            string        `bson:"secao"`
	Numero           string        `bson:"numero"`
	Andar            int           `bson:"andar"`
	Metragem         float64       `bson:"metragem"`
	Tipo             string        `bson:"tipo"`
	Valor            *float64      `bson:"valor,omitempty"`
	ValorTexto       string        `bson:"valor_texto"`
	Status           string        `bson:"status"`
	CreatedAt        time.Time     `bson:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at"`
}

func (d *unidadeDocument) toDomain() domain.Unidade {
	return domain.Unidade{
		ID: d.ID.Hex(), TenantID: d.TenantID.Hex(),
		EmpreendimentoID: d.EmpreendimentoID.Hex(),
		Secao: d.Secao, Numero: d.Numero, Andar: d.Andar,
		Metragem: d.Metragem, Tipo: d.Tipo,
		Valor: d.Valor, ValorTexto: d.ValorTexto, Status: d.Status,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type UnidadeRepo struct {
	col *mongo.Collection
}

func NewUnidadeRepository(db *mongo.Database) domain.UnidadeRepository {
	return &UnidadeRepo{col: db.Collection("unidades")}
}

func (r *UnidadeRepo) Create(ctx context.Context, u *domain.Unidade) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(u.TenantID)
	if err != nil {
		return domain.ErrInvalidID
	}
	empOID, err := bson.ObjectIDFromHex(u.EmpreendimentoID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	doc := &unidadeDocument{
		ID: bson.NewObjectID(), TenantID: tenantOID, EmpreendimentoID: empOID,
		Secao: u.Secao, Numero: u.Numero, Andar: u.Andar,
		Metragem: u.Metragem, Tipo: u.Tipo,
		Valor: u.Valor, ValorTexto: u.ValorTexto, Status: u.Status,
		CreatedAt: now, UpdatedAt: now,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return database.MapError(err, "unidade")
	}

	u.ID = doc.ID.Hex()
	u.CreatedAt = doc.CreatedAt
	u.UpdatedAt = doc.UpdatedAt
	return nil
}

func (r *UnidadeRepo) CreateMany(ctx context.Context, unidades []domain.Unidade) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if len(unidades) == 0 {
		return nil
	}

	now := time.Now().UTC()
	docs := make([]interface{}, len(unidades))
	for i, u := range unidades {
		tenantOID, err := bson.ObjectIDFromHex(u.TenantID)
		if err != nil {
			return domain.ErrInvalidID
		}
		empOID, err := bson.ObjectIDFromHex(u.EmpreendimentoID)
		if err != nil {
			return domain.ErrInvalidID
		}
		docs[i] = &unidadeDocument{
			ID: bson.NewObjectID(), TenantID: tenantOID, EmpreendimentoID: empOID,
			Secao: u.Secao, Numero: u.Numero, Andar: u.Andar,
			Metragem: u.Metragem, Tipo: u.Tipo,
			Valor: u.Valor, ValorTexto: u.ValorTexto, Status: u.Status,
			CreatedAt: now, UpdatedAt: now,
		}
	}

	_, err := r.col.InsertMany(ctx, docs)
	return database.MapError(err, "unidade")
}

func (r *UnidadeRepo) FindByID(ctx context.Context, id string) (*domain.Unidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var doc unidadeDocument
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, database.MapError(err, "unidade")
	}

	result := doc.toDomain()
	return &result, nil
}

func (r *UnidadeRepo) List(ctx context.Context, tenantID string, filter domain.UnidadeFilter) ([]domain.Unidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tenantOID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	empOID, err := bson.ObjectIDFromHex(filter.EmpreendimentoID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	f := bson.M{"tenant_id": tenantOID, "empreendimento_id": empOID}
	if filter.Secao != nil {
		f["secao"] = *filter.Secao
	}
	if filter.Status != nil {
		f["status"] = *filter.Status
	}
	if filter.ValorMin != nil {
		f["valor"] = bson.M{"$gte": *filter.ValorMin}
	}
	if filter.ValorMax != nil {
		if _, ok := f["valor"]; ok {
			f["valor"].(bson.M)["$lte"] = *filter.ValorMax
		} else {
			f["valor"] = bson.M{"$lte": *filter.ValorMax}
		}
	}

	cursor, err := r.col.Find(ctx, f, options.Find().SetSort(bson.D{{Key: "secao", Value: 1}, {Key: "numero", Value: 1}}))
	if err != nil {
		return nil, database.MapError(err, "unidade")
	}
	defer cursor.Close(ctx)

	var docs []unidadeDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, database.MapError(err, "unidade")
	}

	result := make([]domain.Unidade, len(docs))
	for i, doc := range docs {
		result[i] = doc.toDomain()
	}
	return result, nil
}

func (r *UnidadeRepo) Update(ctx context.Context, u *domain.Unidade) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(u.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	now := time.Now().UTC()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"valor":       u.Valor,
			"valor_texto": u.ValorTexto,
			"status":      u.Status,
			"updated_at":  now,
		}},
	)
	return database.MapError(err, "unidade")
}

func (r *UnidadeRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return database.MapError(err, "unidade")
}

func (r *UnidadeRepo) DeleteByEmpreendimento(ctx context.Context, empreendimentoID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	empOID, err := bson.ObjectIDFromHex(empreendimentoID)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.col.DeleteMany(ctx, bson.M{"empreendimento_id": empOID})
	return database.MapError(err, "unidade")
}

func EnsureUnidadeIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("unidades")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "empreendimento_id", Value: 1}, {Key: "secao", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "empreendimento_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}, {Key: "valor", Value: 1}}},
	})
	return err
}
