package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type SubscriptionRepository struct {
	db *database.Database
}

func NewSubscriptionRepository(db *database.Database) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	result, err := r.db.Querier.CreateSubscription(ctx, db.CreateSubscriptionParams{
		UserID:                          stringToPgUUID(sub.UserID),
		RevenuecatUserID:                sub.RevenueCatUserID,
		RevenuecatOriginalTransactionID: stringPtr(sub.RevenueCatOriginalTransactionID),
		IsActive:                        sub.IsActive,
		Plan:                            sub.Plan.String(),
		IsTrial:                         sub.IsTrial,
		CurrentPeriodStart:              toNullTimestamp(sub.CurrentPeriodStart),
		CurrentPeriodEnd:                toNullTimestamp(sub.CurrentPeriodEnd),
		LastEventID:                     sub.LastEventID,
		LastEventType:                   sub.LastEventType,
		LastEventAt:                     toNullTimestamp(sub.LastEventAt),
	})
	if err != nil {
		return nil, err
	}

	return toSubscription(&result), nil
}

func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "GetByUserID")
	defer span.End()

	result, err := r.db.Querier.GetSubscriptionByUserID(ctx, stringToPgUUID(userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toSubscription(&result), nil
}

func (r *SubscriptionRepository) GetByRevenueCatUserID(ctx context.Context, rcUserID string) (*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "GetByRevenueCatUserID")
	defer span.End()

	result, err := r.db.Querier.GetSubscriptionByRevenueCatUserID(ctx, rcUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toSubscription(&result), nil
}

func (r *SubscriptionRepository) Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "Upsert")
	defer span.End()

	result, err := r.db.Querier.UpsertSubscription(ctx, db.UpsertSubscriptionParams{
		UserID:                          stringToPgUUID(sub.UserID),
		RevenuecatUserID:                sub.RevenueCatUserID,
		RevenuecatOriginalTransactionID: stringPtr(sub.RevenueCatOriginalTransactionID),
		IsActive:                        sub.IsActive,
		Plan:                            sub.Plan.String(),
		IsTrial:                         sub.IsTrial,
		CurrentPeriodStart:              toNullTimestamp(sub.CurrentPeriodStart),
		CurrentPeriodEnd:                toNullTimestamp(sub.CurrentPeriodEnd),
		LastEventID:                     sub.LastEventID,
		LastEventType:                   sub.LastEventType,
		LastEventAt:                     toNullTimestamp(sub.LastEventAt),
	})
	if err != nil {
		return nil, err
	}

	return toSubscription(&result), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	result, err := r.db.Querier.UpdateSubscription(ctx, db.UpdateSubscriptionParams{
		UserID:             stringToPgUUID(sub.UserID),
		IsActive:           sub.IsActive,
		Plan:               sub.Plan.String(),
		IsTrial:            sub.IsTrial,
		CurrentPeriodStart: toNullTimestamp(sub.CurrentPeriodStart),
		CurrentPeriodEnd:   toNullTimestamp(sub.CurrentPeriodEnd),
		LastEventID:        sub.LastEventID,
		LastEventType:      sub.LastEventType,
		LastEventAt:        toNullTimestamp(sub.LastEventAt),
	})
	if err != nil {
		return nil, err
	}

	return toSubscription(&result), nil
}

func (r *SubscriptionRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "IsEventProcessed")
	defer span.End()

	exists, err := r.db.Querier.CheckEventProcessed(ctx, stringPtr(eventID))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SubscriptionRepository) ListActive(ctx context.Context) ([]*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "ListActive")
	defer span.End()

	results, err := r.db.Querier.ListActiveSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	subs := make([]*domain.Subscription, len(results))
	for i, result := range results {
		subs[i] = toSubscription(&result)
	}
	return subs, nil
}

func (r *SubscriptionRepository) ListExpired(ctx context.Context) ([]*domain.Subscription, error) {
	tr := otel.Tracer("repositories/subscription_repo.go")
	ctx, span := tr.Start(ctx, "ListExpired")
	defer span.End()

	results, err := r.db.Querier.ListExpiredSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	subs := make([]*domain.Subscription, len(results))
	for i, result := range results {
		subs[i] = toSubscription(&result)
	}
	return subs, nil
}

func toSubscription(s *db.Subscription) *domain.Subscription {
	return &domain.Subscription{
		ID:                              s.ID,
		UserID:                          uuid.UUID(s.UserID.Bytes).String(),
		RevenueCatUserID:                s.RevenuecatUserID,
		RevenueCatOriginalTransactionID: stringValue(s.RevenuecatOriginalTransactionID),
		IsActive:                        s.IsActive,
		Plan:                            enum.SubscriptionPlan(s.Plan),
		IsTrial:                         s.IsTrial,
		CurrentPeriodStart:              fromNullTimestamp(s.CurrentPeriodStart),
		CurrentPeriodEnd:                fromNullTimestamp(s.CurrentPeriodEnd),
		LastEventID:                     s.LastEventID,
		LastEventType:                   s.LastEventType,
		LastEventAt:                     fromNullTimestamp(s.LastEventAt),
		CreatedAt:                       s.CreatedAt.Time,
		UpdatedAt:                       s.UpdatedAt.Time,
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toNullTimestamp(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func fromNullTimestamp(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
