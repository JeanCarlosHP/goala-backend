# Environment Variables Documentation

Documentação completa de todas as variáveis de ambiente utilizadas pelo Calorie AI Backend.

## Configuração Inicial

1. Copie o arquivo de exemplo:
```bash
cp .env.example .env
```

2. Configure as variáveis obrigatórias antes de executar a aplicação.

## Variáveis Obrigatórias

### Server Configuration

| Variável | Descrição | Exemplo | Padrão |
|----------|-----------|---------|--------|
| `HTTP_PORT` | Porta HTTP do servidor | `8080` | `8080` |
| `HTTP_REQUEST_SIZE_LIMIT` | Limite de tamanho da requisição (bytes) | `10485760` | `10485760` (10MB) |
| `HTTP_CORS_ALLOWED_ORIGINS` | Origens permitidas para CORS | `*` ou `http://localhost:3000` | `*` |
| `HTTP_CORS_ALLOWED_HEADERS` | Headers permitidos para CORS | `*` ou `Content-Type,Authorization` | `*` |
| `HTTP_CORS_ALLOWED_METHODS` | Métodos HTTP permitidos | `GET,POST,PUT,DELETE,OPTIONS` | - |

### Environment

| Variável | Descrição | Valores Permitidos | Padrão |
|----------|-----------|-------------------|--------|
| `ENVIRONMENT` | Ambiente de execução | `development`, `staging`, `production` | `development` |
| `LOGGING_LEVEL` | Nível de log | `debug`, `info`, `warn`, `error`, `fatal` | `debug` |
| `LOGGING_JSON_FORMAT` | Formato JSON para logs | `true`, `false` | `true` |

### Database (PostgreSQL)

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `DATABASE_HOST` | Host do PostgreSQL | `localhost` | ✅ |
| `DATABASE_PORT` | Porta do PostgreSQL | `5432` | ✅ |
| `DATABASE_USER` | Usuário do banco | `postgres` | ✅ |
| `DATABASE_PASSWORD` | Senha do banco | `postgres` | ✅ |
| `DATABASE_NAME` | Nome do banco de dados | `calorie_ai` | ✅ |
| `DATABASE_SSL_MODE` | Modo SSL | `disable`, `require`, `verify-full` | ✅ |
| `DATABASE_TRACING` | Habilitar trace de queries | `true`, `false` | ❌ |
| `DATABASE_MIGRATION_PATH` | Caminho das migrations | `file://pkg/database/migrations` | ✅ |
| `DATABASE_MIGRATION_FORCE_SET_LOWER_VERSION` | Forçar versão inferior | `true`, `false` | ❌ |

### Pagination

| Variável | Descrição | Exemplo | Padrão |
|----------|-----------|---------|--------|
| `PAGINATION_MAX_PER_PAGE` | Máximo de itens por página | `100` | `100` |
| `PAGINATION_DEFAULT_PAGE_SIZE` | Tamanho padrão da página | `20` | `20` |

### Firebase Authentication

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `FIREBASE_CREDENTIALS_FILE` | Caminho do arquivo JSON | `./secrets/firebase-credentials.json` | ⚠️ * |
| `FIREBASE_CREDENTIALS_JSON` | JSON inline das credenciais | `{"type":"service_account"...}` | ⚠️ * |

> ⚠️ * Uma das duas variáveis deve ser configurada. Prioridade: `FIREBASE_CREDENTIALS_JSON` > `FIREBASE_CREDENTIALS_FILE`

**Como obter as credenciais do Firebase:**
1. Acesse o [Firebase Console](https://console.firebase.google.com)
2. Selecione seu projeto
3. Vá em "Project Settings" > "Service Accounts"
4. Clique em "Generate New Private Key"
5. Salve o arquivo JSON ou copie o conteúdo

### AWS (S3 e Lambda)

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `AWS_REGION` | Região AWS | `us-east-1` | ✅ |
| `AWS_ACCESS_KEY_ID` | Access Key ID | `AKIAIOSFODNN7EXAMPLE` | ✅ |
| `AWS_SECRET_ACCESS_KEY` | Secret Access Key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` | ✅ |
| `AWS_S3_BUCKET_NAME` | Nome do bucket S3 | `calorieai-meals` | ✅ |
| `AWS_S3_REGION` | Região do bucket S3 | `us-east-1` | ✅ |

### CDN

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `CDN_DOMAIN` | Domínio do CDN para servir assets | `https://cdn.example.com` ou `https://d1234567890abc.cloudfront.net` | ✅ |

**Configuração do CDN:**
- Usado para servir avatares e imagens do usuário
- O backend salva apenas paths relativos no banco de dados (ex: `/avatars/{userId}/avatar.jpg`)
- Ao retornar dados para o frontend, concatena `CDN_DOMAIN + path`
- Suporta CloudFront, Cloudflare, Fastly ou qualquer CDN compatível com S3

**Como configurar AWS:**
1. Crie uma conta AWS ou faça login
2. Crie um bucket S3 para armazenar fotos
3. Crie um usuário IAM com permissões: `s3:PutObject`, `s3:GetObject`, `lambda:InvokeFunction`
4. Gere e copie as credenciais

### RevenueCat

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `REVENUECAT_WEBHOOK_SECRET` | Secret para validar webhooks do RevenueCat | `rc_whsec_xxxxxxxxxxxxx` | ✅ |

**Configuração do RevenueCat:**
- Usado para gerenciar assinaturas e acesso a recursos premium
- O backend valida webhooks usando este secret
- Configure o webhook URL no dashboard do RevenueCat: `https://your-api.com/api/v1/webhooks/revenuecat`
- Copie o webhook secret do dashboard do RevenueCat

### RabbitMQ

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `RABBITMQ_URL` | URL de conexão | `amqp://guest:guest@localhost:5672/` | ✅ |
| `RABBITMQ_QUEUE` | Nome da fila | `ai-photo-analysis` | ✅ |
| `RABBITMQ_EXCHANGE` | Nome do exchange | `calorieai` | ✅ |

### AWS Lambda

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `LAMBDA_FUNCTION_NAME` | Nome da função Lambda | `calorieai-openai` ou `calorieai-gemini` | ✅ |
| `LAMBDA_WEBHOOK_URL` | URL do webhook para resultados | `http://localhost:8080/api/v1/webhook/ai` | ✅ |

## Variáveis por Ambiente

### Development
```bash
ENVIRONMENT=development
LOGGING_LEVEL=debug
LOGGING_JSON_FORMAT=true
DATABASE_SSL_MODE=disable
DATABASE_TRACING=true
HTTP_CORS_ALLOWED_ORIGINS=*
```

### Staging
```bash
ENVIRONMENT=staging
LOGGING_LEVEL=info
LOGGING_JSON_FORMAT=true
DATABASE_SSL_MODE=require
DATABASE_TRACING=false
HTTP_CORS_ALLOWED_ORIGINS=https://staging.calorieai.com
```

### Production
```bash
ENVIRONMENT=production
LOGGING_LEVEL=warn
LOGGING_JSON_FORMAT=true
DATABASE_SSL_MODE=verify-full
DATABASE_TRACING=false
HTTP_CORS_ALLOWED_ORIGINS=https://calorieai.com
```

## Exemplo Completo (.env)

```bash
# Server
HTTP_PORT=8080
HTTP_REQUEST_SIZE_LIMIT=10485760
HTTP_CORS_ALLOWED_ORIGINS=*
HTTP_CORS_ALLOWED_HEADERS=*
HTTP_CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS

# Environment
ENVIRONMENT=development
LOGGING_LEVEL=debug
LOGGING_JSON_FORMAT=true

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=calorie_ai
DATABASE_SSL_MODE=disable
DATABASE_TRACING=false
DATABASE_MIGRATION_PATH=file://pkg/database/migrations
DATABASE_MIGRATION_FORCE_SET_LOWER_VERSION=false

# Pagination
PAGINATION_MAX_PER_PAGE=100
PAGINATION_DEFAULT_PAGE_SIZE=20

# Firebase
FIREBASE_CREDENTIALS_FILE=./secrets/firebase-credentials.json

# AWS
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_S3_BUCKET_NAME=calorieai-meals
AWS_S3_REGION=us-east-1

# CDN
CDN_DOMAIN=https://d1234567890abc.cloudfront.net

# RevenueCat
REVENUECAT_WEBHOOK_SECRET=rc_whsec_xxxxxxxxxxxxx

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=ai-photo-analysis
RABBITMQ_EXCHANGE=calorieai

# Lambda
LAMBDA_FUNCTION_NAME=calorieai-openai
LAMBDA_WEBHOOK_URL=http://localhost:8080/api/v1/webhook/ai
```

## Validação

Para validar se todas as variáveis obrigatórias estão configuradas, execute:

```bash
go run cmd/api/main.go
```

A aplicação falhará na inicialização se alguma variável obrigatória estiver faltando.

## Segurança

⚠️ **NUNCA commite o arquivo `.env` com credenciais reais.**

- Use `.env.example` para documentar as variáveis
- Mantenha `.env` no `.gitignore`
- Use serviços de gestão de secrets em produção (AWS Secrets Manager, HashiCorp Vault, etc.)
- Rotacione credenciais regularmente
- Use diferentes credenciais para cada ambiente

## Troubleshooting

### Erro: "failed to connect to database"
- Verifique se PostgreSQL está rodando
- Confirme `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD`
- Teste a conexão: `psql -h localhost -U postgres -d calorie_ai`

### Erro: "authentication service unavailable"
- Verifique se `FIREBASE_CREDENTIALS_FILE` ou `FIREBASE_CREDENTIALS_JSON` está configurado
- Confirme se o arquivo JSON é válido
- Verifique permissões do arquivo

### Erro: "failed to upload to S3"
- Confirme `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- Verifique se o bucket `AWS_S3_BUCKET` existe
- Confirme permissões IAM do usuário

### Erro: "failed to connect to RabbitMQ"
- Verifique se RabbitMQ está rodando: `docker ps | grep rabbitmq`
- Confirme `RABBITMQ_URL`
- Teste a conexão: `curl http://localhost:15672`
