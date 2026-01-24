package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
)

type MealRepository struct {
	db *database.Database
}

func NewMealRepository(db *database.Database) *MealRepository {
	return &MealRepository{db: db}
}

func (r *MealRepository) Create(ctx context.Context, meal *domain.Meal) error {
	mealDate := pgtype.Date{}
	_ = mealDate.Scan(meal.MealDate)

	var mealTime pgtype.Time
	if meal.MealTime != nil {
		_ = mealTime.Scan(*meal.MealTime)
	}

	result, err := r.db.Querier.CreateMeal(ctx, db.CreateMealParams{
		ID:       pgtype.UUID{Bytes: meal.ID, Valid: true},
		UserID:   pgtype.UUID{Bytes: meal.UserID, Valid: true},
		MealType: stringToPtr(meal.MealType),
		MealDate: mealDate,
		MealTime: mealTime,
		PhotoUrl: meal.PhotoURL,
	})
	if err != nil {
		return err
	}

	meal.CreatedAt = timePtrValue(result.CreatedAt)
	return nil
}

func (r *MealRepository) GetByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]domain.Meal, error) {
	mealDate := pgtype.Date{}
	_ = mealDate.Scan(date)

	results, err := r.db.Querier.GetMealsByUserAndDate(ctx, db.GetMealsByUserAndDateParams{
		UserID:   pgtype.UUID{Bytes: userID, Valid: true},
		MealDate: mealDate,
	})
	if err != nil {
		return nil, err
	}

	meals := make([]domain.Meal, 0, len(results))
	for _, result := range results {
		var mealTime *time.Time
		if result.MealTime.Valid {
			t := result.MealTime.Microseconds / 1000000
			h := int(t / 3600)
			m := int((t % 3600) / 60)
			s := int(t % 60)
			parsed := time.Date(0, 1, 1, h, m, s, 0, time.UTC)
			mealTime = &parsed
		}

		var mealDate time.Time
		if result.MealDate.Valid {
			mealDate = result.MealDate.Time
		}

		meals = append(meals, domain.Meal{
			ID:        result.ID.Bytes,
			UserID:    result.UserID.Bytes,
			MealType:  stringPtrValue(result.MealType),
			MealDate:  mealDate,
			MealTime:  mealTime,
			PhotoURL:  result.PhotoUrl,
			CreatedAt: timePtrValue(result.CreatedAt),
		})
	}

	return meals, nil
}

func (r *MealRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Meal, error) {
	result, err := r.db.Querier.GetMealByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	var mealTime *time.Time
	if result.MealTime.Valid {
		t := result.MealTime.Microseconds / 1000000
		h := int(t / 3600)
		m := int((t % 3600) / 60)
		s := int(t % 60)
		parsed := time.Date(0, 1, 1, h, m, s, 0, time.UTC)
		mealTime = &parsed
	}

	var mealDate time.Time
	if result.MealDate.Valid {
		mealDate = result.MealDate.Time
	}

	return &domain.Meal{
		ID:        result.ID.Bytes,
		UserID:    result.UserID.Bytes,
		MealType:  stringPtrValue(result.MealType),
		MealDate:  mealDate,
		MealTime:  mealTime,
		PhotoURL:  result.PhotoUrl,
		CreatedAt: timePtrValue(result.CreatedAt),
	}, nil
}
