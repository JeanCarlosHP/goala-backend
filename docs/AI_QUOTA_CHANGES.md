# Modificações na Validação de Cotas de IA - Documentação Técnica
As modificações implementadas alteram fundamentalmente como as cotas de recursos de IA são validadas no backend. Anteriormente, as cotas eram baseadas em planos de assinatura (Free, Monthly, Yearly, Trial). Agora, as cotas são **limites diários fixos** por recurso de IA, calculados com base no timezone do usuário.

## Mudanças Principais

### 1. Remoção da Dependência de Planos de Assinatura

**Antes:**
- Cotas variavam por plano: Free (10), Monthly (100), Yearly (100), Trial (50)
- Lógica complexa baseada em `SubscriptionRepository`

**Depois:**
- Cotas fixas diárias por feature, independente do plano
- Limites diários:
  - `food_recognition`: 10 usos/dia
  - `meal_analysis`: 5 usos/dia

### 2. Integração do Timezone do Usuário

**Campo Novo no Usuário:**
```json
{
  "timezone": "America/Sao_Paulo"
}
```

**Comportamento:**
- O início do dia é calculado no timezone do usuário
- Se timezone inválido, usa UTC como fallback
- Reset diário acontece à meia-noite no timezone do usuário

### 3. Mudanças na API

#### Endpoints Afetados

##### `GET /api/v1/ai/usage`
Retorna lista de uso de IA do usuário.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "feature": "food_recognition",
      "used": 3,
      "quota": 10,
      "remaining": 7
    },
    {
      "feature": "meal_analysis",
      "used": 1,
      "quota": 5,
      "remaining": 4
    }
  ]
}
```

##### `GET /api/v1/ai/usage/:feature`
Verifica quota para um feature específico.

**Parâmetros:**
- `feature`: `food_recognition` ou `meal_analysis`

**Response:**
```json
{
  "success": true,
  "data": {
    "has_quota": true,
    "used": 3,
    "quota": 10,
    "remaining": 7
  }
}
```

### 4. Mudanças no Perfil do Usuário

#### `PUT /api/v1/user/profile`
Agora aceita campo `timezone` no request body.

**Request Body:**
```json
{
  "id": "uuid-do-usuario",
  "displayName": "Nome do Usuário",
  "email": "usuario@email.com",
  "language": "pt-BR",
  "timezone": "America/Sao_Paulo",
  "weight": 70,
  "height": 175,
  "age": 30,
  "gender": "male",
  "activityLevel": "moderate",
  "notificationsEnabled": true,
  "dailyCalorieGoal": 2500,
  "dailyProteinGoal": 150,
  "dailyCarbsGoal": 300,
  "dailyFatGoal": 80
}
```

#### `GET /api/v1/user/profile`
Agora retorna `timezone` no response.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid-do-usuario",
    "displayName": "Nome do Usuário",
    "email": "usuario@email.com",
    "photo": "https://cdn.exemplo.com/avatar.jpg",
    "language": "pt-BR",
    "timezone": "America/Sao_Paulo",
    "notificationsEnabled": true,
    "dailyCalorieGoal": 2500,
    "dailyProteinGoal": 150,
    "dailyCarbsGoal": 300,
    "dailyFatGoal": 80,
    "weight": 70,
    "height": 175,
    "age": 30,
    "gender": "male",
    "activityLevel": "moderate",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

## Impactos no Frontend

### 1. Atualização do Perfil do Usuário

**Necessário:**
- Incluir `timezone` no request de atualização de perfil

**Timezones Suportados:**
- Qualquer timezone válido do IANA (ex: "America/Sao_Paulo", "Europe/London", "Asia/Tokyo")
- Valor padrão: "UTC"

### 2. Tratamento de Cotas de IA

**Antes:**
- Cotas baseadas em plano de assinatura
- Lógica complexa de verificação

**Depois:**
- Cotas diárias simples por feature
- Reset automático à meia-noite no timezone do usuário
- Mensagens de erro mais claras em português

### 3. Tratamento de Erros

**Ações recomendadas no frontend:**
- Desabilitar botões de features de IA quando quota excedida
- Mostrar modal/toast explicativo sobre limites diários
- Sugerir aguardar até o próximo reset

## Considerações Técnicas

### Timezone Handling
- Frontend deve enviar timezone válido do IANA
- Backend valida e usa UTC como fallback
- Reset acontece exatamente à meia-noite no timezone do usuário

### Backward Compatibility
- Endpoints existentes mantêm compatibilidade
- Campos novos são opcionais (timezone defaults to "UTC")
- Assinaturas continuam existindo mas não afetam cotas de IA

### Performance
- Verificação de quota é feita por usuário/feature/dia
- Cache pode ser implementado no frontend para reduzir requests
- Reset automático não requer ação do usuário

## Testes Recomendados

1. **Testar diferentes timezones**
2. **Verificar reset diário correto**
3. **Testar limite de cotas**
4. **Validar mensagens de erro**
5. **Testar integração com perfil do usuário**
