package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"github.com/rs/zerolog/log"
)

type AchievementService struct {
	achievementRepo *repositories.AchievementRepository
	statsRepo       *repositories.StatsRepository
}

func NewAchievementService(achievementRepo *repositories.AchievementRepository, statsRepo *repositories.StatsRepository) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		statsRepo:       statsRepo,
	}
}

func (s *AchievementService) GetUserAchievements(ctx context.Context, userID uuid.UUID) (*domain.AchievementsResponse, error) {
	achievements, err := s.achievementRepo.GetUserAchievements(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get user achievements")
		return nil, err
	}

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to get user stats, using defaults")
		stats = &domain.UserStats{}
	}

	var avgCaloriesPerDay int32
	if stats.TotalDaysLogged > 0 {
		avgCaloriesPerDay = stats.TotalCaloriesConsumed / stats.TotalDaysLogged
	}

	statsResponse := domain.UserStatsResponse{
		CurrentStreak:         stats.CurrentStreak,
		BestStreak:            stats.LongestStreak,
		TotalMealsLogged:      stats.TotalMealsLogged,
		TotalCaloriesLogged:   stats.TotalCaloriesConsumed,
		TotalDaysLogged:       stats.TotalDaysLogged,
		AverageCaloriesPerDay: avgCaloriesPerDay,
	}

	response := &domain.AchievementsResponse{
		Achievements: achievements,
		Stats:        statsResponse,
	}

	return response, nil
}

func (s *AchievementService) SyncAchievements(ctx context.Context, userID uuid.UUID) (*domain.AchievementsResponse, error) {
	allAchievements, err := s.achievementRepo.GetAllAchievements(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all achievements")
		return nil, err
	}

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to get user stats")
		stats = &domain.UserStats{}
	}

	for _, achievement := range allAchievements {
		progress, unlocked := s.calculateAchievementProgress(achievement, stats)

		var unlockedAt *time.Time
		if unlocked {
			now := time.Now()
			unlockedAt = &now
		}

		err := s.achievementRepo.UpsertUserAchievement(ctx, userID, achievement.ID, unlocked, progress, unlockedAt)
		if err != nil {
			log.Error().Err(err).
				Str("user_id", userID.String()).
				Str("achievement_id", achievement.ID.String()).
				Msg("Failed to upsert user achievement")
		}
	}

	return s.GetUserAchievements(ctx, userID)
}

func (s *AchievementService) calculateAchievementProgress(achievement domain.Achievement, stats *domain.UserStats) (int32, bool) {
	var progress int32

	switch achievement.ID.String() {
	case "first_meal":
		if stats.TotalMealsLogged >= 1 {
			progress = 1
		}
	case "streak_3":
		progress = stats.CurrentStreak
		if progress > achievement.Target {
			progress = achievement.Target
		}
	case "streak_7":
		progress = stats.CurrentStreak
		if progress > achievement.Target {
			progress = achievement.Target
		}
	case "streak_30":
		progress = stats.CurrentStreak
		if progress > achievement.Target {
			progress = achievement.Target
		}
	case "meals_10":
		progress = stats.TotalMealsLogged
		if progress > achievement.Target {
			progress = achievement.Target
		}
	case "meals_50":
		progress = stats.TotalMealsLogged
		if progress > achievement.Target {
			progress = achievement.Target
		}
	case "meals_100":
		progress = stats.TotalMealsLogged
		if progress > achievement.Target {
			progress = achievement.Target
		}
	default:
		progress = 0
	}

	unlocked := progress >= achievement.Target
	return progress, unlocked
}
