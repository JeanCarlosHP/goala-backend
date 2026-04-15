package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"
)

type UserService struct {
	userRepo  *repositories.UserRepository
	goalRepo  *repositories.GoalRepository
	cdnDomain string
}

func NewUserService(userRepo *repositories.UserRepository, goalRepo *repositories.GoalRepository, cdnDomain string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		goalRepo:  goalRepo,
		cdnDomain: cdnDomain,
	}
}

func (s *UserService) buildPhotoURL(ctx context.Context, photoPath *string) *string {
	tr := otel.Tracer("services/user_service.go")
	_, span := tr.Start(ctx, "buildPhotoURL")
	defer span.End()

	if photoPath == nil || *photoPath == "" {
		return nil
	}

	if strings.HasPrefix(*photoPath, "http://") || strings.HasPrefix(*photoPath, "https://") {
		externalURL := *photoPath
		return &externalURL
	}

	baseURL := strings.TrimRight(s.cdnDomain, "/")
	path := "/" + strings.TrimLeft(*photoPath, "/")
	fullURL := baseURL + path
	return &fullURL
}

func (s *UserService) RegisterUser(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "RegisterUser")
	defer span.End()

	exists, err := s.userRepo.ExistsByFirebaseUID(ctx, req.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		ID:          uuid.New(),
		FirebaseUID: req.FirebaseUID,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		PhotoURL:    req.PhotoURL,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	defaultGoal := &domain.UserGoal{
		UserID:           user.ID,
		DailyCalorieGoal: 2000,
		DailyProteinGoal: 150,
		DailyCarbsGoal:   200,
		DailyFatGoal:     65,
	}

	if err := s.goalRepo.Upsert(ctx, defaultGoal); err != nil {
		return nil, fmt.Errorf("failed to create default goal: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "GetUserByFirebaseUID")
	defer span.End()

	return s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
}

func (s *UserService) GetUserGoal(ctx context.Context, userID uuid.UUID) (*domain.UserGoal, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "GetUserGoal")
	defer span.End()

	return s.goalRepo.GetByUserID(ctx, userID)
}

func (s *UserService) UpdateUserGoal(ctx context.Context, userID uuid.UUID, req domain.UpdateGoalRequest) (*domain.UserProfileResponse, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "UpdateUserGoal")
	defer span.End()

	goal := &domain.UserGoal{
		UserID:           userID,
		DailyCalorieGoal: req.DailyCalorieGoal,
		DailyProteinGoal: req.DailyProteinGoal,
		DailyCarbsGoal:   req.DailyCarbsGoal,
		DailyFatGoal:     req.DailyFatGoal,
	}

	if err := s.goalRepo.Upsert(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return s.GetUserProfile(ctx, userID)
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfileResponse, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "GetUserProfile")
	defer span.End()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	goal, err := s.goalRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	carbsGoal := int32(goal.DailyCarbsGoal)
	fatGoal := int32(goal.DailyFatGoal)

	return &domain.UserProfileResponse{
		ID:                   user.ID.String(),
		DisplayName:          user.DisplayName,
		Email:                user.Email,
		Photo:                s.buildPhotoURL(ctx, user.PhotoURL),
		DailyCalorieGoal:     int32(goal.DailyCalorieGoal),
		DailyProteinGoal:     int32(goal.DailyProteinGoal),
		DailyCarbsGoal:       &carbsGoal,
		DailyFatGoal:         &fatGoal,
		Weight:               user.Weight,
		Height:               user.Height,
		Age:                  user.Age,
		Gender:               user.Gender,
		ActivityLevel:        user.ActivityLevel,
		Language:             user.Language,
		Timezone:             user.Timezone,
		NotificationsEnabled: user.NotificationsEnabled,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileRequest) (*domain.UserProfileResponse, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "UpdateUserProfile")
	defer span.End()

	user := &domain.User{
		ID:                   userID,
		DisplayName:          req.DisplayName,
		Email:                req.Email,
		PhotoURL:             req.Photo,
		Weight:               req.Weight,
		Height:               req.Height,
		Age:                  req.Age,
		Gender:               req.Gender,
		ActivityLevel:        req.ActivityLevel,
		Language:             req.Language,
		Timezone:             req.Timezone,
		NotificationsEnabled: req.NotificationsEnabled,
	}

	if err := s.userRepo.UpdateProfile(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	proteinGoalInt := int(0)
	carbsGoalInt := int(0)
	fatGoalInt := int(0)
	if req.DailyProteinGoal != nil {
		proteinGoalInt = int(*req.DailyProteinGoal)
	}
	if req.DailyCarbsGoal != nil {
		carbsGoalInt = int(*req.DailyCarbsGoal)
	}
	if req.DailyFatGoal != nil {
		fatGoalInt = int(*req.DailyFatGoal)
	}

	goal := &domain.UserGoal{
		UserID:           userID,
		DailyCalorieGoal: int(req.DailyCalorieGoal),
		DailyProteinGoal: proteinGoalInt,
		DailyCarbsGoal:   carbsGoalInt,
		DailyFatGoal:     fatGoalInt,
	}

	if err := s.goalRepo.Upsert(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return s.GetUserProfile(ctx, userID)
}

func (s *UserService) PatchUserPreferences(ctx context.Context, userID uuid.UUID, req domain.PatchUserPreferencesRequest) (*domain.UserProfileResponse, error) {
	tr := otel.Tracer("services/user_service.go")
	ctx, span := tr.Start(ctx, "PatchUserPreferences")
	defer span.End()

	if req.DisplayName != nil {
		if err := s.userRepo.UpdateDisplayName(ctx, userID, req.DisplayName); err != nil {
			return nil, fmt.Errorf("failed to update display name: %w", err)
		}
	}

	if req.PhotoURL != nil {
		if err := s.userRepo.UpdateAvatar(ctx, userID, req.PhotoURL); err != nil {
			return nil, fmt.Errorf("failed to update avatar: %w", err)
		}
	}

	if req.NotificationsEnabled != nil {
		if err := s.userRepo.UpdateNotifications(ctx, userID, req.NotificationsEnabled); err != nil {
			return nil, fmt.Errorf("failed to update notifications: %w", err)
		}
	}

	return s.GetUserProfile(ctx, userID)
}
