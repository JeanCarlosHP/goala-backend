package domain

import (
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain/enum"
)

type Subscription struct {
	ID                              int64                 `json:"id"`
	UserID                          string                `json:"user_id"`
	RevenueCatUserID                string                `json:"revenuecat_user_id"`
	RevenueCatOriginalTransactionID string                `json:"revenuecat_original_transaction_id,omitempty"`
	IsActive                        bool                  `json:"is_active"`
	Plan                            enum.SubscriptionPlan `json:"plan"`
	IsTrial                         bool                  `json:"is_trial"`
	CurrentPeriodStart              *time.Time            `json:"current_period_start,omitempty"`
	CurrentPeriodEnd                *time.Time            `json:"current_period_end,omitempty"`
	LastEventID                     *string               `json:"last_event_id,omitempty"`
	LastEventType                   *string               `json:"last_event_type,omitempty"`
	LastEventAt                     *time.Time            `json:"last_event_at,omitempty"`
	CreatedAt                       time.Time             `json:"created_at"`
	UpdatedAt                       time.Time             `json:"updated_at"`
}

func (s *Subscription) IsExpired() bool {
	if s.CurrentPeriodEnd == nil {
		return false
	}
	return s.CurrentPeriodEnd.Before(time.Now())
}

func (s *Subscription) HasAccess() bool {
	return s.IsActive && !s.IsExpired()
}

type AIUsage struct {
	ID          int64          `json:"id"`
	UserID      string         `json:"user_id"`
	Feature     enum.AIFeature `json:"feature"`
	UsageCount  int32          `json:"usage_count"`
	Quota       int32          `json:"quota"`
	PeriodStart time.Time      `json:"period_start"`
	PeriodEnd   time.Time      `json:"period_end"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (u *AIUsage) RemainingQuota() int32 {
	remaining := u.Quota - u.UsageCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (u *AIUsage) HasQuota() bool {
	return u.UsageCount < u.Quota
}

func (u *AIUsage) IsExpired() bool {
	return u.PeriodEnd.Before(time.Now())
}

type RevenueCatWebhook struct {
	Event             RevenueCatEvent `json:"event"`
	APIVersion        string          `json:"api_version"`
	AppUserID         string          `json:"app_user_id,omitempty"`
	OriginalAppUserID string          `json:"original_app_user_id,omitempty"`
}

type RevenueCatEvent struct {
	ID                    string                   `json:"id"`
	Type                  enum.RevenueCatEventType `json:"type"`
	AppUserID             string                   `json:"app_user_id"`
	OriginalAppUserID     string                   `json:"original_app_user_id"`
	ProductID             string                   `json:"product_id"`
	EntitlementID         string                   `json:"entitlement_id,omitempty"`
	EntitlementIDs        []string                 `json:"entitlement_ids,omitempty"`
	PeriodType            string                   `json:"period_type"`
	PurchasedAtMs         int64                    `json:"purchased_at_ms"`
	ExpirationAtMs        int64                    `json:"expiration_at_ms,omitempty"`
	Environment           string                   `json:"environment"`
	IsTrial               bool                     `json:"is_trial_period"`
	OriginalTransactionID string                   `json:"original_transaction_id"`
	TransactionID         string                   `json:"transaction_id"`
}

func (e *RevenueCatEvent) PurchasedAt() time.Time {
	return time.UnixMilli(e.PurchasedAtMs)
}

func (e *RevenueCatEvent) ExpirationAt() *time.Time {
	if e.ExpirationAtMs == 0 {
		return nil
	}
	t := time.UnixMilli(e.ExpirationAtMs)
	return &t
}
