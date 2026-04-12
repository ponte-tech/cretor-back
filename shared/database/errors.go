package database

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

func MapError(err error, entity string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%s: %w", entity, ErrNotFound)
	}
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, we := range writeErr.WriteErrors {
			if we.Code == 11000 {
				return fmt.Errorf("%s: %w", entity, ErrDuplicateKey)
			}
		}
	}
	return fmt.Errorf("%s: %w", entity, err)
}
