# Referencia IA - Go + MongoDB: Boas Praticas de Desenvolvimento

Este documento serve como referencia para a IA auxiliar no desenvolvimento do backend Cretor usando Go com MongoDB.

---

## 1. Estrutura do Projeto

O projeto segue **DDD (Domain-Driven Design)** com **Clean Architecture**, adaptado para Go + MongoDB (migrando de DynamoDB/Lambda para MongoDB standalone).

```
cretor-back/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point, wiring, server startup
│
├── internal/
│   ├── domain/                        # Entidades puras de negocio (ZERO imports externos)
│   │   ├── cliente.go
│   │   ├── imovel.go
│   │   ├── projeto.go
│   │   ├── pipeline.go
│   │   ├── conversa.go
│   │   ├── match.go
│   │   ├── lead.go
│   │   ├── usuario.go
│   │   ├── errors.go                  # Erros de dominio
│   │   └── repository.go             # Interfaces dos repositorios
│   │
│   ├── application/                   # Casos de uso / services
│   │   ├── cliente_service.go
│   │   ├── imovel_service.go
│   │   ├── projeto_service.go
│   │   ├── pipeline_service.go
│   │   ├── matching_service.go
│   │   ├── whatsapp_service.go
│   │   └── auth_service.go
│   │
│   ├── infrastructure/
│   │   ├── mongo/
│   │   │   ├── client.go             # Conexao MongoDB singleton
│   │   │   ├── indexes.go            # Criacao de indices no startup
│   │   │   ├── cliente_repo.go       # Implementacao ClienteRepository
│   │   │   ├── imovel_repo.go
│   │   │   ├── projeto_repo.go
│   │   │   ├── pipeline_repo.go
│   │   │   ├── conversa_repo.go
│   │   │   ├── match_repo.go
│   │   │   ├── lead_repo.go
│   │   │   └── usuario_repo.go
│   │   │
│   │   └── config/
│   │       └── config.go             # Carregamento de env vars
│   │
│   ├── interfaces/
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── cliente_handler.go
│   │       │   ├── imovel_handler.go
│   │       │   ├── projeto_handler.go
│   │       │   ├── pipeline_handler.go
│   │       │   ├── matching_handler.go
│   │       │   ├── whatsapp_handler.go
│   │       │   └── auth_handler.go
│   │       │
│   │       ├── middleware/
│   │       │   ├── auth.go
│   │       │   ├── tenant.go
│   │       │   ├── cors.go
│   │       │   └── logger.go
│   │       │
│   │       ├── dto/
│   │       │   ├── cliente_dto.go
│   │       │   ├── imovel_dto.go
│   │       │   └── ...
│   │       │
│   │       └── router.go             # chi router config
│   │
│   └── pkg/
│       ├── mongoutil/
│       │   └── filters.go            # Helpers de filtro reutilizaveis
│       └── validator/
│           └── validator.go
│
├── docs/
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── Dockerfile
```

---

## 2. Conexao MongoDB

### Singleton - Um unico `mongo.Client` por aplicacao

```go
package mongo

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
    opts := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second).
        SetConnectTimeout(10 * time.Second).
        SetServerSelectionTimeout(5 * time.Second).
        SetRetryWrites(true).
        SetRetryReads(true)

    client, err := mongo.Connect(opts)
    if err != nil {
        return nil, fmt.Errorf("mongo connect: %w", err)
    }

    pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
        return nil, fmt.Errorf("mongo ping: %w", err)
    }
    return client, nil
}
```

### Regras de conexao

| Regra | Detalhe |
|---|---|
| Um `mongo.Client` por app | E goroutine-safe, compartilhe entre handlers |
| Sempre passe `context.Context` | Use contextos com timeout por operacao |
| Feche cursors | `defer cursor.Close(ctx)` apos obter cursor |
| Verifique `SingleResult.Err()` | `FindOne` nunca retorna nil |
| Use projections | Busque apenas os campos necessarios |

---

## 3. BSON Tags e Design de Structs

### Tags de referencia

| Tag | Efeito |
|---|---|
| `bson:"nome_campo"` | Mapeia campo Go para campo BSON |
| `bson:"_id,omitempty"` | Chave primaria MongoDB |
| `bson:",omitempty"` | Omite campo se valor zero |
| `bson:",inline"` | Achata struct embeddada no documento pai |
| `bson:"-"` | Nunca persiste este campo |

### Struct base reutilizavel

```go
type BaseDocument struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

type ClienteDocument struct {
    BaseDocument `bson:",inline"`
    TenantID     bson.ObjectID `bson:"tenant_id"`
    Nome         string        `bson:"nome"`
    Email        string        `bson:"email"`
    // ... demais campos
}
```

### Regras de struct

- Use `bson.Decimal128` para valores monetarios (NUNCA float64)
- Use `time.Time` para datas
- Use `bson.ObjectID` para IDs
- Separe structs de dominio (com tags JSON) de structs de persistencia (com tags BSON)
- Mapeie entre elas no repositorio

---

## 4. Padrao Repository

### Interface no dominio (zero imports do MongoDB)

```go
// domain/repository.go
package domain

import "context"

type ClienteRepository interface {
    FindByID(ctx context.Context, id string) (*Cliente, error)
    FindByEmail(ctx context.Context, tenantID, email string) (*Cliente, error)
    List(ctx context.Context, filter ClienteFilter, page Pagination) ([]Cliente, int64, error)
    Create(ctx context.Context, cliente *Cliente) error
    Update(ctx context.Context, cliente *Cliente) error
    Delete(ctx context.Context, id string) error
}

type Pagination struct {
    Page     int64
    PageSize int64
}
```

### Implementacao no infrastructure

```go
// infrastructure/mongo/cliente_repo.go
package mongo

type clienteRepo struct {
    col *mongo.Collection
}

func NewClienteRepository(db *mongo.Database) domain.ClienteRepository {
    return &clienteRepo{col: db.Collection("clientes")}
}

func (r *clienteRepo) FindByID(ctx context.Context, id string) (*domain.Cliente, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    oid, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return nil, domain.ErrInvalidID
    }

    var doc clienteDocument
    err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
    if errors.Is(err, mongo.ErrNoDocuments) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("find cliente %s: %w", id, err)
    }
    return doc.toDomain(), nil
}
```

---

## 5. Tratamento de Erros

### Erros de dominio

```go
// domain/errors.go
package domain

var (
    ErrNotFound      = errors.New("not found")
    ErrDuplicateKey  = errors.New("duplicate key")
    ErrInvalidEntity = errors.New("invalid entity")
    ErrInvalidID     = errors.New("invalid id")
    ErrUnauthorized  = errors.New("unauthorized")
)
```

### Mapeamento de erros MongoDB para dominio

```go
func mapError(err error, entity string) error {
    if err == nil {
        return nil
    }
    if errors.Is(err, mongo.ErrNoDocuments) {
        return fmt.Errorf("%s: %w", entity, domain.ErrNotFound)
    }
    var writeErr mongo.WriteException
    if errors.As(err, &writeErr) {
        for _, we := range writeErr.WriteErrors {
            if we.Code == 11000 {
                return fmt.Errorf("%s: %w", entity, domain.ErrDuplicateKey)
            }
        }
    }
    return fmt.Errorf("%s: %w", entity, err)
}
```

---

## 6. Indices - Criar no Startup

```go
// infrastructure/mongo/indexes.go
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
    // Indices sao idempotentes - seguros para chamar no startup
    
    // Clientes
    clientesCol := db.Collection("clientes")
    clientesCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{"tenant_id", 1}, {"email", 1}}, Options: options.Index().SetUnique(true)},
        {Keys: bson.D{{"tenant_id", 1}, {"status", 1}, {"created_at", -1}}},
        {Keys: bson.D{{"tenant_id", 1}, {"nome", "text"}, {"email", "text"}}},
    })

    // Imoveis
    imoveisCol := db.Collection("imoveis")
    imoveisCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{"tenant_id", 1}, {"status", 1}, {"preco", 1}}},
        {Keys: bson.D{{"tenant_id", 1}, {"tipo", 1}, {"cidade", 1}}},
        {Keys: bson.D{{"tenant_id", 1}, {"titulo", "text"}, {"descricao", "text"}}},
    })

    // Pipeline
    pipelineCol := db.Collection("pipeline")
    pipelineCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{"tenant_id", 1}, {"cliente_id", 1}}},
        {Keys: bson.D{{"tenant_id", 1}, {"etapa", 1}, {"data_criacao", -1}}},
    })

    // Match Scores
    matchCol := db.Collection("match_scores")
    matchCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{"tenant_id", 1}, {"cliente_id", 1}, {"score", -1}}},
        {Keys: bson.D{{"tenant_id", 1}, {"imovel_id", 1}}},
    })

    return nil
}
```

### Regra ESR para indices compostos

**Equality -> Sort -> Range** (nesta ordem):
```go
// Query: status == "ativo" ORDER BY created_at DESC WHERE preco > 500000
// Indice ideal:
{Keys: bson.D{{"status", 1}, {"created_at", -1}, {"preco", 1}}}
```

---

## 7. Paginacao

### Offset-based (simples, lento para paginas profundas)

```go
opts := options.Find().
    SetSkip((page - 1) * pageSize).
    SetLimit(pageSize).
    SetSort(bson.D{{"created_at", -1}})
```

### Cursor-based (eficiente em qualquer profundidade)

```go
filter := bson.M{"_id": bson.M{"$gt": lastID}}
opts := options.Find().SetLimit(pageSize).SetSort(bson.D{{"_id", 1}})
```

---

## 8. Aggregation Pipeline

```go
pipeline := mongo.Pipeline{
    // 1. $match primeiro - usa indices
    {{Key: "$match", Value: bson.M{"tenant_id": tenantID, "status": "ativo"}}},
    // 2. $project cedo - reduz memoria
    {{Key: "$project", Value: bson.M{"nome": 1, "preco": 1, "tipo": 1}}},
    // 3. $group
    {{Key: "$group", Value: bson.M{"_id": "$tipo", "total": bson.M{"$sum": 1}}}},
    // 4. $sort
    {{Key: "$sort", Value: bson.D{{"total", -1}}}},
}
```

**Regras:**
- `$match` e `$sort` primeiro (usa indices)
- `$project/$unset` cedo (reduz memoria)
- Evite `$lookup` em caminhos quentes (denormalize com Extended Reference Pattern)
- Use `bson.D` (ordenado) para stages de pipeline, `bson.M` para filtros

---

## 9. Transacoes

```go
session, err := client.StartSession()
if err != nil {
    return err
}
defer session.EndSession(ctx)

_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (any, error) {
    // operacoes atomicas aqui
    return nil, nil
})
```

**Regras:**
- Transacoes requerem **replica set** (Atlas ja e replica set)
- Use `WithTransaction` - faz retry automatico em erros transientes
- Mantenha transacoes curtas (< 60s)
- Prefira operacoes atomicas em documento unico (`$set`, `$push`, `$inc` juntos) a transacoes

---

## 10. Patterns de Schema Design

### Embedding vs Referencing

| Estrategia | Quando usar |
|---|---|
| **Embed** | 1:1 ou 1:poucos, dados lidos juntos, filho sem ciclo de vida proprio |
| **Reference** | 1:muitos ou N:N, entidade acessada independentemente, crescimento ilimitado |

### Extended Reference Pattern
Copie campos frequentemente acessados de documentos referenciados:
```go
type PipelineDocument struct {
    ClienteID    bson.ObjectID `bson:"cliente_id"`
    ClienteNome  string        `bson:"cliente_nome"`   // copia
    ClienteEmail string        `bson:"cliente_email"`  // copia
}
```

### Subset Pattern
Para arrays grandes, armazene subconjunto no pai e lista completa em collection separada.

---

## 11. Regras Gerais para a IA

1. **Sempre use `context.Context`** com timeout em toda operacao do driver
2. **Separe dominio de persistencia** - mapeie entre entidades e documentos BSON
3. **Um `mongo.Client` por aplicacao**, criado no startup
4. **Mapeie erros do driver para erros de dominio** na fronteira do repositorio
5. **Denormalize estrategicamente** com Extended Reference Pattern
6. **Crie indices no startup** seguindo regra ESR
7. **Use paginacao cursor-based** para listas que podem ser profundas
8. **Projete apenas campos necessarios** em toda query
9. **Use chi como router HTTP** (ja definido no projeto)
10. **Multi-tenant**: todo documento tem `tenant_id`, todo indice comeca com `tenant_id`
11. **Valores monetarios**: sempre `bson.Decimal128`, nunca float64
12. **Bulk operations**: use `BulkWrite` com `SetOrdered(false)` para writes em lote
