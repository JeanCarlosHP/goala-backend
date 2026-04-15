package domain

type Config struct {
	OtelCollectorURL                   string   `mapstructure:"OTEL_COLLECTOR_URL"`
	HTTPPort                           string   `mapstructure:"HTTP_PORT"`
	HTTPRequestSizeLimit               string   `mapstructure:"HTTP_REQUEST_SIZE_LIMIT"`
	HTTPCorsAllowedOrigins             []string `mapstructure:"HTTP_CORS_ALLOWED_ORIGINS"`
	HTTPCorsAllowedHeaders             []string `mapstructure:"HTTP_CORS_ALLOWED_HEADERS"`
	HTTPCorsAllowedMethods             []string `mapstructure:"HTTP_CORS_ALLOWED_METHODS"`
	Environment                        string   `mapstructure:"ENVIRONMENT"`
	LoggingLevel                       string   `mapstructure:"LOGGING_LEVEL"`
	DatabaseMigrationPath              string   `mapstructure:"DATABASE_MIGRATION_PATH"`
	DatabaseURL                        string   `mapstructure:"DATABASE_URL"`
	DatabaseMigrationURL               string   `mapstructure:"DATABASE_MIGRATION_URL"`
	PaginationMaxPerPage               int      `mapstructure:"PAGINATION_MAX_PER_PAGE"`
	PaginationDefaultPageSize          int      `mapstructure:"PAGINATION_DEFAULT_PAGE_SIZE"`
	DatabaseTracing                    bool     `mapstructure:"DATABASE_TRACING"`
	DatabaseShouldForceSetLowerVersion bool     `mapstructure:"DATABASE_MIGRATION_FORCE_SET_LOWER_VERSION"`
	LoggingJSONFormat                  bool     `mapstructure:"LOGGING_JSON_FORMAT"`
	FirebaseCredentialsFile            string   `mapstructure:"FIREBASE_CREDENTIALS_FILE"`
	FirebaseCredentialsJSON            string   `mapstructure:"FIREBASE_CREDENTIALS_JSON"`
	StorageBucketName                  string   `mapstructure:"BUCKET"`
	StorageRegion                      string   `mapstructure:"REGION"`
	StorageEndpoint                    string   `mapstructure:"ENDPOINT"`
	StorageAccessKeyID                 string   `mapstructure:"ACCESS_KEY_ID"`
	StorageSecretAccessKey             string   `mapstructure:"SECRET_ACCESS_KEY"`
	GeminiAPIKey                       string   `mapstructure:"GEMINI_API_KEY"`
	GeminiModel                        string   `mapstructure:"GEMINI_MODEL"`
	OpenAIAPIKey                       string   `mapstructure:"OPENAI_API_KEY"`
	OpenAIModel                        string   `mapstructure:"OPENAI_MODEL"`
	AIProvider                         string   `mapstructure:"AI_PROVIDER"`
	OpenFoodFactsAPIURL                string   `mapstructure:"OPENFOODFACTS_API_URL"`
	RedisURL                           string   `mapstructure:"REDIS_URL"`
	MeiliSearchURL                     string   `mapstructure:"MEILISEARCH_URL"`
	MeiliSearchAPIKey                  string   `mapstructure:"MEILISEARCH_API_KEY"`
	MeiliSearchFoodsIndex              string   `mapstructure:"MEILISEARCH_FOODS_INDEX"`
	CDNDomain                          string   `mapstructure:"CDN_DOMAIN"`
	RevenueCatWebhookSecret            string   `mapstructure:"REVENUECAT_WEBHOOK_SECRET"`
}
