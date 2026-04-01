package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type User struct {
	ID                      uuid.UUID               `json:"id"`
	FirebaseUID             string                  `json:"firebaseUid"`
	Email                   string                  `json:"email"`
	DisplayName             string                  `json:"displayName"`
	PhotoURL                *string                 `json:"photoUrl,omitempty"`
	Weight                  *int32                  `json:"weight,omitempty"`
	Height                  *int32                  `json:"height,omitempty"`
	Age                     *int32                  `json:"age,omitempty"`
	Gender                  *string                 `json:"gender,omitempty"`
	ActivityLevel           *string                 `json:"activityLevel,omitempty"`
	Language                string                  `json:"language"`
	Timezone                string                  `json:"timezone"` // Novo campo para timezone
	NotificationsEnabled    bool                    `json:"notificationsEnabled"`
	NotificationPreferences NotificationPreferences `json:"notificationPreferences"`
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
}

type UserGoal struct {
	UserID           uuid.UUID `json:"userId"`
	DailyCalorieGoal int       `json:"dailyCalorieGoal"`
	DailyProteinGoal int       `json:"dailyProteinGoal"`
	DailyCarbsGoal   int       `json:"dailyCarbsGoal"`
	DailyFatGoal     int       `json:"dailyFatGoal"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type RegisterRequest struct {
	FirebaseUID string  `json:"firebaseUid"`
	Email       string  `json:"email"`
	DisplayName string  `json:"displayName"`
	PhotoURL    *string `json:"photoUrl,omitempty"`
}

type UpdateGoalRequest struct {
	DailyCalorieGoal int `json:"dailyCalorieGoal" validate:"required,min=500,max=10000"`
	DailyProteinGoal int `json:"dailyProteinGoal" validate:"required,min=0"`
	DailyCarbsGoal   int `json:"dailyCarbsGoal" validate:"required,min=0"`
	DailyFatGoal     int `json:"dailyFatGoal" validate:"required,min=0"`
}

type UserProfileResponse struct {
	ID                      string                  `json:"id"`
	DisplayName             string                  `json:"displayName"`
	Email                   string                  `json:"email"`
	Photo                   *string                 `json:"photo,omitempty"`
	DailyCalorieGoal        int32                   `json:"dailyCalorieGoal"`
	DailyProteinGoal        int32                   `json:"dailyProteinGoal"`
	DailyCarbsGoal          *int32                  `json:"dailyCarbsGoal,omitempty"`
	DailyFatGoal            *int32                  `json:"dailyFatGoal,omitempty"`
	Weight                  *int32                  `json:"weight,omitempty"`
	Height                  *int32                  `json:"height,omitempty"`
	Age                     *int32                  `json:"age,omitempty"`
	Gender                  *string                 `json:"gender,omitempty"`
	ActivityLevel           *string                 `json:"activityLevel,omitempty"`
	Language                string                  `json:"language"`
	Timezone                string                  `json:"timezone"` // Novo campo
	NotificationsEnabled    bool                    `json:"notificationsEnabled"`
	NotificationPreferences NotificationPreferences `json:"notificationPreferences"`
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
}

type UpdateProfileRequest struct {
	ID                   string  `json:"id" validate:"required,uuid"`
	DisplayName          string  `json:"displayName" validate:"required,min=2,max=255"`
	Email                string  `json:"email" validate:"required,email"`
	Photo                *string `json:"photo" validate:"omitempty,startswith=/"`
	DailyCalorieGoal     int32   `json:"dailyCalorieGoal" validate:"required,gte=0,lte=10000"`
	DailyProteinGoal     *int32  `json:"dailyProteinGoal" validate:"omitempty,gte=0,lte=1000"`
	DailyCarbsGoal       *int32  `json:"dailyCarbsGoal" validate:"omitempty,gte=0,lte=2000"`
	DailyFatGoal         *int32  `json:"dailyFatGoal" validate:"omitempty,gte=0,lte=1000"`
	Weight               *int32  `json:"weight" validate:"omitempty,gt=0,lte=1000"`
	Height               *int32  `json:"height" validate:"omitempty,gt=0,lte=300"`
	Age                  *int32  `json:"age" validate:"omitempty,gt=0,lte=150"`
	Gender               *string `json:"gender" validate:"omitempty,oneof=male female other"`
	ActivityLevel        *string `json:"activityLevel" validate:"omitempty,oneof=sedentary light moderate active very_active"`
	Language             string  `json:"language" validate:"required,oneof=en-US pt-BR"`
	Timezone             string  `json:"timezone" validate:"required"` // Novo campo
	NotificationsEnabled bool    `json:"notificationsEnabled"`
}

type AvatarUploadRequest struct {
	ContentType string `json:"contentType" validate:"required,oneof=image/jpeg image/jpg image/png image/webp"`
	FileSize    int64  `json:"fileSize" validate:"required,gt=0,lte=5242880"`
}

type AvatarUploadResponse struct {
	UploadURL string `json:"uploadUrl"`
	PhotoURL  string `json:"photoUrl"`
	ExpiresIn int    `json:"expiresIn"`
}

type PatchUserPreferencesRequest struct {
	DisplayName             *string                              `json:"displayName" validate:"omitempty,min=2,max=255"`
	PhotoURL                *string                              `json:"photoUrl" validate:"omitempty,startswith=/"`
	NotificationsEnabled    *bool                                `json:"notificationsEnabled" validate:"omitempty"`
	NotificationPreferences *PatchNotificationPreferencesRequest `json:"notificationPreferences" validate:"omitempty"`
}
