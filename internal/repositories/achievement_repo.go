package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type AchievementRepository struct {
	queries *db.Queries
}

func NewAchievementRepository(queries *db.Queries) *AchievementRepository {
	return &AchievementRepository{
		queries: queries,
	}
}

func (r *AchievementRepository) GetAllAchievements(ctx context.Context) ([]domain.Achievement, error) {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "GetAllAchievements")
	defer span.End()

	dbAchievements, err := r.queries.GetAllAchievements(ctx)
	if err != nil {
		return nil, err
	}

	achievements := make([]domain.Achievement, 0, len(dbAchievements))
	for _, dbAch := range dbAchievements {
		achievement := domain.Achievement{
			ID:             dbAch.ID,
			NameKey:        dbAch.NameKey,
			DescriptionKey: dbAch.DescriptionKey,
			Icon:           dbAch.Icon,
			Target:         int32(dbAch.Target),
			Category:       dbAch.Category,
			CreatedAt:      timePtrValue(dbAch.CreatedAt),
		}

		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (r *AchievementRepository) GetAchievementByID(ctx context.Context, achievementID uuid.UUID) (*domain.Achievement, error) {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "GetAchievementByID")
	defer span.End()

	dbAch, err := r.queries.GetAchievementByID(ctx, achievementID)
	if err != nil {
		return nil, err
	}

	achievement := &domain.Achievement{
		ID:             dbAch.ID,
		NameKey:        dbAch.NameKey,
		DescriptionKey: dbAch.DescriptionKey,
		Icon:           dbAch.Icon,
		Target:         int32(dbAch.Target),
		Category:       dbAch.Category,
		CreatedAt:      timePtrValue(dbAch.CreatedAt),
	}

	return achievement, nil
}

func (r *AchievementRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]domain.AchievementResponse, error) {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "GetUserAchievements")
	defer span.End()

	dbUserAchievements, err := r.queries.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, err
	}

	achievements := make([]domain.AchievementResponse, 0, len(dbUserAchievements))
	for _, dbUA := range dbUserAchievements {
		achievement := domain.AchievementResponse{
			ID:             dbUA.AchievementID.String(),
			NameKey:        dbUA.NameKey,
			DescriptionKey: dbUA.DescriptionKey,
			Icon:           dbUA.Icon,
			Unlocked:       dbUA.Unlocked,
			UnlockedAt:     dbUA.UnlockedAt,
			Progress:       int32(dbUA.Progress),
			Target:         int32(dbUA.Target),
		}

		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (r *AchievementRepository) UpsertUserAchievement(ctx context.Context, userID, achievementID uuid.UUID, unlocked bool, progress int32, unlockedAt *time.Time) error {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "UpsertUserAchievement")
	defer span.End()

	params := db.UpsertUserAchievementParams{
		UserID:        userID,
		AchievementID: achievementID,
		Unlocked:      unlocked,
		Progress:      int(progress),
		UnlockedAt:    unlockedAt,
	}

	return r.queries.UpsertUserAchievement(ctx, params)
}

func (r *AchievementRepository) UpdateAchievementProgress(ctx context.Context, userID, achievementID uuid.UUID, progress int32) error {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "UpdateAchievementProgress")
	defer span.End()

	params := db.UpdateAchievementProgressParams{
		UserID:        userID,
		AchievementID: achievementID,
		Progress:      int(progress),
	}

	return r.queries.UpdateAchievementProgress(ctx, params)
}

func (r *AchievementRepository) GetUserAchievement(ctx context.Context, userID, achievementID uuid.UUID) (*domain.UserAchievement, error) {
	tr := otel.Tracer("services/achievement_repo.go")
	ctx, span := tr.Start(ctx, "GetUserAchievement")
	defer span.End()

	params := db.GetUserAchievementParams{
		UserID:        userID,
		AchievementID: achievementID,
	}

	dbUA, err := r.queries.GetUserAchievement(ctx, params)
	if err != nil {
		return nil, err
	}

	userAchievement := &domain.UserAchievement{
		ID:            dbUA.ID,
		UserID:        dbUA.UserID,
		AchievementID: dbUA.AchievementID,
		Unlocked:      dbUA.Unlocked,
		Progress:      int32(dbUA.Progress),
		UnlockedAt:    dbUA.UnlockedAt,
		CreatedAt:     timePtrValue(dbUA.CreatedAt),
		UpdatedAt:     timePtrValue(dbUA.UpdatedAt),
	}

	return userAchievement, nil
}
