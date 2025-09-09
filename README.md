# Magazin POS/ERP (Go)

[![Release](https://img.shields.io/github/v/release/aslon1213/g4h_pos_erp?style=flat-square)](https://github.com/aslon1213/g4h_pos_erp/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/aslon1213/g4h_pos_erp/release.yaml?style=flat-square)](https://github.com/aslon1213/g4h_pos_erp/actions/workflows/release.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aslon1213/g4h_pos_erp?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/github/license/aslon1213/g4h_pos_erp?style=flat-square)](LICENSE)

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
- **Authentication**: JWT/BasicAuth
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
   git clone https://github.com/aslon1213/g4h_pos_erp.git
   cd g4h_pos_erp
   ```

2. **Prepare configuration**

   ```bash
   cp example.config.yaml config.yaml
   # Edit config.yaml to match your MongoDB/Redis/S3 settings
   ```

3. **Create required docker networks (first run only)**

   ```bash
   docker network create mongoCluster || true
   docker network create caddy || true
   ```

4. **Start optional infrastructure (reverse proxy, redis UI)**

   ```bash
   # Caddy (for domain/HTTPS via labels). Optional for local use
   docker compose -f deploy/docker-compose-caddy.yaml up -d

   # Redis (redis-stack with UI). Optional
   docker compose -f deploy/docker-compose-db.yml up -d
   ```

5. **Start the backend**

   ```bash
   docker compose -f deploy/docker-compose.yml up -d --build
   ```

### Manual Setup

1. **Install dependencies**

   ```bash
   go mod download
   ```

2. **Setup MongoDB**

   - Use an existing MongoDB replica set, or update `config.yaml` to point to your MongoDB URL. Ensure `replica_set` aligns with your cluster if replication is enabled.

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
http://localhost:12000/docs/index.html
```

Note: Access to `/docs` is protected with BasicAuth. Default credentials are configured in `config.yaml` under `server.admin_docs_users`.

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
├── deploy/                         # Deployment configurations - docker-compose*.yml
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
# Start backend
docker compose -f deploy/docker-compose.yml up

# View logs
docker compose -f deploy/docker-compose.yml logs -f

# Stop services
docker compose -f deploy/docker-compose.yml down
```

## 📦 Releases & CI/CD

### Automated Releases

This project uses GitHub Actions with [GoReleaser](https://goreleaser.com/) for automated releases:

- **Trigger**: Releases are automatically triggered when you push a git tag starting with `v` (e.g., `v1.0.0`, `v2.1.3`)
- **Build Matrix**: Cross-platform builds for Linux, macOS, and Windows (amd64, arm64)
- **Artifacts**: Pre-built binaries, archives, and checksums
- **Distribution**: Automatically published to GitHub Releases

### Release Workflow

1. **Create and push a tag**:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Builds binaries for all supported platforms
   - Runs tests and quality checks
   - Creates GitHub release with artifacts
   - Generates changelog

### Build Platforms

| OS      | Architecture | Status       |
| ------- | ------------ | ------------ |
| Linux   | amd64, arm64 | ✅ Supported |
| macOS   | amd64, arm64 | ✅ Supported |
| Windows | amd64        | ✅ Supported |

### Local Development Build

```bash
# Build snapshot for testing (without releasing)
goreleaser build --snapshot --clean

# Build for current platform only
go build -o bin/pos-erp cmd/main.go
```

### Download Latest Release

Visit the [Releases page](https://github.com/aslon1213/g4h_pos_erp/releases) to download pre-built binaries for your platform.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/<some-feature>`)
3. Commit your changes (`git commit -m 'Add some amazing -- <some-feature>'`)
4. Push to the branch (`git push origin feature/<some-feature>`)
5. Open a Pull Request

## 🆘 Support

For support and questions:

- Create an issue in the GitHub repository
- Contact: hamidovaslon13@gmail.com

## 🚧 Development Status

This project is actively maintained and under development. Please check the [issues](https://github.com/aslon1213/g4h_pos_erp/issues) for current development status and upcoming features.
