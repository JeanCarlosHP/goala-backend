package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
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
	dbAchievements, err := r.queries.GetAllAchievements(ctx)
	if err != nil {
		return nil, err
	}

	achievements := make([]domain.Achievement, 0, len(dbAchievements))
	for _, dbAch := range dbAchievements {
		achievement := domain.Achievement{
			ID:             uuidFromPgtype(dbAch.ID),
			NameKey:        dbAch.NameKey,
			DescriptionKey: dbAch.DescriptionKey,
			Icon:           dbAch.Icon,
			Target:         int32(dbAch.Target),
			Category:       dbAch.Category,
		}

		if dbAch.CreatedAt != nil {
			achievement.CreatedAt = *dbAch.CreatedAt
		}

		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (r *AchievementRepository) GetAchievementByID(ctx context.Context, achievementID uuid.UUID) (*domain.Achievement, error) {
	var pgAchievementID pgtype.UUID
	if err := pgAchievementID.Scan(achievementID.String()); err != nil {
		return nil, err
	}

	dbAch, err := r.queries.GetAchievementByID(ctx, pgAchievementID)
	if err != nil {
		return nil, err
	}

	achievement := &domain.Achievement{
		ID:             uuidFromPgtype(dbAch.ID),
		NameKey:        dbAch.NameKey,
		DescriptionKey: dbAch.DescriptionKey,
		Icon:           dbAch.Icon,
		Target:         int32(dbAch.Target),
		Category:       dbAch.Category,
	}

	if dbAch.CreatedAt != nil {
		achievement.CreatedAt = *dbAch.CreatedAt
	}

	return achievement, nil
}

func (r *AchievementRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]domain.AchievementResponse, error) {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return nil, err
	}

	dbUserAchievements, err := r.queries.GetUserAchievements(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	achievements := make([]domain.AchievementResponse, 0, len(dbUserAchievements))
	for _, dbUA := range dbUserAchievements {
		achievement := domain.AchievementResponse{
			ID:             uuidFromPgtype(dbUA.AchievementID).String(),
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
	var pgUserID, pgAchievementID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return err
	}
	if err := pgAchievementID.Scan(achievementID.String()); err != nil {
		return err
	}

	params := db.UpsertUserAchievementParams{
		UserID:        pgUserID,
		AchievementID: pgAchievementID,
		Unlocked:      unlocked,
		Progress:      int(progress),
		UnlockedAt:    unlockedAt,
	}

	return r.queries.UpsertUserAchievement(ctx, params)
}

func (r *AchievementRepository) UpdateAchievementProgress(ctx context.Context, userID, achievementID uuid.UUID, progress int32) error {
	var pgUserID, pgAchievementID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return err
	}
	if err := pgAchievementID.Scan(achievementID.String()); err != nil {
		return err
	}

	params := db.UpdateAchievementProgressParams{
		UserID:        pgUserID,
		AchievementID: pgAchievementID,
		Progress:      int(progress),
	}

	return r.queries.UpdateAchievementProgress(ctx, params)
}

func (r *AchievementRepository) GetUserAchievement(ctx context.Context, userID, achievementID uuid.UUID) (*domain.UserAchievement, error) {
	var pgUserID, pgAchievementID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return nil, err
	}
	if err := pgAchievementID.Scan(achievementID.String()); err != nil {
		return nil, err
	}

	params := db.GetUserAchievementParams{
		UserID:        pgUserID,
		AchievementID: pgAchievementID,
	}

	dbUA, err := r.queries.GetUserAchievement(ctx, params)
	if err != nil {
		return nil, err
	}

	userAchievement := &domain.UserAchievement{
		ID:            uuidFromPgtype(dbUA.ID),
		UserID:        uuidFromPgtype(dbUA.UserID),
		AchievementID: uuidFromPgtype(dbUA.AchievementID),
		Unlocked:      dbUA.Unlocked,
		Progress:      int32(dbUA.Progress),
		UnlockedAt:    dbUA.UnlockedAt,
	}

	if dbUA.CreatedAt != nil {
		userAchievement.CreatedAt = *dbUA.CreatedAt
	}

	if dbUA.UpdatedAt != nil {
		userAchievement.UpdatedAt = *dbUA.UpdatedAt
	}

	return userAchievement, nil
}
