package domain

import (
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain/enum"
)

type DeviceInfo struct {
	Platform   string `json:"platform" validate:"required"`
	OsVersion  string `json:"osVersion" validate:"required"`
	AppVersion string `json:"appVersion" validate:"required"`
}

type CreateFeedbackRequest struct {
	Type        string      `json:"type" validate:"required,oneof=problem improvement"`
	Title       string      `json:"title" validate:"required,min=3,max=255"`
	Description string      `json:"description" validate:"required,min=10,max=5000"`
	UserEmail   string      `json:"userEmail" validate:"required,email"`
	DeviceInfo  *DeviceInfo `json:"deviceInfo" validate:"omitempty"`
}

type Feedback struct {
	ID          string            `json:"id"`
	UserID      *string           `json:"userId,omitempty"`
	Type        enum.FeedbackType `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	UserEmail   string            `json:"userEmail"`
	Platform    *string           `json:"platform,omitempty"`
	OsVersion   *string           `json:"osVersion,omitempty"`
	AppVersion  *string           `json:"appVersion,omitempty"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
}
