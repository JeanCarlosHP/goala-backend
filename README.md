# Calorie AI Backend

Backend API for AI-powered calorie tracking application.

## Stack

- **Framework**: Fiber (Go)
- **Database**: PostgreSQL with pgxpool
- **Authentication**: Firebase Auth
- **Storage**: AWS S3
- **AI Processing**: OpenAI GPT-4o / Google Gemini Flash
- **Architecture**: Clean Architecture + CQRS
- **Logging**: Slog

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Firebase project credentials
- AWS account (for S3 and Lambda)
- OpenAI API key or Google Gemini API key

## Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Add Firebase credentials:
   - Download your Firebase service account JSON from Firebase Console
   - Save it as `firebase-credentials.json` in the project root
   - Or set `FIREBASE_CREDENTIALS_JSON` in `.env`

3. Start services:
```bash
docker-compose up -d postgres rabbitmq
```

4. Run migrations:
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable" up
```

5. Install dependencies:
```bash
go mod download
```

6. Run the application:
```bash
go run cmd/api/main.go
```

## API Documentation

Complete API documentation is available in OpenAPI/Swagger format:
- **[Swagger YAML](./docs/swagger.yaml)** - OpenAPI 3.0 specification
- **[Environment Variables](./docs/ENVIRONMENT_VARIABLES.md)** - Complete environment variables documentation

### Quick API Reference

## API Endpoints

### Health & Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `GET /api/v1/auth/me` - Get current user (protected)

### Goals
- `PUT /api/v1/goals` - Update user goals (protected)

### Meals
- `GET /api/v1/meals?date=2025-11-24` - Get meals by date (protected)
- `POST /api/v1/meals` - Create new meal (protected)
- `GET /api/v1/summary/daily?date=2025-11-24` - Get daily summary (protected)

### Foods
- `GET /api/v1/foods/search?q=arroz` - Search food database (protected)
- `GET /api/v1/foods/recent` - Get recently used foods (protected)
- `POST /api/v1/ai/autocomplete` - AI-powered food macro autocomplete (protected)

### AI Photo Analysis
- `POST /api/v1/upload/sign` - Generate S3 signed URL for photo upload (protected)
- `POST /api/v1/ai/analyze` - Trigger AI analysis of uploaded photo (protected)
- `POST /api/v1/webhook/ai` - Webhook to receive AI analysis results (internal)

## Photo Analysis Flow

See [PHOTO_FLOW.md](./PHOTO_FLOW.md) for detailed documentation on the AI photo analysis flow.

## Testing & Monitoring

See [TESTING_MONITORING.md](./TESTING_MONITORING.md) for complete testing and monitoring guide.

### Quick Start

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Start monitoring stack (Prometheus + Grafana)
docker-compose up -d prometheus grafana

# Access Grafana at http://localhost:3000 (admin/admin)
# Access Prometheus at http://localhost:9090
```

## Docker

Build and run with Docker Compose:
```bash
docker-compose up --build
```

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain entities
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Data access
│   ├── middleware/      # Auth, logging
│   └── infrastructure/  # DB, Firebase, S3, RabbitMQ
├── lambda/              # AWS Lambda functions
│   ├── openai_analyzer.py   # OpenAI GPT-4o Vision
│   ├── gemini_analyzer.py   # Google Gemini Flash
│   ├── requirements.txt
│   ├── deploy.sh
│   └── README.md
├── migrations/          # Database migrations
└── docker-compose.yml
```

## Development

Install migrate tool:
```bash
brew install golang-migrate
```

Create new migration:
```bash
migrate create -ext sql -dir migrations -seq migration_name
```

## AWS Lambda Deployment

See [lambda/README.md](./lambda/README.md) for Lambda deployment instructions.

Quick deploy:
```bash
cd lambda
chmod +x deploy.sh
./deploy.sh calorieai-openai openai_analyzer.py
./deploy.sh calorieai-gemini gemini_analyzer.py
```

## Environment Variables

All environment variables are documented in [docs/ENVIRONMENT_VARIABLES.md](./docs/ENVIRONMENT_VARIABLES.md).

Required environment variables are documented in [.env.example](./.env.example).

Key variables:
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_S3_BUCKET` - AWS S3 configuration
- `RABBITMQ_URL` - RabbitMQ connection string
- `LAMBDA_FUNCTION_NAME` - AWS Lambda function name
- `LAMBDA_WEBHOOK_URL` - Backend webhook URL for AI results

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make test-coverage
```

See [TESTING_MONITORING.md](./TESTING_MONITORING.md) for detailed testing documentation.

## License

MIT
