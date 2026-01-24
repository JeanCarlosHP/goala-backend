package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
)

type StatsRepository struct {
	queries *db.Queries
}

func NewStatsRepository(queries *db.Queries) *StatsRepository {
	return &StatsRepository{
		queries: queries,
	}
}

func (r *StatsRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*domain.UserStats, error) {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return nil, err
	}

	dbStats, err := r.queries.GetUserStats(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	stats := &domain.UserStats{
		ID:                    uuidFromPgtype(dbStats.ID),
		UserID:                uuidFromPgtype(dbStats.UserID),
		CurrentStreak:         int32(dbStats.CurrentStreak),
		LongestStreak:         int32(dbStats.LongestStreak),
		TotalMealsLogged:      int32(dbStats.TotalMealsLogged),
		TotalDaysLogged:       int32(dbStats.TotalDaysLogged),
		TotalCaloriesConsumed: int32(dbStats.TotalCaloriesConsumed),
		TotalProteinConsumed:  int32(dbStats.TotalProteinConsumed),
		TotalCarbsConsumed:    int32(dbStats.TotalCarbsConsumed),
		TotalFatConsumed:      int32(dbStats.TotalFatConsumed),
	}

	if dbStats.LastLogDate.Valid {
		lastLogDate := dbStats.LastLogDate.Time
		stats.LastLogDate = &lastLogDate
	}

	if dbStats.CreatedAt != nil {
		stats.CreatedAt = *dbStats.CreatedAt
	}

	if dbStats.UpdatedAt != nil {
		stats.UpdatedAt = *dbStats.UpdatedAt
	}

	return stats, nil
}

func (r *StatsRepository) CreateUserStats(ctx context.Context, userID uuid.UUID) (*domain.UserStats, error) {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return nil, err
	}

	dbStats, err := r.queries.CreateUserStats(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	stats := &domain.UserStats{
		ID:                    uuidFromPgtype(dbStats.ID),
		UserID:                uuidFromPgtype(dbStats.UserID),
		CurrentStreak:         int32(dbStats.CurrentStreak),
		LongestStreak:         int32(dbStats.LongestStreak),
		TotalMealsLogged:      int32(dbStats.TotalMealsLogged),
		TotalDaysLogged:       int32(dbStats.TotalDaysLogged),
		TotalCaloriesConsumed: int32(dbStats.TotalCaloriesConsumed),
		TotalProteinConsumed:  int32(dbStats.TotalProteinConsumed),
		TotalCarbsConsumed:    int32(dbStats.TotalCarbsConsumed),
		TotalFatConsumed:      int32(dbStats.TotalFatConsumed),
	}

	if dbStats.LastLogDate.Valid {
		lastLogDate := dbStats.LastLogDate.Time
		stats.LastLogDate = &lastLogDate
	}

	if dbStats.CreatedAt != nil {
		stats.CreatedAt = *dbStats.CreatedAt
	}

	if dbStats.UpdatedAt != nil {
		stats.UpdatedAt = *dbStats.UpdatedAt
	}

	return stats, nil
}

func (r *StatsRepository) UpdateUserStats(ctx context.Context, stats *domain.UserStats) error {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(stats.UserID.String()); err != nil {
		return err
	}

	var lastLogDate pgtype.Date
	if stats.LastLogDate != nil {
		lastLogDate.Time = *stats.LastLogDate
		lastLogDate.Valid = true
	}

	params := db.UpdateUserStatsParams{
		UserID:                pgUserID,
		CurrentStreak:         int(stats.CurrentStreak),
		LongestStreak:         int(stats.LongestStreak),
		TotalMealsLogged:      int(stats.TotalMealsLogged),
		TotalDaysLogged:       int(stats.TotalDaysLogged),
		TotalCaloriesConsumed: int(stats.TotalCaloriesConsumed),
		TotalProteinConsumed:  int(stats.TotalProteinConsumed),
		TotalCarbsConsumed:    int(stats.TotalCarbsConsumed),
		TotalFatConsumed:      int(stats.TotalFatConsumed),
		LastLogDate:           lastLogDate,
	}

	return r.queries.UpdateUserStats(ctx, params)
}

func (r *StatsRepository) IncrementMealCount(ctx context.Context, userID uuid.UUID) error {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return err
	}

	return r.queries.IncrementMealCount(ctx, pgUserID)
}

func (r *StatsRepository) UpdateStreakAndLastLogDate(ctx context.Context, userID uuid.UUID, currentStreak int32, lastLogDate time.Time) error {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return err
	}

	var pgLastLogDate pgtype.Date
	pgLastLogDate.Time = lastLogDate
	pgLastLogDate.Valid = true

	params := db.UpdateStreakAndLastLogDateParams{
		UserID:        pgUserID,
		CurrentStreak: int(currentStreak),
		LastLogDate:   pgLastLogDate,
	}

	return r.queries.UpdateStreakAndLastLogDate(ctx, params)
}

func (r *StatsRepository) AddNutritionToStats(ctx context.Context, userID uuid.UUID, calories, protein, carbs, fat int32) error {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID.String()); err != nil {
		return err
	}

	params := db.AddNutritionToStatsParams{
		UserID:                pgUserID,
		TotalCaloriesConsumed: int(calories),
		TotalProteinConsumed:  int(protein),
		TotalCarbsConsumed:    int(carbs),
		TotalFatConsumed:      int(fat),
	}

	return r.queries.AddNutritionToStats(ctx, params)
}
