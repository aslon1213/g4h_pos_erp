# Go POS ERP System

A modern Point of Sale (POS) and Enterprise Resource Planning (ERP) system built with Go, designed for retail businesses and stores.

## 🚀 Features

- **Sales Management** - Complete sales transaction processing
- **Product Management** - Inventory tracking and product catalog
- **Customer Management** - Customer data and relationship management
- **Supplier Management** - Vendor and supplier relationship tracking
- **Financial Management** - Transaction processing and financial reporting
- **Journal Entries** - Accounting and bookkeeping functionality
- **Internal Expenses** - Business expense tracking and management
- **REST API** - Full RESTful API with Swagger documentation
- **Real-time Data** - Redis caching for optimal performance
- **Observability** - OpenTelemetry integration for monitoring

## 🛠️ Tech Stack

- **Backend**: Go 1.24.3
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Database**: MongoDB
- **Cache**: Redis
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT/BasicAuthc
- **Observability**: OpenTelemetry
- **Logging**: Zerolog
- **Configuration**: Viper (YAML-based)

## 📋 Prerequisites

- Go 1.24.3 or higher
- MongoDB 4.4+ cluster with replica set
- Redis 6.0+
- Docker & Docker Compose (for containerized setup)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/aslon1213/go-pos-erp.git
   cd go-pos-erp
   ```

2. **Start the services**
   ```bash
   docker-compose -f deployments/docker-compose.yml up -d
   ```

3. **Run the application**
   ```bash
   go run cmd/main.go
   ```

### Manual Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup MongoDB**
   ```bash
   # Use the provided MongoDB setup script
   chmod +x deploy/mongo.sh
   ./deploy/mongo.sh
   ```

3. **Configure the application**
   ```bash
   cp example.config.yaml config.yaml
   # Edit config.yaml with your database and Redis settings
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

## ⚙️ Configuration

The application uses YAML configuration. Copy `example.config.yaml` to `config.yaml` and modify as needed:

```yaml
database:
  host: "localhost"
  port: "27017"
  username: "admin"
  password: "admin"
  database: "store"
  max_connections: 20
  min_pool_size: 10
  auth: false
  replica_set: "rs0"
  url: "mongodb://localhost:27017/?replicaSet=rs0"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  database: 0

server:
  port: ":12000"
```

## 📖 API Documentation

Once the application is running, access the Swagger documentation at:
```
http://localhost:12000/swagger/
```

The API provides endpoints for:
- `/sales` - Sales transactions
- `/products` - Product management
- `/customers` - Customer management
- `/suppliers` - Supplier management
- `/transactions` - Financial transactions
- `/journals` - Journal entries
- `/expenses` - Internal expenses
- `/finance` - Financial operations

## 🔧 Development

### Project Structure
```
.
├── cmd/                            # Application entry points - main.go    
├── pkg/                            # Application packages
│   ├── app/                        # Main application logic 
│   ├── controllers/                # Business logic controllers
│   │   ├── customers/              # Customer management
│   │   ├── finance/                # Financial operations
│   │   ├── internalExpenses/       # Internal expenses
│   │   ├── journals/               # Journal entries for daily financial operations
│   │   ├── products/               # Product management
│   │   ├── sales/                  # Sales transactions
│   │   ├── suppliers/              # Supplier management
│   │   └── transactions/           # Financial transactions
│   ├── routes/                     # API route definitions
│   ├── middleware/                 # HTTP middleware
│   ├── repository/                 # Data access layer, models
│   ├── utils/                      # Utility functions
│   └── configs/                    # Configuration management
├── docs/                           # API documentation
├── deploy/                         # Deployment configurations - docker-compose.yml
├── test/                           # Test files
├── web/                            # Frontend assets (if any)
└── scripts/                        # Build and deployment scripts
```
### Building the Application

```bash
# Build for current platform
go build -o bin/pos-erp cmd/main.go

# Build for multiple platforms using GoReleaser
goreleaser build --snapshot --clean
```

### Running Tests

```bash
go test ./test
```

### Generate Swagger Documentation

```bash
swag init -g app/app.go --dir pkg
```

## 🐳 Docker Deployment

The project includes Docker Compose configuration for easy deployment:

```bash
# Start all services
docker-compose -f deployments/docker-compose.yml up -d

# View logs
docker-compose -f deployments/docker-compose.yml logs -f

# Stop services
docker-compose -f deployments/docker-compose.yml down
```

## 📦 Release

This project uses [GoReleaser](https://goreleaser.com/) for automated releases. Releases are built for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

```bash 
goreleaser build --snapshot --clean
```


## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/<some-feature>`)
3. Commit your changes (`git commit -m 'Add some amazing -- <some-feature>'`)
4. Push to the branch (`git push origin feature/<some-feature>`)
5. Open a Pull Request

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Contact: support@swagger.io

## 🚧 Development Status

This project is actively maintained and under development. Please check the [issues](https://github.com/aslon1213/go-pos-erp/issues) for current development status and upcoming features.
