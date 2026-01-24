package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func stringToPtr(s string) *string {
	return &s
}

func timePtrValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func stringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(f)
	return n
}

func intToPtr(i int) *int {
	return &i
}

func intPtrValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func int32PtrToIntPtr(i *int32) *int {
	if i == nil {
		return nil
	}
	val := int(*i)
	return &val
}

func intPtrToInt32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}
	val := int32(*i)
	return &val
}

func boolPtrValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func boolToPtr(b bool) *bool {
	return &b
}

func uuidFromPgtype(pguuid pgtype.UUID) uuid.UUID {
	uid, _ := uuid.FromBytes(pguuid.Bytes[:])
	return uid
}
