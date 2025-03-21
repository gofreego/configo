package configo

import (
	"context"
	"net/http"
	"time"
)

type registeredConfigsMap map[string]config

type configManagerImpl struct {
	repository        Repository
	config            *ConfigManagerConfig
	registeredConfigs registeredConfigsMap
}

func newConfigManagerImpl(ctx context.Context, cfg *ConfigManagerConfig, repository Repository) (*configManagerImpl, error) {
	cfg.withDefault()
	manager := &configManagerImpl{
		repository:        repository,
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
func (manager *configManagerImpl) RegisterConfig(ctx context.Context, cfg config) error {
	// validate config
	cfgStr, err := marshal(ctx, cfg)
	if err != nil {
		return err
	}

	// check if config is already present in the repository
	value, err := manager.repository.GetConfig(ctx, cfg.Key())
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

		if err := manager.repository.SaveConfig(ctx, &value); err != nil {
			return err
		}
	}
	manager.addConfigToMap(ctx, cfg)
	// save the config in manager
	return nil
}

func (manager *configManagerImpl) Get(ctx context.Context, cfg config) error {
	dbCfg, err := manager.repository.GetConfig(ctx, cfg.Key())
	if err != nil {
		return err
	}
	return unmarshal(ctx, dbCfg.Value, cfg)
}

// RegisterRoute registers routes for the configuration manager.
// register routes with /configs/* endpoints
func (c *configManagerImpl) RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error {

	//setup swagger
	if err := registerFunc(http.MethodGet, "/configo/swagger/*any", c.handleSwagger); err != nil {
		return err
	}

	// setup ui
	if err := registerFunc(http.MethodGet, "/configo/web/*any", c.handleUI); err != nil {
		return err
	}
	// setup get config
	if err := registerFunc(http.MethodGet, "/configo/config/{key}", c.handleGetConfig); err != nil {
		return err
	}

	//setup save config
	if err := registerFunc(http.MethodPost, "/configo/config", c.handleSaveConfig); err != nil {
		return err
	}

	// setup get all configs
	if err := registerFunc(http.MethodGet, "/configo/metadata", c.handleGetConfigMetadata); err != nil {
		return err
	}
	return nil
}
