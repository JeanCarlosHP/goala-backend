package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	filePath      string
	databaseURL   string
	batchSize     int
	limit         int
	truncate      bool
	dryRun        bool
	progressEvery int
	meiliURL      string
	meiliAPIKey   string
	meiliIndex    string
	indexMeili    bool
	meiliTimeout  time.Duration
	meiliRetries  int
	meiliBackoff  time.Duration
}

type foodRow struct {
	ExternalID      string
	Barcode         string
	Name            string
	Brand           *string
	CaloriesPer100g int
	ProteinPer100g  float64
	CarbsPer100g    float64
	FatPer100g      float64
}

type portionRow struct {
	ExternalID string
	Name       string
	Grams      float64
}

type stats struct {
	ReadRows         int
	ImportedFoods    int64
	ImportedPortions int64
	IndexedFoods     int
	SkippedEmpty     int
	SkippedInvalid   int
	Batches          int
}

type meiliDocument struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Brand    *string `json:"brand,omitempty"`
	Source   string  `json:"source"`
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}

type meiliIndexer struct {
	baseURL   string
	apiKey    string
	indexName string
	client    *http.Client
	timeout   time.Duration
	retries   int
	backoff   time.Duration
}

var servingPattern = regexp.MustCompile(`(?i)^\s*([0-9]+(?:[.,][0-9]+)?)\s*(g|ml)\s*$`)

func main() {
	cfg := parseFlags()
	if cfg.databaseURL == "" && !cfg.dryRun {
		fmt.Fprintln(os.Stderr, "DATABASE_URL or -database-url is required")
		os.Exit(1)
	}

	ctx := context.Background()
	start := time.Now()
	indexer := newMeiliIndexer(cfg)

	file, err := os.Open(cfg.filePath)
	if err != nil {
		fail("open input file", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		fail("open gzip reader", err)
	}
	defer gzReader.Close()

	reader := csv.NewReader(gzReader)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.ReuseRecord = true

	headers, err := reader.Read()
	if err != nil {
		fail("read header", err)
	}

	columnIndex := indexColumns(headers)
	required := []string{
		"code", "product_name", "brands", "energy-kcal_100g",
		"proteins_100g", "carbohydrates_100g", "fat_100g", "serving_size",
	}
	for _, column := range required {
		if _, ok := columnIndex[column]; !ok {
			fail("validate header", fmt.Errorf("missing required column %q", column))
		}
	}

	var pool *pgxpool.Pool
	var conn *pgxpool.Conn
	if !cfg.dryRun {
		pool, err = pgxpool.New(ctx, cfg.databaseURL)
		if err != nil {
			fail("connect database", err)
		}
		defer pool.Close()

		if err := ensureSchema(ctx, pool); err != nil {
			fail("ensure schema", err)
		}

		if cfg.truncate {
			if err := truncateImportedData(ctx, pool); err != nil {
				fail("truncate imported data", err)
			}
		}

		conn, err = pool.Acquire(ctx)
		if err != nil {
			fail("acquire connection", err)
		}
		defer conn.Release()

		if err := createStagingTables(ctx, conn.Conn()); err != nil {
			fail("create staging tables", err)
		}
	}

	foodBatch := make([]foodRow, 0, cfg.batchSize)
	portionBatch := make([]portionRow, 0, cfg.batchSize)
	runStats := stats{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			runStats.SkippedInvalid++
			continue
		}
		runStats.ReadRows++

		food, portion, ok := parseRecord(record, columnIndex)
		if !ok {
			runStats.SkippedEmpty++
			if cfg.limit > 0 && runStats.ReadRows >= cfg.limit {
				break
			}
			continue
		}

		foodBatch = append(foodBatch, food)
		if portion != nil {
			portionBatch = append(portionBatch, *portion)
		}

		if len(foodBatch) >= cfg.batchSize {
			if err := processBatch(ctx, cfg, dbConn(conn), indexer, foodBatch, portionBatch, &runStats, start); err != nil {
				fail("process batch", err)
			}
			foodBatch = foodBatch[:0]
			portionBatch = portionBatch[:0]
		}

		if cfg.progressEvery > 0 && runStats.ReadRows%cfg.progressEvery == 0 {
			printProgress(runStats, start)
		}

		if cfg.limit > 0 && runStats.ReadRows >= cfg.limit {
			break
		}
	}

	if len(foodBatch) > 0 {
		if err := processBatch(ctx, cfg, dbConn(conn), indexer, foodBatch, portionBatch, &runStats, start); err != nil {
			fail("process final batch", err)
		}
	}

	fmt.Printf("Arquivo analisado: %s\n", cfg.filePath)
	fmt.Printf("Colunas: %d\n", len(headers))
	fmt.Printf("Linhas lidas: %d\n", runStats.ReadRows)
	fmt.Printf("Linhas ignoradas por falta de dados úteis: %d\n", runStats.SkippedEmpty)
	fmt.Printf("Linhas inválidas: %d\n", runStats.SkippedInvalid)
	if cfg.dryRun {
		fmt.Println("Modo dry-run: nenhum dado foi inserido.")
	} else {
		fmt.Printf("Foods importados/atualizados: %d\n", runStats.ImportedFoods)
		fmt.Printf("Porções importadas/atualizadas: %d\n", runStats.ImportedPortions)
		if cfg.indexMeili {
			fmt.Printf("Foods indexados no Meilisearch: %d\n", runStats.IndexedFoods)
		}
	}
	fmt.Printf("Lotes processados: %d\n", runStats.Batches)
	fmt.Printf("Tempo total: %s\n", time.Since(start).Round(time.Millisecond))
}

func parseFlags() config {
	cfg := config{}
	flag.StringVar(&cfg.filePath, "file", "../en.openfoodfacts.org.products.csv.gz", "Path to en.openfoodfacts.org.products.csv.gz")
	flag.StringVar(&cfg.databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	flag.IntVar(&cfg.batchSize, "batch-size", 2000, "Rows per in-memory batch before flushing to staging")
	flag.IntVar(&cfg.limit, "limit", 0, "Maximum number of CSV rows to read (0 = all)")
	flag.IntVar(&cfg.progressEvery, "progress-every", 10000, "Print progress every N CSV rows read")
	flag.BoolVar(&cfg.truncate, "truncate", false, "Delete previously imported Open Food Facts data before importing")
	flag.BoolVar(&cfg.dryRun, "dry-run", false, "Analyze and parse the file without writing to the database")
	flag.StringVar(&cfg.meiliURL, "meili-url", os.Getenv("MEILISEARCH_URL"), "Meilisearch base URL")
	flag.StringVar(&cfg.meiliAPIKey, "meili-api-key", os.Getenv("MEILISEARCH_API_KEY"), "Meilisearch API key")
	flag.StringVar(&cfg.meiliIndex, "meili-index", envOrDefault("MEILISEARCH_FOODS_INDEX", "foods"), "Meilisearch foods index name")
	flag.BoolVar(&cfg.indexMeili, "index-meili", false, "Index imported foods in Meilisearch after each batch")
	flag.DurationVar(&cfg.meiliTimeout, "meili-timeout", envDurationOrDefault("MEILISEARCH_TIMEOUT", 2*time.Minute), "Timeout per Meilisearch indexing request")
	flag.IntVar(&cfg.meiliRetries, "meili-retries", envIntOrDefault("MEILISEARCH_RETRIES", 3), "Number of retries for Meilisearch indexing")
	flag.DurationVar(&cfg.meiliBackoff, "meili-backoff", envDurationOrDefault("MEILISEARCH_RETRY_BACKOFF", 3*time.Second), "Backoff between Meilisearch retries")
	flag.Parse()
	return cfg
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		ALTER TABLE food_database
			ADD COLUMN IF NOT EXISTS external_id VARCHAR(100);
		ALTER TABLE food_items
			ADD COLUMN IF NOT EXISTS food_database_id UUID REFERENCES food_database(id) ON DELETE SET NULL;
		CREATE TABLE IF NOT EXISTS food_portions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			food_id UUID NOT NULL REFERENCES food_database(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			grams DECIMAL(10,2) NOT NULL CHECK (grams > 0),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE (food_id, name)
		);
		CREATE TABLE IF NOT EXISTS favorite_foods (
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			food_id UUID NOT NULL REFERENCES food_database(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, food_id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_food_database_external_id_unique
			ON food_database(external_id)
			WHERE external_id IS NOT NULL;
	`)
	return err
}

func truncateImportedData(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		DELETE FROM favorite_foods
		WHERE food_id IN (
			SELECT id FROM food_database WHERE source = 'openfoodfacts'
		);
		DELETE FROM food_portions
		WHERE food_id IN (
			SELECT id FROM food_database WHERE source = 'openfoodfacts'
		);
		DELETE FROM food_database
		WHERE source = 'openfoodfacts';
	`)
	return err
}

func createStagingTables(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS off_import_foods (
			external_id TEXT,
			barcode TEXT,
			name TEXT,
			brand TEXT,
			calories_per_100g INT,
			protein_per_100g DOUBLE PRECISION,
			carbs_per_100g DOUBLE PRECISION,
			fat_per_100g DOUBLE PRECISION
		);

		CREATE TEMP TABLE IF NOT EXISTS off_import_portions (
			external_id TEXT,
			name TEXT,
			grams DOUBLE PRECISION
		);
	`)
	return err
}

func flushBatch(ctx context.Context, conn *pgx.Conn, foods []foodRow, portions []portionRow, dryRun bool) error {
	if dryRun {
		return nil
	}

	if len(foods) > 0 {
		rows := make([][]any, 0, len(foods))
		for _, food := range foods {
			rows = append(rows, []any{
				food.ExternalID,
				food.Barcode,
				food.Name,
				stringOrNil(food.Brand),
				food.CaloriesPer100g,
				food.ProteinPer100g,
				food.CarbsPer100g,
				food.FatPer100g,
			})
		}

		_, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{"off_import_foods"},
			[]string{"external_id", "barcode", "name", "brand", "calories_per_100g", "protein_per_100g", "carbs_per_100g", "fat_per_100g"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return err
		}
	}

	if len(portions) > 0 {
		rows := make([][]any, 0, len(portions))
		for _, portion := range portions {
			rows = append(rows, []any{portion.ExternalID, portion.Name, portion.Grams})
		}

		_, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{"off_import_portions"},
			[]string{"external_id", "name", "grams"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeStagingTables(ctx context.Context, conn *pgx.Conn) (int64, int64, []meiliDocument, error) {
	commandTag, err := conn.Exec(ctx, `
		INSERT INTO food_database (
			external_id,
			barcode,
			name,
			brand,
			calories_per_100g,
			protein_per_100g,
			carbs_per_100g,
			fat_per_100g,
			source,
			verified,
			updated_at
		)
		SELECT DISTINCT ON (external_id)
			external_id,
			barcode,
			name,
			NULLIF(brand, ''),
			calories_per_100g,
			protein_per_100g,
			carbs_per_100g,
			fat_per_100g,
			'openfoodfacts',
			false,
			NOW()
		FROM off_import_foods
		WHERE external_id IS NOT NULL
		ON CONFLICT (external_id) WHERE external_id IS NOT NULL
		DO UPDATE SET
			barcode = EXCLUDED.barcode,
			name = EXCLUDED.name,
			brand = EXCLUDED.brand,
			calories_per_100g = EXCLUDED.calories_per_100g,
			protein_per_100g = EXCLUDED.protein_per_100g,
			carbs_per_100g = EXCLUDED.carbs_per_100g,
			fat_per_100g = EXCLUDED.fat_per_100g,
			source = EXCLUDED.source,
			updated_at = NOW()
	`)
	if err != nil {
		return 0, 0, nil, err
	}

	portionTag, err := conn.Exec(ctx, `
		INSERT INTO food_portions (food_id, name, grams)
		SELECT
			fd.id,
			ip.name,
			ip.grams
		FROM off_import_portions ip
		JOIN food_database fd ON fd.external_id = ip.external_id
		ON CONFLICT (food_id, name)
		DO UPDATE SET grams = EXCLUDED.grams
	`)
	if err != nil {
		return 0, 0, nil, err
	}

	documents, err := fetchMeiliDocuments(ctx, conn)
	if err != nil {
		return 0, 0, nil, err
	}

	if _, err := conn.Exec(ctx, `TRUNCATE TABLE off_import_foods, off_import_portions`); err != nil {
		return 0, 0, nil, err
	}

	return commandTag.RowsAffected(), portionTag.RowsAffected(), documents, nil
}

func parseRecord(record []string, index map[string]int) (foodRow, *portionRow, bool) {
	name := strings.TrimSpace(valueAt(record, index, "product_name"))
	if name == "" {
		return foodRow{}, nil, false
	}

	code := strings.TrimSpace(valueAt(record, index, "code"))
	if code == "" {
		return foodRow{}, nil, false
	}

	calories, okCalories := parseInt(valueAt(record, index, "energy-kcal_100g"))
	protein, okProtein := parseFloat(valueAt(record, index, "proteins_100g"))
	carbs, okCarbs := parseFloat(valueAt(record, index, "carbohydrates_100g"))
	fat, okFat := parseFloat(valueAt(record, index, "fat_100g"))

	if !okCalories && !okProtein && !okCarbs && !okFat {
		return foodRow{}, nil, false
	}

	food := foodRow{
		ExternalID:      code,
		Barcode:         code,
		Name:            name,
		Brand:           normalizedOptional(valueAt(record, index, "brands")),
		CaloriesPer100g: calories,
		ProteinPer100g:  protein,
		CarbsPer100g:    carbs,
		FatPer100g:      fat,
	}

	portion := parseServingSize(code, valueAt(record, index, "serving_size"))
	return food, portion, true
}

func parseServingSize(code, raw string) *portionRow {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return nil
	}

	matches := servingPattern.FindStringSubmatch(raw)
	if len(matches) != 3 {
		return nil
	}

	grams, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ",", "."), 64)
	if err != nil || grams <= 0 {
		return nil
	}

	return &portionRow{
		ExternalID: code,
		Name:       "serving",
		Grams:      round(grams, 2),
	}
}

func parseFloat(value string) (float64, bool) {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", "."))
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return 0, false
	}
	return round(parsed, 2), true
}

func parseInt(value string) (int, bool) {
	parsed, ok := parseFloat(value)
	if !ok {
		return 0, false
	}
	return int(math.Round(parsed)), true
}

func normalizedOptional(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	first := strings.TrimSpace(parts[0])
	if first == "" {
		return nil
	}
	return &first
}

func valueAt(record []string, index map[string]int, column string) string {
	position, ok := index[column]
	if !ok || position >= len(record) {
		return ""
	}
	return record[position]
}

func indexColumns(headers []string) map[string]int {
	index := make(map[string]int, len(headers))
	for i, header := range headers {
		index[header] = i
	}
	return index
}

func round(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}

func stringOrNil(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func fail(action string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", action, err)
	os.Exit(1)
}

func dbConn(conn *pgxpool.Conn) *pgx.Conn {
	if conn == nil {
		return nil
	}
	return conn.Conn()
}

func processBatch(
	ctx context.Context,
	cfg config,
	conn *pgx.Conn,
	indexer *meiliIndexer,
	foods []foodRow,
	portions []portionRow,
	runStats *stats,
	start time.Time,
) error {
	if err := flushBatch(ctx, conn, foods, portions, cfg.dryRun); err != nil {
		return err
	}

	runStats.Batches++
	if cfg.dryRun {
		fmt.Printf(
			"[batch %d] read=%d valid=%d skipped=%d elapsed=%s\n",
			runStats.Batches,
			runStats.ReadRows,
			len(foods),
			runStats.SkippedEmpty+runStats.SkippedInvalid,
			time.Since(start).Round(time.Millisecond),
		)
		return nil
	}

	importedFoods, importedPortions, documents, err := mergeStagingTables(ctx, conn)
	if err != nil {
		return err
	}
	runStats.ImportedFoods += importedFoods
	runStats.ImportedPortions += importedPortions

	if cfg.indexMeili && indexer != nil && len(documents) > 0 {
		indexed, err := indexer.IndexDocuments(ctx, documents)
		if err != nil {
			return err
		}
		runStats.IndexedFoods += indexed
	}

	fmt.Printf(
		"[batch %d] read=%d imported_foods=%d imported_portions=%d indexed=%d elapsed=%s\n",
		runStats.Batches,
		runStats.ReadRows,
		importedFoods,
		importedPortions,
		runStats.IndexedFoods,
		time.Since(start).Round(time.Millisecond),
	)
	return nil
}

func fetchMeiliDocuments(ctx context.Context, conn *pgx.Conn) ([]meiliDocument, error) {
	rows, err := conn.Query(ctx, `
		SELECT DISTINCT
			fd.id::text,
			fd.name,
			fd.brand,
			COALESCE(fd.source, 'openfoodfacts'),
			COALESCE(fd.calories_per_100g, 0),
			COALESCE(fd.protein_per_100g, 0),
			COALESCE(fd.carbs_per_100g, 0),
			COALESCE(fd.fat_per_100g, 0)
		FROM off_import_foods oif
		JOIN food_database fd ON fd.external_id = oif.external_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]meiliDocument, 0)
	for rows.Next() {
		var doc meiliDocument
		if err := rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.Brand,
			&doc.Source,
			&doc.Calories,
			&doc.Protein,
			&doc.Carbs,
			&doc.Fat,
		); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func newMeiliIndexer(cfg config) *meiliIndexer {
	if !cfg.indexMeili {
		return nil
	}
	if cfg.meiliURL == "" {
		fail("init meilisearch", fmt.Errorf("MEILISEARCH_URL or -meili-url is required when -index-meili is enabled"))
	}
	return &meiliIndexer{
		baseURL:   strings.TrimRight(cfg.meiliURL, "/"),
		apiKey:    cfg.meiliAPIKey,
		indexName: cfg.meiliIndex,
		client:    &http.Client{},
		timeout:   cfg.meiliTimeout,
		retries:   cfg.meiliRetries,
		backoff:   cfg.meiliBackoff,
	}
}

func (m *meiliIndexer) IndexDocuments(ctx context.Context, documents []meiliDocument) (int, error) {
	if m == nil || len(documents) == 0 {
		return 0, nil
	}

	body, err := json.Marshal(documents)
	if err != nil {
		return 0, err
	}

	var lastErr error
	for attempt := 1; attempt <= max(1, m.retries); attempt++ {
		requestCtx := ctx
		cancel := func() {}
		if m.timeout > 0 {
			requestCtx, cancel = context.WithTimeout(ctx, m.timeout)
		}

		req, err := http.NewRequestWithContext(
			requestCtx,
			http.MethodPost,
			fmt.Sprintf("%s/indexes/%s/documents", m.baseURL, m.indexName),
			bytes.NewReader(body),
		)
		if err != nil {
			cancel()
			return 0, err
		}
		req.Header.Set("Content-Type", "application/json")
		if m.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+m.apiKey)
		}

		resp, err := m.client.Do(req)
		cancel()
		if err != nil {
			lastErr = err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode < 300 {
				return len(documents), nil
			}
			payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			lastErr = fmt.Errorf("meilisearch returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(payload)))
		}

		if attempt < max(1, m.retries) {
			time.Sleep(m.backoff)
		}
	}

	return 0, lastErr
}

func printProgress(runStats stats, start time.Time) {
	fmt.Printf(
		"[progress] read=%d batches=%d skipped_empty=%d skipped_invalid=%d elapsed=%s\n",
		runStats.ReadRows,
		runStats.Batches,
		runStats.SkippedEmpty,
		runStats.SkippedInvalid,
		time.Since(start).Round(time.Millisecond),
	)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
