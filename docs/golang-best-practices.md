# Guia de Boas Práticas em Go

Este documento serve como referência para desenvolvimento em Go seguindo princípios de Clean Code, SOLID, DRY e outras melhores práticas da linguagem.

---

## Índice

1. [Princípios Fundamentais](#princípios-fundamentais)
2. [Clean Code em Go](#clean-code-em-go)
3. [Princípios SOLID](#princípios-solid)
4. [Estrutura de Projetos](#estrutura-de-projetos)
5. [Nomenclatura e Convenções](#nomenclatura-e-convenções)
6. [Tratamento de Erros](#tratamento-de-erros)
7. [Concorrência](#concorrência)
8. [Testes](#testes)
9. [Performance e Otimização](#performance-e-otimização)
10. [Segurança](#segurança)

---

## Princípios Fundamentais

### DRY (Don't Repeat Yourself) - NÃO REPITA CÓDIGO

**REGRA FUNDAMENTAL: JAMAIS DUPLIQUE CÓDIGO**

A duplicação de código é um dos maiores problemas em manutenção de software. Sempre que você encontrar código repetido, extraia-o para uma função reutilizável.

#### Por que evitar repetição?
- **Manutenção**: Bugs precisam ser corrigidos em apenas um lugar
- **Consistência**: Mudanças são aplicadas uniformemente
- **Legibilidade**: Código mais limpo e fácil de entender
- **Testabilidade**: Testa-se uma vez, funciona em todos os lugares

#### Princípios para evitar duplicação:
1. **Extraia lógica comum** em funções/métodos reutilizáveis
2. **Use composição** ao invés de herança
3. **Crie utilitários** para operações repetidas
4. **Abstraia padrões** recorrentes em interfaces
5. **Refatore imediatamente** quando detectar duplicação

```go
// ❌ RUIM - código duplicado (NUNCA FAÇA ISSO!)
func ProcessUser(user User) error {
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email")
    }
    if len(user.Email) > 255 {
        return errors.New("email too long")
    }
    // processar usuário
    return nil
}

func ProcessAdmin(admin Admin) error {
    if admin.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(admin.Email, "@") {
        return errors.New("invalid email")
    }
    if len(admin.Email) > 255 {
        return errors.New("email too long")
    }
    // processar admin
    return nil
}

func ProcessCustomer(customer Customer) error {
    if customer.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(customer.Email, "@") {
        return errors.New("invalid email")
    }
    if len(customer.Email) > 255 {
        return errors.New("email too long")
    }
    // processar customer
    return nil
}

// ✅ BOM - extrai validação comum (SEMPRE FAÇA ISSO!)
func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email")
    }
    if len(email) > 255 {
        return errors.New("email too long")
    }
    return nil
}

func ProcessUser(user User) error {
    if err := ValidateEmail(user.Email); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }
    // processar usuário
    return nil
}

func ProcessAdmin(admin Admin) error {
    if err := ValidateEmail(admin.Email); err != nil {
        return fmt.Errorf("admin validation failed: %w", err)
    }
    // processar admin
    return nil
}

func ProcessCustomer(customer Customer) error {
    if err := ValidateEmail(customer.Email); err != nil {
        return fmt.Errorf("customer validation failed: %w", err)
    }
    // processar customer
    return nil
}
```

#### Exemplo Avançado: Extraindo lógica de negócio comum

```go
// ❌ RUIM - lógica de desconto duplicada
func CalculateOrderDiscount(order Order) float64 {
    if order.Total >= 1000 {
        return order.Total * 0.15
    } else if order.Total >= 500 {
        return order.Total * 0.10
    } else if order.Total >= 100 {
        return order.Total * 0.05
    }
    return 0
}

func CalculateInvoiceDiscount(invoice Invoice) float64 {
    if invoice.Total >= 1000 {
        return invoice.Total * 0.15
    } else if invoice.Total >= 500 {
        return invoice.Total * 0.10
    } else if invoice.Total >= 100 {
        return invoice.Total * 0.05
    }
    return 0
}

// ✅ BOM - função única reutilizável
func CalculateDiscountByAmount(amount float64) float64 {
    switch {
    case amount >= 1000:
        return amount * 0.15
    case amount >= 500:
        return amount * 0.10
    case amount >= 100:
        return amount * 0.05
    default:
        return 0
    }
}

func ProcessOrder(order Order) {
    order.Discount = CalculateDiscountByAmount(order.Total)
}

func ProcessInvoice(invoice Invoice) {
    invoice.Discount = CalculateDiscountByAmount(invoice.Total)
}
```

#### Detectando código duplicado

**SINAIS DE ALERTA:**
- Você está copiando e colando código
- Mesma lógica aparece em múltiplas funções
- Mudanças exigem atualização em vários lugares
- Testes similares para funções diferentes

**AÇÃO IMEDIATA:**
1. Identifique o código comum
2. Extraia para uma função separada
3. Substitua todas as ocorrências
4. Teste para garantir que funciona
5. Delete o código duplicado

**LEMBRE-SE: A única duplicação aceitável é ZERO duplicação!**

### KISS (Keep It Simple, Stupid)
- Prefira soluções simples e diretas
- Evite over-engineering
- Use as features da linguagem de forma idiomática

```go
// ❌ Ruim - complexo desnecessariamente
func IsEven(n int) bool {
    return n&1 == 0
}

// ✅ Bom - simples e claro
func IsEven(n int) bool {
    return n%2 == 0
}
```

### YAGNI (You Aren't Gonna Need It)
- Não implemente funcionalidades que você "pode precisar no futuro"
- Implemente apenas o necessário agora
- Refatore quando realmente precisar

---

## Clean Code em Go

### 1. Funções Pequenas e Focadas

```go
// ❌ Ruim - função fazendo muitas coisas
func ProcessOrder(order Order) error {
    // validação
    if order.Total <= 0 {
        return errors.New("invalid total")
    }

    // calcular desconto
    discount := 0.0
    if order.Total > 100 {
        discount = order.Total * 0.1
    }

    // salvar no banco
    db.Save(order)

    // enviar email
    sendEmail(order.CustomerEmail)

    return nil
}

// ✅ Bom - separado em funções específicas
func ProcessOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    order = applyDiscount(order)

    if err := saveOrder(order); err != nil {
        return fmt.Errorf("failed to save order: %w", err)
    }

    if err := notifyCustomer(order); err != nil {
        return fmt.Errorf("failed to notify customer: %w", err)
    }

    return nil
}

func validateOrder(order Order) error {
    if order.Total <= 0 {
        return errors.New("total must be positive")
    }
    return nil
}

func applyDiscount(order Order) Order {
    if order.Total > 100 {
        order.Discount = order.Total * 0.1
    }
    return order
}
```

### 2. Nomes Significativos

```go
// ❌ Ruim
func calc(a, b int) int {
    return a + b
}

var d int // dias
var u User // usuário

// ✅ Bom
func calculateTotal(price, quantity int) int {
    return price * quantity
}

var daysSinceLastLogin int
var currentUser User
```

### 3. Comentários Úteis

```go
// ❌ Ruim - comentário óbvio
// Adiciona 1 ao contador
counter++

// ✅ Bom - explica o "porquê", não o "o quê"
// Incrementa o contador de tentativas para implementar rate limiting
// após 5 tentativas, o usuário será temporariamente bloqueado
attemptCounter++

// ✅ Bom - documentação de pacote/função pública
// CalculateMonthlyInterest calcula os juros mensais baseado na taxa anual.
// A taxa deve ser fornecida como decimal (ex: 0.05 para 5%).
// Retorna erro se a taxa for negativa.
func CalculateMonthlyInterest(principal float64, annualRate float64) (float64, error) {
    if annualRate < 0 {
        return 0, errors.New("annual rate cannot be negative")
    }
    monthlyRate := annualRate / 12
    return principal * monthlyRate, nil
}
```

---

## Princípios SOLID

### S - Single Responsibility Principle (SRP)

Cada tipo/função deve ter apenas uma responsabilidade.

```go
// ❌ Ruim - UserService faz muitas coisas
type UserService struct {
    db *sql.DB
}

func (s *UserService) CreateUser(user User) error {
    // validação
    if user.Email == "" {
        return errors.New("email required")
    }

    // salvar no banco
    _, err := s.db.Exec("INSERT INTO users ...")

    // enviar email
    smtp.SendEmail(user.Email, "Welcome!")

    // log
    log.Printf("User created: %s", user.Email)

    return err
}

// ✅ Bom - responsabilidades separadas
type UserValidator struct{}

func (v *UserValidator) Validate(user User) error {
    if user.Email == "" {
        return errors.New("email required")
    }
    return nil
}

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(user User) error {
    _, err := r.db.Exec("INSERT INTO users ...")
    return err
}

type EmailNotifier struct {
    smtpClient *smtp.Client
}

func (n *EmailNotifier) SendWelcomeEmail(email string) error {
    return n.smtpClient.SendEmail(email, "Welcome!")
}

type UserService struct {
    validator  *UserValidator
    repository *UserRepository
    notifier   *EmailNotifier
}

func (s *UserService) CreateUser(user User) error {
    if err := s.validator.Validate(user); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if err := s.repository.Save(user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }

    if err := s.notifier.SendWelcomeEmail(user.Email); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}
```

### O - Open/Closed Principle (OCP)

Aberto para extensão, fechado para modificação.

```go
// ❌ Ruim - precisa modificar o código para adicionar novos tipos
func CalculateArea(shape string, dimensions map[string]float64) float64 {
    switch shape {
    case "circle":
        return math.Pi * dimensions["radius"] * dimensions["radius"]
    case "rectangle":
        return dimensions["width"] * dimensions["height"]
    default:
        return 0
    }
}

// ✅ Bom - usa interfaces para extensão
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Adicionar novo shape não requer modificar código existente
type Triangle struct {
    Base   float64
    Height float64
}

func (t Triangle) Area() float64 {
    return 0.5 * t.Base * t.Height
}

func CalculateTotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, shape := range shapes {
        total += shape.Area()
    }
    return total
}
```

### L - Liskov Substitution Principle (LSP)

Subtipos devem ser substituíveis por seus tipos base.

```go
// ✅ Bom - todos os tipos implementam a interface completamente
type Storage interface {
    Save(key string, value []byte) error
    Load(key string) ([]byte, error)
}

type DiskStorage struct {
    basePath string
}

func (d *DiskStorage) Save(key string, value []byte) error {
    return os.WriteFile(filepath.Join(d.basePath, key), value, 0644)
}

func (d *DiskStorage) Load(key string) ([]byte, error) {
    return os.ReadFile(filepath.Join(d.basePath, key))
}

type MemoryStorage struct {
    data map[string][]byte
}

func (m *MemoryStorage) Save(key string, value []byte) error {
    m.data[key] = value
    return nil
}

func (m *MemoryStorage) Load(key string) ([]byte, error) {
    if val, ok := m.data[key]; ok {
        return val, nil
    }
    return nil, errors.New("key not found")
}

// Ambas implementações podem ser usadas de forma intercambiável
func ProcessData(storage Storage, key string, data []byte) error {
    if err := storage.Save(key, data); err != nil {
        return err
    }

    loaded, err := storage.Load(key)
    if err != nil {
        return err
    }

    // processar loaded data
    _ = loaded
    return nil
}
```

### I - Interface Segregation Principle (ISP)

Clientes não devem ser forçados a depender de interfaces que não usam.

```go
// ❌ Ruim - interface muito grande
type Worker interface {
    Work()
    Eat()
    Sleep()
    GetSalary() float64
    TakePaidLeave()
}

// ✅ Bom - interfaces segregadas
type Workable interface {
    Work()
}

type Eatable interface {
    Eat()
}

type Sleepable interface {
    Sleep()
}

type Payable interface {
    GetSalary() float64
}

type LeaveManager interface {
    TakePaidLeave()
}

// Tipos implementam apenas o que precisam
type Robot struct{}

func (r Robot) Work() {
    // robots trabalham
}

type Human struct {
    salary float64
}

func (h Human) Work() {}
func (h Human) Eat() {}
func (h Human) Sleep() {}
func (h Human) GetSalary() float64 { return h.salary }
func (h Human) TakePaidLeave() {}

// Funções usam apenas interfaces necessárias
func AssignTask(w Workable) {
    w.Work()
}

func ProcessPayroll(p Payable) {
    salary := p.GetSalary()
    // processar pagamento
    _ = salary
}
```

### D - Dependency Inversion Principle (DIP)

Dependa de abstrações, não de implementações concretas.

```go
// ❌ Ruim - depende de implementação concreta
type EmailService struct {
    smtpHost string
    smtpPort int
}

func (e *EmailService) Send(to, message string) error {
    // lógica específica de SMTP
    return nil
}

type UserService struct {
    emailService *EmailService // dependência concreta
}

func (u *UserService) Register(user User) error {
    // registrar usuário
    return u.emailService.Send(user.Email, "Welcome!")
}

// ✅ Bom - depende de abstração
type Notifier interface {
    Notify(recipient, message string) error
}

type EmailNotifier struct {
    smtpHost string
    smtpPort int
}

func (e *EmailNotifier) Notify(recipient, message string) error {
    // lógica específica de SMTP
    return nil
}

type SMSNotifier struct {
    apiKey string
}

func (s *SMSNotifier) Notify(recipient, message string) error {
    // lógica específica de SMS
    return nil
}

type UserService struct {
    notifier Notifier // dependência abstrata
}

func (u *UserService) Register(user User) error {
    // registrar usuário
    return u.notifier.Notify(user.Email, "Welcome!")
}

// Facilita testes e troca de implementação
func NewUserService(notifier Notifier) *UserService {
    return &UserService{notifier: notifier}
}
```

---

## Estrutura de Projetos

### Layout Padrão de Projeto Go

```
projeto/
├── cmd/
│   └── app/
│       └── main.go              # Entry point da aplicação
├── internal/
│   ├── domain/                  # Entidades de domínio
│   │   ├── user.go
│   │   └── order.go
│   ├── handler/                 # Handlers HTTP/gRPC
│   │   ├── user_handler.go
│   │   └── order_handler.go
│   ├── service/                 # Lógica de negócio
│   │   ├── user_service.go
│   │   └── order_service.go
│   └── repository/              # Acesso a dados
│       ├── user_repository.go
│       └── order_repository.go
├── pkg/                         # Bibliotecas reutilizáveis
│   ├── logger/
│   └── validator/
├── config/                      # Arquivos de configuração
│   └── config.go
├── migrations/                  # Migrações de banco de dados
├── scripts/                     # Scripts auxiliares
├── docs/                        # Documentação
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Organização por Camadas

```go
// domain/user.go - Entidades de domínio
package domain

import "time"

type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// repository/user_repository.go - Interface de persistência
package repository

import "myapp/internal/domain"

type UserRepository interface {
    Create(user *domain.User) error
    FindByID(id string) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id string) error
}

// repository/postgres/user_repository.go - Implementação concreta
package postgres

import (
    "database/sql"
    "myapp/internal/domain"
    "myapp/internal/repository"
)

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    query := `INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.CreatedAt)
    return err
}

// service/user_service.go - Lógica de negócio
package service

import (
    "errors"
    "myapp/internal/domain"
    "myapp/internal/repository"
    "time"
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(email, name string) (*domain.User, error) {
    // Validação
    if email == "" {
        return nil, errors.New("email is required")
    }

    // Verificar duplicidade
    existing, _ := s.repo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    // Criar usuário
    user := &domain.User{
        ID:        generateID(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
    }

    if err := s.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// handler/user_handler.go - HTTP handlers
package handler

import (
    "encoding/json"
    "myapp/internal/service"
    "net/http"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service: service}
}

type CreateUserRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.service.CreateUser(req.Email, req.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

---

## Nomenclatura e Convenções

### 1. Nomes de Pacotes

```go
// ❌ Ruim
package userPackage
package user_utils
package myUserStuff

// ✅ Bom
package user
package http
package postgres
```

### 2. Nomes de Variáveis

```go
// ❌ Ruim
var MyVar int
var my_var int
var myVeryLongVariableNameThatDescribesEverything int

// ✅ Bom
var count int
var userID string
var maxRetries int

// Contexto curto permite nomes curtos
for i := 0; i < 10; i++ {
    // i é aceitável aqui
}

// Contexto longo requer nomes descritivos
var activeUserCount int
```

### 3. Nomes de Funções e Métodos

```go
// ✅ Funções exportadas começam com maiúscula
func CreateUser(user User) error

// ✅ Funções internas começam com minúscula
func validateEmail(email string) error

// ✅ Getters não usam "Get" prefix
type User struct {
    name string
}

func (u *User) Name() string {
    return u.name
}

// ✅ Setters usam "Set" prefix
func (u *User) SetName(name string) {
    u.name = name
}

// ✅ Métodos booleanos usam "Is", "Has", "Can"
func (u *User) IsActive() bool
func (u *User) HasPermission(perm string) bool
func (u *User) CanEdit() bool
```

### 4. Nomes de Interfaces

```go
// ✅ Interfaces de um método usam sufixo "-er"
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// ✅ Interfaces maiores usam nomes substantivos
type UserRepository interface {
    Create(user User) error
    FindByID(id string) (User, error)
    Update(user User) error
    Delete(id string) error
}
```

### 5. Constantes e Enumerações

```go
// ✅ Constantes exportadas
const MaxConnections = 100
const DefaultTimeout = 30 * time.Second

// ✅ Enumerações usando iota
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
    StatusDeleted
)

// ✅ Grupo de constantes relacionadas
const (
    StatusPending  = "pending"
    StatusActive   = "active"
    StatusInactive = "inactive"
)
```

---

## Tratamento de Erros

### 1. Sempre Verifique Erros

```go
// ❌ Ruim - ignora erro
file, _ := os.Open("file.txt")

// ✅ Bom - trata erro apropriadamente
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
```

### 2. Erros Customizados

```go
// ✅ Bom - erro customizado com contexto
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func ValidateUser(user User) error {
    if user.Email == "" {
        return &ValidationError{
            Field:   "email",
            Message: "email is required",
        }
    }
    return nil
}

// Verificar tipo de erro
err := ValidateUser(user)
if err != nil {
    var valErr *ValidationError
    if errors.As(err, &valErr) {
        // tratar erro de validação especificamente
        log.Printf("Validation failed on field: %s", valErr.Field)
    }
}
```

### 3. Wrapping de Erros

```go
// ✅ Bom - usa %w para wrapping
func ProcessOrder(orderID string) error {
    order, err := fetchOrder(orderID)
    if err != nil {
        return fmt.Errorf("process order failed: %w", err)
    }

    if err := validateOrder(order); err != nil {
        return fmt.Errorf("order validation failed: %w", err)
    }

    return nil
}

// Permite verificar erro original
err := ProcessOrder("123")
if errors.Is(err, sql.ErrNoRows) {
    // tratar caso específico de não encontrado
}
```

### 4. Erros de Sentinela

```go
// ✅ Bom - erros de sentinela para casos específicos
var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrInvalidInput  = errors.New("invalid input")
)

func GetUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

// Uso
user, err := GetUser("123")
if errors.Is(err, ErrNotFound) {
    // tratar não encontrado
}
```

### 5. Defer para Cleanup

```go
// ✅ Bom - usar defer para cleanup
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // sempre será executado

    // processar arquivo
    return nil
}

// ✅ Bom - capturar erro de Close
func ProcessFile(filename string) (err error) {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := file.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()

    // processar arquivo
    return nil
}
```

---

## Concorrência

### 1. Goroutines com Contexto

```go
// ✅ Bom - usa context para cancelamento
func ProcessData(ctx context.Context, data []string) error {
    errChan := make(chan error, len(data))

    for _, item := range data {
        go func(item string) {
            select {
            case <-ctx.Done():
                errChan <- ctx.Err()
                return
            default:
                // processar item
                errChan <- processItem(item)
            }
        }(item)
    }

    // coletar erros
    for i := 0; i < len(data); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }

    return nil
}
```

### 2. WaitGroup para Sincronização

```go
// ✅ Bom - usa WaitGroup
func ProcessConcurrently(items []string) {
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(item string) {
            defer wg.Done()
            processItem(item)
        }(item)
    }

    wg.Wait()
}
```

### 3. Channels para Comunicação

```go
// ✅ Bom - produtor-consumidor com channels
func ProcessPipeline(input <-chan string) <-chan Result {
    output := make(chan Result)

    go func() {
        defer close(output)
        for item := range input {
            result := processItem(item)
            output <- result
        }
    }()

    return output
}

// Uso
inputChan := make(chan string)
resultChan := ProcessPipeline(inputChan)

// Enviar dados
go func() {
    defer close(inputChan)
    for _, item := range items {
        inputChan <- item
    }
}()

// Receber resultados
for result := range resultChan {
    // usar resultado
}
```

### 4. Mutex para Proteção de Dados

```go
// ✅ Bom - usa mutex para acesso concorrente
type SafeCounter struct {
    mu    sync.RWMutex
    count map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count[key]++
}

func (c *SafeCounter) Value(key string) int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count[key]
}
```

### 5. Worker Pool Pattern

```go
// ✅ Bom - worker pool para limitar concorrência
func WorkerPool(ctx context.Context, jobs <-chan Job, numWorkers int) <-chan Result {
    results := make(chan Result)

    var wg sync.WaitGroup

    // Criar workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case <-ctx.Done():
                    return
                case results <- processJob(job):
                }
            }
        }()
    }

    // Fechar results quando todos workers terminarem
    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

---

## Testes

### 1. Testes Unitários

```go
// user_service_test.go
package service

import (
    "errors"
    "testing"
)

// Mock do repository
type mockUserRepository struct {
    users map[string]*User
}

func (m *mockUserRepository) Create(user *User) error {
    if _, exists := m.users[user.Email]; exists {
        return errors.New("user already exists")
    }
    m.users[user.Email] = user
    return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*User, error) {
    if user, ok := m.users[email]; ok {
        return user, nil
    }
    return nil, ErrNotFound
}

// Teste
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name      string
        email     string
        userName  string
        wantErr   bool
        errMsg    string
    }{
        {
            name:     "valid user",
            email:    "test@example.com",
            userName: "Test User",
            wantErr:  false,
        },
        {
            name:     "empty email",
            email:    "",
            userName: "Test User",
            wantErr:  true,
            errMsg:   "email is required",
        },
        {
            name:     "duplicate email",
            email:    "duplicate@example.com",
            userName: "Test User",
            wantErr:  true,
            errMsg:   "email already exists",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mockUserRepository{
                users: make(map[string]*User),
            }

            // Preparar dados para teste de duplicidade
            if tt.name == "duplicate email" {
                repo.users[tt.email] = &User{Email: tt.email}
            }

            service := NewUserService(repo)

            user, err := service.CreateUser(tt.email, tt.userName)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
                    t.Errorf("expected error message %q, got %q", tt.errMsg, err.Error())
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if user == nil {
                    t.Error("expected user but got nil")
                }
                if user != nil && user.Email != tt.email {
                    t.Errorf("expected email %q, got %q", tt.email, user.Email)
                }
            }
        })
    }
}
```

### 2. Testes de Integração

```go
// integration_test.go
// +build integration

package integration

import (
    "database/sql"
    "testing"
)

func TestUserRepository_Integration(t *testing.T) {
    // Setup database de teste
    db, err := sql.Open("postgres", "postgres://localhost/testdb")
    if err != nil {
        t.Fatalf("failed to connect to database: %v", err)
    }
    defer db.Close()

    // Limpar database antes do teste
    if _, err := db.Exec("TRUNCATE TABLE users"); err != nil {
        t.Fatalf("failed to truncate table: %v", err)
    }

    repo := NewUserRepository(db)

    // Teste real de criação
    user := &User{
        ID:    "123",
        Email: "test@example.com",
        Name:  "Test User",
    }

    if err := repo.Create(user); err != nil {
        t.Errorf("failed to create user: %v", err)
    }

    // Teste real de busca
    found, err := repo.FindByEmail("test@example.com")
    if err != nil {
        t.Errorf("failed to find user: %v", err)
    }

    if found.Email != user.Email {
        t.Errorf("expected email %q, got %q", user.Email, found.Email)
    }
}
```

### 3. Testes de Tabela

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"empty email", "", true},
        {"no @ symbol", "userexample.com", true},
        {"no domain", "user@", true},
        {"no local part", "@example.com", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 4. Benchmarks

```go
func BenchmarkProcessData(b *testing.B) {
    data := generateTestData(1000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ProcessData(data)
    }
}

func BenchmarkProcessDataParallel(b *testing.B) {
    data := generateTestData(1000)

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ProcessData(data)
        }
    })
}
```

### 5. Testify para Assertions

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserCreation(t *testing.T) {
    user, err := CreateUser("test@example.com", "Test")

    require.NoError(t, err) // para se houver erro
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    assert.NotEmpty(t, user.ID)
}
```

---

## Performance e Otimização

### 1. Use Strings Builder para Concatenação

```go
// ❌ Ruim - concatenação ineficiente
func BuildString(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ", "
    }
    return result
}

// ✅ Bom - usa strings.Builder
func BuildString(items []string) string {
    var builder strings.Builder
    for i, item := range items {
        builder.WriteString(item)
        if i < len(items)-1 {
            builder.WriteString(", ")
        }
    }
    return builder.String()
}
```

### 2. Pré-alocar Slices

```go
// ❌ Ruim - sem pré-alocação
func ProcessItems(n int) []int {
    var results []int
    for i := 0; i < n; i++ {
        results = append(results, i*2)
    }
    return results
}

// ✅ Bom - pré-aloca capacidade
func ProcessItems(n int) []int {
    results := make([]int, 0, n)
    for i := 0; i < n; i++ {
        results = append(results, i*2)
    }
    return results
}
```

### 3. Use Ponteiros Apropriadamente

```go
// ✅ Bom - usa ponteiro para structs grandes
type LargeStruct struct {
    Data [1000]int
    // muitos campos...
}

func ProcessLargeStruct(ls *LargeStruct) {
    // evita cópia de toda a struct
}

// ✅ Bom - não usa ponteiro para tipos pequenos
func ProcessInt(n int) int {
    return n * 2
}
```

### 4. Sync.Pool para Reuso de Objetos

```go
// ✅ Bom - usa sync.Pool para buffers
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    buf.Write(data)
    // processar...

    return buf.Bytes()
}
```

### 5. Evite Alocações Desnecessárias

```go
// ❌ Ruim - cria slice a cada iteração
for i := 0; i < 1000; i++ {
    items := []int{1, 2, 3}
    process(items)
}

// ✅ Bom - reusa slice
items := []int{1, 2, 3}
for i := 0; i < 1000; i++ {
    process(items)
}
```

---

## Segurança

### 1. Validação de Entrada

```go
import "html"

// ✅ Bom - valida e sanitiza entrada
func CreatePost(title, content string) error {
    // Validação
    if len(title) == 0 || len(title) > 200 {
        return errors.New("title must be between 1 and 200 characters")
    }

    // Sanitização
    title = html.EscapeString(title)
    content = html.EscapeString(content)

    // processar...
    return nil
}
```

### 2. Prepared Statements para SQL

```go
// ❌ Ruim - vulnerável a SQL injection
func GetUser(email string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
    row := db.QueryRow(query)
    // ...
}

// ✅ Bom - usa prepared statement
func GetUser(email string) (*User, error) {
    query := "SELECT id, email, name FROM users WHERE email = $1"
    row := db.QueryRow(query, email)

    var user User
    err := row.Scan(&user.ID, &user.Email, &user.Name)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 3. Gestão de Secrets

```go
// ❌ Ruim - hardcoded secrets
const apiKey = "sk_live_123456789"

// ✅ Bom - usa variáveis de ambiente
func GetAPIKey() string {
    key := os.Getenv("API_KEY")
    if key == "" {
        log.Fatal("API_KEY environment variable not set")
    }
    return key
}

// ✅ Melhor ainda - usa configuração estruturada
type Config struct {
    APIKey      string `env:"API_KEY,required"`
    DatabaseURL string `env:"DATABASE_URL,required"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

### 4. Rate Limiting

```go
import "golang.org/x/time/rate"

// ✅ Bom - implementa rate limiting
type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
    }
}

func (rl *RateLimiter) Allow() bool {
    return rl.limiter.Allow()
}

// Middleware HTTP
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### 5. Timeout para Operações

```go
// ✅ Bom - sempre usa timeout
func FetchData(url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

---

## Referências e Recursos

### Documentação Oficial
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)

### Livros Recomendados
- "The Go Programming Language" - Alan Donovan & Brian Kernighan
- "Clean Code" - Robert C. Martin
- "Design Patterns in Go" - Mario Castro Contreras

### Ferramentas Essenciais
- `gofmt` - Formatação de código
- `golangci-lint` - Linter agregado
- `go vet` - Análise estática
- `go test -race` - Detector de race conditions
- `pprof` - Profiling de performance

### Convenções de Commit
```
feat: adiciona nova funcionalidade de autenticação
fix: corrige bug no cálculo de impostos
refactor: reorganiza estrutura do handler
test: adiciona testes para UserService
docs: atualiza documentação da API
chore: atualiza dependências
```

---

## Checklist de Revisão de Código

- [ ] Código segue as convenções de nomenclatura do Go
- [ ] Erros são tratados apropriadamente
- [ ] Funções têm responsabilidade única
- [ ] Interfaces são pequenas e focadas
- [ ] Testes cobrem casos importantes
- [ ] Não há código duplicado
- [ ] Documentação está presente para APIs públicas
- [ ] Contextos são usados para cancelamento
- [ ] Recursos são liberados corretamente (defer)
- [ ] SQL usa prepared statements
- [ ] Secrets não estão hardcoded
- [ ] Operações têm timeout apropriado
- [ ] Código é thread-safe quando necessário
- [ ] Performance é aceitável (benchmark se crítico)

---

## Conclusão

Este guia deve ser usado como referência contínua durante o desenvolvimento. Lembre-se:

1. **Simplicidade é fundamental** - Go valoriza código simples e direto
2. **Erros são valores** - Trate-os explicitamente
3. **Interfaces são contratos** - Mantenha-as pequenas
4. **Composição sobre herança** - Use embedding
5. **Concorrência não é paralelismo** - Use goroutines e channels apropriadamente
6. **Documentação é código** - Mantenha atualizada
7. **Teste é essencial** - Escreva testes desde o início

A consistência é mais importante que a perfeição. Siga estas práticas e adapte conforme necessário para seu contexto específico.
