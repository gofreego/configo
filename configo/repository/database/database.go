package database

import (
	"context"
	"fmt"

	"github.com/gofreego/configo/configo/models"
)

type Config struct {
}

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) GetConfig(ctx context.Context, key string) (*models.Config, error) {
	return nil, nil
}

func (r *Repository) SaveConfig(ctx context.Context, cfg *models.Config) error {
	return nil
}
