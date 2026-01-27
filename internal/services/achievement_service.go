package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
)

type AchievementService struct {
	achievementRepo *repositories.AchievementRepository
	statsRepo       *repositories.StatsRepository
	logger          domain.Logger
}

func NewAchievementService(achievementRepo *repositories.AchievementRepository, statsRepo *repositories.StatsRepository, logger domain.Logger) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		statsRepo:       statsRepo,
		logger:          logger,
	}
}

func (s *AchievementService) GetUserAchievements(ctx context.Context, userID uuid.UUID) (*domain.AchievementsResponse, error) {
	achievements, err := s.achievementRepo.GetUserAchievements(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user achievements", "user_id", userID.String(), "error", err)
		return nil, err
	}

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.Warn("Failed to get user stats", "user_id", userID.String(), "error", err)
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
		s.logger.Error("Failed to get all achievements", "error", err)
		return nil, err
	}

	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.Warn("Failed to get user stats", "user_id", userID.String(), "error", err)
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
			s.logger.Error("Failed to upsert user achievement", "user_id", userID.String(), "achievement_id", achievement.ID.String(), "error", err)
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
