# Notifications v1 Backend Contract

This document is the canonical backend contract for notifications v1. It freezes the JSON shape, persistence decision, migration defaults, and compatibility rules so backend and app work can proceed in parallel.

## Decision Summary

- Keep the existing profile surface: `GET /api/v1/user/profile` and `PATCH /api/v1/user/profile`.
- Keep top-level `notificationsEnabled` as the legacy global kill switch for backward compatibility.
- Add a nested `notificationPreferences` object for v1 category-level settings.
- Keep `timezone` as the source of truth for local-time execution context.
- Normalize reminder time as zero-padded 24-hour `HH:mm`.
- Persist v1 preferences as explicit columns on `users`, not as an opaque JSON blob.

## Canonical Response Shape

`GET /api/v1/user/profile` must continue returning the existing top-level fields and add the stable v1 payload below.

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "displayName": "Jean",
    "email": "jean@example.com",
    "photo": "https://cdn.example.com/avatars/550e8400-e29b-41d4-a716-446655440000/avatar.jpg",
    "dailyCalorieGoal": 2000,
    "dailyProteinGoal": 150,
    "dailyCarbsGoal": 200,
    "dailyFatGoal": 65,
    "weight": 75,
    "height": 175,
    "age": 30,
    "gender": "male",
    "activityLevel": "moderate",
    "language": "pt-BR",
    "timezone": "America/Sao_Paulo",
    "notificationsEnabled": true,
    "notificationPreferences": {
      "dailyReminder": {
        "enabled": true,
        "time": "09:00"
      },
      "streakAtRisk": {
        "enabled": false
      },
      "achievementUnlocked": {
        "enabled": true
      }
    },
    "createdAt": "2026-04-01T00:00:00Z",
    "updatedAt": "2026-04-01T00:00:00Z"
  },
  "message": "profile retrieved successfully"
}
```

## Canonical Patch Shape

`PATCH /api/v1/user/profile` is the preferred partial update path for notifications v1.

Example payload:

```json
{
    "notificationsEnabled": true,
    "notificationPreferences": {
      "dailyReminder": {
        "enabled": true,
        "time": "09:00"
      },
      "streakAtRisk": {
        "enabled": false
    },
    "achievementUnlocked": {
      "enabled": true
    }
  }
}
```

### Patch Rules

- Omitted top-level fields are unchanged.
- Omitted nested objects inside `notificationPreferences` are unchanged.
- Omitted fields inside a provided nested object are unchanged.
- `notificationsEnabled` remains independently writable and acts as the global kill switch.
- `notificationPreferences` stores category intent even when `notificationsEnabled` is `false`.
- When `notificationsEnabled` is `false`, all categories are effectively disabled for app behavior regardless of stored per-category flags.
- `reminderTime` is interpreted in the user's existing `timezone`; no timezone is stored inside `notificationPreferences`.

## Validation Rules

- `timezone` remains the existing IANA timezone string already stored on the user profile.
- `notificationPreferences.dailyReminder.time` must match `HH:mm` in 24-hour format.
- `notificationPreferences.dailyReminder.time` is required when that field is sent; omitted means "leave unchanged".
- `notificationPreferences.streakAtRisk.enabled` is required when `streakAtRisk` is sent.
- `notificationPreferences.achievementUnlocked.enabled` is required when `achievementUnlocked` is sent.

Valid `time` examples:

- `00:00`
- `08:30`
- `23:59`

Invalid `time` examples:

- `8:30`
- `24:00`
- `25:99`
- `08:30:00`

## Persistence Decision

Persist v1 preferences on the existing `users` row with explicit columns:

- `notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE` (existing)
- `timezone VARCHAR(50) NOT NULL DEFAULT 'UTC'` (existing)
- `notification_daily_reminder_enabled BOOLEAN NOT NULL DEFAULT FALSE`
- `notification_daily_reminder_time VARCHAR(5) NOT NULL DEFAULT '09:00'`
- `notification_streak_at_risk_enabled BOOLEAN NOT NULL DEFAULT FALSE`
- `notification_achievement_unlocked_enabled BOOLEAN NOT NULL DEFAULT FALSE`

### Why explicit columns

- The v1 domain is still small and fixed.
- Validation stays strict and simple in Go and SQL.
- `GET /user/profile` stays on the current bootstrap path without an extra join.
- Future schema refactors can happen behind the same JSON contract if phase 2 needs a dedicated notification model.

## Migration and Default Rules

Existing users must be migrated without silently opting them into new categories.

### Backfill rules

- Keep the current `notifications_enabled` value unchanged.
- Set `notification_daily_reminder_enabled = notifications_enabled`.
- Set `notification_daily_reminder_time = '09:00'`.
- Set `notification_streak_at_risk_enabled = notifications_enabled`.
- Set `notification_achievement_unlocked_enabled = notifications_enabled`.

### Why this default set

- It preserves the legacy global toggle.
- It preserves the legacy on/off behavior across the initial v1 categories.
- It returns a fully materialized, stable shape for existing users immediately after migration.
- It avoids null handling in the app bootstrap path.

### Materialized response for migrated users

After the migration, `GET /api/v1/user/profile` must always return `notificationPreferences`, even for users created before the new columns existed.

If a migrated user previously had `notificationsEnabled = true`, the response shape must look like this:

```json
{
  "timezone": "America/Sao_Paulo",
  "notificationsEnabled": true,
  "notificationPreferences": {
    "dailyReminder": {
      "enabled": true,
      "time": "09:00"
    },
    "streakAtRisk": {
      "enabled": true
    },
    "achievementUnlocked": {
      "enabled": true
    }
  }
}
```

## Compatibility Rules

- Do not remove top-level `notificationsEnabled` in v1.
- Old clients may continue sending only `notificationsEnabled`; that must remain supported.
- New clients should read and write `notificationPreferences` while still respecting top-level `notificationsEnabled`.
- `GET /api/v1/user/profile` remains the one-read bootstrap surface for notification settings.
- `PATCH /api/v1/user/profile` is merge-style, not replace-style, for nested notification updates.
- Internal persistence changes must not change the external JSON field names defined here.

## Parallel Work Guidance

- Backend implementation can land later as long as it preserves this JSON contract.
- App work can proceed immediately using this exact shape for types, local storage, and scheduling logic.
- If repository internals change during implementation, keep these field names stable and confine changes to handlers, services, repositories, and migrations.

## Implementation Notes for GOA-6

- Extend `internal/domain/user.go` with the nested response/request types defined here.
- Update `internal/handlers/user.go` and `internal/services/user_service.go` so `PATCH /user/profile` accepts partial nested notification updates.
- Add migration files under `pkg/database/migrations` and sqlc query updates under `pkg/database/queries/users.sql`.
- Add validation tests for invalid `HH:mm`, omitted nested fields, and migrated users with legacy `notifications_enabled` data.
