# Guia de Boas Práticas: Chi Framework

Este documento serve como referência para desenvolvimento de APIs usando o framework Chi em Go, com foco em AWS Lambda e aplicações serverless.

---

## Índice

1. [Por que Chi?](#por-que-chi)
2. [Estrutura do Projeto](#estrutura-do-projeto)
3. [Setup Básico](#setup-básico)
4. [Routing e Organização](#routing-e-organização)
5. [Middlewares](#middlewares)
6. [Handlers](#handlers)
7. [Integração com AWS Lambda](#integração-com-aws-lambda)
8. [Validação e Error Handling](#validação-e-error-handling)
9. [Testes](#testes)
10. [Performance e Otimização](#performance-e-otimização)

---

## Por que Chi?

### Vantagens do Chi

✅ **Minimalista e Leve**
- Apenas ~100KB de overhead
- Bundle size pequeno (~13-20 MB em Lambda)
- Cold start mínimo (150-350ms)

✅ **Compatível com stdlib**
- Usa `net/http` padrão
- Fácil integração com código existente
- Sem dependências pesadas

✅ **Idiomático Go**
- Segue convenções da linguagem
- Context-aware
- Composição sobre configuração

✅ **Roteamento Eficiente**
- Radix tree para performance
- Suporta path parameters
- Regex patterns quando necessário

✅ **Middlewares Ricos**
- Biblioteca oficial de middlewares
- Fácil criar middlewares customizados
- Composição limpa

---

## Estrutura do Projeto

### Layout Recomendado para Lambda com Chi

```
casa-digital-back/
├── cmd/
│   └── lambda/
│       ├── main.go                    # Entry point Lambda
│       └── routes.go                  # Definição de rotas
├── internal/
│   ├── handler/                       # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── auth_handler.go
│   │   └── health_handler.go
│   ├── middleware/                    # Middlewares customizados
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── service/                       # Lógica de negócio
│   │   └── user_service.go
│   ├── repository/                    # Acesso a dados
│   │   └── user_repository.go
│   └── domain/                        # Modelos de domínio
│       └── user.go
├── pkg/
│   ├── response/                      # Helpers de resposta
│   │   └── json.go
│   └── validator/                     # Validação
│       └── validator.go
└── go.mod
```

---

## Setup Básico

### Instalação

```bash
# Chi router
go get -u github.com/go-chi/chi/v5

# Middlewares oficiais (opcional)
go get -u github.com/go-chi/cors
go get -u github.com/go-chi/httprate

# Adapter para Lambda
go get -u github.com/awslabs/aws-lambda-go-api-proxy/chi
go get -u github.com/aws/aws-lambda-go/lambda
```

### main.go Básico

```go
// cmd/lambda/main.go
package main

import (
    "context"
    "log"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

var chiLambda *chiadapter.ChiLambda

// init executa antes do handler ser chamado (warm-up)
func init() {
    log.Println("Cold start - initializing router")

    r := chi.NewRouter()

    // Middlewares globais
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Definir rotas
    setupRoutes(r)

    // Criar adapter Lambda
    chiLambda = chiadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
    lambda.Start(Handler)
}
```

---

## Routing e Organização

### 1. Organização de Rotas

```go
// cmd/lambda/routes.go
package main

import (
    "github.com/go-chi/chi/v5"
    "myapp/internal/handler"
    "myapp/internal/middleware"
)

func setupRoutes(r chi.Router) {
    // Rotas públicas
    r.Group(func(r chi.Router) {
        r.Get("/health", handler.HealthCheck)
        r.Post("/auth/login", handler.Login)
        r.Post("/auth/register", handler.Register)
    })

    // API v1 - rotas autenticadas
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(middleware.Auth) // middleware de autenticação

        // Rotas de usuários
        r.Route("/users", func(r chi.Router) {
            r.Get("/", handler.ListUsers)
            r.Post("/", handler.CreateUser)
            r.Get("/{userID}", handler.GetUser)
            r.Put("/{userID}", handler.UpdateUser)
            r.Delete("/{userID}", handler.DeleteUser)
        })

        // Rotas de produtos
        r.Route("/products", func(r chi.Router) {
            r.Get("/", handler.ListProducts)
            r.Post("/", handler.CreateProduct)
            r.Get("/{productID}", handler.GetProduct)
        })

        // Rotas admin (requer permissão especial)
        r.Route("/admin", func(r chi.Router) {
            r.Use(middleware.RequireAdmin)
            r.Get("/stats", handler.GetStats)
            r.Get("/users", handler.AdminListUsers)
        })
    })
}
```

### 2. Rotas com Subrouter (Melhor para Lambdas Separadas)

```go
// internal/handler/user_routes.go
package handler

import (
    "github.com/go-chi/chi/v5"
)

// NewUserRouter cria um subrouter para usuários
func NewUserRouter() chi.Router {
    r := chi.NewRouter()

    r.Get("/", listUsers)
    r.Post("/", createUser)
    r.Get("/{id}", getUser)
    r.Put("/{id}", updateUser)
    r.Delete("/{id}", deleteUser)

    return r
}

// Usar no main.go
func setupRoutes(r chi.Router) {
    r.Mount("/api/v1/users", handler.NewUserRouter())
    r.Mount("/api/v1/products", handler.NewProductRouter())
}
```

### 3. Path Parameters

```go
// ✅ Bom - extrair parâmetros de URL
func GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")

    if userID == "" {
        http.Error(w, "user ID is required", http.StatusBadRequest)
        return
    }

    // usar userID
}

// ✅ Bom - múltiplos parâmetros
r.Get("/tenants/{tenantID}/users/{userID}", handler.GetTenantUser)

func GetTenantUser(w http.ResponseWriter, r *http.Request) {
    tenantID := chi.URLParam(r, "tenantID")
    userID := chi.URLParam(r, "userID")
    // processar...
}
```

### 4. Query Parameters

```go
// ✅ Bom - parsing de query params
func ListUsers(w http.ResponseWriter, r *http.Request) {
    // GET /users?page=1&limit=10&status=active

    page := r.URL.Query().Get("page")
    limit := r.URL.Query().Get("limit")
    status := r.URL.Query().Get("status")

    // converter e validar
    pageNum, _ := strconv.Atoi(page)
    if pageNum < 1 {
        pageNum = 1
    }

    // processar...
}
```

---

## Middlewares

### 1. Middlewares Built-in do Chi

```go
import "github.com/go-chi/chi/v5/middleware"

r.Use(middleware.RequestID)      // Adiciona request ID
r.Use(middleware.RealIP)          // Detecta IP real
r.Use(middleware.Logger)          // Log de requests
r.Use(middleware.Recoverer)       // Recover de panics
r.Use(middleware.Timeout(60 * time.Second)) // Timeout
r.Use(middleware.Compress(5))     // Compressão gzip
r.Use(middleware.StripSlashes)    // Remove trailing slash
r.Use(middleware.Heartbeat("/ping")) // Health check
```

### 2. Middleware Customizado: Autenticação

```go
// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Auth middleware de autenticação JWT
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        // Extrair token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
            return
        }

        token := parts[1]

        // Validar token (exemplo simplificado)
        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Adicionar userID ao context
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// GetUserID extrai userID do context
func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(UserIDKey).(string)
    return userID, ok
}

func validateToken(token string) (string, error) {
    // Implementar validação JWT aqui
    // Por exemplo, usando github.com/golang-jwt/jwt
    return "user-123", nil
}
```

### 3. Middleware de CORS

```go
// internal/middleware/cors.go
package middleware

import (
    "github.com/go-chi/cors"
    "net/http"
)

func CORS() func(http.Handler) http.Handler {
    return cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    })
}

// Uso
r.Use(middleware.CORS())
```

### 4. Middleware de Logging Customizado

```go
// internal/middleware/logger.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

type responseWriter struct {
    http.ResponseWriter
    status int
    size   int
}

func (rw *responseWriter) WriteHeader(status int) {
    rw.status = status
    rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    size, err := rw.ResponseWriter.Write(b)
    rw.size += size
    return size, err
}

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rw := &responseWriter{
            ResponseWriter: w,
            status:         200,
        }

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        log.Printf(
            "[%s] %s %s - Status: %d - Duration: %v - Size: %d bytes",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            rw.status,
            duration,
            rw.size,
        )
    })
}
```

### 5. Middleware de Rate Limiting

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "github.com/go-chi/httprate"
    "net/http"
    "time"
)

// RateLimit limita requests por IP
func RateLimit() func(http.Handler) http.Handler {
    return httprate.LimitByIP(100, 1*time.Minute)
}

// RateLimitByUser limita por usuário autenticado
func RateLimitByUser() func(http.Handler) http.Handler {
    return httprate.Limit(
        100,                    // requests
        1*time.Minute,          // janela de tempo
        httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
            userID, ok := GetUserID(r.Context())
            if !ok {
                return "", httprate.ErrKeyNotFound
            }
            return userID, nil
        }),
    )
}

// Uso
r.Route("/api/v1", func(r chi.Router) {
    r.Use(middleware.RateLimit())
    // rotas...
})
```

### 6. Middleware de Tenant Isolation

```go
// internal/middleware/tenant.go
package middleware

import (
    "context"
    "net/http"
)

const TenantIDKey contextKey = "tenantID"

// TenantIsolation extrai e valida tenant ID
func TenantIsolation(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := r.Header.Get("X-Tenant-ID")

        if tenantID == "" {
            http.Error(w, "Tenant ID required", http.StatusBadRequest)
            return
        }

        // Validar tenant existe
        if !isValidTenant(tenantID) {
            http.Error(w, "Invalid tenant", http.StatusForbidden)
            return
        }

        ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetTenantID(ctx context.Context) (string, bool) {
    tenantID, ok := ctx.Value(TenantIDKey).(string)
    return tenantID, ok
}
```

---

## Handlers

### 1. Estrutura de Handler

```go
// internal/handler/user_handler.go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "myapp/internal/service"
    "myapp/pkg/response"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// ListUsers retorna lista de usuários
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    users, err := h.userService.ListUsers(ctx)
    if err != nil {
        response.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response.JSON(w, users, http.StatusOK)
}

// GetUser retorna um usuário específico
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := chi.URLParam(r, "userID")

    user, err := h.userService.GetUser(ctx, userID)
    if err != nil {
        response.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    response.JSON(w, user, http.StatusOK)
}

// CreateUser cria novo usuário
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var input service.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.userService.CreateUser(ctx, input)
    if err != nil {
        response.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    response.JSON(w, user, http.StatusCreated)
}
```

### 2. Response Helpers

```go
// pkg/response/json.go
package response

import (
    "encoding/json"
    "net/http"
)

// JSON envia resposta JSON
func JSON(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Error envia erro JSON
func Error(w http.ResponseWriter, message string, status int) {
    JSON(w, map[string]string{
        "error": message,
    }, status)
}

// Success envia mensagem de sucesso
func Success(w http.ResponseWriter, message string) {
    JSON(w, map[string]string{
        "message": message,
    }, http.StatusOK)
}

// ErrorWithDetails envia erro com detalhes
func ErrorWithDetails(w http.ResponseWriter, message string, details interface{}, status int) {
    JSON(w, map[string]interface{}{
        "error":   message,
        "details": details,
    }, status)
}
```

### 3. Handler com Validação

```go
// internal/handler/user_handler.go

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=3"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validar
    if err := validator.Validate(req); err != nil {
        response.ErrorWithDetails(w, "Validation failed", err, http.StatusBadRequest)
        return
    }

    user, err := h.userService.CreateUser(ctx, req.Email, req.Name, req.Password)
    if err != nil {
        response.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response.JSON(w, user, http.StatusCreated)
}
```

---

## Integração com AWS Lambda

### 1. Main.go para Lambda

```go
// cmd/lambda/main.go
package main

import (
    "context"
    "log"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "myapp/internal/handler"
    "myapp/internal/repository"
    "myapp/internal/service"
    custommw "myapp/internal/middleware"
)

var chiLambda *chiadapter.ChiLambda

func init() {
    log.Println("Initializing Lambda function...")

    // Inicializar dependências
    db := initDatabase()

    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // Criar router
    r := chi.NewRouter()

    // Middlewares
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(custommw.CORS())

    // Rotas
    r.Get("/health", handler.HealthCheck)

    r.Route("/api/v1", func(r chi.Router) {
        r.Use(custommw.Auth)
        r.Mount("/users", userHandler.Routes())
    })

    // Criar adapter
    chiLambda = chiadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
    // Detectar se está em Lambda ou local
    if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
        lambda.Start(Handler)
    } else {
        // Modo local para desenvolvimento
        runLocalServer()
    }
}

func runLocalServer() {
    r := chi.NewRouter()
    setupRoutes(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on port %s", port)
    http.ListenAndServe(":"+port, r)
}
```

### 2. Múltiplas Lambdas (Microserviços)

```
casa-digital-back/
├── cmd/
│   ├── auth-lambda/
│   │   └── main.go          # Lambda de autenticação
│   ├── users-lambda/
│   │   └── main.go          # Lambda de usuários
│   └── products-lambda/
│       └── main.go          # Lambda de produtos
```

```go
// cmd/users-lambda/main.go
package main

import (
    "context"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"

    "myapp/internal/handler"
)

var chiLambda *chiadapter.ChiLambda

func init() {
    userHandler := handler.NewUserHandler()
    chiLambda = chiadapter.New(userHandler.Routes())
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
    lambda.Start(Handler)
}
```

---

## Validação e Error Handling

### 1. Validação com go-playground/validator

```go
// pkg/validator/validator.go
package validator

import (
    "fmt"
    "strings"

    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func Validate(s interface{}) []ValidationError {
    err := validate.Struct(s)
    if err == nil {
        return nil
    }

    var errors []ValidationError
    for _, err := range err.(validator.ValidationErrors) {
        errors = append(errors, ValidationError{
            Field:   strings.ToLower(err.Field()),
            Message: getErrorMsg(err),
        })
    }

    return errors
}

func getErrorMsg(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return fmt.Sprintf("%s is required", err.Field())
    case "email":
        return "Invalid email format"
    case "min":
        return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
    case "max":
        return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
    default:
        return fmt.Sprintf("%s is invalid", err.Field())
    }
}
```

### 2. Error Handler Centralizado

```go
// pkg/errors/errors.go
package errors

import "net/http"

type AppError struct {
    Message    string `json:"message"`
    StatusCode int    `json:"-"`
    Code       string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Construtores de erros comuns
func NewBadRequest(message string) *AppError {
    return &AppError{
        Message:    message,
        StatusCode: http.StatusBadRequest,
        Code:       "BAD_REQUEST",
    }
}

func NewNotFound(resource string) *AppError {
    return &AppError{
        Message:    resource + " not found",
        StatusCode: http.StatusNotFound,
        Code:       "NOT_FOUND",
    }
}

func NewUnauthorized(message string) *AppError {
    return &AppError{
        Message:    message,
        StatusCode: http.StatusUnauthorized,
        Code:       "UNAUTHORIZED",
    }
}

func NewInternal(message string) *AppError {
    return &AppError{
        Message:    message,
        StatusCode: http.StatusInternalServerError,
        Code:       "INTERNAL_ERROR",
    }
}
```

### 3. Middleware de Error Handling

```go
// internal/middleware/error_handler.go
package middleware

import (
    "log"
    "net/http"

    "myapp/pkg/errors"
    "myapp/pkg/response"
)

func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                response.Error(w, "Internal server error", http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

---

## Testes

### 1. Teste de Handler

```go
// internal/handler/user_handler_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock do service
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUserHandler_GetUser(t *testing.T) {
    // Setup
    mockService := new(MockUserService)
    handler := NewUserHandler(mockService)

    user := &User{
        ID:    "123",
        Email: "test@example.com",
        Name:  "Test User",
    }

    mockService.On("GetUser", mock.Anything, "123").Return(user, nil)

    // Request
    req := httptest.NewRequest("GET", "/users/123", nil)
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("userID", "123")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

    // Response recorder
    w := httptest.NewRecorder()

    // Execute
    handler.GetUser(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var response User
    json.NewDecoder(w.Body).Decode(&response)
    assert.Equal(t, user.ID, response.ID)
    assert.Equal(t, user.Email, response.Email)

    mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser(t *testing.T) {
    mockService := new(MockUserService)
    handler := NewUserHandler(mockService)

    input := CreateUserRequest{
        Email:    "new@example.com",
        Name:     "New User",
        Password: "password123",
    }

    expectedUser := &User{
        ID:    "new-123",
        Email: input.Email,
        Name:  input.Name,
    }

    mockService.On("CreateUser", mock.Anything, input.Email, input.Name, input.Password).
        Return(expectedUser, nil)

    body, _ := json.Marshal(input)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    w := httptest.NewRecorder()

    handler.CreateUser(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### 2. Teste de Integração com Router

```go
// internal/handler/integration_test.go
package handler

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
)

func TestRouter_Integration(t *testing.T) {
    // Setup router completo
    r := chi.NewRouter()

    userHandler := NewUserHandler(mockUserService)
    r.Mount("/users", userHandler.Routes())

    // Teste GET
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Teste POST
    req = httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### 3. Teste de Middleware

```go
// internal/middleware/auth_test.go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAuth_ValidToken(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID, ok := GetUserID(r.Context())
        assert.True(t, ok)
        assert.Equal(t, "user-123", userID)
        w.WriteHeader(http.StatusOK)
    })

    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Bearer valid-token")

    w := httptest.NewRecorder()

    Auth(handler).ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuth_MissingToken(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("Handler should not be called")
    })

    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()

    Auth(handler).ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

---

## Performance e Otimização

### 1. Minimizar Bundle Size

```bash
# Build otimizado para Lambda
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -ldflags="-s -w" \
  -trimpath \
  -o bootstrap \
  cmd/lambda/main.go

# Comprimir
zip function.zip bootstrap
```

### 2. Otimizar Cold Start

```go
// ❌ Ruim - inicialização dentro do handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    db := connectDatabase() // NÃO FAÇA ISSO!
    // processar...
}

// ✅ Bom - inicialização no init()
var db *sql.DB

func init() {
    db = connectDatabase() // Uma vez só
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // usar db global
}
```

### 3. Connection Pooling

```go
// ✅ Bom - configurar pool de conexões
func initDatabase() *sql.DB {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }

    // Configurar para ambiente serverless
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(2)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(1 * time.Minute)

    return db
}
```

### 4. Caching

```go
// pkg/cache/cache.go
package cache

import (
    "sync"
    "time"
)

type Cache struct {
    mu    sync.RWMutex
    items map[string]*item
}

type item struct {
    value      interface{}
    expiration int64
}

func New() *Cache {
    return &Cache{
        items: make(map[string]*item),
    }
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = &item{
        value:      value,
        expiration: time.Now().Add(duration).UnixNano(),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.items[key]
    if !exists {
        return nil, false
    }

    if time.Now().UnixNano() > item.expiration {
        delete(c.items, key)
        return nil, false
    }

    return item.value, true
}
```

### 5. Métricas e Monitoring

```go
// internal/middleware/metrics.go
package middleware

import (
    "net/http"
    "time"
)

func Metrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrapper para capturar status
        rw := &responseWriter{ResponseWriter: w, status: 200}

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        // Enviar métricas para CloudWatch, DataDog, etc.
        recordMetric("http.request.duration", duration.Milliseconds())
        recordMetric("http.request.count", 1)
        recordMetric("http.request.status."+string(rw.status), 1)
    })
}
```

---

## Checklist de Boas Práticas

- [ ] Rotas organizadas com `Route()` e `Group()`
- [ ] Middlewares aplicados apropriadamente (global vs específico)
- [ ] Path parameters validados
- [ ] Request body validado
- [ ] Erros tratados com status codes corretos
- [ ] Responses padronizadas (JSON helpers)
- [ ] Autenticação implementada corretamente
- [ ] CORS configurado apropriadamente
- [ ] Rate limiting em rotas públicas
- [ ] Logging estruturado
- [ ] Testes unitários para handlers
- [ ] Testes para middlewares
- [ ] init() usado para warm-up
- [ ] Connection pooling configurado
- [ ] Build otimizado (-ldflags="-s -w")
- [ ] Métricas e monitoring implementados

---

## Recursos Adicionais

### Documentação Oficial
- [Chi Router](https://github.com/go-chi/chi)
- [Chi Middlewares](https://github.com/go-chi/chi/tree/master/middleware)
- [AWS Lambda Go Adapter](https://github.com/awslabs/aws-lambda-go-api-proxy)

### Pacotes Úteis
- `github.com/go-chi/cors` - CORS
- `github.com/go-chi/httprate` - Rate limiting
- `github.com/go-chi/render` - JSON rendering
- `github.com/go-playground/validator/v10` - Validação
- `github.com/golang-jwt/jwt/v5` - JWT

### Exemplos
```bash
# Instalar dependências
go get -u github.com/go-chi/chi/v5
go get -u github.com/go-chi/cors
go get -u github.com/go-chi/httprate
go get -u github.com/awslabs/aws-lambda-go-api-proxy/chi
go get -u github.com/aws/aws-lambda-go/lambda
go get -u github.com/go-playground/validator/v10
```

---

## Conclusão

Chi é uma excelente escolha para APIs em AWS Lambda por ser:
- **Leve**: Mínimo overhead e bundle size pequeno
- **Idiomático**: Segue padrões Go e usa stdlib
- **Flexível**: Fácil compor middlewares e rotas
- **Performático**: Cold start rápido e execução eficiente

Siga estas práticas para construir APIs robustas, manuteníveis e performáticas!
