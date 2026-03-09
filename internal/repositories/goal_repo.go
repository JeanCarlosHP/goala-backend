package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type GoalRepository struct {
	db *database.Database
}

func NewGoalRepository(db *database.Database) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Upsert(ctx context.Context, goal *domain.UserGoal) error {
	tr := otel.Tracer("repositories/goal_repo.go")
	ctx, span := tr.Start(ctx, "Upsert")
	defer span.End()

	result, err := r.db.Querier.UpsertUserGoal(ctx, db.UpsertUserGoalParams{
		UserID:        goal.UserID,
		DailyCalories: goal.DailyCalorieGoal,
		Protein:       new(goal.DailyProteinGoal),
		Carbs:         new(goal.DailyCarbsGoal),
		Fat:           new(goal.DailyFatGoal),
	})
	if err != nil {
		return err
	}

	goal.UpdatedAt = timePtrValue(result.UpdatedAt)
	return nil
}

func (r *GoalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserGoal, error) {
	tr := otel.Tracer("repositories/goal_repo.go")
	ctx, span := tr.Start(ctx, "GetByUserID")
	defer span.End()

	result, err := r.db.Querier.GetUserGoalByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserGoal{
		UserID:           result.UserID,
		DailyCalorieGoal: result.DailyCalories,
		DailyProteinGoal: intPtrValue(result.Protein),
		DailyCarbsGoal:   intPtrValue(result.Carbs),
		DailyFatGoal:     intPtrValue(result.Fat),
		UpdatedAt:        timePtrValue(result.UpdatedAt),
	}, nil
}
