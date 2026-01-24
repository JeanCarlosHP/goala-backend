# Subscription & AI Usage API Documentation

Complete API reference for subscription management and AI feature usage tracking.

---

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Subscription Endpoints](#subscription-endpoints)
- [AI Usage Endpoints](#ai-usage-endpoints)
- [Webhook Endpoints](#webhook-endpoints)
- [Error Codes](#error-codes)
- [Usage Examples](#usage-examples)

---

## Overview

The subscription system integrates with RevenueCat to manage user subscriptions and enforce AI feature quotas.

### Key Features
- **Subscription Management**: Track user subscription status, plan type, and trial access
- **AI Usage Quotas**: Enforce daily limits for AI-powered features
- **RevenueCat Integration**: Process webhook events for subscription lifecycle
- **Security**: Webhook signature verification and idempotent event processing

### Plan Types

| Plan | Code | Food Recognition | Meal Analysis | Nutrition Advice |
|------|------|------------------|---------------|------------------|
| Free | `free` | 10/day | 5/day | 3/day |
| Trial | `free` (with `is_trial: true`) | 50/day | 25/day | 15/day |
| Monthly | `monthly` | 100/day | 50/day | 30/day |
| Yearly | `yearly` | 100/day | 50/day | 30/day |

---

## Authentication

All protected endpoints require:
- `Authorization: Bearer <firebase_token>` header
- Valid Firebase ID token
- User must exist in the database

The middleware extracts:
- `firebase_uid` from the token
- `user_id` from the database lookup

---

## Subscription Endpoints

### Get Subscription Status

Get the current user's subscription status and access information.

**Endpoint:** `GET /api/v1/subscription/status`

**Authentication:** Required

**Response:**

```json
{
  "success": true,
  "data": {
    "is_active": true,
    "plan": "monthly",
    "is_trial": false,
    "current_period_end": "2026-02-21T10:30:00Z",
    "has_access": true
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `is_active` | boolean | Whether subscription is currently active |
| `plan` | string | Current plan: `free`, `monthly`, or `yearly` |
| `is_trial` | boolean | Whether user is on trial period |
| `current_period_end` | string (ISO 8601) | When current billing period ends (null for free) |
| `has_access` | boolean | Whether user has active access (considers expiration) |

**Example cURL:**

```bash
curl -X GET https://api.calorieai.com/api/v1/subscription/status \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIs..."
```

---

## AI Usage Endpoints

### Get AI Usage Summary

Retrieve all AI feature usage for the current user.

**Endpoint:** `GET /api/v1/ai/usage`

**Authentication:** Required

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "feature": "food_recognition",
      "used": 15,
      "quota": 100,
      "remaining": 85
    },
    {
      "feature": "meal_analysis",
      "used": 3,
      "quota": 50,
      "remaining": 47
    }
  ]
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `feature` | string | AI feature identifier |
| `used` | integer | Number of times used in current period |
| `quota` | integer | Total allowed uses for current period |
| `remaining` | integer | Remaining quota (`quota - used`) |

**Example cURL:**

```bash
curl -X GET https://api.calorieai.com/api/v1/ai/usage \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIs..."
```

---

### Check Feature Quota

Check quota status for a specific AI feature.

**Endpoint:** `GET /api/v1/ai/usage/:feature`

**Authentication:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `feature` | string | One of: `food_recognition`, `meal_analysis`, `nutrition_advice` |

**Response:**

```json
{
  "success": true,
  "data": {
    "has_quota": true,
    "used": 15,
    "quota": 100,
    "remaining": 85
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `has_quota` | boolean | Whether user has remaining quota |
| `used` | integer | Number of times used |
| `quota` | integer | Total quota |
| `remaining` | integer | Remaining uses |

**Error Response (Invalid Feature):**

```json
{
  "success": false,
  "message": "invalid feature"
}
```

**Example cURL:**

```bash
curl -X GET https://api.calorieai.com/api/v1/ai/usage/food_recognition \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIs..."
```

---

## Webhook Endpoints

### RevenueCat Webhook

Process subscription events from RevenueCat.

**Endpoint:** `POST /api/v1/webhooks/revenuecat`

**Authentication:** Webhook signature verification

**Headers:**

| Header | Description |
|--------|-------------|
| `X-Revenuecat-Signature` | HMAC-SHA256 signature of request body |
| `Authorization` | Alternative header for signature (if X-Revenuecat-Signature not present) |

**Request Body:**

```json
{
  "api_version": "1.0",
  "event": {
    "id": "unique-event-id-12345",
    "type": "INITIAL_PURCHASE",
    "app_user_id": "user_firebase_uid",
    "original_app_user_id": "user_firebase_uid",
    "product_id": "monthly_premium",
    "entitlement_id": "premium",
    "period_type": "normal",
    "purchased_at_ms": 1737453000000,
    "expiration_at_ms": 1740045000000,
    "environment": "PRODUCTION",
    "is_trial_period": false,
    "original_transaction_id": "GPA.1234-5678-9012-34567",
    "transaction_id": "GPA.1234-5678-9012-34567"
  }
}
```

**Supported Event Types:**

| Event Type | Description | Subscription State |
|------------|-------------|-------------------|
| `INITIAL_PURCHASE` | First subscription purchase | Active |
| `RENEWAL` | Subscription renewed | Active |
| `CANCELLATION` | User canceled subscription | Inactive |
| `EXPIRATION` | Subscription expired | Inactive |
| `BILLING_ISSUE` | Payment failed | Inactive |
| `UNCANCELLATION` | User reactivated subscription | Active |
| `PRODUCT_CHANGE` | User changed plan | Active |

**Response:**

```json
{
  "success": true,
  "message": "webhook processed"
}
```

**Error Responses:**

Invalid signature (401):
```json
{
  "success": false,
  "message": "invalid signature"
}
```

Invalid payload (400):
```json
{
  "success": false,
  "message": "invalid webhook format"
}
```

Processing error (500):
```json
{
  "success": false,
  "message": "failed to process event"
}
```

**Webhook Processing:**

1. Verify signature using `REVENUECAT_WEBHOOK_SECRET`
2. Check if event already processed (idempotency)
3. Map product ID to subscription plan
4. Update subscription in database
5. Return success response

**Example cURL (for testing):**

```bash
curl -X POST https://api.calorieai.com/api/v1/webhooks/revenuecat \
  -H "Content-Type: application/json" \
  -H "X-Revenuecat-Signature: sha256=abc123..." \
  -d '{
    "api_version": "1.0",
    "event": {
      "id": "evt_test_123",
      "type": "INITIAL_PURCHASE",
      "app_user_id": "user_123",
      "product_id": "monthly_premium",
      "purchased_at_ms": 1737453000000,
      "expiration_at_ms": 1740045000000
    }
  }'
```

---

## Error Codes

### HTTP Status Codes

| Code | Description | When It Occurs |
|------|-------------|----------------|
| `200` | Success | Request completed successfully |
| `400` | Bad Request | Invalid input, malformed JSON, or invalid parameters |
| `401` | Unauthorized | Missing or invalid authentication token |
| `402` | Payment Required | Subscription inactive or quota exceeded |
| `403` | Forbidden | User does not have permission |
| `404` | Not Found | Resource not found |
| `500` | Internal Server Error | Server-side error |

### Application Error Codes

| Code | Message | Description |
|------|---------|-------------|
| `SUBSCRIPTION_REQUIRED` | "active subscription required" | User subscription is inactive or expired |
| `QUOTA_EXCEEDED` | "quota exceeded for this feature" | User has reached AI usage limit |

**Quota Exceeded Response:**

```json
{
  "success": false,
  "message": "quota exceeded for this feature",
  "code": "QUOTA_EXCEEDED",
  "feature": "food_recognition"
}
```

**Subscription Required Response:**

```json
{
  "success": false,
  "message": "active subscription required",
  "code": "SUBSCRIPTION_REQUIRED"
}
```

---

## Usage Examples

### Frontend Integration

#### Check Subscription Before Showing Premium Features

```typescript
async function checkSubscription(): Promise<boolean> {
  const response = await fetch('https://api.calorieai.com/api/v1/subscription/status', {
    headers: {
      'Authorization': `Bearer ${firebaseToken}`
    }
  });
  
  const data = await response.json();
  return data.data.has_access;
}
```

#### Handle Quota Exceeded Error

```typescript
async function recognizeFood(imageUrl: string) {
  try {
    const response = await fetch('https://api.calorieai.com/api/v1/food/recognize', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${firebaseToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ image_url: imageUrl })
    });

    if (response.status === 402) {
      const error = await response.json();
      if (error.code === 'QUOTA_EXCEEDED') {
        // Show upgrade prompt
        showUpgradeDialog('You have reached your daily limit for food recognition');
        return;
      }
      if (error.code === 'SUBSCRIPTION_REQUIRED') {
        // Show subscription prompt
        showSubscriptionDialog('This feature requires an active subscription');
        return;
      }
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to recognize food:', error);
  }
}
```

#### Display Usage Stats

```typescript
async function displayUsageStats() {
  const response = await fetch('https://api.calorieai.com/api/v1/ai/usage', {
    headers: {
      'Authorization': `Bearer ${firebaseToken}`
    }
  });
  
  const { data } = await response.json();
  
  data.forEach(usage => {
    console.log(`${usage.feature}: ${usage.used}/${usage.quota} (${usage.remaining} remaining)`);
  });
}
```

### RevenueCat Setup

#### Configure Webhook in RevenueCat Dashboard

1. Go to RevenueCat Dashboard → Project Settings → Webhooks
2. Add new webhook:
   - URL: `https://your-api.com/api/v1/webhooks/revenuecat`
   - Copy the generated webhook secret
3. Set `REVENUECAT_WEBHOOK_SECRET` environment variable
4. Test webhook with sample event

#### Product ID Mapping

Configure these product IDs in RevenueCat:

| Product ID | Maps To |
|------------|---------|
| `monthly_premium` or `premium_monthly` | `monthly` plan |
| `yearly_premium` or `premium_yearly` | `yearly` plan |

---

## Notes

### Quota Reset Logic

- Quotas reset based on the user's billing cycle
- Free users: reset daily at midnight UTC
- Paid users: reset at `current_period_start` (beginning of billing cycle)
- Trial users: reset daily at midnight UTC

### Idempotency

- Webhook events are processed idempotently using `event.id`
- Duplicate webhooks are safely ignored
- Check `last_event_id` in database before processing

### Security Best Practices

- Always verify webhook signatures
- Use HTTPS for webhook endpoints
- Rotate `REVENUECAT_WEBHOOK_SECRET` periodically
- Monitor webhook failures
- Implement rate limiting for webhook endpoint

### Testing

- Use RevenueCat sandbox environment for testing
- Test all event types: purchase, renewal, cancellation, expiration
- Verify quota enforcement with different plans
- Test webhook signature verification
- Simulate quota exceeded scenarios

---

## Support

For issues or questions:
- Backend API errors: Check application logs
- Subscription sync issues: Verify webhook delivery in RevenueCat dashboard
- Quota calculation issues: Check `ai_usage` table in database
