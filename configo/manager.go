package configo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/internal/repository"
)

// RouteRegistrar defines a generic function type for registering routes.
type RouteRegistrar func(method, path string, handler http.HandlerFunc) error

type config interface {
	// Key is the unique key by which you want to save the configuration.
	Key() string
}

type ConfigManager interface {
	// It will register the config in manager, it will save the config in repo if not already present with default values.
	RegisterConfig(ctx context.Context, cfg config) error
	// RegisterRoute will register the necesory endpoints for the configuration manager. using the provided RouteRegistrar.
	RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error
}

func NewConfigManager(ctx context.Context, cfg *configs.ConfigManagerConfig, repo repository.Repository) (ConfigManager, error) {
	if cfg == nil {
		cfg = &configs.ConfigManagerConfig{}
	}
	if repo == nil {
		return nil, fmt.Errorf("repository is required, got nil")
	}
	return newConfigManagerImpl(ctx, cfg, repo)
}
