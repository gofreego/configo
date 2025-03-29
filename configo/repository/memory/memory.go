package memory

import (
	"context"

	"github.com/gofreego/configo/configo/models"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/cache/memory"
)

type Repository struct {
	cache cache.Cache
}

func NewRepository() (*Repository, error) {
	return &Repository{
		cache: memory.NewCache(),
	}, nil
}

func (r *Repository) GetConfig(ctx context.Context, key string) (*models.Config, error) {
	var cfg models.Config
	err := r.cache.GetV(ctx, key, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Key == "" {
		return nil, nil
	}
	return &cfg, nil
}

func (r *Repository) SaveConfig(ctx context.Context, cfg *models.Config) error {
	return r.cache.Set(ctx, cfg.Key, cfg)
}
