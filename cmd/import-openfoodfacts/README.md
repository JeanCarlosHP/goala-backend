# Open Food Facts Importer

Importador em Go para carregar o dump `en.openfoodfacts.org.products.csv.gz` na base local do GOALA.

O comando lê o arquivo em streaming, filtra produtos sem dados úteis, faz `COPY` para tabelas temporárias, faz merge por lote em `food_database` e `food_portions`, e pode indexar automaticamente no Meilisearch.

## O que é importado

- `food_database`
  - `external_id` <- `code`
  - `barcode` <- `code`
  - `name` <- `product_name`
  - `brand` <- primeira marca em `brands`
  - `calories_per_100g` <- `energy-kcal_100g`
  - `protein_per_100g` <- `proteins_100g`
  - `carbs_per_100g` <- `carbohydrates_100g`
  - `fat_per_100g` <- `fat_100g`
  - `source` <- `openfoodfacts`

- `food_portions`
  - cria a porção `serving` quando `serving_size` é parseável, por exemplo `30 g` ou `250 ml`

## Pré-requisitos

- migrations da API já aplicadas
- PostgreSQL acessível via `DATABASE_URL`
- opcional: Meilisearch acessível via `MEILISEARCH_URL`

## Uso direto

```bash
cd backend

GOCACHE=$(pwd)/.gocache go run ./cmd/import-openfoodfacts \
  -file ../en.openfoodfacts.org.products.csv.gz \
  -database-url "$DATABASE_URL"
```

## Dry-run

Valida parsing e filtros sem escrever no banco:

```bash
cd backend

GOCACHE=$(pwd)/.gocache go run ./cmd/import-openfoodfacts \
  -dry-run \
  -file ../en.openfoodfacts.org.products.csv.gz \
  -limit 5000 \
  -batch-size 500 \
  -progress-every 1000
```

## Com Meilisearch

```bash
cd backend

GOCACHE=$(pwd)/.gocache go run ./cmd/import-openfoodfacts \
  -file ../en.openfoodfacts.org.products.csv.gz \
  -database-url "$DATABASE_URL" \
  -index-meili \
  -meili-url "$MEILISEARCH_URL" \
  -meili-api-key "$MEILISEARCH_API_KEY" \
  -meili-index "${MEILISEARCH_FOODS_INDEX:-foods}"
```

## Via Makefile

```bash
cd backend

make import-openfoodfacts DATABASE_URL='postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable'
```

Exemplos:

```bash
make import-openfoodfacts DRY_RUN=1 LIMIT=5000 BATCH_SIZE=500 PROGRESS_EVERY=1000

make import-openfoodfacts \
  DATABASE_URL='postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable' \
  TRUNCATE=1

make import-openfoodfacts \
  DATABASE_URL='postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable' \
  INDEX_MEILI=1 \
  MEILISEARCH_URL='http://localhost:7700' \
  MEILISEARCH_API_KEY='masterKey' \
  MEILISEARCH_FOODS_INDEX='foods'
```

## Flags

- `-file`: caminho do `.csv.gz`
- `-database-url`: conexão PostgreSQL
- `-batch-size`: tamanho do lote para `COPY` + merge
- `-progress-every`: imprime progresso a cada N linhas lidas
- `-limit`: limita quantas linhas do CSV serão lidas
- `-truncate`: remove dados anteriores de `openfoodfacts` antes do import
- `-dry-run`: não escreve no banco
- `-index-meili`: envia documentos para o Meilisearch após cada lote
- `-meili-url`: URL base do Meilisearch
- `-meili-api-key`: chave da API do Meilisearch
- `-meili-index`: nome do índice de alimentos

## Comportamento em produção

- o import é incremental por lote
- cada lote faz `COPY` para staging temporário, merge no catálogo, indexação opcional e limpeza do staging
- isso reduz uso de memória e permite acompanhar progresso real no log
- `-truncate` deve ser usado com cuidado em produção, apenas quando você quiser reconstruir completamente o catálogo Open Food Facts
