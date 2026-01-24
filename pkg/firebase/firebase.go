package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

func New(ctx context.Context, cfg *domain.Config) (*firebase.App, error) {
	if cfg.FirebaseCredentialsFile == "" && cfg.FirebaseCredentialsJSON == "" {
		log.Info().Msg("Firebase credentials not provided, skipping Firebase initialization")
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

	log.Info().Msg("Firebase initialized")
	return app, nil
}

func GetAuthClient(ctx context.Context, app *firebase.App) (*auth.Client, error) {
	return app.Auth(ctx)
}
