package repository

import (
	"context"

	"github.com/gofreego/configo/configo/models"
)

type Repository interface {
	// GetConfig will return the value of the configuration for the given key. It will return nil if the configuration is not found.
	// it will return an error if there is an issue with the repository.
	GetConfig(ctx context.Context, key string) (*models.Config, error)
	// SaveConfig will save the configuration with the given key and value. It will return an error if there is an issue with the repository.
	SaveConfig(ctx context.Context, cfg *models.Config) error
}
