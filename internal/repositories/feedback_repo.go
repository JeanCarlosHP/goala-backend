package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type FeedbackRepository struct {
	db *database.Database
}

func NewFeedbackRepository(db *database.Database) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateFeedbackRequest) (*domain.Feedback, error) {
	tr := otel.Tracer("repositories/feedback_repo.go")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	var platform *string
	var osVersion *string
	var appVersion *string

	if req.DeviceInfo != nil {
		platform = new(req.DeviceInfo.Platform)
		osVersion = new(req.DeviceInfo.OsVersion)
		appVersion = new(req.DeviceInfo.AppVersion)
	}

	result, err := r.db.Querier.CreateFeedback(ctx, db.CreateFeedbackParams{
		UserID:      userID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		UserEmail:   req.UserEmail,
		Platform:    platform,
		OsVersion:   osVersion,
		AppVersion:  appVersion,
	})
	if err != nil {
		return nil, err
	}

	userIDStr := result.UserID.String()
	return &domain.Feedback{
		ID:          uuid.UUID(result.ID).String(),
		UserID:      new(userIDStr),
		Type:        enum.FeedbackType(result.Type),
		Title:       result.Title,
		Description: result.Description,
		UserEmail:   result.UserEmail,
		Platform:    result.Platform,
		OsVersion:   result.OsVersion,
		AppVersion:  result.AppVersion,
		Status:      stringPtrValue(result.Status),
		CreatedAt:   timePtrValue(result.CreatedAt),
	}, nil
}

func (r *FeedbackRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Feedback, error) {
	tr := otel.Tracer("repositories/feedback_repo.go")
	ctx, span := tr.Start(ctx, "GetByID")
	defer span.End()

	result, err := r.db.Querier.GetFeedback(ctx, id)
	if err != nil {
		return nil, err
	}

	userIDStr := result.UserID.String()
	return &domain.Feedback{
		ID:          uuid.UUID(result.ID).String(),
		UserID:      new(userIDStr),
		Type:        enum.FeedbackType(result.Type),
		Title:       result.Title,
		Description: result.Description,
		UserEmail:   result.UserEmail,
		Platform:    result.Platform,
		OsVersion:   result.OsVersion,
		AppVersion:  result.AppVersion,
		Status:      stringPtrValue(result.Status),
		CreatedAt:   timePtrValue(result.CreatedAt),
	}, nil
}

func (r *FeedbackRepository) List(ctx context.Context, limit, offset int32) ([]domain.Feedback, error) {
	tr := otel.Tracer("repositories/feedback_repo.go")
	ctx, span := tr.Start(ctx, "List")
	defer span.End()

	results, err := r.db.Querier.ListFeedback(ctx, db.ListFeedbackParams{
		Limit:  int(limit),
		Offset: int(offset),
	})
	if err != nil {
		return nil, err
	}

	feedbacks := make([]domain.Feedback, len(results))
	for i, result := range results {
		userIDStr := result.UserID.String()
		feedbacks[i] = domain.Feedback{
			ID:          uuid.UUID(result.ID).String(),
			UserID:      new(userIDStr),
			Type:        enum.FeedbackType(result.Type),
			Title:       result.Title,
			Description: result.Description,
			UserEmail:   result.UserEmail,
			Platform:    result.Platform,
			OsVersion:   result.OsVersion,
			AppVersion:  result.AppVersion,
			Status:      stringPtrValue(result.Status),
			CreatedAt:   timePtrValue(result.CreatedAt),
		}
	}

	return feedbacks, nil
}

func (r *FeedbackRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Feedback, error) {
	tr := otel.Tracer("repositories/feedback_repo.go")
	ctx, span := tr.Start(ctx, "GetByUser")
	defer span.End()

	results, err := r.db.Querier.GetFeedbackByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	feedbacks := make([]domain.Feedback, len(results))
	for i, result := range results {
		feedbacks[i] = domain.Feedback{
			ID:          uuid.UUID(result.ID).String(),
			UserID:      new(result.UserID.String()),
			Type:        enum.FeedbackType(result.Type),
			Title:       result.Title,
			Description: result.Description,
			UserEmail:   result.UserEmail,
			Platform:    result.Platform,
			OsVersion:   result.OsVersion,
			AppVersion:  result.AppVersion,
			Status:      stringPtrValue(result.Status),
			CreatedAt:   timePtrValue(result.CreatedAt),
		}
	}

	return feedbacks, nil
}
