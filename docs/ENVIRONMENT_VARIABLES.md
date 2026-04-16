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
| `TRACING_ENABLED` | Habilitar tracing global da aplicação | `true`, `false` | `true` |

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

### Redis

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `REDIS_URL` | URL de conexão do Redis para cache de busca | `redis://localhost:6379` | ❌ |

### Meilisearch

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `MEILISEARCH_URL` | URL base do Meilisearch | `http://localhost:7700` | ❌ |
| `MEILISEARCH_API_KEY` | Chave de API do Meilisearch | `masterKey` | ❌ |
| `MEILISEARCH_FOODS_INDEX` | Nome do índice de alimentos | `foods` | ❌ |

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

### Railway Storage Bucket

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `BUCKET` | Nome real do bucket para a API S3 do Railway | `calorieai-meals-ab12cd34` | ✅ |
| `ENDPOINT` | Endpoint S3-compatible do Railway | `https://storage.railway.app` | ✅ |
| `REGION` | Região exposta pelo bucket Railway | `auto` | ✅ |
| `ACCESS_KEY_ID` | Access Key ID do bucket Railway | `rwk_xxxxxxxxxxxxx` | ✅ |
| `SECRET_ACCESS_KEY` | Secret Access Key do bucket Railway | `rwks_xxxxxxxxxxxxx` | ✅ |

### CDN

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `CDN_DOMAIN` | Base pública usada para servir imagens via proxy do backend | `https://api.example.com` | ✅ |

**Configuração de entrega de imagens:**
- Usado para servir avatares e imagens do usuário
- O backend salva apenas paths relativos no banco de dados (ex: `/avatars/{userId}/avatar.jpg`)
- Ao retornar dados para o frontend, concatena `CDN_DOMAIN + path`
- As imagens são servidas por rotas públicas do backend, porque o bucket do Railway é privado
- Rotas públicas atuais:
  - `GET /avatars/:firebaseUID/:filename`
  - `GET /users/:userID/food_images/:filename`

**Como configurar o bucket Railway:**
1. Crie um Storage Bucket no projeto Railway
2. Abra a aba de credenciais do bucket
3. Copie `BUCKET`, `ENDPOINT`, `REGION`, `ACCESS_KEY_ID` e `SECRET_ACCESS_KEY` para o serviço do backend
4. Configure CORS no bucket para permitir uploads diretos do frontend por `presigned PUT URL`
5. Configure `CDN_DOMAIN` com o domínio público do backend

### RevenueCat

| Variável | Descrição | Exemplo | Obrigatório |
|----------|-----------|---------|-------------|
| `REVENUECAT_WEBHOOK_SECRET` | Secret para validar webhooks do RevenueCat | `rc_whsec_xxxxxxxxxxxxx` | ✅ |

**Configuração do RevenueCat:**
- Usado para gerenciar assinaturas e acesso a recursos premium
- O backend valida webhooks usando este secret
- Configure o webhook URL no dashboard do RevenueCat: `https://your-api.com/api/v1/webhooks/revenuecat`
- Copie o webhook secret do dashboard do RevenueCat

### AI Providers

| Variável | Descrição | Valores Permitidos | Padrão | Obrigatório |
|----------|-----------|-------------------|--------|-------------|
| `AI_PROVIDER` | Provedor de IA para reconhecimento de alimentos | `gemini`, `openai` | `gemini` | ❌ |
| `GEMINI_API_KEY` | Chave da API do Google Gemini | - | - | ✅ (se AI_PROVIDER=gemini) |
| `GEMINI_MODEL` | Modelo do Gemini a usar | `gemini-3-flash-preview`, `gemini-2.5-flash` | `gemini-3-flash-preview` | ❌ |
| `OPENAI_API_KEY` | Chave da API do OpenAI | - | - | ✅ (se AI_PROVIDER=openai) |

**Configuração do Gemini:**
- Obtenha a API key no [Google AI Studio](https://aistudio.google.com/app/apikey)
- O modelo padrão é otimizado para velocidade e custo
- Suporta análise multimodal (texto + imagem)

**Configuração do OpenAI:**
- Obtenha a API key no [OpenAI Platform](https://platform.openai.com/api-keys)
- Atualmente não implementado, usa Gemini como fallback

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
TRACING_ENABLED=false
DATABASE_SSL_MODE=disable
DATABASE_TRACING=true
HTTP_CORS_ALLOWED_ORIGINS=*
```

### Staging
```bash
ENVIRONMENT=staging
LOGGING_LEVEL=info
LOGGING_JSON_FORMAT=true
TRACING_ENABLED=false
DATABASE_SSL_MODE=require
DATABASE_TRACING=false
HTTP_CORS_ALLOWED_ORIGINS=https://staging.calorieai.com
```

### Production
```bash
ENVIRONMENT=production
LOGGING_LEVEL=warn
LOGGING_JSON_FORMAT=true
TRACING_ENABLED=false
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
TRACING_ENABLED=false

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

# Railway Storage Bucket
BUCKET=calorieai-meals-ab12cd34
ENDPOINT=https://storage.railway.app
REGION=auto
ACCESS_KEY_ID=your-railway-access-key-id
SECRET_ACCESS_KEY=your-railway-secret-access-key

# Public image base URL
CDN_DOMAIN=https://api.example.com

# RevenueCat
REVENUECAT_WEBHOOK_SECRET=rc_whsec_xxxxxxxxxxxxx

# AI Providers
AI_PROVIDER=gemini
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-3-flash-preview

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
- Confirme `BUCKET`, `ENDPOINT`, `REGION`, `ACCESS_KEY_ID`, `SECRET_ACCESS_KEY`
- Verifique se o bucket Railway existe e está ativo
- Confirme se o frontend está usando a `presigned PUT URL` antes da expiração
- Verifique se o bucket tem CORS configurado para o domínio do frontend

### Erro: "failed to connect to RabbitMQ"
- Verifique se RabbitMQ está rodando: `docker ps | grep rabbitmq`
- Confirme `RABBITMQ_URL`
- Teste a conexão: `curl http://localhost:15672`
