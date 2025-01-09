// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ImageStatus string

const (
	ImageStatusStarted   ImageStatus = "started"
	ImageStatusCancelled ImageStatus = "cancelled"
	ImageStatusFailed    ImageStatus = "failed"
	ImageStatusSucceeded ImageStatus = "succeeded"
)

func (e *ImageStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ImageStatus(s)
	case string:
		*e = ImageStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ImageStatus: %T", src)
	}
	return nil
}

type NullImageStatus struct {
	ImageStatus ImageStatus
	Valid       bool // Valid is true if ImageStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullImageStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ImageStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ImageStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullImageStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ImageStatus), nil
}

type Account struct {
	ID       int32
	UserID   pgtype.UUID
	Username string
}

type Image struct {
	ID             int32
	ProviderID     string
	OwnerID        pgtype.Int4
	Status         ImageStatus
	Prompt         string
	NegativePrompt string
	Model          string
	Url            pgtype.Text
	CreatedAt      pgtype.Timestamptz
}
