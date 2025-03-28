package service

import (
	"context"
	"time"

	"github.com/gofreego/configo/configo/internal/parser"
	"github.com/gofreego/goutils/logger"
)

type ConfigDetails struct {
	Key          string
	configObject any
	notifyFunc   func()
	updatedOn    int64
}
type Notify func()

func NewConfigDetails(key string, cfg any, notifyFunc ...Notify) *ConfigDetails {
	var f func() = nil
	if len(notifyFunc) > 0 {
		f = notifyFunc[0]
	}
	return &ConfigDetails{
		Key:          key,
		configObject: cfg,
		notifyFunc:   f,
		updatedOn:    time.Now().UnixMilli(),
	}
}

type registeredConfigsMap map[string]*ConfigDetails

func (manager *Service) refreshConfigs(ctx context.Context) {
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
			if value.UpdatedAt <= cfg.updatedOn {
				continue
			}
			if value == nil {
				continue
			}
			err = parser.Unmarshal(ctx, value.Value, cfg.configObject)
			if err != nil {
				continue
			}
			cfg.updatedOn = value.UpdatedAt
			if cfg.notifyFunc != nil {
				cfg.notifyFunc()
			}
		}
	}

}
