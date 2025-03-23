package configo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofreego/goutils/logger"
)

func (manager *configManagerImpl) updateConfig(ctx context.Context, req *UpdateConfigRequest) error {

	if manager.registeredConfigs[req.Key] == nil {
		logger.Error(ctx, "config not registered: %v", req.Key)
		return ErrConfigNotFound
	}

	config, err := manager.repository.GetConfig(ctx, req.Key)
	if err != nil {
		logger.Error(ctx, "failed to get config: %v", err)
		return NewInternalServerErr("failed to get config, Err: %v", err)
	}

	config.UpdatedAt = time.Now().UnixMilli()
	config.UpdatedBy = req.UpdatedBy
	bytes, err := json.Marshal(req.Value)
	if err != nil {
		logger.Error(ctx, "failed to marshal config: %v", err)
		return ErrInvalidConfig
	}
	config.Value = string(bytes)
	err = manager.repository.SaveConfig(ctx, config)
	if err != nil {
		logger.Error(ctx, "failed to save config: %v", err)
		return NewInternalServerErr("failed to save config, Err: %v", err)
	}
	return nil
}

func (manager *configManagerImpl) getConfigByKey(ctx context.Context, key string) ([]ConfigObject, error) {
	config, err := manager.repository.GetConfig(ctx, key)
	if err != nil {
		logger.Error(ctx, "failed to get config: %v", err)
		return nil, NewInternalServerErr("failed to get config, Err: %v", err)
	}
	if config == nil {
		logger.Error(ctx, "config not found: %v", key)
		return nil, ErrConfigNotFound
	}
	var obj []ConfigObject
	err = json.Unmarshal([]byte(config.Value), &obj)
	if err != nil {
		logger.Error(ctx, "failed to unmarshal config: %v", err)
		return nil, ErrInvalidConfig
	}

	return obj, nil
}

func (manager *configManagerImpl) getConfigsMetadata(_ context.Context) (*configMetadataResponse, error) {
	keys := make([]string, 0, len(manager.registeredConfigs))
	for k := range manager.registeredConfigs {
		keys = append(keys, k)
	}
	return &configMetadataResponse{
		ServiceInfo: ServiceInfo{
			Name:        manager.config.ServiceName,
			Description: manager.config.ServiceDescription,
		},
		ConfigInfo: ConfigInfo{
			ConfigKeys: keys,
		},
	}, nil
}

func (manager *configManagerImpl) addConfigToMap(_ context.Context, cfg config) {
	manager.registeredConfigs[cfg.Key()] = cfg
}
