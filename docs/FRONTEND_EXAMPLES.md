# Exemplos de Integração Frontend

## 1. Reconhecimento de Alimentos com Streaming

### Requisição
```typescript
async function recognizeFood(
  imageFile: File,
  mealData: {
    name: string;
    type: string;
    mealLocation: string;
  }
) {
  const formData = new FormData();
  formData.append('image', imageFile);
  formData.append('name', mealData.name);
  formData.append('type', mealData.type);
  formData.append('mealLocation', mealData.mealLocation);
  formData.append('uri', imageFile.name); // URI opcional

  const response = await fetch('http://localhost:8080/api/v1/food/recognize', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${getToken()}`
    },
    body: formData
  });

  if (!response.ok) {
    throw new Error('Failed to recognize food');
  }

  return response;
}
```

### Consumindo o Streaming SSE
```typescript
async function recognizeFoodWithProgress(
  imageFile: File,
  mealData: {
    name: string;
    type: string;
    mealLocation: string;
  },
  onProgress: (progress: ProgressUpdate) => void
): Promise<FoodRecognitionResponse> {
  const response = await recognizeFood(imageFile, mealData);
  
  const reader = response.body!.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  return new Promise((resolve, reject) => {
    const processText = async () => {
      while (true) {
        const { done, value } = await reader.read();
        
        if (done) break;
        
        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        
        buffer = lines.pop() || '';
        
        for (const line of lines) {
          if (line.startsWith('event:')) {
            const eventType = line.slice(6).trim();
            continue;
          }
          
          if (line.startsWith('data:')) {
            const data = line.slice(5).trim();
            
            try {
              const parsed = JSON.parse(data);
              
              if (parsed.success === false) {
                reject(new Error(parsed.message));
                return;
              }
              
              if (parsed.stage) {
                onProgress(parsed as ProgressUpdate);
              } else if (parsed.data) {
                resolve(parsed.data);
                return;
              }
            } catch (e) {
              console.error('Failed to parse SSE data:', e);
            }
          }
        }
      }
    };
    
    processText().catch(reject);
  });
}
```

### Tipos TypeScript
```typescript
interface ProgressUpdate {
  stage: string;
  percentage: number;
  message: string;
}

interface RecognizedFoodItem {
  name: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  quantity: number;
  unit: string;
  confidence: number;
}

interface FoodRecognitionResponse {
  foodItems: RecognizedFoodItem[];
  processingTime: number;
}
```

### Exemplo de Uso React
```typescript
import { useState } from 'react';

function FoodRecognition() {
  const [progress, setProgress] = useState<ProgressUpdate | null>(null);
  const [result, setResult] = useState<FoodRecognitionResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const handleRecognize = async (file: File) => {
    setLoading(true);
    setProgress(null);
    setResult(null);

    try {
      const data = await recognizeFoodWithProgress(
        file,
        {
          name: 'Foto do prato',
          type: 'lunch',
          mealLocation: 'Casa'
        },
        (progressUpdate) => {
          setProgress(progressUpdate);
          console.log(`${progressUpdate.stage}: ${progressUpdate.percentage}%`);
        }
      );

      setResult(data);
    } catch (error) {
      console.error('Recognition failed:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        accept="image/*"
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) handleRecognize(file);
        }}
      />

      {progress && (
        <div>
          <p>{progress.message}</p>
          <progress value={progress.percentage} max={100} />
        </div>
      )}

      {result && (
        <div>
          <h3>Alimentos Reconhecidos ({result.processingTime}ms)</h3>
          {result.foodItems.map((item, i) => (
            <div key={i}>
              <h4>{item.name} ({item.confidence * 100}% confiança)</h4>
              <p>
                Quantidade: {item.quantity} {item.unit} | 
                Calorias: {item.calories}kcal | 
                P: {item.protein}g | C: {item.carbs}g | G: {item.fat}g
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

## 2. Estimativa de Quantidade com Streaming

### Requisição
```typescript
async function estimateQuantity(
  imageFile: File,
  foodData: {
    name: string;
    type: string;
    mealLocation: string;
    referenceServingSize?: string;
    referenceServingUnit?: string;
  },
  onProgress: (progress: ProgressUpdate) => void
): Promise<EstimateQuantityResponse> {
  const formData = new FormData();
  formData.append('image', imageFile);
  formData.append('name', foodData.name);
  formData.append('type', foodData.type);
  formData.append('mealLocation', foodData.mealLocation);
  formData.append('uri', imageFile.name);
  
  if (foodData.referenceServingSize) {
    formData.append('referenceServingSize', foodData.referenceServingSize);
  }
  if (foodData.referenceServingUnit) {
    formData.append('referenceServingUnit', foodData.referenceServingUnit);
  }

  const response = await fetch('http://localhost:8080/api/v1/food/estimate-quantity', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${getToken()}`
    },
    body: formData
  });

  // Processar SSE igual ao exemplo anterior
  // ...
}
```

### Tipos
```typescript
interface EstimateQuantityResponse {
  estimatedQuantity: number;
  unit: string;
  confidence: number;
  reasoning?: string;
}
```

## 3. Busca por Código de Barras (Sem Streaming)

### Requisição
```typescript
async function getFoodByBarcode(barcode: string): Promise<FoodBarcodeResponse> {
  const response = await fetch(`http://localhost:8080/api/v1/food/barcode/${barcode}`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`
    }
  });

  const data = await response.json();
  
  if (!data.success) {
    throw new Error(data.message);
  }

  return data.data;
}
```

### Tipos
```typescript
interface FoodBarcodeResponse {
  barcode: string;
  name: string;
  brand?: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  servingSize?: number;
  servingUnit?: string;
  source?: string;
}
```

### Exemplo de Uso
```typescript
const barcodeData = await getFoodByBarcode('7891234567890');
console.log(`${barcodeData.name} - ${barcodeData.calories}kcal`);
```

## Eventos SSE do Backend

O backend envia eventos Server-Sent Events no seguinte formato:

### Evento de Progresso
```
event: progress
data: {"stage":"upload","percentage":20,"message":"Uploading to S3..."}
```

### Evento de Sucesso
```
event: complete
data: {"success":true,"data":{...},"message":"food recognized successfully"}
```

### Evento de Erro
```
event: error
data: {"success":false,"message":"failed to recognize food"}
```

## Estágios de Progresso

- `upload` (0-25%): Upload e processamento da imagem
- `ai_analysis` (25-90%): Análise pela IA
- `complete` (100%): Processamento finalizado
