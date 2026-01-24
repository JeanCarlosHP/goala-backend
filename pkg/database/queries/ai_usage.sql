-- name: IncrementAIUsage :one
INSERT INTO ai_usage (
    user_id,
    feature,
    usage_count,
    quota,
    period_start,
    period_end
) VALUES (
    $1, $2, 1, $3, $4, $5
)
ON CONFLICT (user_id, feature, period_start) DO UPDATE SET
    usage_count = ai_usage.usage_count + 1,
    updated_at = NOW()
RETURNING *;

-- name: GetAIUsage :one
SELECT * FROM ai_usage
WHERE user_id = $1
AND feature = $2
AND period_end > NOW()
ORDER BY period_start DESC
LIMIT 1;

-- name: GetAIUsageByPeriod :one
SELECT * FROM ai_usage
WHERE user_id = $1
AND feature = $2
AND period_start = $3
LIMIT 1;

-- name: ResetAIUsage :exec
DELETE FROM ai_usage
WHERE period_end < NOW();

-- name: ListUserAIUsage :many
SELECT * FROM ai_usage
WHERE user_id = $1
AND period_end > NOW()
ORDER BY feature ASC;

-- name: CreateOrResetAIUsage :one
INSERT INTO ai_usage (
    user_id,
    feature,
    usage_count,
    quota,
    period_start,
    period_end
) VALUES (
    $1, $2, 0, $3, $4, $5
)
ON CONFLICT (user_id, feature, period_start) DO UPDATE SET
    usage_count = 0,
    quota = EXCLUDED.quota,
    updated_at = NOW()
RETURNING *;
