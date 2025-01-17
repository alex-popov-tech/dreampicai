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

type ReplicateModel = string

const (
	// playgroundai/playground-v2.5-1024px-aesthetic:a45f82a1382bed5c7aeb861dac7c7d191b0fdf74d8d57c4a0e6ed7d4d0bf7d24
	REPLICATE_MODEL_PLAYGROUND ReplicateModel = "a45f82a1382bed5c7aeb861dac7c7d191b0fdf74d8d57c4a0e6ed7d4d0bf7d24"
	// ai-forever/kandinsky-2.2
	REPLICATE_MODEL_KADNINSKY ReplicateModel = "ad9d7879fbffa2874e1d909d1d37d9bc682889cc65b31f7bb00d2362619f194a"
	// datacte / proteus-v0.3
	REPLICATE_MODEL_PROTEUS ReplicateModel = "b28b79d725c8548b173b6a19ff9bffd16b9b80df5b18b8dc5cb9e1ee471bfa48"
)

var MODELS = []ReplicateModel{
	REPLICATE_MODEL_PLAYGROUND,
	REPLICATE_MODEL_KADNINSKY,
	REPLICATE_MODEL_PROTEUS,
}

const (
	DEFAULT_MODEL           = REPLICATE_MODEL_PLAYGROUND
	DEFAULT_PROMPT          = "Neon cyberpunk Ukrainian woman in yellow-blue clothes, 8k"
	DEFAULT_NEGATIVE_PROMPT = "ugly, deformed, noisy, blurry, distorted, worst quality, low quality"
	DEFAULT_COUNT           = "1"
)
