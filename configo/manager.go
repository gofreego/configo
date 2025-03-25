package configo

import (
	"context"
	"fmt"
	"net/http"
)

// ConfigTag is a type for configuration tags.
type ConfigTag string

const (
	// CONFIG_TAG_NAME is the tag for the name of the configuration. It is required tag.
	CONFIG_TAG_NAME ConfigTag = "name"
	// CONFIG_TAG_DESCRIPTION is the tag for the description of the configuration. It is optional tag.
	CONFIG_TAG_DESCRIPTION ConfigTag = "description"
	// CONFIG_TAG_TYPE is the tag for the type of the configuration. It is required
	CONFIG_TAG_TYPE ConfigTag = "type"
	// CONFIG_TAG_REQUIRED is the tag for the required value of the configuration. It should be true or false. It will be false by default.
	CONFIG_TAG_REQUIRED ConfigTag = "required"
	// CONFIG_TAG_CHOICES is the tag for the choices of the configuration. It is required if the type is choice.
	CONFIG_TAG_CHOICES ConfigTag = "choices"
)

type ConfigType string

const (
	// CONFIG_TYPE_STRING is the type for string configuration, it will show a textbox on ui.
	CONFIG_TYPE_STRING ConfigType = "string"
	// CONFIG_TYPE_NUMBER is the type for integer configuration, it will show a number input on ui.
	CONFIG_TYPE_NUMBER ConfigType = "number"
	// CONFIG_TYPE_BOOLEAN is the type for boolean configuration, it will show a checkbox on ui.
	CONFIG_TYPE_BOOLEAN ConfigType = "boolean"
	// CONFIG_TYPE_JSON is the type for json configuration, it will show a textarea on ui which will have json formatting.
	CONFIG_TYPE_JSON ConfigType = "json"
	// CONFIG_TYPE_FLOAT is the type for float configuration, it will show a number input on ui.
	CONFIG_TYPE_BIG_TEXT ConfigType = "bigText"
	// CONFIG_TYPE_CHOICE is the type for choice configuration, it will show a dropdown on ui and it should have type string.
	CONFIG_TYPE_CHOICE ConfigType = "choice"
	//CONFIG_TYPE_PARENT
	CONFIG_TYPE_PARENT ConfigType = "parent"
	//CONFIG_TYPE_LIST
	CONFIG_TYPE_LIST ConfigType = "list"
)

type Config struct {
	Key   string
	Value string
	// UpdatedBy is the user who updated the configuration. It will taken from header (X-User-Id) of the request. it will be empty if header is not present.
	UpdatedBy string
	UpdatedAt int64
	CreatedAt int64
}

type Repository interface {
	// GetConfig will return the value of the configuration for the given key. It will return nil if the configuration is not found.
	// it will return an error if there is an issue with the repository.
	GetConfig(ctx context.Context, key string) (*Config, error)
	// SaveConfig will save the configuration with the given key and value. It will return an error if there is an issue with the repository.
	SaveConfig(ctx context.Context, cfg *Config) error
}

// RouteRegistrar defines a generic function type for registering routes.
type RouteRegistrar func(method, path string, handler http.HandlerFunc) error

type config interface {
	// Key is the unique key by which you want to save the configuration.
	Key() string
}

type ConfigManagerConfig struct {
	ServiceName         string `name:"service_name" type:"string" description:"service name" required:"true"`
	ServiceDescription  string `name:"service_description" type:"bigText" description:"This config is to update the description of this service" required:"false"`
	ConfigRefreshInSecs int    `name:"config_refresh_in_secs" type:"number" description:"config refresh in seconds" required:"false"`
}

func (c *ConfigManagerConfig) withDefault() {
	if c.ServiceName == "" {
		c.ServiceName = "Config Manager Service"
	}
}

func (c *ConfigManagerConfig) Key() string {
	return "config-manager-config"
}

type ConfigManager interface {
	// It will register the config in manager, it will save the config in repo if not already present with default values.
	RegisterConfig(ctx context.Context, cfg config) error
	// RegisterRoute will register the necesory endpoints for the configuration manager. using the provided RouteRegistrar.
	RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error
}

func NewConfigManager(ctx context.Context, cfg *ConfigManagerConfig, repo Repository) (ConfigManager, error) {
	if cfg == nil {
		cfg = &ConfigManagerConfig{}
	}
	if repo == nil {
		return nil, fmt.Errorf("repository is required, got nil")
	}
	return newConfigManagerImpl(ctx, cfg, repo)
}
