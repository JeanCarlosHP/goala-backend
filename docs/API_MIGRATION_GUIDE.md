# Guia de Migração de API - Calorie AI Backend

**Data:** 15 de Janeiro de 2026  
**Versão da API:** Atualização conforme documentação Apidog  
**Link da Documentação:** https://mu0v6q54o3.apidog.io/

---

## 📋 Índice

1. [Resumo Executivo](#resumo-executivo)
2. [Fases de Implementação](#fases-de-implementação)
3. [Rotas Novas](#rotas-novas)
4. [Rotas Alteradas](#rotas-alteradas)
5. [Rotas Removidas](#rotas-removidas)
6. [Schemas e Banco de Dados](#schemas-e-banco-de-dados)
7. [Contratos de Validação](#contratos-de-validação)
8. [Checklist de Validação](#checklist-de-validação)

---

## 🎯 Resumo Executivo

### Estatísticas da Migração
- **Total de Rotas Novas:** 8
- **Total de Rotas Alteradas:** 3
- **Total de Rotas Removidas:** 0
- **Novas Tabelas:** 4
- **Tabelas Alteradas:** 2

### Principais Mudanças
1. Sistema de achievements/conquistas
2. Sistema de estatísticas e métricas de usuário
3. Sistema de feedback
4. Reconhecimento de comida por imagem (AI)
5. Busca de alimentos por código de barras
6. Estimativa de quantidade de alimentos por imagem
7. Perfil de usuário expandido com novas preferências

---

## 🔄 Fases de Implementação

### **FASE 1: Infraestrutura Base** ⏱️ Estimativa: 2-3 dias

#### Objetivo
Criar a estrutura de banco de dados necessária e modelos de domínio para suportar as novas features.

#### Tarefas
1. Criar novas migrations
2. Atualizar schemas SQLC
3. Criar novos modelos de domínio
4. Atualizar tabela `users` com novos campos
5. Criar tabela `user_stats`
6. Criar tabela `achievements`
7. Criar tabela `user_achievements`
8. Criar tabela `feedback`

#### Entregáveis
- ✅ Migrations criadas e testadas
- ✅ Modelos de domínio implementados
- ✅ Queries SQLC configuradas

---

### **FASE 2: User Profile & Goals** ⏱️ Estimativa: 2-3 dias

#### Objetivo
Implementar os endpoints de perfil de usuário com os novos campos de configuração.

#### Tarefas
1. Atualizar endpoint `GET /user/profile`
2. Atualizar endpoint `PUT /user/profile`
3. Implementar validações com validator/v10
4. Criar repositório de goals
5. Atualizar service de usuário

#### Entregáveis
- ✅ Profile endpoints atualizados
- ✅ Validações implementadas
- ✅ Testes unitários

---

### **FASE 3: Stats & Achievements** ⏱️ Estimativa: 3-4 dias

#### Objetivo
Implementar sistema de estatísticas e achievements.

#### Tarefas
1. Implementar `GET /stats`
2. Implementar `GET /stats/range`
3. Implementar `GET /achievements`
4. Implementar `POST /achievements/sync`
5. Criar lógica de cálculo de streaks
6. Criar sistema de unlock de achievements
7. Implementar agregações de estatísticas

#### Entregáveis
- ✅ Endpoints de stats funcionando
- ✅ Sistema de achievements completo
- ✅ Cálculo de métricas automático

---

### **FASE 4: Feedback System** ⏱️ Estimativa: 1 dia

#### Objetivo
Implementar sistema de feedback do usuário.

#### Tarefas
1. Implementar `POST /feedback`
2. Criar repositório de feedback
3. Implementar validações
4. Configurar notificações (opcional)

#### Entregáveis
- ✅ Endpoint de feedback funcionando
- ✅ Armazenamento de feedbacks

---

### **FASE 5: Food Recognition & Barcode** ⏱️ Estimativa: 4-5 dias

#### Objetivo
Implementar features de reconhecimento de alimentos por imagem e barcode.

#### Tarefas
1. Implementar `POST /food/recognize`
2. Implementar `GET /food/barcode/{barcode}`
3. Implementar `POST /food/estimate-quantity`
4. Integrar com Lambda (Gemini/OpenAI)
5. Implementar upload de imagens para S3
6. Criar sistema de cache para barcodes
7. Processamento síncrono (S3 dispara Lambda)

#### Entregáveis
- ✅ Reconhecimento de imagem funcionando
- ✅ Busca por barcode implementada
- ✅ Estimativa de quantidade implementada
- ✅ Integração com Lambda

---

### **FASE 6: Food Items CRUD** ⏱️ Estimativa: 2 dias

#### Objetivo
Atualizar endpoints de CRUD de food items.

#### Tarefas
1. Implementar `POST /food-items`
2. Implementar `PUT /food-items/{id}`
3. Atualizar validações
4. Implementar filtros por tipo de refeição

#### Entregáveis
- ✅ CRUD completo de food items
- ✅ Validações atualizadas

---

## 🆕 Rotas Novas

### 1. `GET /user/profile`

**Descrição:** Obter perfil completo do usuário autenticado

**Método:** GET  
**Autenticação:** Bearer Token (Firebase)

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "email",
    "photo": "uri",
    "dailyCalorieGoal": 2000,
    "dailyProteinGoal": 150,
    "dailyCarbsGoal": 200,
    "dailyFatGoal": 65,
    "weight": 70,
    "height": 175,
    "age": 30,
    "gender": "male",
    "activityLevel": "moderate",
    "language": "en-US",
    "notificationsEnabled": true,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  },
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type UserProfileResponse struct {
    ID                     string  `json:"id" validate:"required,uuid"`
    Name                   string  `json:"name" validate:"required"`
    Email                  string  `json:"email" validate:"required,email"`
    Photo                  *string `json:"photo" validate:"omitempty,url"`
    DailyCalorieGoal       int32   `json:"dailyCalorieGoal" validate:"required,gte=0,lte=10000"`
    DailyProteinGoal       int32   `json:"dailyProteinGoal" validate:"required,gte=0,lte=1000"`
    DailyCarbsGoal         *int32  `json:"dailyCarbsGoal" validate:"omitempty,gte=0,lte=2000"`
    DailyFatGoal           *int32  `json:"dailyFatGoal" validate:"omitempty,gte=0,lte=1000"`
    Weight                 *int32  `json:"weight" validate:"omitempty,gt=0,lte=1000"`
    Height                 *int32  `json:"height" validate:"omitempty,gt=0,lte=300"`
    Age                    *int32  `json:"age" validate:"omitempty,gt=0,lte=150"`
    Gender                 *string `json:"gender" validate:"omitempty,oneof=male female other"`
    ActivityLevel          *string `json:"activityLevel" validate:"omitempty,oneof=sedentary light moderate active very_active"`
    Language               string  `json:"language" validate:"required,oneof=en-US pt-BR"`
    NotificationsEnabled   bool    `json:"notificationsEnabled"`
    CreatedAt              *string `json:"createdAt"`
    UpdatedAt              *string `json:"updatedAt"`
}
```

#### Tabelas Envolvidas
- `users` (READ)
- `user_goals` (READ)

---

### 2. `PUT /user/profile`

**Descrição:** Atualizar perfil do usuário

**Método:** PUT  
**Autenticação:** Bearer Token (Firebase)

#### Request Body
```json
{
  "id": "uuid",
  "name": "string",
  "email": "user@example.com",
  "photo": "http://example.com",
  "dailyCalorieGoal": 2000,
  "dailyProteinGoal": 150,
  "dailyCarbsGoal": 200,
  "dailyFatGoal": 65,
  "weight": 70,
  "height": 175,
  "age": 30,
  "gender": "male",
  "activityLevel": "moderate",
  "language": "en-US",
  "notificationsEnabled": true
}
```

#### Validação (go-playground/validator)
```go
type UpdateProfileRequest struct {
    ID                     string  `json:"id" validate:"required,uuid"`
    Name                   string  `json:"name" validate:"required,min=2,max=255"`
    Email                  string  `json:"email" validate:"required,email"`
    Photo                  *string `json:"photo" validate:"omitempty,url"`
    DailyCalorieGoal       int32   `json:"dailyCalorieGoal" validate:"required,gte=0,lte=10000"`
    DailyProteinGoal       int32   `json:"dailyProteinGoal" validate:"required,gte=0,lte=1000"`
    DailyCarbsGoal         *int32  `json:"dailyCarbsGoal" validate:"omitempty,gte=0,lte=2000"`
    DailyFatGoal           *int32  `json:"dailyFatGoal" validate:"omitempty,gte=0,lte=1000"`
    Weight                 *int32  `json:"weight" validate:"omitempty,gt=0,lte=1000"`
    Height                 *int32  `json:"height" validate:"omitempty,gt=0,lte=300"`
    Age                    *int32  `json:"age" validate:"omitempty,gt=0,lte=150"`
    Gender                 *string `json:"gender" validate:"omitempty,oneof=male female other"`
    ActivityLevel          *string `json:"activityLevel" validate:"omitempty,oneof=sedentary light moderate active very_active"`
    Language               string  `json:"language" validate:"required,oneof=en-US pt-BR"`
    NotificationsEnabled   bool    `json:"notificationsEnabled"`
}
```

#### Tabelas Envolvidas
- `users` (UPDATE)
- `user_goals` (UPDATE/INSERT)

---

### 3. `GET /stats/range`

**Descrição:** Obter estatísticas em um intervalo de datas com paginação

**Método:** GET  
**Autenticação:** Bearer Token (Firebase)

#### Query Parameters
- `startDate` (required): string date (formato: yyyy-MM-dd)
- `endDate` (required): string date (formato: yyyy-MM-dd)
- `page` (optional): integer, default: 1
- `limit` (optional): integer, default: 30

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "days": [
      {
        "date": "2024-01-15T00:00:00Z",
        "totalCalories": 2000,
        "totalProtein": 150,
        "totalCarbs": 200,
        "totalFat": 65,
        "meals": [
          {
            "id": "uuid",
            "type": "breakfast",
            "name": "Café da Manhã",
            "calories": 500,
            "protein": 30,
            "carbs": 50,
            "fat": 20,
            "itemCount": 3,
            "items": [
              {
                "id": "uuid",
                "name": "Aveia",
                "calories": 150,
                "protein": 10,
                "carbs": 25,
                "fat": 5,
                "quantity": 50,
                "unit": "g",
                "imageUrl": "http://example.com",
                "date": "2024-01-15T08:00:00Z"
              }
            ]
          }
        ],
        "waterIntake": 2000
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 30,
      "total": 100,
      "totalPages": 4,
      "hasNext": true,
      "hasPrev": false
    },
    "aggregated": {
      "totalCalories": 60000,
      "totalProtein": 4500,
      "totalCarbs": 6000,
      "totalFat": 1950,
      "avgCalories": 2000,
      "avgProtein": 150,
      "avgCarbs": 200,
      "avgFat": 65
    }
  },
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type StatsRangeQuery struct {
    StartDate string `query:"startDate" validate:"required,datetime=2006-01-02"`
    EndDate   string `query:"endDate" validate:"required,datetime=2006-01-02"`
    Page      int    `query:"page" validate:"omitempty,gte=1"`
    Limit     int    `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

type DayStats struct {
    Date          time.Time `json:"date" validate:"required"`
    TotalCalories int32     `json:"totalCalories" validate:"gte=0"`
    TotalProtein  int32     `json:"totalProtein" validate:"gte=0"`
    TotalCarbs    int32     `json:"totalCarbs" validate:"gte=0"`
    TotalFat      int32     `json:"totalFat" validate:"gte=0"`
    Meals         []Meal    `json:"meals" validate:"dive"`
    WaterIntake   int32     `json:"waterIntake" validate:"gte=0"`
}

type Pagination struct {
    Page       int  `json:"page" validate:"required,gte=1"`
    Limit      int  `json:"limit" validate:"required,gte=1,lte=100"`
    Total      int  `json:"total" validate:"gte=0"`
    TotalPages int  `json:"totalPages" validate:"gte=0"`
    HasNext    bool `json:"hasNext"`
    HasPrev    bool `json:"hasPrev"`
}

type AggregatedStats struct {
    TotalCalories int32 `json:"totalCalories" validate:"gte=0"`
    TotalProtein  int32 `json:"totalProtein" validate:"gte=0"`
    TotalCarbs    int32 `json:"totalCarbs" validate:"gte=0"`
    TotalFat      int32 `json:"totalFat" validate:"gte=0"`
    AvgCalories   int32 `json:"avgCalories" validate:"gte=0"`
    AvgProtein    int32 `json:"avgProtein" validate:"gte=0"`
    AvgCarbs      int32 `json:"avgCarbs" validate:"gte=0"`
    AvgFat        int32 `json:"avgFat" validate:"gte=0"`
}
```

#### Tabelas Envolvidas
- `meals` (READ com JOIN)
- `food_items` (READ com JOIN)

---

### 4. `POST /feedback`

**Descrição:** Enviar feedback do usuário

**Método:** POST  
**Autenticação:** Bearer Token (Firebase)

#### Request Body
```json
{
  "type": "problem",
  "title": "Bug no login",
  "description": "Descrição detalhada do problema",
  "userEmail": "user@example.com",
  "deviceInfo": {
    "platform": "iOS",
    "osVersion": "17.0",
    "appVersion": "1.0.0"
  }
}
```

#### Validação (go-playground/validator)
```go
type FeedbackRequest struct {
    Type        string       `json:"type" validate:"required,oneof=problem improvement"`
    Title       string       `json:"title" validate:"required,min=3,max=255"`
    Description string       `json:"description" validate:"required,min=10,max=5000"`
    UserEmail   string       `json:"userEmail" validate:"required,email"`
    DeviceInfo  *DeviceInfo  `json:"deviceInfo" validate:"omitempty"`
}

type DeviceInfo struct {
    Platform   string `json:"platform" validate:"required"`
    OsVersion  string `json:"osVersion" validate:"required"`
    AppVersion string `json:"appVersion" validate:"required"`
}
```

#### Response 200 (Success)
```json
{}
```

#### Tabelas Envolvidas
- `feedback` (INSERT) - **NOVA TABELA**

---

### 5. `GET /achievements`

**Descrição:** Obter achievements do usuário

**Método:** GET  
**Autenticação:** Bearer Token (Firebase)

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "achievements": [
      {
        "id": "first_meal",
        "nameKey": "achievements.first_meal.name",
        "descriptionKey": "achievements.first_meal.description",
        "icon": "🍽️",
        "unlocked": true,
        "unlockedAt": "2024-01-15T00:00:00Z",
        "progress": 1,
        "target": 1
      }
    ],
    "stats": {
      "currentStreak": 5,
      "bestStreak": 10,
      "totalMealsLogged": 100,
      "totalCaloriesLogged": 200000,
      "totalDaysLogged": 30,
      "averageCaloriesPerDay": 2000
    }
  },
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type Achievement struct {
    ID             string  `json:"id" validate:"required"`
    NameKey        string  `json:"nameKey" validate:"required"`
    DescriptionKey string  `json:"descriptionKey" validate:"required"`
    Icon           string  `json:"icon" validate:"required"`
    Unlocked       bool    `json:"unlocked"`
    UnlockedAt     *string `json:"unlockedAt" validate:"omitempty"`
    Progress       int32   `json:"progress" validate:"gte=0"`
    Target         int32   `json:"target" validate:"gte=0"`
}

type UserStats struct {
    CurrentStreak          int32 `json:"currentStreak" validate:"gte=0"`
    BestStreak             int32 `json:"bestStreak" validate:"gte=0"`
    TotalMealsLogged       int32 `json:"totalMealsLogged" validate:"gte=0"`
    TotalCaloriesLogged    int32 `json:"totalCaloriesLogged" validate:"gte=0"`
    TotalDaysLogged        int32 `json:"totalDaysLogged" validate:"gte=0"`
    AverageCaloriesPerDay  int32 `json:"averageCaloriesPerDay" validate:"gte=0"`
}
```

#### Tabelas Envolvidas
- `achievements` (READ) - **NOVA TABELA**
- `user_achievements` (READ) - **NOVA TABELA**
- `user_stats` (READ) - **NOVA TABELA**

---

### 6. `GET /stats`

**Descrição:** Obter estatísticas do usuário

**Método:** GET  
**Autenticação:** Bearer Token (Firebase)

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "currentStreak": 5,
    "bestStreak": 10,
    "totalMealsLogged": 100,
    "totalCaloriesLogged": 200000,
    "totalDaysLogged": 30,
    "averageCaloriesPerDay": 2000
  },
  "message": "string"
}
```

#### Validação
Usar `UserStats` definido anteriormente

#### Tabelas Envolvidas
- `user_stats` (READ) - **NOVA TABELA**

---

### 7. `POST /achievements/sync`

**Descrição:** Sincronizar achievements (calcular progresso e desbloquear)

**Método:** POST  
**Autenticação:** Bearer Token (Firebase)

#### Response 200 (Success)
Mesmo formato de `GET /achievements`

#### Tabelas Envolvidas
- `achievements` (READ)
- `user_achievements` (UPDATE/INSERT)
- `user_stats` (UPDATE)
- `meals` (READ para cálculos)

---

### 8. `POST /food/recognize`

**Descrição:** Reconhecer comida por imagem usando AI

**Método:** POST  
**Content-Type:** multipart/form-data  
**Autenticação:** Bearer Token (Firebase)

#### Request Body (Form Data)
- `uri` (required): string - URL ou base64 da imagem
- `name` (required): string - Nome do arquivo
- `type` (required): string - MIME type
- `mealLocation` (required): string - Localização GPS ou contexto

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "foodItems": [
      {
        "name": "Arroz",
        "calories": 130,
        "protein": 3,
        "carbs": 28,
        "fat": 0,
        "quantity": 100,
        "unit": "g",
        "confidence": 0.95
      }
    ],
    "processingTime": 1500
  },
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type FoodRecognitionRequest struct {
    URI          string `form:"uri" validate:"required"`
    Name         string `form:"name" validate:"required"`
    Type         string `form:"type" validate:"required"`
    MealLocation string `form:"mealLocation" validate:"required"`
}

type FoodItem struct {
    Name       string  `json:"name" validate:"required"`
    Calories   int32   `json:"calories" validate:"gte=0,lte=5000"`
    Protein    int32   `json:"protein" validate:"gte=0,lte=500"`
    Carbs      int32   `json:"carbs" validate:"gte=0,lte=500"`
    Fat        int32   `json:"fat" validate:"gte=0,lte=500"`
    Quantity   int32   `json:"quantity" validate:"gte=1,lte=10000"`
    Unit       string  `json:"unit" validate:"required"`
    Confidence float64 `json:"confidence" validate:"gte=0,lte=1"`
}

type FoodRecognitionResponse struct {
    FoodItems      []FoodItem `json:"foodItems" validate:"required,dive"`
    ProcessingTime int32      `json:"processingTime" validate:"gte=0"`
}
```

#### Integrações
- AWS S3 (upload de imagem)
- Lambda (Gemini/OpenAI Analyzer)

#### Tabelas Envolvidas
Nenhuma tabela direta (processamento assíncrono)

---

### 9. `GET /food/barcode/{barcode}`

**Descrição:** Buscar informações de alimento por código de barras

**Método:** GET  
**Autenticação:** Bearer Token (Firebase)

#### Path Parameters
- `barcode` (required): string - Código de barras (EAN, UPC, etc)

#### Response 200 (Success)
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
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type FoodBarcodeResponse struct {
    Barcode     string  `json:"barcode" validate:"required"`
    Name        string  `json:"name" validate:"required"`
    Brand       *string `json:"brand"`
    Calories    int32   `json:"calories" validate:"required,gte=0,lte=5000"`
    Protein     int32   `json:"protein" validate:"required,gte=0,lte=5000"`
    Carbs       int32   `json:"carbs" validate:"required,gte=0,lte=5000"`
    Fat         int32   `json:"fat" validate:"required,gte=0,lte=5000"`
    ServingSize *int32  `json:"servingSize" validate:"omitempty,gte=1,lte=5000"`
    ServingUnit *string `json:"servingUnit"`
    Source      *string `json:"source"`
}
```

#### Integrações
- OpenFoodFacts API
- Cache (Redis/In-Memory)

#### Tabelas Envolvidas
- `food_database` (READ/INSERT para cache) - **NOVA TABELA**

---

### 10. `POST /food-items`

**Descrição:** Criar novo food item

**Método:** POST  
**Autenticação:** Bearer Token (Firebase)

#### Request Body
```json
{
  "name": "Arroz",
  "calories": 130,
  "protein": 3,
  "carbs": 28,
  "fat": 0,
  "quantity": 100,
  "unit": "g",
  "imageUrl": "http://example.com/image.jpg",
  "mealType": "lunch",
  "date": 1705276800000
}
```

#### Validação (go-playground/validator)
```go
type CreateFoodRequest struct {
    Name     string  `json:"name" validate:"required,min=1,max=255"`
    Calories int32   `json:"calories" validate:"required,gte=0,lte=2500"`
    Protein  int32   `json:"protein" validate:"required,gte=0,lte=300"`
    Carbs    int32   `json:"carbs" validate:"required,gte=0,lte=500"`
    Fat      int32   `json:"fat" validate:"required,gte=0,lte=500"`
    Quantity *int32  `json:"quantity" validate:"omitempty,gte=1,lte=100000"`
    Unit     *string `json:"unit"`
    ImageUrl *string `json:"imageUrl" validate:"omitempty,url"`
    MealType *string `json:"mealType" validate:"omitempty,oneof=breakfast lunch dinner snack"`
    Date     *int64  `json:"date" validate:"omitempty,gte=0"`
}
```

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Arroz",
    "calories": 130,
    "protein": 3,
    "carbs": 28,
    "fat": 0,
    "quantity": 100,
    "unit": "g",
    "imageUrl": "http://example.com/image.jpg",
    "date": "2024-01-15T00:00:00Z"
  },
  "message": "string"
}
```

#### Tabelas Envolvidas
- `meals` (INSERT se não existir)
- `food_items` (INSERT)
- `user_stats` (UPDATE - incrementar contadores)

---

### 11. `POST /food/estimate-quantity`

**Descrição:** Estimar quantidade de alimento por imagem

**Método:** POST  
**Content-Type:** multipart/form-data  
**Autenticação:** Bearer Token (Firebase)

#### Request Body (Form Data)
- `uri` (required): string
- `name` (required): string
- `type` (required): string
- `productName` (optional): string
- `referenceServingSize` (optional): string
- `referenceServingUnit` (optional): string

#### Response 200 (Success)
```json
{
  "success": true,
  "data": {
    "estimatedQuantity": 150,
    "unit": "g",
    "confidence": 0.85,
    "reasoning": "Based on plate size and reference objects"
  },
  "message": "string"
}
```

#### Validação (go-playground/validator)
```go
type EstimateQuantityRequest struct {
    URI                  string  `form:"uri" validate:"required"`
    Name                 string  `form:"name" validate:"required"`
    Type                 string  `form:"type" validate:"required"`
    ProductName          *string `form:"productName"`
    ReferenceServingSize *string `form:"referenceServingSize"`
    ReferenceServingUnit *string `form:"referenceServingUnit"`
}

type EstimateQuantityResponse struct {
    EstimatedQuantity int32   `json:"estimatedQuantity" validate:"required,gte=1,lte=500"`
    Unit              string  `json:"unit" validate:"required,oneof=g ml serving"`
    Confidence        float64 `json:"confidence" validate:"required,gte=0,lte=1"`
    Reasoning         *string `json:"reasoning"`
}
```

#### Integrações
- AWS S3 (upload)
- Lambda (AI estimation)

#### Tabelas Envolvidas
Nenhuma (processamento via AI)

---

### 12. `PUT /food-items/{id}`

**Descrição:** Atualizar food item existente

**Método:** PUT  
**Autenticação:** Bearer Token (Firebase)

#### Path Parameters
- `id` (required): string (UUID)

#### Request Body
Mesmo formato de `POST /food-items`

#### Validação (go-playground/validator)
```go
type UpdateFoodRequest struct {
    Name     string  `json:"name" validate:"required,min=1,max=255"`
    Calories int32   `json:"calories" validate:"required,gte=0,lte=2500"`
    Protein  int32   `json:"protein" validate:"required,gte=0,lte=300"`
    Carbs    int32   `json:"carbs" validate:"required,gte=0,lte=500"`
    Fat      int32   `json:"fat" validate:"required,gte=0,lte=500"`
    Quantity *int32  `json:"quantity" validate:"omitempty,gte=1,lte=100000"`
    Unit     *string `json:"unit"`
    ImageUrl *string `json:"imageUrl" validate:"omitempty,url"`
    MealType *string `json:"mealType" validate:"omitempty,oneof=breakfast lunch dinner snack"`
    Date     *int64  `json:"date" validate:"omitempty,gte=0"`
}
```

#### Response 200 (Success)
Mesmo formato de resposta de `POST /food-items`

#### Tabelas Envolvidas
- `food_items` (UPDATE)
- `meals` (UPDATE totais)

---

## 🔄 Rotas Alteradas

### 1. Tabela `users` - Novos Campos

**Campos Adicionados:**
- `weight` INT (peso em kg)
- `height` INT (altura em cm)
- `age` INT (idade)
- `gender` VARCHAR(20) (male, female, other)
- `activity_level` VARCHAR(20) (sedentary, light, moderate, active, very_active)
- `language` VARCHAR(10) DEFAULT 'en-US' (en-US, pt-BR)
- `notifications_enabled` BOOLEAN DEFAULT false

### 2. Tabela `user_goals` - Novos Campos

**Campos Adicionados:**
- `daily_carbs_g` INT (meta de carboidratos)
- `daily_fat_g` INT (meta de gorduras)

### 3. Tabela `food_items` - Alterações

**Campos Alterados:**
- Adicionar restrições de validação conforme API

---

## 🗑️ Rotas Removidas

**Nenhuma rota foi removida nesta atualização.**

---

## 🗄️ Schemas e Banco de Dados

### Nova Tabela: `user_stats`

```sql
CREATE TABLE user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    best_streak INT NOT NULL DEFAULT 0,
    total_meals_logged INT NOT NULL DEFAULT 0,
    total_calories_logged BIGINT NOT NULL DEFAULT 0,
    total_days_logged INT NOT NULL DEFAULT 0,
    average_calories_per_day INT NOT NULL DEFAULT 0,
    last_meal_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);
```

**SQLC Query:**
```sql
-- name: GetUserStats :one
SELECT * FROM user_stats WHERE user_id = $1;

-- name: UpdateUserStats :exec
UPDATE user_stats SET
    current_streak = $2,
    best_streak = $3,
    total_meals_logged = $4,
    total_calories_logged = $5,
    total_days_logged = $6,
    average_calories_per_day = $7,
    last_meal_date = $8,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpsertUserStats :exec
INSERT INTO user_stats (
    user_id, current_streak, best_streak, total_meals_logged,
    total_calories_logged, total_days_logged, average_calories_per_day, last_meal_date
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id) DO UPDATE SET
    current_streak = EXCLUDED.current_streak,
    best_streak = EXCLUDED.best_streak,
    total_meals_logged = EXCLUDED.total_meals_logged,
    total_calories_logged = EXCLUDED.total_calories_logged,
    total_days_logged = EXCLUDED.total_days_logged,
    average_calories_per_day = EXCLUDED.average_calories_per_day,
    last_meal_date = EXCLUDED.last_meal_date,
    updated_at = NOW();
```

---

### Nova Tabela: `achievements`

```sql
CREATE TABLE achievements (
    id VARCHAR(50) PRIMARY KEY,
    name_key VARCHAR(255) NOT NULL,
    description_key VARCHAR(255) NOT NULL,
    icon VARCHAR(50) NOT NULL,
    target INT NOT NULL,
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO achievements (id, name_key, description_key, icon, target, category) VALUES
('first_meal', 'achievements.first_meal.name', 'achievements.first_meal.description', '🍽️', 1, 'getting_started'),
('streak_3', 'achievements.streak_3.name', 'achievements.streak_3.description', '🔥', 3, 'consistency'),
('streak_7', 'achievements.streak_7.name', 'achievements.streak_7.description', '🔥', 7, 'consistency'),
('streak_30', 'achievements.streak_30.name', 'achievements.streak_30.description', '🔥', 30, 'consistency'),
('meals_10', 'achievements.meals_10.name', 'achievements.meals_10.description', '📊', 10, 'logging'),
('meals_50', 'achievements.meals_50.name', 'achievements.meals_50.description', '📊', 50, 'logging'),
('meals_100', 'achievements.meals_100.name', 'achievements.meals_100.description', '📊', 100, 'logging'),
('goal_achieved', 'achievements.goal_achieved.name', 'achievements.goal_achieved.description', '🎯', 1, 'goals');
```

**SQLC Query:**
```sql
-- name: ListAchievements :many
SELECT * FROM achievements ORDER BY category, target;

-- name: GetAchievement :one
SELECT * FROM achievements WHERE id = $1;
```

---

### Nova Tabela: `user_achievements`

```sql
CREATE TABLE user_achievements (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id VARCHAR(50) NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked BOOLEAN NOT NULL DEFAULT false,
    progress INT NOT NULL DEFAULT 0,
    unlocked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_unlocked ON user_achievements(user_id, unlocked);
```

**SQLC Query:**
```sql
-- name: GetUserAchievements :many
SELECT ua.*, a.name_key, a.description_key, a.icon, a.target, a.category
FROM user_achievements ua
JOIN achievements a ON ua.achievement_id = a.id
WHERE ua.user_id = $1
ORDER BY ua.unlocked DESC, a.category, a.target;

-- name: UpsertUserAchievement :exec
INSERT INTO user_achievements (user_id, achievement_id, unlocked, progress, unlocked_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, achievement_id) DO UPDATE SET
    unlocked = EXCLUDED.unlocked,
    progress = EXCLUDED.progress,
    unlocked_at = EXCLUDED.unlocked_at,
    updated_at = NOW();

-- name: GetUserAchievement :one
SELECT * FROM user_achievements
WHERE user_id = $1 AND achievement_id = $2;
```

---

### Nova Tabela: `feedback`

```sql
CREATE TABLE feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('problem', 'improvement')),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    platform VARCHAR(50),
    os_version VARCHAR(50),
    app_version VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_feedback_user_id ON feedback(user_id);
CREATE INDEX idx_feedback_status ON feedback(status);
CREATE INDEX idx_feedback_created_at ON feedback(created_at DESC);
```

**SQLC Query:**
```sql
-- name: CreateFeedback :one
INSERT INTO feedback (
    user_id, type, title, description, user_email,
    platform, os_version, app_version
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListFeedback :many
SELECT * FROM feedback
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetFeedback :one
SELECT * FROM feedback WHERE id = $1;
```

---

### Nova Tabela: `food_database`

```sql
CREATE TABLE food_database (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    barcode VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    calories INT NOT NULL,
    protein_g DECIMAL(10,2) NOT NULL,
    carbs_g DECIMAL(10,2) NOT NULL,
    fat_g DECIMAL(10,2) NOT NULL,
    serving_size INT,
    serving_unit VARCHAR(50),
    source VARCHAR(50),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_food_database_barcode ON food_database(barcode);
CREATE INDEX idx_food_database_name ON food_database(name);
```

**SQLC Query:**
```sql
-- name: GetFoodByBarcode :one
SELECT * FROM food_database WHERE barcode = $1;

-- name: CreateFoodFromBarcode :one
INSERT INTO food_database (
    barcode, name, brand, calories, protein_g, carbs_g, fat_g,
    serving_size, serving_unit, source
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: SearchFoodByName :many
SELECT * FROM food_database
WHERE LOWER(name) LIKE LOWER($1)
ORDER BY verified DESC, name
LIMIT $2;
```

---

### Alteração: Tabela `users`

```sql
ALTER TABLE users ADD COLUMN weight INT CHECK (weight > 0 AND weight <= 1000);
ALTER TABLE users ADD COLUMN height INT CHECK (height > 0 AND height <= 300);
ALTER TABLE users ADD COLUMN age INT CHECK (age > 0 AND age <= 150);
ALTER TABLE users ADD COLUMN gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other'));
ALTER TABLE users ADD COLUMN activity_level VARCHAR(20) CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'active', 'very_active'));
ALTER TABLE users ADD COLUMN language VARCHAR(10) DEFAULT 'en-US' CHECK (language IN ('en-US', 'pt-BR'));
ALTER TABLE users ADD COLUMN notifications_enabled BOOLEAN DEFAULT false;
```

**SQLC Query Atualizada:**
```sql
-- name: GetUserByFirebaseUID :one
SELECT * FROM users WHERE firebase_uid = $1;

-- name: UpdateUserProfile :exec
UPDATE users SET
    display_name = $2,
    email = $3,
    photo_url = $4,
    weight = $5,
    height = $6,
    age = $7,
    gender = $8,
    activity_level = $9,
    language = $10,
    notifications_enabled = $11,
    updated_at = NOW()
WHERE id = $1;
```

---

### Alteração: Tabela `user_goals`

```sql
ALTER TABLE user_goals ADD COLUMN carbs_g INT CHECK (carbs_g >= 0 AND carbs_g <= 2000);
ALTER TABLE user_goals ADD COLUMN fat_g INT CHECK (fat_g >= 0 AND fat_g <= 1000);
```

**SQLC Query Atualizada:**
```sql
-- name: GetUserGoals :one
SELECT * FROM user_goals WHERE user_id = $1;

-- name: UpsertUserGoals :exec
INSERT INTO user_goals (user_id, daily_calories, protein_g, carbs_g, fat_g)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE SET
    daily_calories = EXCLUDED.daily_calories,
    protein_g = EXCLUDED.protein_g,
    carbs_g = EXCLUDED.carbs_g,
    fat_g = EXCLUDED.fat_g,
    updated_at = NOW();
```

---

## ✅ Contratos de Validação

### Enums para Validação

```go
// internal/domain/enum/gender.go
package enum

type Gender string

const (
    GenderMale   Gender = "male"
    GenderFemale Gender = "female"
    GenderOther  Gender = "other"
)

func (g Gender) IsValid() bool {
    switch g {
    case GenderMale, GenderFemale, GenderOther:
        return true
    }
    return false
}
```

```go
// internal/domain/enum/activity_level.go
package enum

type ActivityLevel string

const (
    ActivityLevelSedentary  ActivityLevel = "sedentary"
    ActivityLevelLight      ActivityLevel = "light"
    ActivityLevelModerate   ActivityLevel = "moderate"
    ActivityLevelActive     ActivityLevel = "active"
    ActivityLevelVeryActive ActivityLevel = "very_active"
)

func (a ActivityLevel) IsValid() bool {
    switch a {
    case ActivityLevelSedentary, ActivityLevelLight, ActivityLevelModerate, 
         ActivityLevelActive, ActivityLevelVeryActive:
        return true
    }
    return false
}
```

```go
// internal/domain/enum/language.go
package enum

type Language string

const (
    LanguageEnUS Language = "en-US"
    LanguagePtBR Language = "pt-BR"
)

func (l Language) IsValid() bool {
    switch l {
    case LanguageEnUS, LanguagePtBR:
        return true
    }
    return false
}
```

```go
// internal/domain/enum/meal_type.go
package enum

type MealType string

const (
    MealTypeBreakfast MealType = "breakfast"
    MealTypeLunch     MealType = "lunch"
    MealTypeDinner    MealType = "dinner"
    MealTypeSnack     MealType = "snack"
)

func (m MealType) IsValid() bool {
    switch m {
    case MealTypeBreakfast, MealTypeLunch, MealTypeDinner, MealTypeSnack:
        return true
    }
    return false
}
```

```go
// internal/domain/enum/feedback_type.go
package enum

type FeedbackType string

const (
    FeedbackTypeProblem     FeedbackType = "problem"
    FeedbackTypeImprovement FeedbackType = "improvement"
)

func (f FeedbackType) IsValid() bool {
    switch f {
    case FeedbackTypeProblem, FeedbackTypeImprovement:
        return true
    }
    return false
}
```

```go
// internal/domain/enum/quantity_unit.go
package enum

type QuantityUnit string

const (
    QuantityUnitGram    QuantityUnit = "g"
    QuantityUnitML      QuantityUnit = "ml"
    QuantityUnitServing QuantityUnit = "serving"
)

func (q QuantityUnit) IsValid() bool {
    switch q {
    case QuantityUnitGram, QuantityUnitML, QuantityUnitServing:
        return true
    }
    return false
}
```

### Custom Validators

```go
// pkg/validator/custom_validators.go
package validator

import (
    "github.com/go-playground/validator/v10"
    "github.com/jeancarloshp/calorieai/internal/domain/enum"
)

func RegisterCustomValidators(v *validator.Validate) error {
    if err := v.RegisterValidation("gender", validateGender); err != nil {
        return err
    }
    if err := v.RegisterValidation("activityLevel", validateActivityLevel); err != nil {
        return err
    }
    if err := v.RegisterValidation("language", validateLanguage); err != nil {
        return err
    }
    if err := v.RegisterValidation("mealType", validateMealType); err != nil {
        return err
    }
    if err := v.RegisterValidation("feedbackType", validateFeedbackType); err != nil {
        return err
    }
    if err := v.RegisterValidation("quantityUnit", validateQuantityUnit); err != nil {
        return err
    }
    return nil
}

func validateGender(fl validator.FieldLevel) bool {
    gender := enum.Gender(fl.Field().String())
    return gender.IsValid()
}

func validateActivityLevel(fl validator.FieldLevel) bool {
    level := enum.ActivityLevel(fl.Field().String())
    return level.IsValid()
}

func validateLanguage(fl validator.FieldLevel) bool {
    lang := enum.Language(fl.Field().String())
    return lang.IsValid()
}

func validateMealType(fl validator.FieldLevel) bool {
    mealType := enum.MealType(fl.Field().String())
    return mealType.IsValid()
}

func validateFeedbackType(fl validator.FieldLevel) bool {
    feedbackType := enum.FeedbackType(fl.Field().String())
    return feedbackType.IsValid()
}

func validateQuantityUnit(fl validator.FieldLevel) bool {
    unit := enum.QuantityUnit(fl.Field().String())
    return unit.IsValid()
}
```

---

## 📝 Checklist de Validação

### Fase 1: Infraestrutura Base
- [X] Migration `000003_add_user_stats.up.sql` criada
- [X] Migration `000003_add_user_stats.down.sql` criada
- [X] Migration `000004_add_achievements.up.sql` criada
- [X] Migration `000004_add_achievements.down.sql` criada
- [X] Migration `000005_add_feedback.up.sql` criada
- [X] Migration `000005_add_feedback.down.sql` criada
- [X] Migration `000006_add_food_database.up.sql` criada
- [X] Migration `000006_add_food_database.down.sql` criada
- [X] Migration `000007_alter_users.up.sql` criada
- [X] Migration `000007_alter_users.down.sql` criada
- [X] Migration `000008_alter_user_goals.up.sql` criada
- [X] Migration `000008_alter_user_goals.down.sql` criada
- [X] SQLC queries atualizadas em `pkg/database/queries/`
- [X] Executar `make sqlc-generate`
- [X] Verificar models gerados em `pkg/database/db/models.go`
- [X] Criar enums em `internal/domain/enum/`
- [X] Registrar validators customizados
- [X] Testar migrations em ambiente local
- [X] Rollback migrations testado

### Fase 2: User Profile & Goals
- [X] Criar `internal/domain/user_profile.go`
- [X] Criar `internal/repositories/user_profile_repo.go`
- [X] Atualizar `internal/services/user_service.go`
- [X] Criar handler `GetProfile` em `internal/handlers/auth.go`
- [X] Criar handler `UpdateProfile` em `internal/handlers/auth.go`
- [X] Adicionar rotas no server
- [X] Testes unitários do repositório
- [X] Testes unitários do service
- [X] Testes de integração dos endpoints
- [X] Validar response schemas
- [X] Testar com Postman/Insomnia

### Fase 3: Stats & Achievements
- [X] Criar `internal/domain/stats.go`
- [X] Criar `internal/domain/achievement.go`
- [X] Criar `internal/repositories/stats_repo.go`
- [X] Criar `internal/repositories/achievement_repo.go`
- [X] Criar `internal/services/stats_service.go`
- [X] Criar `internal/services/achievement_service.go`
- [X] Criar `internal/handlers/stats.go`
- [X] Criar `internal/handlers/achievement.go`
- [X] Implementar cálculo de streaks
- [X] Implementar lógica de unlock de achievements
- [X] Implementar agregações em `GET /stats/range`
- [X] Implementar paginação
- [X] Adicionar rotas no server
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Testar cálculos de agregação
- [ ] Testar unlock automático de achievements

### Fase 4: Feedback System
- [X] Criar `internal/domain/feedback.go`
- [X] Criar `internal/repositories/feedback_repo.go`
- [X] Criar `internal/services/feedback_service.go`
- [X] Criar `internal/handlers/feedback.go`
- [X] Adicionar rota no server
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Configurar notificações (opcional)
- [ ] Testar com diferentes tipos de feedback

### Fase 5: Food Recognition & Barcode
- [X] Criar `internal/domain/food_recognition.go`
- [X] Criar `internal/repositories/food_database_repo.go`
- [X] Criar `internal/services/food_recognition_service.go`
- [X] Criar `internal/handlers/food_recognition.go`
- [X] Implementar upload S3
- [X] Implementar cache de barcodes
- [X] Implementar busca na OpenFoodFacts API
- [X] Adicionar rotas no server
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Testar upload de imagens
- [ ] Testar processamento assíncrono
- [ ] Testar cache de barcodes
- [X] Configurar timeouts adequados

### Fase 6: Food Items CRUD
- [X] Atualizar `internal/domain/food.go`
- [X] Atualizar `internal/repositories/food_repo.go`
- [X] Atualizar `internal/services/food_service.go`
- [X] Atualizar `internal/handlers/food.go`
- [X] Implementar `POST /food-items`
- [X] Implementar `PUT /food-items/{id}`
- [X] Atualizar validações
- [X] Adicionar rotas no server
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Testar validações de limites
- [ ] Testar com diferentes meal types

### Validação Final
- [ ] Documentação Swagger/OpenAPI atualizada
- [ ] README.md atualizado
- [ ] Variáveis de ambiente documentadas
- [ ] Docker Compose atualizado
- [ ] Logs estruturados verificados
- [ ] Métricas Prometheus verificadas
- [ ] Dashboards Grafana atualizados

### Checklist de Qualidade de Código
- [ ] Código segue Clean Architecture
- [ ] Nenhum código comentado desnecessário
- [ ] Logs estruturados com zerolog
- [ ] Tratamento adequado de erros
- [ ] Validações com go-playground/validator
- [ ] Código sem magic numbers
- [ ] Constantes definidas para valores fixos
- [ ] Transactions de banco quando necessário
- [ ] Context propagation adequada
- [ ] Graceful shutdown implementado
- [ ] Rate limiting configurado
- [ ] CORS configurado adequadamente
- [ ] Headers de segurança configurados

### Checklist de Segurança
- [ ] Autenticação Firebase em todas as rotas protegidas
- [ ] Validação de ownership (user só acessa seus próprios dados)
- [ ] SQL injection prevenido (usando SQLC)
- [ ] XSS prevenido
- [ ] CSRF protection configurado
- [ ] Rate limiting por IP/usuário
- [ ] Input sanitization
- [ ] File upload validation
- [ ] Secrets não commitados no git
- [ ] Env vars validadas no startup

### Checklist de Performance
- [ ] Índices de banco otimizados
- [ ] Queries com LIMIT apropriado
- [ ] Cache implementado onde necessário
- [ ] Conexões de banco pooled
- [ ] Timeouts configurados
- [ ] Processamento assíncrono onde aplicável
- [ ] Lazy loading implementado
- [ ] Paginação implementada
- [ ] N+1 queries evitados
- [ ] Compressão de responses habilitada

---

## 📊 Dependências Novas

```bash
# Adicionar ao go.mod se necessário
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
```

---

## 🔗 Links Úteis

- [Documentação API](https://mu0v6q54o3.apidog.io/)
- [Go Playground Validator](https://github.com/go-playground/validator)
- [SQLC](https://docs.sqlc.dev/)
- [Fiber Framework](https://docs.gofiber.io/)
- [OpenFoodFacts API](https://world.openfoodfacts.org/data)

---

## 📞 Contato e Suporte

Em caso de dúvidas durante a implementação, consultar:
- Tech Lead
- Product Owner
- Documentação Apidog

---

**Fim do Documento**
