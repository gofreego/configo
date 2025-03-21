package configo

import (
	"context"
	"time"

	"github.com/gofreego/goutils/logger"
)

// saveConfig saves/update the config in cache and repository
func (manager *configo) saveConfig(ctx context.Context, cfg *Config) error {

	if err := manager.cache.SetWithTimeout(ctx, cfg.Key, cfg, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		logger.Error(ctx, "failed to save config in cache: %v", err)
		return NewInternalServerErr("failed to save config in cache, Err: %v", err)
	}

	if err := manager.repository.SaveConfig(ctx, cfg); err != nil {
		logger.Error(ctx, "failed to save config in repository: %v", err)
		return NewInternalServerErr("failed to save config in repository, Err: %v", err)
	}
	return nil
}

// getConfig gets the config from cache or repository. returns ErrConfigNotFound if config is not found
func (manager *configo) getConfig(ctx context.Context, key string) (*Config, error) {
	var cfg Config
	err := manager.cache.GetV(ctx, key, &cfg)
	if err != nil {
		logger.Error(ctx, "failed to get config from cache: %v", err)
		return nil, NewInternalServerErr("failed to get config from cache, Err: %v", err)
	}
	if cfg.Key != "" {
		return &cfg, nil
	}

	repoCfg, err := manager.repository.GetConfig(ctx, key)
	if err != nil {
		logger.Error(ctx, "failed to get config from repository: %v", err)
		return nil, NewInternalServerErr("failed to get config from repository, Err: %v", err)
	}
	if repoCfg == nil {
		return nil, ErrConfigNotFound
	}

	if err := manager.cache.SetWithTimeout(ctx, key, repoCfg, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		logger.Error(ctx, "failed to save config in cache: %v", err)
		return nil, NewInternalServerErr("failed to save config in cache, Err: %v", err)
	}
	return repoCfg, nil
}

func (manager *configo) addConfigToMap(_ context.Context, cfg config) {
	manager.registeredConfigs[cfg.Key()] = cfg
}
