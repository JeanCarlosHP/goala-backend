# PATCH User Profile Preferences

This document describes the live partial-update route for user preferences and points to the canonical notifications v1 contract.

## Endpoint

```http
PATCH /api/v1/user/profile
```

## Current Live Fields

The handler currently accepts these partial fields:

- `displayName`
- `photoUrl`
- `notificationsEnabled`
- `notificationPreferences`

Example:

```json
{
  "displayName": "Jean",
  "photoUrl": "/avatars/550e8400-e29b-41d4-a716-446655440000/avatar.jpg",
  "notificationsEnabled": true,
  "notificationPreferences": {
    "dailyReminder": {
      "time": "09:00"
    }
  }
}
```

## Response Shape

The route returns the same profile envelope used by `GET /api/v1/user/profile`:

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "displayName": "Jean",
    "email": "jean@example.com",
    "photo": "https://cdn.example.com/avatars/550e8400-e29b-41d4-a716-446655440000/avatar.jpg",
    "timezone": "America/Sao_Paulo",
    "notificationsEnabled": true,
    "createdAt": "2026-04-01T00:00:00Z",
    "updatedAt": "2026-04-01T00:00:00Z"
  },
  "message": "preferences updated successfully"
}
```

## Notifications v1

The canonical contract for the upcoming nested notification payload lives in:

- `docs/NOTIFICATIONS_V1_BACKEND_CONTRACT.md`

That document freezes:

- `notificationPreferences` field names
- `GET /api/v1/user/profile` response shape
- `PATCH /api/v1/user/profile` merge semantics
- migration defaults for existing users
- the persistence strategy for GOA-6

## Notes

- Older references to `/api/v1/user/preferences` are stale; the live route is `PATCH /api/v1/user/profile`.
- Older examples using `avatar` are stale; the live request field is `photoUrl`.
