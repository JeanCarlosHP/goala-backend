package domain

type Config struct {
	DatabaseName                       string `mapstructure:"DATABASE_NAME"`
	DatabaseHost                       string `mapstructure:"DATABASE_HOST"`
	DatabaseUser                       string `mapstructure:"DATABASE_USER"`
	DatabasePassword                   string `mapstructure:"DATABASE_PASSWORD"`
	DatabasePort                       string `mapstructure:"DATABASE_PORT"`
	HTTPPort                           string `mapstructure:"HTTP_PORT"`
	HTTPRequestSizeLimit               string `mapstructure:"HTTP_REQUEST_SIZE_LIMIT"`
	HTTPCorsAllowedOrigins             string `mapstructure:"HTTP_CORS_ALLOWED_ORIGINS"`
	HTTPCorsAllowedHeaders             string `mapstructure:"HTTP_CORS_ALLOWED_HEADERS"`
	HTTPCorsAllowedMethods             string `mapstructure:"HTTP_CORS_ALLOWED_METHODS"`
	Environment                        string `mapstructure:"ENVIRONMENT"`
	LoggingLevel                       string `mapstructure:"LOGGING_LEVEL"`
	DatabaseMigrationPath              string `mapstructure:"DATABASE_MIGRATION_PATH"`
	PaginationMaxPerPage               int    `mapstructure:"PAGINATION_MAX_PER_PAGE"`
	PaginationDefaultPageSize          int    `mapstructure:"PAGINATION_DEFAULT_PAGE_SIZE"`
	DatabaseTracing                    bool   `mapstructure:"DATABASE_TRACING"`
	DatabaseShouldForceSetLowerVersion bool   `mapstructure:"DATABASE_MIGRATION_FORCE_SET_LOWER_VERSION"`
	DatabaseSslMode                    string `mapstructure:"DATABASE_SSL_MODE"`
	LoggingJSONFormat                  bool   `mapstructure:"LOGGING_JSON_FORMAT"`
	FirebaseCredentialsFile            string `mapstructure:"FIREBASE_CREDENTIALS_FILE"`
	FirebaseCredentialsJSON            string `mapstructure:"FIREBASE_CREDENTIALS_JSON"`
	AWSS3BucketName                    string `mapstructure:"AWS_S3_BUCKET_NAME"`
	AWSS3Region                        string `mapstructure:"AWS_S3_REGION"`
	AWSAccessKeyID                     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey                 string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	GeminiAPIKey                       string `mapstructure:"GEMINI_API_KEY"`
	OpenAIAPIKey                       string `mapstructure:"OPENAI_API_KEY"`
	AIProvider                         string `mapstructure:"AI_PROVIDER"`
	OpenFoodFactsAPIURL                string `mapstructure:"OPENFOODFACTS_API_URL"`
	CDNDomain                          string `mapstructure:"CDN_DOMAIN"`
	RevenueCatWebhookSecret            string `mapstructure:"REVENUECAT_WEBHOOK_SECRET"`
}
