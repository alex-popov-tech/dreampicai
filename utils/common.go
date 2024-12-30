package utils

import (
	"github.com/google/uuid"
)

func ToUUIDBytes(str string) ([16]byte, error) {
	parsed, err := uuid.Parse(str)
	if err != nil {
		return [16]byte{}, err
	}
	var bytes [16]byte
	copy(bytes[:], parsed[:])
	return bytes, nil
}
