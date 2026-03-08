package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type FoodRepository struct {
	db *database.Database
}

func NewFoodRepository(db *database.Database) *FoodRepository {
	return &FoodRepository{db: db}
}

func (r *FoodRepository) Create(ctx context.Context, food *domain.FoodItem) error {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	return r.db.Querier.CreateFoodItem(ctx, db.CreateFoodItemParams{
		ID:          food.ID,
		MealID:      food.MealID,
		Name:        food.Name,
		PortionSize: &food.PortionSize,
		PortionUnit: stringToPtr(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     &food.Protein,
		Carbs:       &food.Carbs,
		Fat:         &food.Fat,
		Source:      stringToPtr(food.Source),
	})
}

func (r *FoodRepository) GetByMealID(ctx context.Context, mealID uuid.UUID) ([]domain.FoodItem, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "GetByMealID")
	defer span.End()

	results, err := r.db.Querier.GetFoodItemsByMealID(ctx, mealID)
	if err != nil {
		return nil, err
	}

	foods := make([]domain.FoodItem, 0, len(results))
	for _, result := range results {
		foods = append(foods, domain.FoodItem{
			ID:          result.ID,
			MealID:      result.MealID,
			Name:        result.Name,
			PortionSize: *result.PortionSize,
			PortionUnit: stringPtrValue(result.PortionUnit),
			Calories:    result.Calories,
			Protein:     *result.Protein,
			Carbs:       *result.Carbs,
			Fat:         *result.Fat,
			Source:      stringPtrValue(result.Source),
		})
	}

	return foods, nil
}

func (r *FoodRepository) GetByMealIDs(ctx context.Context, mealIDs []uuid.UUID) (map[uuid.UUID][]domain.FoodItem, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "GetByMealIDs")
	defer span.End()

	results, err := r.db.Querier.GetFoodItemsByMealIDs(ctx, mealIDs)
	if err != nil {
		return nil, err
	}

	foodsByMeal := make(map[uuid.UUID][]domain.FoodItem)
	for _, result := range results {
		food := domain.FoodItem{
			ID:          result.ID,
			MealID:      result.MealID,
			Name:        result.Name,
			PortionSize: *result.PortionSize,
			PortionUnit: *result.PortionUnit,
			Calories:    result.Calories,
			Protein:     *result.Protein,
			Carbs:       *result.Carbs,
			Fat:         *result.Fat,
			Source:      stringPtrValue(result.Source),
		}
		foodsByMeal[food.MealID] = append(foodsByMeal[food.MealID], food)
	}

	return foodsByMeal, nil
}

func (r *FoodRepository) SearchFoodDatabase(ctx context.Context, query string, limit int) ([]domain.FoodDatabase, error) {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "SearchFoodDatabase")
	defer span.End()

	results, err := r.db.Querier.SearchFoodDatabase(ctx, db.SearchFoodDatabaseParams{
		PlaintoTsquery: query,
		Column2:        stringToPtr(query),
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	foods := make([]domain.FoodDatabase, 0, len(results))
	for _, result := range results {
		foods = append(foods, domain.FoodDatabase{
			ID:              result.ID,
			Name:            result.Name,
			Brand:           result.Brand,
			CaloriesPer100g: intPtrValue(result.CaloriesPer100g),
			ProteinPer100g:  *result.ProteinPer100g,
			CarbsPer100g:    *result.CarbsPer100g,
			FatPer100g:      *result.FatPer100g,
			Source:          stringPtrValue(result.Source),
			CreatedAt:       timePtrValue(result.CreatedAt),
		})
	}

	return foods, nil
}

func (r *FoodRepository) GetRecentFoods(ctx context.Context, userID uuid.UUID, limit int) ([]domain.RecentFood, error) {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "GetRecentFoods")
	defer span.End()

	results, err := r.db.Querier.GetRecentFoods(ctx, db.GetRecentFoodsParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	foods := make([]domain.RecentFood, 0, len(results))
	for _, result := range results {
		foods = append(foods, domain.RecentFood{
			Name:        result.Name,
			PortionSize: *result.PortionSize,
			PortionUnit: stringPtrValue(result.PortionUnit),
			Calories:    result.Calories,
			Protein:     *result.Protein,
			Carbs:       *result.Carbs,
			Fat:         *result.Fat,
			LastUsed:    timePtrValue(result.LastUsed),
		})
	}

	return foods, nil
}

func (r *FoodRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "GetByID")
	defer span.End()

	result, err := r.db.Querier.GetFoodItemByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          result.ID,
		MealID:      result.MealID,
		Name:        result.Name,
		PortionSize: *result.PortionSize,
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     *result.Protein,
		Carbs:       *result.Carbs,
		Fat:         *result.Fat,
		Source:      stringPtrValue(result.Source),
	}, nil
}

func (r *FoodRepository) Update(ctx context.Context, id uuid.UUID, food *domain.UpdateFoodItemRequest) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	result, err := r.db.Querier.UpdateFoodItemComplete(ctx, db.UpdateFoodItemCompleteParams{
		ID:          id,
		Name:        food.Name,
		PortionSize: &food.PortionSize,
		PortionUnit: stringToPtr(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     &food.Protein,
		Carbs:       &food.Carbs,
		Fat:         &food.Fat,
		Source:      stringToPtr(food.Source),
	})
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          result.ID,
		MealID:      result.MealID,
		Name:        result.Name,
		PortionSize: *result.PortionSize,
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     *result.Protein,
		Carbs:       *result.Carbs,
		Fat:         *result.Fat,
		Source:      stringPtrValue(result.Source),
	}, nil
}

func (r *FoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	return r.db.Querier.DeleteFoodItem(ctx, id)
}

func (r *FoodRepository) CreateStandalone(ctx context.Context, food *domain.CreateFoodItemRequest) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_repo.go")
	ctx, span := tr.Start(ctx, "CreateStandalone")
	defer span.End()

	id := uuid.New()
	result, err := r.db.Querier.CreateStandaloneFoodItem(ctx, db.CreateStandaloneFoodItemParams{
		ID:          id,
		MealID:      food.MealID,
		Name:        food.Name,
		PortionSize: &food.PortionSize,
		PortionUnit: stringToPtr(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     &food.Protein,
		Carbs:       &food.Carbs,
		Fat:         &food.Fat,
		Source:      stringToPtr(food.Source),
	})
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          result.ID,
		MealID:      result.MealID,
		Name:        result.Name,
		PortionSize: *result.PortionSize,
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     *result.Protein,
		Carbs:       *result.Carbs,
		Fat:         *result.Fat,
		Source:      stringPtrValue(result.Source),
	}, nil
}
