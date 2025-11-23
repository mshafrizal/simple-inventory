# Architecture Documentation

## Clean Architecture Overview

This project follows Clean Architecture principles to ensure maintainability, testability, and scalability.

```
┌─────────────────────────────────────────────────────────┐
│                   HTTP Requests                          │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│              INTERFACE LAYER                             │
│  ┌────────────┐  ┌────────────┐  ┌─────────────┐       │
│  │ Middleware │  │  Handlers  │  │    DTOs     │       │
│  └────────────┘  └────────────┘  └─────────────┘       │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│               USE CASE LAYER                             │
│  ┌────────────────────────────────────────────────┐     │
│  │  Business Logic & Application Rules             │     │
│  │  - AuthUseCase                                  │     │
│  │  - ProductUseCase                               │     │
│  │  - InventoryUseCase                             │     │
│  │  - LocationUseCase                              │     │
│  └────────────────────────────────────────────────┘     │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│               DOMAIN LAYER                               │
│  ┌──────────────────┐  ┌──────────────────────┐        │
│  │    Entities      │  │  Repository Interfaces│        │
│  │  - User          │  │  - UserRepository     │        │
│  │  - Product       │  │  - ProductRepository  │        │
│  │  - Location      │  │  - LocationRepository │        │
│  │  - Transaction   │  │  - TransactionRepo    │        │
│  └──────────────────┘  └──────────────────────┘        │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│           INFRASTRUCTURE LAYER                           │
│  ┌────────────┐  ┌────────────┐  ┌─────────────┐       │
│  │   GORM     │  │   Viper    │  │  PostgreSQL │       │
│  │ Repository │  │   Config   │  │  Database   │       │
│  └────────────┘  └────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Core business logic and rules, independent of frameworks and external concerns.

**Components**:
- **Entities**: Core business objects (User, Product, Location, etc.)
  - Contains business methods (e.g., `IsLowStock()`, `UpdateQuantity()`)
  - No external dependencies

- **Repository Interfaces**: Define contracts for data access
  - Abstractions that infrastructure implements
  - Enables dependency inversion

**Key Files**:
- `entity/user.go` - User domain model with password hashing
- `entity/product.go` - Product model with inventory logic
- `entity/location.go` - Location hierarchy management
- `entity/inventory_transaction.go` - Audit trail for inventory
- `repository/*.go` - Repository interface definitions

**Rules**:
- No dependencies on outer layers
- Pure business logic only
- Framework-agnostic

### 2. Use Case Layer (`internal/usecase/`)

**Purpose**: Application-specific business rules and orchestration.

**Components**:
- **AuthUseCase**: User registration, login, session management
- **ProductUseCase**: Product CRUD, search, validation
- **InventoryUseCase**: Inventory transactions (receive, issue, adjust, transfer)
- **LocationUseCase**: Location management

**Responsibilities**:
- Orchestrate domain entities
- Implement business workflows
- Coordinate repository operations
- Handle business validation

**Example Flow**:
```go
func (uc *InventoryUseCase) ReceiveInventory(...) error {
    // 1. Validate input
    // 2. Fetch product from repository
    // 3. Update product quantity (domain logic)
    // 4. Save to repository
    // 5. Create transaction record
}
```

**Rules**:
- Depends on domain layer only
- No HTTP or database concerns
- Testable in isolation

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: External concerns implementation (database, config, external services).

**Components**:

**Database** (`database/`):
- Connection management
- Automatic migrations
- Connection pooling

**Config** (`config/`):
- Viper-based configuration
- Environment variable loading
- Configuration validation

**Persistence** (`persistence/`):
- GORM repository implementations
- Implements domain repository interfaces
- Database queries and operations

**Key Features**:
```go
// Example: UserRepository implementation
type userRepositoryImpl struct {
    db *gorm.DB
}

func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    return &user, err
}
```

**Rules**:
- Implements domain interfaces
- Contains framework-specific code
- Replaceable without affecting business logic

### 4. Interface Layer (`internal/interfaces/http/`)

**Purpose**: External communication handling (HTTP in this case).

**Components**:

**Handlers** (`handler/`):
- HTTP request/response handling
- Input validation
- Call use cases
- Map entities to DTOs

**DTOs** (`dto/`):
- Request/response models
- Validation tags
- API contracts

**Middleware** (`middleware/`):
- Authentication (JWT validation)
- Logging
- CORS
- Error handling
- Request recovery

**Router** (`router/`):
- Route definitions
- Middleware application
- Handler registration

**Example**:
```go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // 1. Parse request (DTO)
    var req dto.CreateProductRequest
    c.ShouldBindJSON(&req)

    // 2. Convert to entity
    product := &entity.Product{...}

    // 3. Call use case
    err := h.productUseCase.CreateProduct(ctx, product)

    // 4. Return response (DTO)
    c.JSON(http.StatusCreated, response)
}
```

**Rules**:
- Only layer that knows about HTTP
- Converts between DTOs and entities
- No business logic

## Dependency Flow

```
Interfaces → Use Cases → Domain ← Infrastructure
     ↓           ↓          ↑          ↑
    HTTP    Business    Entities   Database
           Rules                   Config
```

**Key Principle**: Dependencies point inward
- Outer layers depend on inner layers
- Inner layers know nothing about outer layers
- Domain is the most stable, interfaces are the most volatile

## Data Flow Example: Create Product

```
1. HTTP POST /api/v1/products
   ↓
2. Router → ProductHandler.CreateProduct()
   ↓
3. Parse CreateProductRequest DTO
   ↓
4. Convert to Product entity
   ↓
5. ProductUseCase.CreateProduct()
   ├─ Validate SKU uniqueness (via repository)
   ├─ Validate barcode (via repository)
   └─ Save product (via repository)
   ↓
6. ProductRepository.Create()
   ↓
7. GORM → PostgreSQL
   ↓
8. Return ProductResponse DTO
   ↓
9. HTTP 201 Created
```

## Security Architecture

### Authentication Flow

```
1. User Login → AuthHandler
   ↓
2. AuthUseCase.Login()
   ├─ Validate credentials
   ├─ Generate token
   └─ Create session
   ↓
3. Return token to client
   ↓
4. Client includes token in Authorization header
   ↓
5. AuthMiddleware validates token
   ├─ Extract token
   ├─ Validate session
   ├─ Check expiration
   └─ Inject user into context
   ↓
6. Protected handler accesses user from context
```

### Security Features

- **Password Hashing**: bcrypt with salt
- **Token-Based Auth**: Secure token generation
- **Session Management**: Expiration and cleanup
- **CORS**: Configurable cross-origin access
- **Input Validation**: Request validation with tags

## Database Schema

```sql
users
├── id (PK)
├── username (unique)
├── email (unique)
├── password (hashed)
├── role
├── is_active
└── timestamps

sessions
├── id (PK)
├── user_id (FK → users)
├── token (unique, indexed)
├── expires_at
└── timestamps

locations
├── id (PK)
├── name
├── code (unique, indexed)
├── building, floor, aisle, shelf
├── is_active
└── timestamps

products
├── id (PK)
├── name
├── sku (unique, indexed)
├── barcode (unique, indexed)
├── description
├── quantity
├── min_quantity
├── price
├── location_id (FK → locations)
└── timestamps

inventory_transactions
├── id (PK)
├── product_id (FK → products)
├── type (IN/OUT/ADJUST/TRANSFER)
├── quantity
├── from_location_id (FK → locations)
├── to_location_id (FK → locations)
├── user_id (FK → users)
├── notes
└── created_at
```

## Testing Strategy

### Unit Tests
- **Domain entities**: Test business logic methods
- **Use cases**: Mock repositories, test workflows
- **Handlers**: Mock use cases, test HTTP logic

### Integration Tests
- **Repository**: Test with real database
- **API**: End-to-end endpoint testing

### Test Structure
```go
// Example: Use case test
func TestAuthUseCase_Login(t *testing.T) {
    mockRepo := new(MockUserRepository)
    useCase := NewAuthUseCase(mockRepo, ...)

    mockRepo.On("GetByUsername", "admin").Return(user, nil)

    session, err := useCase.Login(ctx, "admin", "password")

    assert.NoError(t, err)
    assert.NotEmpty(t, session.Token)
}
```

## Benefits of This Architecture

1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Easy to swap implementations (e.g., change database)
4. **Scalability**: Can add features without affecting existing code
5. **Team Collaboration**: Different teams can work on different layers
6. **Framework Independence**: Not locked into Gin or GORM

## Design Patterns Used

- **Repository Pattern**: Data access abstraction
- **Dependency Injection**: Loose coupling between layers
- **DTO Pattern**: Separate external and internal models
- **Middleware Pattern**: Cross-cutting concerns
- **Factory Pattern**: Object creation (e.g., NewAuthUseCase)

## Best Practices Implemented

- **SOLID Principles**: Single responsibility, open/closed, dependency inversion
- **Clean Code**: Readable, self-documenting code
- **Error Handling**: Consistent error responses
- **Logging**: Structured logging for debugging
- **Configuration**: Environment-based configuration
- **Security**: Authentication, authorization, input validation
