## Meals Endpoints

### 1. POST /meals
Criar refeição com food items.

#### ✅ Sucesso (201 Created)
```json
{
  "success": true,
  "message": "meal created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "meal_type": "breakfast",
    "meal_date": "2026-02-24T00:00:00Z",
    "meal_time": "2000-01-01T08:30:00Z",
    "photo_url": "https://s3.amazonaws.com/bucket/photo.jpg",
    "created_at": "2026-02-24T08:30:15Z",
    "foods": [
      {
        "id": "650e8400-e29b-41d4-a716-446655440001",
        "meal_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Arroz integral",
        "portion_size": 150,
        "portion_unit": "g",
        "calories": 195,
        "protein": 4.2,
        "carbs": 43.0,
        "fat": 0.4,
        "source": "manual"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "meal_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Ovo mexido",
        "portion_size": 2,
        "portion_unit": "unit",
        "calories": 140,
        "protein": 12.0,
        "carbs": 1.0,
        "fat": 10.0,
        "source": "ai_photo"
      }
    ]
  }
}
```

---

### 2. GET /meals?date=2026-02-24
Listar todas as refeições de um dia específico.

**Query Parameters:**
- `date` (opcional): Data no formato `YYYY-MM-DD`. Default: data atual.

#### ✅ Sucesso (200 OK)
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "meal_type": "breakfast",
      "meal_date": "2026-02-24T00:00:00Z",
      "meal_time": "2000-01-01T08:30:00Z",
      "photo_url": "https://s3.amazonaws.com/bucket/photo.jpg",
      "created_at": "2026-02-24T08:30:15Z",
      "foods": [
        {
          "id": "650e8400-e29b-41d4-a716-446655440001",
          "meal_id": "550e8400-e29b-41d4-a716-446655440000",
          "name": "Arroz integral",
          "portion_size": 150,
          "portion_unit": "g",
          "calories": 195,
          "protein": 4.2,
          "carbs": 43.0,
          "fat": 0.4,
          "source": "manual"
        }
      ]
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "meal_type": "lunch",
      "meal_date": "2026-02-24T00:00:00Z",
      "meal_time": "2000-01-01T12:00:00Z",
      "created_at": "2026-02-24T12:05:30Z",
      "foods": [
        {
          "id": "670e8400-e29b-41d4-a716-446655440003",
          "meal_id": "660e8400-e29b-41d4-a716-446655440002",
          "name": "Frango grelhado",
          "portion_size": 150,
          "portion_unit": "g",
          "calories": 247,
          "protein": 46.5,
          "carbs": 0.0,
          "fat": 5.4,
          "source": "barcode"
        }
      ]
    }
  ]
}
```

#### ✅ Sucesso - Lista vazia (200 OK)
```json
{
  "success": true,
  "data": []
}
```

---

## Food Items Endpoints

### 1. POST /food-items
Adicionar food item a uma refeição existente.

#### ✅ Sucesso (201 Created)
```json
{
  "success": true,
  "message": "food item created successfully",
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440003",
    "meal_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Frango grelhado",
    "portion_size": 100,
    "portion_unit": "g",
    "calories": 165,
    "protein": 31.0,
    "carbs": 0.0,
    "fat": 3.6,
    "source": "manual"
  }
}
```

---

### 2. GET /food-items/:id
Buscar um food item específico pelo ID.

#### ✅ Sucesso (200 OK)
```json
{
  "success": true,
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440003",
    "meal_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Frango grelhado",
    "portion_size": 100,
    "portion_unit": "g",
    "calories": 165,
    "protein": 31.0,
    "carbs": 0.0,
    "fat": 3.6,
    "source": "manual"
  }
}
```

---

### 3. PUT /food-items/:id
Atualizar um food item existente.

#### ✅ Sucesso (200 OK)
```json
{
  "success": true,
  "message": "food item updated successfully",
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440003",
    "meal_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Frango grelhado (corrigido)",
    "portion_size": 120,
    "portion_unit": "g",
    "calories": 198,
    "protein": 37.2,
    "carbs": 0.0,
    "fat": 4.3,
    "source": "manual"
  }
}
```

---

### 4. DELETE /food-items/:id
Remover um food item.

#### ✅ Sucesso (200 OK)
```json
{
  "success": true,
  "message": "food item deleted successfully"
}
```

---

## Códigos de Erro

| Código | Descrição |
|--------|-----------|
| `invalid_request_body` | Corpo da requisição malformado ou inválido |
| `validation_failed` | Falha na validação dos campos (detalhes na mensagem) |
| `invalid_date_format` | Formato de data inválido (use YYYY-MM-DD) |
| `invalid_user_id` | UUID de usuário inválido |
| `invalid_food_item_id` | UUID de food item inválido |
| `missing_query_parameter` | Query parameter obrigatório ausente |
| `missing_food_name` | Campo food_name não informado |
| `create_meal_failed` | Falha ao criar refeição no banco |
| `get_meals_failed` | Falha ao buscar refeições |
| `create_food_item_failed` | Falha ao criar food item |
| `food_item_not_found` | Food item não encontrado |
| `update_food_item_failed` | Falha ao atualizar food item |
| `delete_food_item_failed` | Falha ao deletar food item |

---

## Notas

- Todos os timestamps estão em formato ISO 8601 (UTC).
- UUIDs são strings no formato `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
- Campos opcionais podem ser `null` ou omitidos da resposta.
- A propriedade `message` em respostas de sucesso é opcional e pode não estar presente em todos os casos.
