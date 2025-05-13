package configo

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/internal/handlers"
	"github.com/gofreego/configo/configo/internal/parser"
	"github.com/gofreego/configo/configo/internal/service"
	"github.com/gofreego/configo/configo/internal/utils"
	"github.com/gofreego/configo/configo/models"
	"github.com/gofreego/configo/configo/repository"
)

type configManagerImpl struct {
	repository repository.Repository
	service    *service.Service
	handler    *handlers.Handler
	mux        *gin.Engine
}

func newConfigManagerImpl(ctx context.Context, cfg *configs.ConfigManagerConfig, repository repository.Repository, pathPrefix ...string) (*configManagerImpl, error) {
	cfg.WithDefault()
	service, err := service.NewService(ctx, cfg, repository)
	if err != nil {
		return nil, err
	}
	manager := &configManagerImpl{
		repository: repository,
		service:    service,
		handler:    handlers.NewHandler(service),
		mux:        gin.New(),
	}

	manager.registerRoutes(ctx, pathPrefix...)
	err = manager.RegisterConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// RegisterConfig will register config and setup a UI for it. It will also validate the config.
// if the notifyFunc is provided, it will be called after the config is changed.
func (manager *configManagerImpl) RegisterConfig(ctx context.Context, cfg any, notifyFunc ...service.Notify) error {
	// validate config
	cfgStr, err := parser.Marshal(ctx, cfg)
	if err != nil {
		return err
	}
	configName := utils.GetNameOfTheObject(cfg)
	// check if config is already present in the repository
	value, err := manager.repository.GetConfig(ctx, configName)
	if err != nil {
		return err
	}

	// if config is not present in the repository, save it
	if value == nil {
		var value models.Config = models.Config{
			Key:       configName,
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

	manager.service.AddConfigToMap(ctx, service.NewConfigDetails(configName, cfg, notifyFunc...))
	// save the config in manager
	return nil
}

// RegisterRoute registers routes for the configuration manager.
// register routes with /configs/* endpoints
func (c *configManagerImpl) registerRoutes(ctx context.Context, pathPrefix ...string) error {
	pathPrefixStr := ""
	if len(pathPrefix) > 0 {
		pathPrefixStr = pathPrefix[0]
	}
	group := c.mux.Group(pathPrefixStr + "/configo/v1")
	group.GET("/swagger/*any", c.handler.Swagger)
	group.GET("/web/*any", c.handler.UI)
	group.GET("/config", c.handler.GetConfig)
	group.POST("/config", c.handler.SaveConfig)
	group.GET("/metadata", c.handler.GetConfigMetadata)
	return nil
}

func (c *configManagerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}
