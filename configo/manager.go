package configo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/internal/service"
	"github.com/gofreego/configo/configo/repository"
)

// RouteRegistrar defines a generic function type for registering routes.
type RouteRegistrar func(method, path string, handler http.HandlerFunc) error

type ConfigManager interface {
	// It will register the config in manager, it will save the config in repo if not already present with default values.
	// if notifyFunc is provided, it will be called after the config is changed.
	RegisterConfig(ctx context.Context, cfg any, notifyFunc ...service.Notify) error
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
