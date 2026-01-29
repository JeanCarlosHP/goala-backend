# Testes de API - Reconhecimento de Alimentos

## Testando com cURL

### 1. Reconhecimento de Alimentos
```bash
curl -X POST http://localhost:8080/api/v1/food/recognize \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/image.jpg" \
  -F "name=Foto do almoço" \
  -F "type=lunch" \
  -F "mealLocation=Casa" \
  -F "uri=image.jpg" \
  --no-buffer
```

### 2. Estimativa de Quantidade
```bash
curl -X POST http://localhost:8080/api/v1/food/estimate-quantity \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/image.jpg" \
  -F "name=Pizza" \
  -F "type=dinner" \
  -F "mealLocation=Restaurante" \
  -F "uri=pizza.jpg" \
  -F "referenceServingSize=100" \
  -F "referenceServingUnit=g" \
  --no-buffer
```

### 3. Busca por Código de Barras
```bash
curl -X GET http://localhost:8080/api/v1/food/barcode/7891234567890 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Testando com Postman

### Configuração do Postman para Reconhecimento

1. Método: `POST`
2. URL: `http://localhost:8080/api/v1/food/recognize`
3. Headers:
   - `Authorization`: `Bearer YOUR_TOKEN`
4. Body (form-data):
   - `image`: [Selecionar arquivo]
   - `name`: `Foto do prato`
   - `type`: `lunch`
   - `mealLocation`: `Casa`
   - `uri`: `plate.jpg`

**Importante**: Para ver o streaming no Postman, você precisará usar um cliente que suporte SSE ou fazer o parsing manual dos eventos.

## Campos Obrigatórios

### Para `/food/recognize`
- ✅ `image` (file): Imagem do alimento
- ✅ `name` (string): Nome/descrição da foto
- ✅ `type` (string): Tipo da refeição (breakfast, lunch, dinner, snack)
- ✅ `mealLocation` (string): Local da refeição
- ⚪ `uri` (string): URI/nome da imagem (opcional)

### Para `/food/estimate-quantity`
- ✅ `image` (file): Imagem do alimento
- ✅ `name` (string): Nome do alimento
- ✅ `type` (string): Tipo da refeição
- ✅ `mealLocation` (string): Local da refeição
- ⚪ `uri` (string): URI/nome da imagem (opcional)
- ⚪ `referenceServingSize` (string): Tamanho de referência (opcional)
- ⚪ `referenceServingUnit` (string): Unidade de referência (opcional)

### Para `/food/barcode/:barcode`
- ✅ `barcode` (path param): Código de barras

## Formato da Resposta

### Streaming (SSE)
Durante o processamento, você receberá eventos:

```
event: progress
data: {"stage":"upload","percentage":10,"message":"Reading image data..."}

event: progress
data: {"stage":"upload","percentage":20,"message":"Uploading to S3..."}

event: progress
data: {"stage":"ai_analysis","percentage":50,"message":"Analyzing image..."}

event: complete
data: {"success":true,"data":{"foodItems":[...]},"message":"food recognized successfully"}
```

### Resposta Final
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
        "confidence": 0.95
      },
      {
        "name": "Feijão preto",
        "calories": 77,
        "protein": 5,
        "carbs": 14,
        "fat": 0,
        "quantity": 100,
        "unit": "g"
      }
    ]
  },
  "message": "food recognized successfully"
}
```

## Códigos de Status HTTP

- `200 OK`: Sucesso (resposta em streaming)
- `400 Bad Request`: Parâmetros inválidos ou ausentes
- `401 Unauthorized`: Token inválido ou ausente
- `404 Not Found`: Recurso não encontrado (barcode)
- `500 Internal Server Error`: Erro no servidor
