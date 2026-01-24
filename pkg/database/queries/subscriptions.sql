-- name: CreateSubscription :one
INSERT INTO subscriptions (
    user_id,
    revenuecat_user_id,
    revenuecat_original_transaction_id,
    is_active,
    plan,
    is_trial,
    current_period_start,
    current_period_end,
    last_event_id,
    last_event_type,
    last_event_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetSubscriptionByUserID :one
SELECT * FROM subscriptions
WHERE user_id = $1 LIMIT 1;

-- name: GetSubscriptionByRevenueCatUserID :one
SELECT * FROM subscriptions
WHERE revenuecat_user_id = $1 LIMIT 1;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET
    is_active = $2,
    plan = $3,
    is_trial = $4,
    current_period_start = $5,
    current_period_end = $6,
    last_event_id = $7,
    last_event_type = $8,
    last_event_at = $9,
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: UpsertSubscription :one
INSERT INTO subscriptions (
    user_id,
    revenuecat_user_id,
    revenuecat_original_transaction_id,
    is_active,
    plan,
    is_trial,
    current_period_start,
    current_period_end,
    last_event_id,
    last_event_type,
    last_event_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (user_id) DO UPDATE SET
    is_active = EXCLUDED.is_active,
    plan = EXCLUDED.plan,
    is_trial = EXCLUDED.is_trial,
    current_period_start = EXCLUDED.current_period_start,
    current_period_end = EXCLUDED.current_period_end,
    last_event_id = EXCLUDED.last_event_id,
    last_event_type = EXCLUDED.last_event_type,
    last_event_at = EXCLUDED.last_event_at,
    revenuecat_original_transaction_id = EXCLUDED.revenuecat_original_transaction_id,
    updated_at = NOW()
RETURNING *;

-- name: CheckEventProcessed :one
SELECT EXISTS(
    SELECT 1 FROM subscriptions
    WHERE last_event_id = $1
);

-- name: ListActiveSubscriptions :many
SELECT * FROM subscriptions
WHERE is_active = true
ORDER BY created_at DESC;

-- name: ListExpiredSubscriptions :many
SELECT * FROM subscriptions
WHERE is_active = true
AND current_period_end < NOW()
ORDER BY current_period_end ASC;
