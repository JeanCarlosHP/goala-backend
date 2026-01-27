package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
)

type AIUsageRepository struct {
	db *database.Database
}

func NewAIUsageRepository(db *database.Database) *AIUsageRepository {
	return &AIUsageRepository{db: db}
}

func (r *AIUsageRepository) Increment(ctx context.Context, userID string, feature enum.AIFeature, quota int32, periodStart, periodEnd time.Time) (*domain.AIUsage, error) {
	result, err := r.db.Querier.IncrementAIUsage(ctx, db.IncrementAIUsageParams{
		UserID:      stringToPgUUID(userID),
		Feature:     feature.String(),
		Quota:       int(quota),
		PeriodStart: pgtype.Timestamptz{Time: periodStart, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return toAIUsage(&result), nil
}

func (r *AIUsageRepository) Get(ctx context.Context, userID string, feature enum.AIFeature) (*domain.AIUsage, error) {
	result, err := r.db.Querier.GetAIUsage(ctx, db.GetAIUsageParams{
		UserID:  stringToPgUUID(userID),
		Feature: feature.String(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toAIUsage(&result), nil
}

func (r *AIUsageRepository) GetByPeriod(ctx context.Context, userID string, feature enum.AIFeature, periodStart time.Time) (*domain.AIUsage, error) {
	result, err := r.db.Querier.GetAIUsageByPeriod(ctx, db.GetAIUsageByPeriodParams{
		UserID:      stringToPgUUID(userID),
		Feature:     feature.String(),
		PeriodStart: pgtype.Timestamptz{Time: periodStart, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toAIUsage(&result), nil
}

func (r *AIUsageRepository) Reset(ctx context.Context) error {
	return r.db.Querier.ResetAIUsage(ctx)
}

func (r *AIUsageRepository) ListByUser(ctx context.Context, userID string) ([]*domain.AIUsage, error) {
	results, err := r.db.Querier.ListUserAIUsage(ctx, stringToPgUUID(userID))
	if err != nil {
		return nil, err
	}

	usages := make([]*domain.AIUsage, len(results))
	for i, result := range results {
		usages[i] = toAIUsage(&result)
	}
	return usages, nil
}

func (r *AIUsageRepository) CreateOrReset(ctx context.Context, userID string, feature enum.AIFeature, quota int32, periodStart, periodEnd time.Time) (*domain.AIUsage, error) {
	result, err := r.db.Querier.CreateOrResetAIUsage(ctx, db.CreateOrResetAIUsageParams{
		UserID:      stringToPgUUID(userID),
		Feature:     feature.String(),
		Quota:       int(quota),
		PeriodStart: pgtype.Timestamptz{Time: periodStart, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return toAIUsage(&result), nil
}

func toAIUsage(u *db.AiUsage) *domain.AIUsage {
	return &domain.AIUsage{
		ID:          u.ID,
		UserID:      uuid.UUID(u.UserID.Bytes).String(),
		Feature:     enum.AIFeature(u.Feature),
		UsageCount:  int32(u.UsageCount),
		Quota:       int32(u.Quota),
		PeriodStart: u.PeriodStart.Time,
		PeriodEnd:   u.PeriodEnd.Time,
		CreatedAt:   u.CreatedAt.Time,
		UpdatedAt:   u.UpdatedAt.Time,
	}
}

func stringToPgUUID(s string) pgtype.UUID {
	u, _ := uuid.Parse(s)
	return pgtype.UUID{Bytes: u, Valid: true}
}
