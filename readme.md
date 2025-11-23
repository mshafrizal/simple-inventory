# Simple Inventory API

A clean architecture REST API for inventory management built with Go, following best practices and SOLID principles.

## Features

- **Inventory Management**: Track products, quantities, and stock levels
- **Barcode Support**: Scan and manage product barcodes
- **Location Tracking**: Organize inventory across multiple locations
- **Transaction History**: Complete audit trail of all inventory movements
- **User Authentication**: Secure JWT-based authentication system
- **Low Stock Alerts**: Track products below minimum quantity
- **Search & Filter**: Powerful search across products and locations

## Tech Stack

- **Go 1.21**: Backend language
- **Gin**: HTTP web framework
- **GORM**: ORM for database operations
- **PostgreSQL**: Relational database
- **Viper**: Configuration management
- **JWT**: Authentication tokens

## Project Structure (Clean Architecture)

```
simple-inventory/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   ├── entity/              # Domain entities
│   │   └── repository/          # Repository interfaces
│   ├── usecase/                 # Business logic layer
│   ├── infrastructure/
│   │   ├── config/              # Configuration
│   │   ├── database/            # Database setup
│   │   └── persistence/         # Repository implementations
│   └── interfaces/
│       └── http/
│           ├── dto/             # Data Transfer Objects
│           ├── handler/         # HTTP handlers
│           ├── middleware/      # HTTP middleware
│           └── router/          # Route definitions
├── .env.example                 # Environment variables template
├── docker-compose.yml           # Docker setup
├── Dockerfile                   # Container image
├── Makefile                     # Build commands
└── API_DOCUMENTATION.md         # API documentation

```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd simple-inventory
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your configuration:
```env
APP_NAME=simple-inventory
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inventory_db
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
```

4. Install dependencies:
```bash
make install-deps
```

### Running with Docker

```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Running Locally

1. Start PostgreSQL database

2. Install dependencies:
```bash
make install-deps
```

3. Run the application:
```bash
make run
```

Or build and run:
```bash
make build
./bin/api
```

## API Documentation

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for detailed API endpoints and usage examples.

### Quick Start Example

1. Register a user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "secure123"
  }'
```

2. Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure123"
  }'
```

3. Create a product (use token from login):
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "name": "Laptop",
    "sku": "LAP-001",
    "barcode": "1234567890",
    "quantity": 10,
    "min_quantity": 5,
    "price": 999.99
  }'
```

## Development

### Available Make Commands

- `make help` - Show available commands
- `make install-deps` - Install Go dependencies
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business entities and repository interfaces
2. **Use Case Layer** (`internal/usecase/`): Business logic and application rules
3. **Infrastructure Layer** (`internal/infrastructure/`): Database, config, external services
4. **Interface Layer** (`internal/interfaces/`): HTTP handlers, DTOs, middleware

### Key Design Patterns

- **Repository Pattern**: Abstracts data access
- **Dependency Injection**: Decouples components
- **Clean Architecture**: Separates concerns by layer
- **DTO Pattern**: Separates external and internal models

## Testing

Run tests with:
```bash
make test
```

## Database Migrations

The application automatically runs migrations on startup. All tables are created from GORM entity definitions.

## Security Features

- Password hashing with bcrypt
- JWT-based authentication
- Token expiration handling
- Session management
- CORS middleware
- Input validation

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.