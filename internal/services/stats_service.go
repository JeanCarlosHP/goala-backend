package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"
)

type StatsService struct {
	statsRepo *repositories.StatsRepository
	mealRepo  *repositories.MealRepository
	foodRepo  *repositories.FoodRepository
	logger    domain.Logger
}

func NewStatsService(
	statsRepo *repositories.StatsRepository,
	mealRepo *repositories.MealRepository,
	foodRepo *repositories.FoodRepository,
	logger domain.Logger,
) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		mealRepo:  mealRepo,
		foodRepo:  foodRepo,
		logger:    logger,
	}
}

func (s *StatsService) GetUserStats(ctx context.Context, userID uuid.UUID) (*domain.UserStatsResponse, error) {
	tr := otel.Tracer("services/stats_service.go")
	ctx, span := tr.Start(ctx, "GetUserStats")
	defer span.End()

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user stats", "user_id", userID.String(), "error", err)
		return nil, err
	}

	var avgCaloriesPerDay int32
	if stats.TotalDaysLogged > 0 {
		avgCaloriesPerDay = stats.TotalCaloriesConsumed / stats.TotalDaysLogged
	}

	response := &domain.UserStatsResponse{
		CurrentStreak:         stats.CurrentStreak,
		BestStreak:            stats.LongestStreak,
		TotalMealsLogged:      stats.TotalMealsLogged,
		TotalCaloriesLogged:   stats.TotalCaloriesConsumed,
		TotalDaysLogged:       stats.TotalDaysLogged,
		AverageCaloriesPerDay: avgCaloriesPerDay,
	}

	return response, nil
}

func (s *StatsService) GetStatsRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, page, limit int) (*domain.StatsRangeResponse, error) {
	tr := otel.Tracer("services/stats_service.go")
	ctx, span := tr.Start(ctx, "GetStatsRange")
	defer span.End()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	allDays := []domain.DayStats{}
	currentDate := startDate

	for !currentDate.After(endDate) {
		meals, err := s.mealRepo.GetByUserAndDate(ctx, userID, currentDate)
		if err != nil {
			s.logger.Error("Failed to get meals for date", "user_id", userID.String(), "date", currentDate.Format("2006-01-02"), "error", err)
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		mealIDs := make([]uuid.UUID, len(meals))
		for i, meal := range meals {
			mealIDs[i] = meal.ID
		}

		foodItems, err := s.foodRepo.GetByMealIDs(ctx, mealIDs)
		if err != nil {
			s.logger.Error("Failed to get food items for meals", "user_id", userID.String(), "date", currentDate.Format("2006-01-02"), "error", err)
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		for i := range meals {
			if foods, ok := foodItems[meals[i].ID]; ok {
				meals[i].Foods = foods
			}
		}

		var totalCalories, totalProtein, totalCarbs, totalFat int32

		for _, meal := range meals {
			mealCalories := 0
			mealProtein := 0.0
			mealCarbs := 0.0
			mealFat := 0.0

			for _, food := range meal.Foods {
				mealCalories += food.Calories
				mealProtein += food.Protein
				mealCarbs += food.Carbs
				mealFat += food.Fat
			}

			totalCalories += int32(mealCalories)
			totalProtein += int32(mealProtein)
			totalCarbs += int32(mealCarbs)
			totalFat += int32(mealFat)
		}

		dayStats := domain.DayStats{
			Date:          currentDate,
			TotalCalories: totalCalories,
			TotalProtein:  totalProtein,
			TotalCarbs:    totalCarbs,
			TotalFat:      totalFat,
			Meals:         meals,
		}

		allDays = append(allDays, dayStats)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	total := len(allDays)
	totalPages := (total + limit - 1) / limit

	offset := (page - 1) * limit
	end := min(offset+limit, total)
	if offset > total {
		offset = total
	}

	paginatedDays := allDays[offset:end]

	var aggTotalCalories, aggTotalProtein, aggTotalCarbs, aggTotalFat int32
	for _, day := range allDays {
		aggTotalCalories += day.TotalCalories
		aggTotalProtein += day.TotalProtein
		aggTotalCarbs += day.TotalCarbs
		aggTotalFat += day.TotalFat
	}

	var avgCalories, avgProtein, avgCarbs, avgFat int32
	if total > 0 {
		avgCalories = aggTotalCalories / int32(total)
		avgProtein = aggTotalProtein / int32(total)
		avgCarbs = aggTotalCarbs / int32(total)
		avgFat = aggTotalFat / int32(total)
	}

	response := &domain.StatsRangeResponse{
		Days: paginatedDays,
		Pagination: domain.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Aggregated: domain.AggregatedStats{
			TotalCalories: aggTotalCalories,
			TotalProtein:  aggTotalProtein,
			TotalCarbs:    aggTotalCarbs,
			TotalFat:      aggTotalFat,
			AvgCalories:   avgCalories,
			AvgProtein:    avgProtein,
			AvgCarbs:      avgCarbs,
			AvgFat:        avgFat,
		},
	}

	return response, nil
}

func (s *StatsService) UpdateStreakForUser(ctx context.Context, userID uuid.UUID, mealDate time.Time) error {
	tr := otel.Tracer("services/stats_service.go")
	ctx, span := tr.Start(ctx, "UpdateStreakForUser")
	defer span.End()

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		stats, err = s.statsRepo.CreateUserStats(ctx, userID)
		if err != nil {
			return err
		}
	}

	today := time.Date(mealDate.Year(), mealDate.Month(), mealDate.Day(), 0, 0, 0, 0, time.UTC)

	if stats.LastLogDate == nil {
		stats.CurrentStreak = 1
		stats.LastLogDate = &today
	} else {
		lastLog := time.Date(stats.LastLogDate.Year(), stats.LastLogDate.Month(), stats.LastLogDate.Day(), 0, 0, 0, 0, time.UTC)
		daysSinceLastLog := int(today.Sub(lastLog).Hours() / 24)

		if daysSinceLastLog == 0 {
			return nil
		} else if daysSinceLastLog == 1 {
			stats.CurrentStreak++
		} else {
			stats.CurrentStreak = 1
		}
		stats.LastLogDate = &today
	}

	if stats.CurrentStreak > stats.LongestStreak {
		stats.LongestStreak = stats.CurrentStreak
	}

	return s.statsRepo.UpdateStreakAndLastLogDate(ctx, userID, stats.CurrentStreak, today)
}
