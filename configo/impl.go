package configo

import (
	"context"
	"net/http"
	"time"

	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/cache/memory"
)

type registeredConfigsMap map[string]config

type configo struct {
	repository        Repository
	cache             cache.Cache
	config            *ConfigManagerConfig
	registeredConfigs registeredConfigsMap
}

func newconfigo(ctx context.Context, cfg *ConfigManagerConfig, repository Repository) (*configo, error) {
	cfg.withDefault()
	manager := &configo{
		repository:        repository,
		cache:             memory.NewCache(),
		config:            cfg,
		registeredConfigs: make(registeredConfigsMap),
	}

	err := manager.RegisterConfig(ctx, manager.config)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// RegisterConfig will register config and setup a UI for it. It will also validate the config.
func (manager *configo) RegisterConfig(ctx context.Context, cfg config) error {
	// validate config
	cfgStr, err := marshal(ctx, cfg)
	if err != nil {
		return err
	}

	// check if config is already present in the repository
	value, err := manager.getConfig(ctx, cfg.Key())
	if err != nil && err != ErrConfigNotFound {
		return err
	}

	// if config is not present in the repository, save it
	if value == nil {
		var value Config = Config{
			Key:       cfg.Key(),
			Value:     cfgStr,
			UpdatedBy: "",
			UpdatedAt: time.Now().UnixMilli(),
			CreatedAt: time.Now().UnixMilli(),
		}

		if err := manager.saveConfig(ctx, &value); err != nil {
			return err
		}
	}
	manager.addConfigToMap(ctx, cfg)
	// save the config in manager
	return nil
}

func (manager *configo) Get(ctx context.Context, cfg config) error {
	dbCfg, err := manager.getConfig(ctx, cfg.Key())
	if err != nil {
		return err
	}
	return unmarshal(ctx, dbCfg.Value, cfg)
}

// RegisterRoute registers routes for the configuration manager.
// register routes with /configs/* endpoints
func (c *configo) RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error {

	//setup swagger
	if err := registerFunc(http.MethodGet, "/configs/swagger/*any", c.handleSwagger); err != nil {
		return err
	}

	// setup ui
	if err := registerFunc(http.MethodGet, "/configs/ui", c.handleUI); err != nil {
		return err
	}
	// setup get config
	if err := registerFunc(http.MethodGet, "/configs/config/{key}", c.handleGetConfig); err != nil {
		return err
	}

	//setup save config
	if err := registerFunc(http.MethodPost, "/configs/config", c.handleSaveConfig); err != nil {
		return err
	}

	// setup get all configs
	if err := registerFunc(http.MethodGet, "/configs/metadata", c.handleGetConfigMetadata); err != nil {
		return err
	}
	return nil
}
