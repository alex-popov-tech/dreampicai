package domain

import (
	"time"

	"dreampicai/pkg/db"
)

const (
	ImageStatusStarted   ImageStatus = "started"
	ImageStatusCancelled ImageStatus = "cancelled"
	ImageStatusFailed    ImageStatus = "failed"
	ImageStatusSucceeded ImageStatus = "succeeded"
)

type (
	ImageStatus = db.ImageStatus
	Image       struct {
		ID        int32
		OwnerID   int32
		Status    ImageStatus
		Prompt    string
		Url       string
		CreatedAt time.Time
	}
)
