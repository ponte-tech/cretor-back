# Casa Digital Backend - Arquitetura Go + AWS Lambda

## 📋 Visão Geral

Backend modular baseado em **Domain-Driven Design (DDD)** usando **Go** e **AWS Lambda**, otimizado para escalabilidade, manutenção e crescimento de equipe.

## 🏗️ Estrutura de Diretórios

```
casa-digital-back/
│
├── shared/                          # 📦 Código compartilhado (Lambda Layer)
│   ├── auth/
│   │   ├── jwt.go                   # Geração e validação JWT
│   │   ├── password.go              # Hash bcrypt
│   │   └── claims.go                # JWT claims customizados
│   │
│   ├── database/
│   │   ├── dynamodb.go              # Cliente DynamoDB singleton
│   │   ├── connection.go            # Gerenciamento de conexões
│   │   └── base_repository.go       # CRUD genérico
│   │
│   ├── middleware/
│   │   ├── auth.go                  # Middleware de autenticação
│   │   ├── tenant.go                # Middleware multi-tenant
│   │   ├── cors.go                  # CORS handler
│   │   ├── logger.go                # Request logger
│   │   └── error_handler.go         # Error handling middleware
│   │
│   ├── config/
│   │   ├── config.go                # Configurações (env vars)
│   │   └── aws.go                   # AWS clients setup
│   │
│   ├── logger/
│   │   └── logger.go                # Logger estruturado (zap)
│   │
│   ├── validator/
│   │   ├── validator.go             # Validador principal
│   │   ├── cpf.go                   # Validação CPF/CNPJ
│   │   └── email.go                 # Validação email
│   │
│   ├── response/
│   │   ├── response.go              # Builders de resposta HTTP
│   │   └── error.go                 # Estruturas de erro padronizadas
│   │
│   ├── errors/
│   │   └── errors.go                # Erros de negócio customizados
│   │
│   └── utils/
│       ├── date.go                  # Manipulação de datas
│       ├── pagination.go            # Helpers de paginação
│       └── strings.go               # Utilitários de string
│
├── modules/                         # 🎯 Módulos de Domínio (Lambdas)
│   │
│   ├── auth/                        # Lambda 1: Autenticação
│   │   ├── main.go                  # Entry point da Lambda
│   │   │
│   │   ├── domain/                  # Entidades de negócio
│   │   │   ├── user.go              # Agregado User
│   │   │   ├── session.go           # Value Object Session
│   │   │   └── password_reset.go    # Value Object PasswordReset
│   │   │
│   │   ├── handler/                 # HTTP Handlers
│   │   │   ├── signup.go            # POST /auth/signup
│   │   │   ├── login.go             # POST /auth/login
│   │   │   ├── logout.go            # POST /auth/logout
│   │   │   ├── refresh.go           # POST /auth/refresh
│   │   │   ├── password_reset.go    # POST /auth/password-reset
│   │   │   ├── password_confirm.go  # POST /auth/password-reset/confirm
│   │   │   └── oauth.go             # POST /auth/oauth/google
│   │   │
│   │   ├── service/                 # Lógica de Negócio
│   │   │   ├── signup_service.go
│   │   │   ├── login_service.go
│   │   │   ├── password_service.go
│   │   │   └── oauth_service.go
│   │   │
│   │   ├── repository/              # Acesso a Dados
│   │   │   ├── repository.go        # Interface UserRepository
│   │   │   └── dynamodb_repo.go     # Implementação DynamoDB
│   │   │
│   │   └── dto/                     # Data Transfer Objects
│   │       ├── signup_request.go
│   │       ├── login_request.go
│   │       ├── auth_response.go
│   │       └── oauth_request.go
│   │
│   │
│   └── integrations/                # Lambda 8: Integrações
│       ├── main.go
│       │
│       ├── domain/
│       │   ├── integration.go       # Agregado Integration
│       │   ├── webhook.go           # Agregado Webhook
│       │   └── api_key.go           # Value Object APIKey
│       │
│       ├── handler/
│       │   ├── whatsapp_send.go     # POST /integrations/whatsapp/send
│       │   ├── email_send.go        # POST /integrations/email/send
│       │   ├── csv_import.go        # POST /integrations/csv/import
│       │   ├── webhook_create.go    # POST /webhooks
│       │   ├── webhook_list.go      # GET /webhooks
│       │   └── webhook_trigger.go   # POST /webhooks/{id}/trigger
│       │
│       ├── service/
│       │   ├── whatsapp_service.go  # WhatsApp Business API
│       │   ├── email_service.go     # Amazon SES
│       │   ├── csv_service.go       # CSV parsing
│       │   └── webhook_service.go
│       │
│       ├── repository/
│       │   ├── integration_repository.go
│       │   ├── webhook_repository.go
│       │   └── dynamodb_repo.go
│       │
│       └── dto/
│           ├── whatsapp_dto.go
│           ├── email_dto.go
│           └── webhook_dto.go
│
├── tests/                           # 🧪 Testes
│   ├── unit/
│   │   ├── auth/
│   │   │   ├── service_test.go
│   │   │   └── handler_test.go
│   │   ├── leads/
│   │   │   ├── service_test.go
│   │   │   └── scoring_test.go
│   │   └── ...
│   │
│   ├── integration/
│   │   ├── auth_integration_test.go
│   │   ├── leads_integration_test.go
│   │   └── ...
│   │
│   └── fixtures/
│       ├── users.go                 # Dados de teste para users
│       ├── leads.go                 # Dados de teste para leads
│       └── ...
│
├── scripts/                         # 🛠️ Scripts auxiliares
│   ├── build.sh                     # Build de todas lambdas
│   ├── deploy.sh                    # Deploy via AWS CLI
│   ├── seed-db.sh                   # Popular DynamoDB com dados de teste
│   ├── test-local.sh                # Testes locais com SAM CLI
│   └── generate-mocks.sh            # Gerar mocks para testes
│
├── docs/                            # 📚 Documentação
│   ├── api/
│   │   └── openapi.yaml             # Especificação OpenAPI 3.0
│   ├── architecture.md              # Este arquivo
│   ├── deployment.md                # Guia de deployment
│   └── development.md               # Guia de desenvolvimento
│
├── .github/
│   └── workflows/
│       ├── ci.yml                   # CI: tests, lint, build
│       ├── cd-dev.yml               # CD: deploy dev
│       └── cd-prod.yml              # CD: deploy prod
│
├── go.mod                           # Dependências Go
├── go.sum
├── Makefile                         # Comandos make
├── .env.example                     # Exemplo de variáveis de ambiente
├── .gitignore
├── .golangci.yml                    # Configuração do linter
└── README.md
```

---

## 🎯 Padrões de Arquitetura

### Domain-Driven Design (DDD)

Cada módulo segue a estrutura DDD:

```
module/
├── domain/          # Entidades, Value Objects, Aggregates
├── handler/         # Camada de apresentação (HTTP)
├── service/         # Lógica de negócio (Use Cases)
├── repository/      # Acesso a dados (Interfaces + Implementações)
└── dto/             # Data Transfer Objects
```

### Separação de Responsabilidades

1. **Domain Layer** (`domain/`)
   - Entidades de negócio puras
   - Regras de validação do domínio
   - Zero dependências externas

2. **Service Layer** (`service/`)
   - Orquestração de casos de uso
   - Lógica de negócio complexa
   - Depende de `repository` (interfaces)

3. **Repository Layer** (`repository/`)
   - Interfaces: definem contratos
   - Implementações: DynamoDB, S3, etc.
   - Isolamento da infraestrutura

4. **Handler Layer** (`handler/`)
   - Recebe requests HTTP
   - Valida entrada (usando DTOs)
   - Chama services
   - Formata responses

5. **DTO Layer** (`dto/`)
   - Request/Response schemas
   - Validações de entrada
   - Serialização JSON

---

## 🔧 Código Compartilhado (shared/)

### Princípios

- **Reutilização**: Código comum a todas lambdas
- **Lambda Layer**: Empacotado como Lambda Layer para reduzir tamanho
- **Zero Lógica de Negócio**: Apenas utilitários e infraestrutura

### Componentes

#### 1. **Auth** (`shared/auth/`)
```go
// jwt.go
func GenerateToken(userID, tenantID string) (string, error)
func ValidateToken(token string) (*Claims, error)

// password.go
func HashPassword(password string) (string, error)
func ComparePassword(hashed, password string) bool
```

#### 2. **Database** (`shared/database/`)
```go
// dynamodb.go - Singleton
func GetDynamoDBClient() *dynamodb.Client

// base_repository.go - CRUD genérico
type BaseRepository struct {
    client *dynamodb.Client
    table  string
}
```

#### 3. **Middleware** (`shared/middleware/`)
```go
// auth.go
func RequireAuth(next http.Handler) http.Handler

// tenant.go
func RequireTenant(next http.Handler) http.Handler

// logger.go
func RequestLogger(next http.Handler) http.Handler
```

#### 4. **Response** (`shared/response/`)
```go
// response.go
func Success(w http.ResponseWriter, data interface{})
func Error(w http.ResponseWriter, err error)
func Paginated(w http.ResponseWriter, data interface{}, pagination Pagination)
```

#### 5. **Validator** (`shared/validator/`)
```go
// validator.go
func Validate(v interface{}) error

// cpf.go
func ValidateCPF(cpf string) bool
func ValidateCNPJ(cnpj string) bool
```

---

## 🚀 Lambda Handlers

### Pattern de Entry Point

```go
// modules/auth/main.go
package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/chi"
    gochi "github.com/go-chi/chi/v5"

    "github.com/ponte-tech/casa-digital-back/modules/auth/handler"
    "github.com/ponte-tech/casa-digital-back/shared/middleware"
    "github.com/ponte-tech/casa-digital-back/shared/config"
    "github.com/ponte-tech/casa-digital-back/shared/logger"
)

var chiLambda *chiadapter.ChiLambda

func init() {
    // Carrega configurações
    cfg := config.Load()

    // Inicializa logger
    logger.Init(cfg.Environment)

    // Cria router
    r := gochi.NewRouter()

    // Middlewares globais
    r.Use(middleware.RequestLogger)
    r.Use(middleware.CORS)
    r.Use(middleware.ErrorHandler)

    // Rotas públicas
    r.Post("/auth/signup", handler.Signup)
    r.Post("/auth/login", handler.Login)
    r.Post("/auth/password-reset", handler.PasswordReset)
    r.Post("/auth/password-reset/confirm", handler.PasswordResetConfirm)
    r.Post("/auth/oauth/google", handler.OAuthGoogle)

    // Rotas protegidas
    r.Group(func(r gochi.Router) {
        r.Use(middleware.RequireAuth)

        r.Post("/auth/logout", handler.Logout)
        r.Post("/auth/refresh", handler.RefreshToken)
    })

    chiLambda = chiadapter.New(r)
}

func main() {
    lambda.Start(chiLambda.ProxyWithContext)
}
```

---

## 📦 Dependências (go.mod)

```go
module github.com/ponte-tech/casa-digital-back

go 1.23

require (
    // AWS Lambda
    github.com/aws/aws-lambda-go v1.47.0
    github.com/awslabs/aws-lambda-go-api-proxy v0.16.2

    // AWS SDK v2
    github.com/aws/aws-sdk-go-v2 v1.30.5
    github.com/aws/aws-sdk-go-v2/config v1.27.33
    github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.9
    github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.14.12
    github.com/aws/aws-sdk-go-v2/service/s3 v1.61.2
    github.com/aws/aws-sdk-go-v2/service/ses v1.27.2

    // HTTP Router
    github.com/go-chi/chi/v5 v5.1.0
    github.com/go-chi/cors v1.2.1

    // Validação
    github.com/go-playground/validator/v10 v10.22.1

    // JWT
    github.com/golang-jwt/jwt/v5 v5.2.1

    // Utilitários
    github.com/google/uuid v1.6.0

    // Logs
    go.uber.org/zap v1.27.0

    // Crypto
    golang.org/x/crypto v0.28.0
)
```

---

## 🔨 Makefile

```makefile
.PHONY: help build test lint deploy clean deps

help: ## Mostra este help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Instala dependências
	go mod download
	go mod tidy

build: ## Build de todas as lambdas
	@echo "Building all lambdas..."
	@for dir in modules/*/main.go; do \
		MODULE=$$(dirname $$dir); \
		echo "Building $$MODULE..."; \
		cd $$MODULE && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go; \
		zip -j function.zip bootstrap; \
		rm bootstrap; \
		cd ../..; \
	done

build-auth: ## Build apenas lambda auth
	cd modules/auth && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go
	cd modules/auth && zip -j function.zip bootstrap && rm bootstrap

build-leads: ## Build apenas lambda leads
	cd modules/leads && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go
	cd modules/leads && zip -j function.zip bootstrap && rm bootstrap

test: ## Roda todos os testes
	go test -v -race -coverprofile=coverage.out ./...

test-unit: ## Roda testes unitários
	go test -v -race ./tests/unit/...

test-integration: ## Roda testes de integração
	go test -v -race ./tests/integration/...

test-coverage: ## Gera relatório de coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Roda linter
	golangci-lint run ./...

fmt: ## Formata código
	go fmt ./...
	goimports -w .

clean: ## Remove arquivos de build
	find modules -name "bootstrap" -delete
	find modules -name "function.zip" -delete
	rm -f coverage.out coverage.html

local-auth: ## Roda lambda auth localmente
	sam local start-api -t sam-templates/auth.yaml

deploy-dev: build ## Deploy para dev
	@echo "Deploy dev via Terraform..."
	# Terraform está em repo separado

deploy-prod: build ## Deploy para prod
	@echo "Deploy prod via Terraform..."
	# Terraform está em repo separado

generate-mocks: ## Gera mocks para testes
	@echo "Generating mocks..."
	mockgen -source=modules/auth/repository/repository.go -destination=tests/mocks/auth_repository_mock.go
	mockgen -source=modules/leads/repository/repository.go -destination=tests/mocks/lead_repository_mock.go
```

---

## 🧪 Testes

### Estrutura de Testes

```
tests/
├── unit/
│   ├── auth/
│   │   ├── service_test.go          # Testa AuthService
│   │   ├── password_test.go         # Testa hash de senha
│   │   └── jwt_test.go              # Testa JWT
│   │
│   ├── leads/
│   │   ├── service_test.go          # Testa LeadService
│   │   └── scoring_test.go          # Testa scoring
│   │
│   └── shared/
│       ├── validator_test.go
│       └── response_test.go
│
├── integration/
│   ├── auth_integration_test.go     # Testa fluxo completo auth
│   └── leads_integration_test.go    # Testa fluxo completo leads
│
└── fixtures/
    ├── users.go                     # Mock data users
    └── leads.go                     # Mock data leads
```

### Exemplo de Teste Unitário

```go
// tests/unit/auth/service_test.go
package auth_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/ponte-tech/casa-digital-back/modules/auth/service"
)

func TestSignupService(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    svc := service.NewSignupService(mockRepo)

    // Act
    user, err := svc.Signup("test@example.com", "password123")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

---

## 🔐 Autenticação e Autorização

### Fluxo de Autenticação

1. **Signup**: `POST /auth/signup`
2. **Login**: `POST /auth/login` → retorna JWT
3. **Request protegido**: Header `Authorization: Bearer <token>`
4. **Middleware valida JWT** → extrai `user_id` e `tenant_id`
5. **Context propagado** para handlers

### Multi-Tenancy

```go
// Middleware adiciona tenant_id ao context
ctx := context.WithValue(r.Context(), "tenant_id", tenantID)

// Repository filtra por tenant_id automaticamente
func (r *LeadRepository) List(ctx context.Context) ([]*Lead, error) {
    tenantID := ctx.Value("tenant_id").(string)
    // Query DynamoDB com filtro tenant_id
}
```

---

## 📊 DynamoDB Schema

### Single Table Design

```
PK                    SK                       Type        Attributes
----------------------------------------------------------------------
TENANT#123            TENANT#123               Tenant      name, plan, ...
TENANT#123            USER#456                 User        email, name, ...
TENANT#123            LEAD#789                 Lead        name, email, status, ...
TENANT#123            OPPORTUNITY#101          Opportunity name, value, ...
USER#456              SESSION#abc              Session     token, expires_at
LEAD#789              ACTIVITY#xyz             Activity    type, description
```

### GSI (Global Secondary Indexes)

1. **GSI-Email**: Para login por email
   - PK: `email`
   - SK: `TENANT#<id>`

2. **GSI-Status**: Para filtros por status
   - PK: `TENANT#<id>#STATUS#<status>`
   - SK: `created_at`

---

## 🚀 Build e Deploy

### Build Local

```bash
# Build todas lambdas
make build

# Build lambda específica
make build-auth
```

### Deploy (via Terraform - repo separado)

```bash
# Dev
make deploy-dev

# Prod
make deploy-prod
```

### Lambda Layer (Shared Code)

```bash
# Criar layer com código shared
cd shared
GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o layer.so
zip layer.zip layer.so

# Deploy layer via Terraform
```

---

## 🔄 CI/CD

### GitHub Actions - CI

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install dependencies
        run: make deps

      - name: Run tests
        run: make test

      - name: Run linter
        run: make lint

      - name: Build
        run: make build

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
```

### GitHub Actions - CD

```yaml
# .github/workflows/cd-dev.yml
name: Deploy Dev

on:
  push:
    branches: [develop]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build
        run: make build

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Deploy to Lambda
        run: |
          for module in modules/*/function.zip; do
            MODULE_NAME=$(basename $(dirname $module))
            aws lambda update-function-code \
              --function-name casa-digital-$MODULE_NAME-dev \
              --zip-file fileb://$module
          done
```

---

## 🌍 Variáveis de Ambiente

### .env.example

```bash
# Application
ENVIRONMENT=development
AWS_REGION=us-east-1

# Database
DYNAMODB_TABLE_PREFIX=casa-digital

# Auth
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# AWS Services
S3_BUCKET=casa-digital-files

# Integrations
WHATSAPP_API_KEY=your-whatsapp-api-key
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
```

---

## 📈 Observabilidade

### CloudWatch Logs

- Cada Lambda tem seu log group
- Logs estruturados em JSON (zap)
- Retention: 30 dias (dev), 90 dias (prod)

### CloudWatch Metrics

- Invocations
- Duration
- Errors
- Throttles
- Cold starts

### X-Ray

- Tracing distribuído
- Performance por operação
- Mapa de serviços

---

## 🎯 Boas Práticas

### 1. **Cold Start Optimization**

```go
// ✅ BOM - Inicializa fora do handler
var (
    db     *dynamodb.Client
    logger *zap.Logger
)

func init() {
    db = setupDynamoDB()
    logger = setupLogger()
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) {
    // Usa db e logger já inicializados
}
```

### 2. **Error Handling**

```go
// ✅ BOM - Errors tipados e estruturados
if err := service.CreateLead(lead); err != nil {
    if errors.Is(err, ErrLeadAlreadyExists) {
        return response.Error(w, errors.Conflict("Lead already exists"))
    }
    logger.Error("failed to create lead", zap.Error(err))
    return response.Error(w, errors.InternalServer("Failed to create lead"))
}
```

### 3. **Context Propagation**

```go
// ✅ BOM - Propaga context para logs e queries
func (s *LeadService) Create(ctx context.Context, lead *Lead) error {
    userID := ctx.Value("user_id").(string)

    logger.Info("creating lead",
        zap.String("user_id", userID),
        zap.String("lead_id", lead.ID))

    return s.repo.Create(ctx, lead)
}
```

### 4. **Validação**

```go
// ✅ BOM - Valida no DTO
type CreateLeadRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=100"`
    Email string `json:"email" validate:"required,email"`
    Phone string `json:"phone" validate:"required,e164"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateLeadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, errors.BadRequest("Invalid request body"))
        return
    }

    if err := validator.Validate(req); err != nil {
        response.Error(w, errors.ValidationError(err))
        return
    }

    // ...
}
```

### 5. **Repository Pattern**

```go
// ✅ BOM - Interface + Implementação
// repository/repository.go
type LeadRepository interface {
    Create(ctx context.Context, lead *Lead) error
    GetByID(ctx context.Context, id string) (*Lead, error)
    List(ctx context.Context, filters map[string]interface{}) ([]*Lead, error)
}

// repository/dynamodb_repo.go
type dynamoLeadRepository struct {
    client *dynamodb.Client
    table  string
}

func NewLeadRepository(client *dynamodb.Client) LeadRepository {
    return &dynamoLeadRepository{
        client: client,
        table:  "casa-digital-leads",
    }
}
```

---

## 🔍 Exemplo Completo: Módulo Leads

### 1. Domain Entity

```go
// modules/leads/domain/lead.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type Lead struct {
    ID         string
    TenantID   string
    Name       string
    Email      string
    Phone      string
    Status     LeadStatus
    Source     LeadSource
    Score      int
    AssignedTo string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type LeadStatus string

const (
    StatusNew        LeadStatus = "new"
    StatusContacted  LeadStatus = "contacted"
    StatusQualified  LeadStatus = "qualified"
    StatusConverted  LeadStatus = "converted"
    StatusLost       LeadStatus = "lost"
)

type LeadSource string

const (
    SourceWebsite    LeadSource = "website"
    SourceFacebook   LeadSource = "facebook"
    SourceGoogle     LeadSource = "google"
    SourceReferral   LeadSource = "referral"
    SourceManual     LeadSource = "manual"
)

// NewLead cria um novo lead
func NewLead(tenantID, name, email, phone string, source LeadSource) *Lead {
    return &Lead{
        ID:        uuid.New().String(),
        TenantID:  tenantID,
        Name:      name,
        Email:     email,
        Phone:     phone,
        Status:    StatusNew,
        Source:    source,
        Score:     0,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// Validate valida regras de negócio
func (l *Lead) Validate() error {
    if l.Name == "" {
        return errors.New("name is required")
    }
    if l.Email == "" {
        return errors.New("email is required")
    }
    return nil
}
```

### 2. Repository Interface

```go
// modules/leads/repository/repository.go
package repository

import (
    "context"
    "github.com/ponte-tech/casa-digital-back/modules/leads/domain"
)

type LeadRepository interface {
    Create(ctx context.Context, lead *domain.Lead) error
    GetByID(ctx context.Context, tenantID, id string) (*domain.Lead, error)
    List(ctx context.Context, tenantID string, filters map[string]interface{}) ([]*domain.Lead, error)
    Update(ctx context.Context, lead *domain.Lead) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 3. DynamoDB Implementation

```go
// modules/leads/repository/dynamodb_repo.go
package repository

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/ponte-tech/casa-digital-back/modules/leads/domain"
)

type dynamoLeadRepository struct {
    client *dynamodb.Client
    table  string
}

func NewDynamoDBRepository(client *dynamodb.Client, table string) LeadRepository {
    return &dynamoLeadRepository{
        client: client,
        table:  table,
    }
}

func (r *dynamoLeadRepository) Create(ctx context.Context, lead *domain.Lead) error {
    item, err := attributevalue.MarshalMap(lead)
    if err != nil {
        return err
    }

    _, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: &r.table,
        Item:      item,
    })

    return err
}

func (r *dynamoLeadRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.Lead, error) {
    // Implementação GetItem
}
```

### 4. Service Layer

```go
// modules/leads/service/lead_service.go
package service

import (
    "context"
    "github.com/ponte-tech/casa-digital-back/modules/leads/domain"
    "github.com/ponte-tech/casa-digital-back/modules/leads/repository"
    "github.com/ponte-tech/casa-digital-back/shared/errors"
    "github.com/ponte-tech/casa-digital-back/shared/logger"
    "go.uber.org/zap"
)

type LeadService struct {
    repo repository.LeadRepository
}

func NewLeadService(repo repository.LeadRepository) *LeadService {
    return &LeadService{repo: repo}
}

func (s *LeadService) Create(ctx context.Context, lead *domain.Lead) error {
    // Validação de negócio
    if err := lead.Validate(); err != nil {
        return errors.BadRequest(err.Error())
    }

    // Verifica duplicação (regra de negócio)
    existing, _ := s.repo.GetByEmail(ctx, lead.TenantID, lead.Email)
    if existing != nil {
        return errors.Conflict("Lead with this email already exists")
    }

    // Cria lead
    if err := s.repo.Create(ctx, lead); err != nil {
        logger.Error("failed to create lead", zap.Error(err))
        return errors.InternalServer("Failed to create lead")
    }

    logger.Info("lead created", zap.String("lead_id", lead.ID))
    return nil
}
```

### 5. DTO

```go
// modules/leads/dto/create_lead.go
package dto

type CreateLeadRequest struct {
    Name   string `json:"name" validate:"required,min=3,max=100"`
    Email  string `json:"email" validate:"required,email"`
    Phone  string `json:"phone" validate:"required,e164"`
    Source string `json:"source" validate:"required,oneof=website facebook google referral manual"`
}

type LeadResponse struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Phone     string `json:"phone"`
    Status    string `json:"status"`
    Source    string `json:"source"`
    Score     int    `json:"score"`
    CreatedAt string `json:"created_at"`
}

func ToLeadResponse(lead *domain.Lead) *LeadResponse {
    return &LeadResponse{
        ID:        lead.ID,
        Name:      lead.Name,
        Email:     lead.Email,
        Phone:     lead.Phone,
        Status:    string(lead.Status),
        Source:    string(lead.Source),
        Score:     lead.Score,
        CreatedAt: lead.CreatedAt.Format(time.RFC3339),
    }
}
```

### 6. Handler

```go
// modules/leads/handler/create.go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/ponte-tech/casa-digital-back/modules/leads/domain"
    "github.com/ponte-tech/casa-digital-back/modules/leads/dto"
    "github.com/ponte-tech/casa-digital-back/modules/leads/service"
    "github.com/ponte-tech/casa-digital-back/shared/response"
    "github.com/ponte-tech/casa-digital-back/shared/validator"
)

type CreateLeadHandler struct {
    service *service.LeadService
}

func NewCreateLeadHandler(service *service.LeadService) *CreateLeadHandler {
    return &CreateLeadHandler{service: service}
}

func (h *CreateLeadHandler) Handle(w http.ResponseWriter, r *http.Request) {
    // Extrai tenant do context (adicionado pelo middleware)
    tenantID := r.Context().Value("tenant_id").(string)

    // Parse request
    var req dto.CreateLeadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, errors.BadRequest("Invalid request body"))
        return
    }

    // Valida request
    if err := validator.Validate(req); err != nil {
        response.Error(w, errors.ValidationError(err))
        return
    }

    // Cria domain entity
    lead := domain.NewLead(
        tenantID,
        req.Name,
        req.Email,
        req.Phone,
        domain.LeadSource(req.Source),
    )

    // Chama service
    if err := h.service.Create(r.Context(), lead); err != nil {
        response.Error(w, err)
        return
    }

    // Retorna response
    response.Created(w, dto.ToLeadResponse(lead))
}
```

### 7. Main (Lambda Entry Point)

```go
// modules/leads/main.go
package main

import (
    "context"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/chi"
    gochi "github.com/go-chi/chi/v5"

    "github.com/ponte-tech/casa-digital-back/modules/leads/handler"
    "github.com/ponte-tech/casa-digital-back/modules/leads/repository"
    "github.com/ponte-tech/casa-digital-back/modules/leads/service"
    "github.com/ponte-tech/casa-digital-back/shared/config"
    "github.com/ponte-tech/casa-digital-back/shared/database"
    "github.com/ponte-tech/casa-digital-back/shared/logger"
    "github.com/ponte-tech/casa-digital-back/shared/middleware"
)

var chiLambda *chiadapter.ChiLambda

func init() {
    // Configurações
    cfg := config.Load()
    logger.Init(cfg.Environment)

    // Database
    db := database.GetDynamoDBClient()

    // Repository
    leadRepo := repository.NewDynamoDBRepository(db, cfg.DynamoDBTablePrefix+"-leads")

    // Service
    leadService := service.NewLeadService(leadRepo)

    // Handlers
    createHandler := handler.NewCreateLeadHandler(leadService)
    listHandler := handler.NewListLeadHandler(leadService)
    getHandler := handler.NewGetLeadHandler(leadService)
    updateHandler := handler.NewUpdateLeadHandler(leadService)
    deleteHandler := handler.NewDeleteLeadHandler(leadService)

    // Router
    r := gochi.NewRouter()

    // Middlewares
    r.Use(middleware.RequestLogger)
    r.Use(middleware.CORS)
    r.Use(middleware.ErrorHandler)
    r.Use(middleware.RequireAuth)
    r.Use(middleware.RequireTenant)

    // Routes
    r.Post("/leads", createHandler.Handle)
    r.Get("/leads", listHandler.Handle)
    r.Get("/leads/{id}", getHandler.Handle)
    r.Put("/leads/{id}", updateHandler.Handle)
    r.Delete("/leads/{id}", deleteHandler.Handle)
    r.Post("/leads/{id}/convert", convertHandler.Handle)
    r.Get("/leads/{id}/scoring", scoringHandler.Handle)

    chiLambda = chiadapter.New(r)
}

func main() {
    lambda.Start(chiLambda.ProxyWithContext)
}
```

---

## 🎓 Resumo das Vantagens

### ✅ Modularidade
- Cada Lambda é independente
- Fácil escalar equipes (cada time cuida de um módulo)

### ✅ DDD
- Código organizado por domínio
- Regras de negócio isoladas
- Fácil manutenção

### ✅ Testabilidade
- Interfaces facilitam mocks
- Testes unitários simples
- Testes de integração isolados

### ✅ Performance
- Cold start otimizado (init() fora do handler)
- Binários pequenos (go build -ldflags="-s -w")
- Lambda Layer para código compartilhado

### ✅ Escalabilidade
- Lambdas escalam automaticamente
- DynamoDB single-table design
- Multi-tenancy nativo

---

## 📚 Próximos Passos

1. ✅ Criar estrutura de diretórios
2. ✅ Implementar módulo shared
3. ✅ Implementar módulo auth
4. ✅ Implementar módulo leads
5. ⏳ Implementar demais módulos
6. ⏳ Configurar Terraform (repo separado)
7. ⏳ Configurar CI/CD
8. ⏳ Deploy dev
9. ⏳ Testes end-to-end
10. ⏳ Deploy prod

---

**Documentação criada para Casa Digital Backend**
*Arquitetura modular Go + Lambda + DDD*
