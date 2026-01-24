package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
)

type FeedbackRepository struct {
	db *database.Database
}

func NewFeedbackRepository(db *database.Database) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateFeedbackRequest) (*domain.Feedback, error) {
	var platform *string
	var osVersion *string
	var appVersion *string

	if req.DeviceInfo != nil {
		platform = &req.DeviceInfo.Platform
		osVersion = &req.DeviceInfo.OsVersion
		appVersion = &req.DeviceInfo.AppVersion
	}

	result, err := r.db.Querier.CreateFeedback(ctx, db.CreateFeedbackParams{
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
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

	return &domain.Feedback{
		ID:          uuid.UUID(result.ID.Bytes).String(),
		UserID:      uuidToStringPtr(result.UserID),
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
	result, err := r.db.Querier.GetFeedback(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &domain.Feedback{
		ID:          uuid.UUID(result.ID.Bytes).String(),
		UserID:      uuidToStringPtr(result.UserID),
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
	results, err := r.db.Querier.ListFeedback(ctx, db.ListFeedbackParams{
		Limit:  int(limit),
		Offset: int(offset),
	})
	if err != nil {
		return nil, err
	}

	feedbacks := make([]domain.Feedback, len(results))
	for i, result := range results {
		feedbacks[i] = domain.Feedback{
			ID:          uuid.UUID(result.ID.Bytes).String(),
			UserID:      uuidToStringPtr(result.UserID),
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
	results, err := r.db.Querier.GetFeedbackByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	feedbacks := make([]domain.Feedback, len(results))
	for i, result := range results {
		feedbacks[i] = domain.Feedback{
			ID:          uuid.UUID(result.ID.Bytes).String(),
			UserID:      uuidToStringPtr(result.UserID),
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

func uuidToStringPtr(pguuid pgtype.UUID) *string {
	if !pguuid.Valid {
		return nil
	}
	id := uuid.UUID(pguuid.Bytes)
	str := id.String()
	return &str
}
