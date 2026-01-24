# Fase 5: Food Recognition & Barcode - Implementação Concluída

## Resumo
Implementação completa dos endpoints de reconhecimento de alimentos por imagem, busca por barcode e estimativa de quantidade.

## Alterações Realizadas

### 1. Domain Layer (`internal/domain/`)
- ✅ Criado `food_recognition.go` com structs:
  - `FoodRecognitionRequest` - request para reconhecimento de imagem
  - `RecognizedFoodItem` - item de alimento reconhecido pela IA
  - `FoodRecognitionResponse` - response com lista de alimentos reconhecidos
  - `FoodBarcodeResponse` - response com informações de barcode
  - `EstimateQuantityRequest` - request para estimativa de quantidade
  - `EstimateQuantityResponse` - response com quantidade estimada

- ✅ Atualizado `config.go` com novos campos:
  - `AWSS3BucketName` - nome do bucket S3
  - `AWSS3Region` - região do S3
  - `AWSAccessKeyID` - credenciais AWS
  - `AWSSecretAccessKey` - credenciais AWS
  - `LambdaFoodRecognitionURL` - URL da Lambda de reconhecimento
  - `LambdaQuantityEstimationURL` - URL da Lambda de estimativa
  - `OpenFoodFactsAPIURL` - URL da API OpenFoodFacts

### 2. Service Layer (`internal/services/`)

#### S3Service (`s3_service.go`)
- ✅ Upload de imagens para S3
- ✅ Geração automática de nomes de arquivo únicos
- ✅ Suporte a múltiplos formatos de imagem (JPEG, PNG, WEBP, HEIC)
- ✅ Retorno da URL pública da imagem

#### FoodRecognitionService (`food_recognition_service.go`)
- ✅ `RecognizeFood` - reconhece alimentos em uma imagem
  - Upload da imagem para S3
  - Chamada à Lambda de reconhecimento (Gemini/OpenAI)
  - Retorno de lista de alimentos com macros e confiança
  - Medição de tempo de processamento

- ✅ `EstimateQuantity` - estima quantidade de um alimento
  - Upload da imagem para S3
  - Chamada à Lambda de estimativa
  - Considera contexto (tipo de refeição, localização)
  - Suporte a porções de referência

#### BarcodeService (`barcode_service.go`)
- ✅ `GetFoodByBarcode` - busca alimento por código de barras
  - Cache em banco de dados (food_database)
  - Integração com OpenFoodFacts API
  - Fallback automático para API externa
  - Armazenamento de resultados para cache

- ✅ Funções auxiliares:
  - `fetchFromOpenFoodFacts` - busca na API externa
  - `cacheFoodInDB` - armazena resultado no cache
  - `mapDBToResponse` - converte modelo DB para response
  - Conversões de tipos (int32, pgtype.Numeric)

### 3. Handler Layer (`internal/handlers/food_recognition.go`)
- ✅ Criado `FoodRecognitionHandler` com 3 endpoints:

#### `RecognizeFood` - POST /food/recognize
- Recebe imagem via multipart/form-data
- Valida presença do arquivo
- Chama service de reconhecimento
- Retorna lista de alimentos reconhecidos

#### `GetFoodByBarcode` - GET /food/barcode/{barcode}
- Recebe barcode via path parameter
- Busca em cache ou OpenFoodFacts
- Retorna informações nutricionais

#### `EstimateQuantity` - POST /food/estimate-quantity
- Recebe imagem e metadados
- Estima quantidade em gramas/ml/porções
- Retorna confiança da estimativa

### 4. Rotas (`cmd/api/main.go`)
- ✅ Instanciado `S3Service`
- ✅ Instanciado `FoodRecognitionService`
- ✅ Instanciado `BarcodeService`
- ✅ Instanciado `FoodRecognitionHandler`
- ✅ Adicionado rotas protegidas:
  - `POST /api/v1/food/recognize`
  - `GET /api/v1/food/barcode/:barcode`
  - `POST /api/v1/food/estimate-quantity`

### 5. Dependências
- ✅ Adicionado AWS SDK v2:
  - `github.com/aws/aws-sdk-go-v2/config`
  - `github.com/aws/aws-sdk-go-v2/service/s3`
  - `github.com/aws/aws-sdk-go-v2/credentials`
  - `github.com/aws/aws-sdk-go-v2/aws`

## Validações Implementadas

### FoodRecognitionRequest
- `uri`: required
- `name`: required
- `type`: required
- `mealLocation`: required

### RecognizedFoodItem
- `name`: required
- `calories`: gte=0, lte=5000
- `protein`: gte=0, lte=500
- `carbs`: gte=0, lte=500
- `fat`: gte=0, lte=500
- `quantity`: gte=1, lte=10000
- `unit`: required
- `confidence`: gte=0, lte=1

### EstimateQuantityRequest
- `uri`: required
- `name`: required
- `type`: required
- `mealLocation`: required
- `referenceServingSize`: optional
- `referenceServingUnit`: optional

### EstimateQuantityResponse
- `estimatedQuantity`: required, gte=1, lte=500
- `unit`: required, oneof=g ml serving
- `confidence`: required, gte=0, lte=1
- `reasoning`: optional

## Integrações Externas

### AWS S3
- Upload de imagens de alimentos
- Armazenamento permanente
- URLs públicas para acesso

### Lambda Functions
1. **Food Recognition Lambda**
   - Análise de imagem com Gemini/OpenAI
   - Identificação de múltiplos alimentos
   - Cálculo de macronutrientes
   - Score de confiança

2. **Quantity Estimation Lambda**
   - Estimativa de porções
   - Análise de contexto (prato, local)
   - Comparação com referências
   - Raciocínio explicado

### OpenFoodFacts API
- Base de dados de alimentos
- Busca por código de barras
- Informações nutricionais padronizadas
- Cache local para performance

## Response Padrão

### Sucesso (Food Recognition)
```json
{
  "success": true,
  "data": {
    "foodItems": [
      {
        "name": "Arroz branco",
        "calories": 130,
        "protein": 3,
        "carbs": 28,
        "fat": 0,
        "quantity": 150,
        "unit": "g",
        "confidence": 0.92
      }
    ],
    "processingTime": 1500
  },
  "message": "food recognized successfully"
}
```

### Sucesso (Barcode)
```json
{
  "success": true,
  "data": {
    "barcode": "7894900011517",
    "name": "Leite Integral",
    "brand": "Marca X",
    "calories": 60,
    "protein": 3,
    "carbs": 5,
    "fat": 3,
    "servingSize": 200,
    "servingUnit": "ml",
    "source": "OpenFoodFacts"
  },
  "message": "food retrieved successfully"
}
```

### Erro
```json
{
  "success": false,
  "message": "error message"
}
```

## Segurança
- ✅ Autenticação via Firebase Token (middleware `AuthRequired`)
- ✅ Validação de tipos de arquivo para upload
- ✅ Validação de tamanho de arquivo (via configuração HTTP)
- ✅ URLs assinadas do S3 para acesso controlado
- ✅ Rate limiting via middleware (já existente)

## Logging
- ✅ Logs estruturados com domain.Logger
- ✅ Rastreamento de erros de upload
- ✅ Rastreamento de chamadas à Lambda
- ✅ Rastreamento de cache hits/misses
- ✅ Logs de tempo de processamento

## Tabelas Utilizadas
- `food_database` (READ/WRITE para cache de barcode)

## Testes de Compilação
- ✅ Código compilado sem erros
- ✅ Sem erros de lint detectados
- ✅ Dependências organizadas com go mod tidy

## Variáveis de Ambiente Necessárias

Adicionar ao `.env`:
```env
# AWS S3
AWS_S3_BUCKET_NAME=calorie-ai-images
AWS_S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# Lambda Functions
LAMBDA_FOOD_RECOGNITION_URL=https://your-lambda-url.amazonaws.com/recognize
LAMBDA_QUANTITY_ESTIMATION_URL=https://your-lambda-url.amazonaws.com/estimate

# OpenFoodFacts (opcional, usa default se não configurado)
OPENFOODFACTS_API_URL=https://world.openfoodfacts.org/api/v2
```

## Próximos Passos
- Fase 6: Food Items CRUD
- Implementar testes unitários para os services
- Adicionar métricas de performance
- Implementar retry logic para chamadas externas
- Adicionar suporte a mais providers de IA

## Notas Técnicas

### Diferenciação de FoodItem
- `domain.FoodItem` (meal.go) - representa item em uma refeição (com IDs, datas)
- `domain.RecognizedFoodItem` (food_recognition.go) - representa item reconhecido pela IA (sem IDs)

### Conversões de Tipos
- Criadas funções auxiliares para conversão entre:
  - `int32` ↔ `int` (para SQLC)
  - `float64` ↔ `pgtype.Numeric` (para PostgreSQL)
  - Tratamento de ponteiros nulos

### Cache Strategy
- Primeiro, busca no banco de dados local
- Se não encontrar, busca na API OpenFoodFacts
- Armazena resultado no banco para próximas consultas
- Logs indicam se foi cache hit ou miss

### Timeouts
- HTTP Client: 30s para Lambda calls
- HTTP Client: 10s para OpenFoodFacts
- Configurável via service se necessário
