package configo

import (
	"context"
	"net/http"
	"time"

	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/internal/errors"
	"github.com/gofreego/configo/configo/internal/parser"
	"github.com/gofreego/configo/configo/internal/repository"
	"github.com/gofreego/configo/configo/models"
	"github.com/gofreego/goutils/logger"
)

const (
	// DefaultConfigRefreshInSecs is the default time in seconds after which the config manager will refresh the configs.
	DefaultConfigRefreshInSecs = 10
)

type registeredConfigsMap map[string]config

type ConfigManagerImpl struct {
	repository        repository.Repository
	config            *configs.ConfigManagerConfig
	registeredConfigs registeredConfigsMap
}

func newConfigManagerImpl(ctx context.Context, cfg *configs.ConfigManagerConfig, repository repository.Repository) (*ConfigManagerImpl, error) {
	cfg.WithDefault()
	manager := &ConfigManagerImpl{
		repository:        repository,
		config:            cfg,
		registeredConfigs: make(registeredConfigsMap),
	}

	err := manager.RegisterConfig(ctx, manager.config)
	if err != nil {
		return nil, err
	}
	go manager.refreshConfigs(ctx)
	return manager, nil
}

func (manager *ConfigManagerImpl) refreshConfigs(ctx context.Context) {
	for {
		if manager.config.ConfigRefreshInSecs == 0 {
			manager.config.ConfigRefreshInSecs = DefaultConfigRefreshInSecs
		}
		logger.Debug(ctx, "refreshing configs after %v seconds", manager.config.ConfigRefreshInSecs)
		time.Sleep(time.Duration(manager.config.ConfigRefreshInSecs) * time.Second)
		for key, cfg := range manager.registeredConfigs {
			value, err := manager.repository.GetConfig(ctx, key)
			if err != nil {
				continue
			}
			if value == nil {
				continue
			}
			err = parser.Unmarshal(ctx, value.Value, cfg)
			if err != nil {
				continue
			}
		}
	}

}

// RegisterConfig will register config and setup a UI for it. It will also validate the config.
func (manager *ConfigManagerImpl) RegisterConfig(ctx context.Context, cfg config) error {
	// validate config
	cfgStr, err := parser.Marshal(ctx, cfg)
	if err != nil {
		return err
	}

	// check if config is already present in the repository
	value, err := manager.repository.GetConfig(ctx, cfg.Key())
	if err != nil && err != errors.ErrConfigNotFound {
		return err
	}

	// if config is not present in the repository, save it
	if value == nil {
		var value models.Config = models.Config{
			Key:       cfg.Key(),
			Value:     cfgStr,
			UpdatedBy: "",
			UpdatedAt: time.Now().UnixMilli(),
			CreatedAt: time.Now().UnixMilli(),
		}
		if err := manager.repository.SaveConfig(ctx, &value); err != nil {
			return err
		}
	} else {
		// if config is present in the repository, unmarshal it
		err = parser.Unmarshal(ctx, value.Value, cfg)
		if err != nil {
			return err
		}
	}

	manager.addConfigToMap(ctx, cfg)
	// save the config in manager
	return nil
}

// RegisterRoute registers routes for the configuration manager.
// register routes with /configs/* endpoints
func (c *ConfigManagerImpl) RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error {

	//setup swagger
	if err := registerFunc(http.MethodGet, "/configo/swagger/*any", c.handleSwagger); err != nil {
		return err
	}

	// setup ui
	if err := registerFunc(http.MethodGet, "/configo/web/*any", c.handleUI); err != nil {
		return err
	}
	// setup get config
	if err := registerFunc(http.MethodGet, "/configo/config", c.handleGetConfig); err != nil {
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
