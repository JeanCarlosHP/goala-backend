package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"github.com/rs/zerolog/log"
)

type StatsService struct {
	statsRepo *repositories.StatsRepository
	mealRepo  *repositories.MealRepository
}

func NewStatsService(statsRepo *repositories.StatsRepository, mealRepo *repositories.MealRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		mealRepo:  mealRepo,
	}
}

func (s *StatsService) GetUserStats(ctx context.Context, userID uuid.UUID) (*domain.UserStatsResponse, error) {
	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get user stats")
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
			log.Warn().Err(err).Time("date", currentDate).Msg("Failed to get meals for date")
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		var totalCalories, totalProtein, totalCarbs, totalFat int32
		mealResponses := make([]domain.Meal, 0)

		for _, meal := range meals {
			mealCalories := 0
			mealProtein := 0.0
			mealCarbs := 0.0
			mealFat := 0.0

			for _, food := range meal.Foods {
				mealCalories += food.Calories
				mealProtein += food.ProteinG
				mealCarbs += food.CarbsG
				mealFat += food.FatG
			}

			totalCalories += int32(mealCalories)
			totalProtein += int32(mealProtein)
			totalCarbs += int32(mealCarbs)
			totalFat += int32(mealFat)

			mealResponses = append(mealResponses, meal)
		}

		dayStats := domain.DayStats{
			Date:          currentDate,
			TotalCalories: totalCalories,
			TotalProtein:  totalProtein,
			TotalCarbs:    totalCarbs,
			TotalFat:      totalFat,
			Meals:         mealResponses,
			WaterIntake:   0,
		}

		allDays = append(allDays, dayStats)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	total := len(allDays)
	totalPages := (total + limit - 1) / limit

	offset := (page - 1) * limit
	end := offset + limit
	if end > total {
		end = total
	}
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
