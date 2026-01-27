package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/jeancarloshp/calorieai/pkg/config"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"github.com/jeancarloshp/calorieai/pkg/firebase"
	"github.com/jeancarloshp/calorieai/pkg/server"
	"github.com/jeancarloshp/calorieai/pkg/server/middleware"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/handlers"
	"github.com/jeancarloshp/calorieai/internal/observability"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"github.com/jeancarloshp/calorieai/internal/services"
	log "github.com/jeancarloshp/calorieai/pkg/logger"
)

var (
	err        error
	configurer *domain.Config
	logger     domain.Logger
)

func init() {
	configurer = config.New()

	logger = log.New(configurer)
}

func main() {
	ctx := context.Background()
	tp, err := observability.InitTracer(ctx, "calorieai-backend", configurer)
	if err != nil {
		logger.Fatal("failed to initialize tracing:", err)
	}
	defer func() { _ = tp.Shutdown(ctx) }()

	database := database.New(logger)
	if err != nil {
		logger.Fatal("failed to initialize database:", err)
	}

	err = database.NewConnection(configurer)
	if err != nil {
		logger.Fatal("failed to connect to database:", err)
	}

	firebaseApp, err := firebase.New(ctx, configurer, logger)
	if err != nil {
		logger.Fatal("failed to initialize firebase app:", err)
	}

	userRepo := repositories.NewUserRepository(database)
	goalRepo := repositories.NewGoalRepository(database)
	mealRepo := repositories.NewMealRepository(database)
	foodRepo := repositories.NewFoodRepository(database)
	statsRepo := repositories.NewStatsRepository(database.Querier.(*db.Queries))
	achievementRepo := repositories.NewAchievementRepository(database.Querier.(*db.Queries))
	feedbackRepo := repositories.NewFeedbackRepository(database)
	subscriptionRepo := repositories.NewSubscriptionRepository(database)
	aiUsageRepo := repositories.NewAIUsageRepository(database)

	userService := services.NewUserService(userRepo, goalRepo, configurer.CDNDomain)
	mealService := services.NewMealService(mealRepo, foodRepo)
	foodService := services.NewFoodService(foodRepo)
	statsService := services.NewStatsService(statsRepo, mealRepo, logger)
	achievementService := services.NewAchievementService(achievementRepo, statsRepo, logger)
	feedbackService := services.NewFeedbackService(feedbackRepo, logger)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, logger)
	aiUsageService := services.NewAIUsageService(aiUsageRepo, subscriptionRepo, logger)
	revenueCatService := services.NewRevenueCatService(configurer.RevenueCatWebhookSecret, logger)

	s3Service, err := services.NewS3Service(configurer, logger)
	if err != nil {
		logger.Fatal("failed to initialize S3 service:", err)
	}

	foodRecognitionService := services.NewFoodRecognitionService(s3Service, configurer, logger)
	barcodeService := services.NewBarcodeService(database.Querier.(*db.Queries), configurer, logger)

	authHandler := handlers.NewAuthHandler(userService, firebaseApp, logger)
	userHandler := handlers.NewUserHandler(userService, s3Service, logger)
	mealHandler := handlers.NewMealHandler(mealService, userService, logger)
	foodHandler := handlers.NewFoodHandler(foodService, validator.New(), logger)
	statsHandler := handlers.NewStatsHandler(statsService, logger)
	achievementHandler := handlers.NewAchievementHandler(achievementService, logger)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService, userService, logger)
	foodRecognitionHandler := handlers.NewFoodRecognitionHandler(foodRecognitionService, barcodeService, aiUsageService, logger)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, revenueCatService, logger)
	aiUsageHandler := handlers.NewAIUsageHandler(aiUsageService, logger)

	httpServer := server.New(configurer, logger)
	app := httpServer.GetApp()

	api := app.Group("/api/v1")

	api.Post("/auth/register", authHandler.Register)
	api.Post("/webhooks/revenuecat", subscriptionHandler.HandleWebhook)

	protected := api.Group("", middleware.AuthRequired(firebaseApp, logger), middleware.UserContext(userRepo, logger))

	protected.Get("/auth/me", authHandler.GetMe)
	protected.Put("/user/goals", authHandler.UpdateGoals)

	protected.Get("/subscription/status", subscriptionHandler.GetStatus)
	protected.Get("/ai/usage", aiUsageHandler.GetUsage)
	protected.Get("/ai/usage/:feature", aiUsageHandler.CheckFeatureQuota)

	protected.Get("/user/profile", userHandler.GetProfile)
	protected.Put("/user/profile", userHandler.UpdateProfile)
	protected.Patch("/user/profile", userHandler.PatchUserPreferences)
	protected.Post("/user/avatar/presigned-url", userHandler.GenerateAvatarUploadURL)

	protected.Get("/meals", mealHandler.GetMeals)
	protected.Post("/meals", mealHandler.CreateMeal)
	protected.Get("/summary/daily", mealHandler.GetDailySummary)

	protected.Get("/foods/search", foodHandler.SearchFoods)
	protected.Get("/foods/recent", foodHandler.GetRecentFoods)
	protected.Post("/ai/autocomplete", foodHandler.AutocompleteFoodMacros)

	protected.Post("/food-items", foodHandler.CreateFoodItem)
	protected.Get("/food-items/:id", foodHandler.GetFoodItem)
	protected.Put("/food-items/:id", foodHandler.UpdateFoodItem)
	protected.Delete("/food-items/:id", foodHandler.DeleteFoodItem)

	protected.Get("/stats", statsHandler.GetStats)
	protected.Get("/stats/range", statsHandler.GetStatsRange)

	protected.Get("/achievements", achievementHandler.GetAchievements)
	protected.Post("/achievements/sync", achievementHandler.SyncAchievements)

	protected.Post("/feedback", feedbackHandler.CreateFeedback)

	protected.Post("/food/recognize", foodRecognitionHandler.RecognizeFood)
	protected.Get("/food/barcode/:barcode", foodRecognitionHandler.GetFoodByBarcode)
	protected.Post("/food/estimate-quantity",
		middleware.AIQuotaCheck(aiUsageService, enum.FeatureMealAnalysis, logger),
		foodRecognitionHandler.EstimateQuantity)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", configurer.HTTPPort)); err != nil {
			logger.Fatal("failed to start server:", err)
		}
	}()

	logger.Info("server started", "port", configurer.HTTPPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Fatal("server forced to shutdown:", err)
	}

	logger.Info("server exited properly")
}
