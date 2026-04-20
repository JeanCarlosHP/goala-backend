#!/bin/sh
set -eu

DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${POSTGRES_USER:-postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD:-postgres}"
DB_NAME="${POSTGRES_DB:-calorie_ai}"

MAX_ATTEMPTS="${MAX_ATTEMPTS:-180}"
SLEEP_SECONDS="${SLEEP_SECONDS:-2}"
ATTEMPT=1

echo "[seed-mock-user] aguardando tabela users..."

while [ "$ATTEMPT" -le "$MAX_ATTEMPTS" ]; do
  if PGPASSWORD="$DB_PASSWORD" psql \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    -tAc "SELECT to_regclass('public.users') IS NOT NULL;" | grep -q "t"; then
    break
  fi

  ATTEMPT=$((ATTEMPT + 1))
  sleep "$SLEEP_SECONDS"
done

if [ "$ATTEMPT" -gt "$MAX_ATTEMPTS" ]; then
  echo "[seed-mock-user] tabela users não encontrada após $MAX_ATTEMPTS tentativas" >&2
  exit 1
fi

PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d "$DB_NAME" \
  -v ON_ERROR_STOP=1 \
  -c "
INSERT INTO users (
  id,
  firebase_uid,
  email,
  display_name,
  photo_url,
  created_at,
  updated_at,
  weight,
  height,
  age,
  gender,
  activity_level,
  language,
  notifications_enabled,
  timezone
)
VALUES (
  '12493671-992b-4a23-b91c-c953a658e4c9',
  'U037isOCwFPM4XYOFnsTmSyvbyf2',
  'test@example.com',
  'test',
  '/users/12493671-992b-4a23-b91c-c953a658e4c9/avatars/avatar.png',
  '2026-04-20 20:45:52.951723',
  '2026-04-20 20:46:00.01217',
  NULL,
  NULL,
  NULL,
  NULL,
  NULL,
  'en-US',
  false,
  'UTC'
)
ON CONFLICT (firebase_uid)
DO UPDATE SET
  id = EXCLUDED.id,
  email = EXCLUDED.email,
  display_name = EXCLUDED.display_name,
  photo_url = EXCLUDED.photo_url,
  created_at = EXCLUDED.created_at,
  updated_at = EXCLUDED.updated_at,
  weight = EXCLUDED.weight,
  height = EXCLUDED.height,
  age = EXCLUDED.age,
  gender = EXCLUDED.gender,
  activity_level = EXCLUDED.activity_level,
  language = EXCLUDED.language,
  notifications_enabled = EXCLUDED.notifications_enabled,
  timezone = EXCLUDED.timezone;
"

echo "[seed-mock-user] usuário mock inserido/atualizado com sucesso"
