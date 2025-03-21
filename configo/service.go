package configo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofreego/goutils/logger"
)

func (manager *configo) updateConfig(ctx context.Context, req *UpdateConfigRequest) error {

	if manager.registeredConfigs[req.Key] == nil {
		logger.Error(ctx, "config not registered: %v", req.Key)
		return ErrConfigNotFound
	}

	config, err := manager.getConfig(ctx, req.Key)
	if err != nil {
		return err
	}

	config.UpdatedAt = time.Now().UnixMilli()
	config.UpdatedBy = req.UpdatedBy
	bytes, err := json.Marshal(req.Value)
	if err != nil {
		logger.Error(ctx, "failed to marshal config: %v", err)
		return ErrInvalidConfig
	}
	config.Value = string(bytes)
	return manager.saveConfig(ctx, config)
}

func (manager *configo) getConfigByKey(ctx context.Context, key string) (*ConfigObject, error) {
	config, err := manager.getConfig(ctx, key)
	if err != nil {
		return nil, err
	}
	var obj ConfigObject
	err = json.Unmarshal([]byte(config.Value), &obj)
	if err != nil {
		logger.Error(ctx, "failed to unmarshal config: %v", err)
		return nil, ErrInvalidConfig
	}
	return &obj, nil
}

func (manager *configo) getConfigsMetadata(_ context.Context) (*configMetadata, error) {
	keys := make([]string, 0, len(manager.registeredConfigs))
	for k := range manager.registeredConfigs {
		keys = append(keys, k)
	}
	return &configMetadata{Keys: keys}, nil
}
