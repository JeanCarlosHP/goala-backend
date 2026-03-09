package repositories

import (
	"time"
)

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

func intPtrValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func int32PtrValue(i *int32) int32 {
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

func valueOrZero[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
