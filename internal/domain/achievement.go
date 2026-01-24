package domain

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID             uuid.UUID `json:"id" db:"id"`
	NameKey        string    `json:"nameKey" db:"name_key" validate:"required"`
	DescriptionKey string    `json:"descriptionKey" db:"description_key" validate:"required"`
	Icon           string    `json:"icon" db:"icon" validate:"required"`
	Target         int32     `json:"target" db:"target" validate:"gte=0"`
	Category       string    `json:"category" db:"category"`
	CreatedAt      time.Time `json:"-" db:"created_at"`
}

type UserAchievement struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID  `json:"achievement_id" db:"achievement_id"`
	Unlocked      bool       `json:"unlocked" db:"unlocked"`
	Progress      int32      `json:"progress" db:"progress" validate:"gte=0"`
	UnlockedAt    *time.Time `json:"unlockedAt,omitempty" db:"unlocked_at"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	UpdatedAt     time.Time  `json:"-" db:"updated_at"`
}

type AchievementResponse struct {
	ID             string     `json:"id" validate:"required"`
	NameKey        string     `json:"nameKey" validate:"required"`
	DescriptionKey string     `json:"descriptionKey" validate:"required"`
	Icon           string     `json:"icon" validate:"required"`
	Unlocked       bool       `json:"unlocked"`
	UnlockedAt     *time.Time `json:"unlockedAt,omitempty" validate:"omitempty"`
	Progress       int32      `json:"progress" validate:"gte=0"`
	Target         int32      `json:"target" validate:"gte=0"`
}

type AchievementsResponse struct {
	Achievements []AchievementResponse `json:"achievements" validate:"required,dive"`
	Stats        UserStatsResponse     `json:"stats"`
}
