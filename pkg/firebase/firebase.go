package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

func New(ctx context.Context, cfg *domain.Config, logger domain.Logger) (*firebase.App, error) {
	if cfg.FirebaseCredentialsFile == "" && cfg.FirebaseCredentialsJSON == "" {
		logger.Info("Firebase credentials not provided, skipping Firebase initialization")
		return nil, nil
	}

	var opt option.ClientOption
	if cfg.FirebaseCredentialsFile != "" {
		opt = option.WithAuthCredentialsFile("service_account", cfg.FirebaseCredentialsFile)
	} else {
		opt = option.WithAuthCredentialsJSON("service_account", []byte(cfg.FirebaseCredentialsJSON))
	}

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	logger.Info("Firebase app initialized successfully")
	return app, nil
}

func GetAuthClient(ctx context.Context, app *firebase.App) (*auth.Client, error) {
	return app.Auth(ctx)
}
