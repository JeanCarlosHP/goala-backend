package repositories

import (
	"context"
	"fmt"
	"strings"

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
		PortionSize: new(food.PortionSize),
		PortionUnit: new(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     new(food.Protein),
		Carbs:       new(food.Carbs),
		Fat:         new(food.Fat),
		Source:      new(food.Source),
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
			PortionSize: valueOrZero(result.PortionSize),
			PortionUnit: stringPtrValue(result.PortionUnit),
			Calories:    result.Calories,
			Protein:     valueOrZero(result.Protein),
			Carbs:       valueOrZero(result.Carbs),
			Fat:         valueOrZero(result.Fat),
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
			PortionSize: valueOrZero(result.PortionSize),
			PortionUnit: stringPtrValue(result.PortionUnit),
			Calories:    result.Calories,
			Protein:     valueOrZero(result.Protein),
			Carbs:       valueOrZero(result.Carbs),
			Fat:         valueOrZero(result.Fat),
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
		Column2:        new(query),
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
			ProteinPer100g:  valueOrZero(result.ProteinPer100g),
			CarbsPer100g:    valueOrZero(result.CarbsPer100g),
			FatPer100g:      valueOrZero(result.FatPer100g),
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
			PortionSize: valueOrZero(result.PortionSize),
			PortionUnit: stringPtrValue(result.PortionUnit),
			Calories:    result.Calories,
			Protein:     valueOrZero(result.Protein),
			Carbs:       valueOrZero(result.Carbs),
			Fat:         valueOrZero(result.Fat),
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
		PortionSize: valueOrZero(result.PortionSize),
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     valueOrZero(result.Protein),
		Carbs:       valueOrZero(result.Carbs),
		Fat:         valueOrZero(result.Fat),
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
		PortionSize: new(food.PortionSize),
		PortionUnit: new(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     new(food.Protein),
		Carbs:       new(food.Carbs),
		Fat:         new(food.Fat),
		Source:      new(food.Source),
	})
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          result.ID,
		MealID:      result.MealID,
		Name:        result.Name,
		PortionSize: valueOrZero(result.PortionSize),
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     valueOrZero(result.Protein),
		Carbs:       valueOrZero(result.Carbs),
		Fat:         valueOrZero(result.Fat),
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
		PortionSize: new(food.PortionSize),
		PortionUnit: new(food.PortionUnit),
		Calories:    food.Calories,
		Protein:     new(food.Protein),
		Carbs:       new(food.Carbs),
		Fat:         new(food.Fat),
		Source:      new(food.Source),
	})
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          result.ID,
		MealID:      result.MealID,
		Name:        result.Name,
		PortionSize: valueOrZero(result.PortionSize),
		PortionUnit: stringPtrValue(result.PortionUnit),
		Calories:    result.Calories,
		Protein:     valueOrZero(result.Protein),
		Carbs:       valueOrZero(result.Carbs),
		Fat:         valueOrZero(result.Fat),
		Source:      stringPtrValue(result.Source),
	}, nil
}

func (r *FoodRepository) SearchFoodsByIDs(
	ctx context.Context,
	foodIDs []uuid.UUID,
	userID uuid.UUID,
) ([]domain.SearchFood, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "SearchFoodsByIDs")
	defer span.End()

	if len(foodIDs) == 0 {
		return []domain.SearchFood{}, nil
	}

	query := `
		SELECT
			fd.id,
			fd.external_id,
			fd.name,
			fd.brand,
			COALESCE(fd.calories_per_100g, fd.calories, 0) AS calories,
			COALESCE(fd.protein_per_100g, fd.protein, 0) AS protein,
			COALESCE(fd.carbs_per_100g, fd.carbs, 0) AS carbs,
			COALESCE(fd.fat_per_100g, fd.fat, 0) AS fat,
			COALESCE(fd.source, 'internal') AS source,
			EXISTS(
				SELECT 1
				FROM favorite_foods ff
				WHERE ff.user_id = $2
				  AND ff.food_id = fd.id
			) AS is_favorite
		FROM food_database fd
		WHERE fd.id = ANY($1::uuid[])
	`

	rows, err := r.db.Pool.Query(ctx, query, foodIDs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSearchFoods(ctx, rows)
}

func (r *FoodRepository) SearchFoodsForAutocomplete(
	ctx context.Context,
	query string,
	limit int,
	userID uuid.UUID,
) ([]domain.SearchFood, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "SearchFoodsForAutocomplete")
	defer span.End()

	sqlQuery := `
		SELECT
			fd.id,
			fd.external_id,
			fd.name,
			fd.brand,
			COALESCE(fd.calories_per_100g, fd.calories, 0) AS calories,
			COALESCE(fd.protein_per_100g, fd.protein, 0) AS protein,
			COALESCE(fd.carbs_per_100g, fd.carbs, 0) AS carbs,
			COALESCE(fd.fat_per_100g, fd.fat, 0) AS fat,
			COALESCE(fd.source, 'internal') AS source,
			EXISTS(
				SELECT 1
				FROM favorite_foods ff
				WHERE ff.user_id = $2
				  AND ff.food_id = fd.id
			) AS is_favorite
		FROM food_database fd
		WHERE fd.name ILIKE '%' || $1 || '%'
		   OR COALESCE(fd.brand, '') ILIKE '%' || $1 || '%'
		ORDER BY
			CASE
				WHEN LOWER(fd.name) = LOWER($1) THEN 0
				WHEN LOWER(fd.name) LIKE LOWER($1) || '%' THEN 1
				ELSE 2
			END,
			COALESCE(fd.verified, false) DESC,
			fd.name ASC
		LIMIT $3
	`

	rows, err := r.db.Pool.Query(ctx, sqlQuery, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSearchFoods(ctx, rows)
}

func (r *FoodRepository) GetRecentSearchFoods(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]domain.SearchFood, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "GetRecentSearchFoods")
	defer span.End()

	sqlQuery := `
		SELECT
			fd.id,
			fd.external_id,
			fd.name,
			fd.brand,
			COALESCE(fd.calories_per_100g, fd.calories, 0) AS calories,
			COALESCE(fd.protein_per_100g, fd.protein, 0) AS protein,
			COALESCE(fd.carbs_per_100g, fd.carbs, 0) AS carbs,
			COALESCE(fd.fat_per_100g, fd.fat, 0) AS fat,
			COALESCE(fd.source, 'internal') AS source,
			EXISTS(
				SELECT 1
				FROM favorite_foods ff
				WHERE ff.user_id = $1
				  AND ff.food_id = fd.id
			) AS is_favorite
		FROM food_items fi
		JOIN meals m ON m.id = fi.meal_id
		JOIN food_database fd ON fd.id = fi.food_database_id
		WHERE m.user_id = $1
		ORDER BY fi.id DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, sqlQuery, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSearchFoods(ctx, rows)
}

func (r *FoodRepository) GetFavoriteFoods(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]domain.SearchFood, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "GetFavoriteFoods")
	defer span.End()

	sqlQuery := `
		SELECT
			fd.id,
			fd.external_id,
			fd.name,
			fd.brand,
			COALESCE(fd.calories_per_100g, fd.calories, 0) AS calories,
			COALESCE(fd.protein_per_100g, fd.protein, 0) AS protein,
			COALESCE(fd.carbs_per_100g, fd.carbs, 0) AS carbs,
			COALESCE(fd.fat_per_100g, fd.fat, 0) AS fat,
			COALESCE(fd.source, 'internal') AS source,
			true AS is_favorite
		FROM favorite_foods ff
		JOIN food_database fd ON fd.id = ff.food_id
		WHERE ff.user_id = $1
		ORDER BY ff.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, sqlQuery, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSearchFoods(ctx, rows)
}

func (r *FoodRepository) ToggleFavorite(
	ctx context.Context,
	userID uuid.UUID,
	foodID uuid.UUID,
	favorite bool,
) error {
	if favorite {
		_, err := r.db.Pool.Exec(ctx, `
			INSERT INTO favorite_foods (user_id, food_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, food_id) DO NOTHING
		`, userID, foodID)
		return err
	}

	_, err := r.db.Pool.Exec(ctx, `
		DELETE FROM favorite_foods
		WHERE user_id = $1 AND food_id = $2
	`, userID, foodID)
	return err
}

func (r *FoodRepository) UpsertFoodCatalogEntry(
	ctx context.Context,
	food domain.SearchFood,
) (*domain.SearchFood, error) {
	tr := otel.Tracer("repositories/food_repo.go")
	ctx, span := tr.Start(ctx, "UpsertFoodCatalogEntry")
	defer span.End()

	sqlQuery := `
		INSERT INTO food_database (
			id, external_id, name, brand, calories_per_100g, protein_per_100g,
			carbs_per_100g, fat_per_100g, source, verified, updated_at
		)
		VALUES (
			$1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, false, NOW()
		)
		ON CONFLICT (external_id) WHERE external_id IS NOT NULL
		DO UPDATE SET
			name = EXCLUDED.name,
			brand = EXCLUDED.brand,
			calories_per_100g = EXCLUDED.calories_per_100g,
			protein_per_100g = EXCLUDED.protein_per_100g,
			carbs_per_100g = EXCLUDED.carbs_per_100g,
			fat_per_100g = EXCLUDED.fat_per_100g,
			source = EXCLUDED.source,
			updated_at = NOW()
		RETURNING id
	`

	id := uuid.New()
	externalID := stringValue(food.ExternalID)
	source := food.Source
	if source == "" {
		source = "internal"
	}

	var storedID uuid.UUID
	err := r.db.Pool.QueryRow(
		ctx,
		sqlQuery,
		id,
		externalID,
		food.Name,
		food.Brand,
		int(food.Calories+0.5),
		food.Protein,
		food.Carbs,
		food.Fat,
		source,
	).Scan(&storedID)
	if err != nil {
		if externalID == "" {
			err = r.db.Pool.QueryRow(
				ctx,
				`SELECT id FROM food_database WHERE LOWER(name) = LOWER($1) AND COALESCE(LOWER(brand), '') = COALESCE(LOWER($2), '') LIMIT 1`,
				food.Name,
				stringValue(food.Brand),
			).Scan(&storedID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := r.replaceFoodPortions(ctx, storedID, food.Portions); err != nil {
		return nil, err
	}

	food.ID = &storedID
	return &food, nil
}

func (r *FoodRepository) CreateLoggedFoodItem(
	ctx context.Context,
	mealID uuid.UUID,
	foodID *uuid.UUID,
	name string,
	quantity float64,
	portionUnit string,
	calories int,
	protein float64,
	carbs float64,
	fat float64,
) (*domain.FoodItem, error) {
	id := uuid.New()
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO food_items (
			id, meal_id, food_database_id, name, portion_size, portion_unit,
			calories, protein, carbs, fat, source
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'manual')
	`, id, mealID, foodID, name, quantity, portionUnit, calories, protein, carbs, fat)
	if err != nil {
		return nil, err
	}

	return &domain.FoodItem{
		ID:          id,
		MealID:      mealID,
		Name:        name,
		PortionSize: quantity,
		PortionUnit: portionUnit,
		Calories:    calories,
		Protein:     protein,
		Carbs:       carbs,
		Fat:         fat,
		Source:      "manual",
	}, nil
}

func (r *FoodRepository) scanSearchFoods(ctx context.Context, rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}) ([]domain.SearchFood, error) {
	foods := make([]domain.SearchFood, 0)
	ids := make([]uuid.UUID, 0)

	for rows.Next() {
		var item domain.SearchFood
		var foodID uuid.UUID
		if err := rows.Scan(
			&foodID,
			&item.ExternalID,
			&item.Name,
			&item.Brand,
			&item.Calories,
			&item.Protein,
			&item.Carbs,
			&item.Fat,
			&item.Source,
			&item.IsFavorite,
		); err != nil {
			return nil, err
		}
		item.ID = &foodID
		foods = append(foods, item)
		ids = append(ids, foodID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	portionsByFoodID, err := r.getFoodPortions(ctx, ids)
	if err != nil {
		return nil, err
	}

	for i := range foods {
		if foods[i].ID != nil {
			foods[i].Portions = append(defaultPortions(), portionsByFoodID[*foods[i].ID]...)
		}
	}

	return foods, nil
}

func (r *FoodRepository) getFoodPortions(ctx context.Context, foodIDs []uuid.UUID) (map[uuid.UUID][]domain.FoodPortion, error) {
	result := make(map[uuid.UUID][]domain.FoodPortion)
	if len(foodIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.Pool.Query(ctx, `
		SELECT food_id, name, grams
		FROM food_portions
		WHERE food_id = ANY($1::uuid[])
		ORDER BY name ASC
	`, foodIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var foodID uuid.UUID
		var portion domain.FoodPortion
		if err := rows.Scan(&foodID, &portion.Name, &portion.Grams); err != nil {
			return nil, err
		}
		result[foodID] = append(result[foodID], portion)
	}

	return result, rows.Err()
}

func (r *FoodRepository) replaceFoodPortions(ctx context.Context, foodID uuid.UUID, portions []domain.FoodPortion) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM food_portions WHERE food_id = $1`, foodID)
	if err != nil {
		return err
	}

	for _, portion := range portions {
		name := strings.TrimSpace(portion.Name)
		if name == "" || strings.EqualFold(name, "g") || portion.Grams <= 0 {
			continue
		}
		if _, err := r.db.Pool.Exec(ctx, `
			INSERT INTO food_portions (food_id, name, grams)
			VALUES ($1, $2, $3)
			ON CONFLICT (food_id, name)
			DO UPDATE SET grams = EXCLUDED.grams
		`, foodID, name, portion.Grams); err != nil {
			return fmt.Errorf("insert food portion: %w", err)
		}
	}

	return nil
}

func defaultPortions() []domain.FoodPortion {
	return []domain.FoodPortion{{Name: "g", Grams: 1}}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
