package svc

import (
	"time"

	"github.com/google/uuid"
)

func genUUID() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// timeToUnix converts *time.Time to int64 Unix timestamp.
// Returns 0 if t is nil.
func timeToUnix(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

// timePtrToInt64Ptr converts *time.Time to *int64 Unix timestamp.
func timePtrToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	v := t.Unix()
	return &v
}
