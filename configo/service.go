package configo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofreego/configo/configo/internal/errors"
	"github.com/gofreego/configo/configo/internal/models"
	"github.com/gofreego/configo/configo/internal/parser"
	"github.com/gofreego/goutils/logger"
)

func (manager *ConfigManagerImpl) updateConfig(ctx context.Context, req *models.UpdateConfigRequest) error {

	if manager.registeredConfigs[req.Id] == nil {
		logger.Error(ctx, "config not registered: %v", req.Id)
		return errors.ErrConfigNotFound
	}

	config, err := manager.repository.GetConfig(ctx, req.Id)
	if err != nil {
		logger.Error(ctx, "failed to get config: %v", err)
		return errors.NewInternalServerErr("failed to get config, Err: %v", err)
	}

	config.UpdatedAt = time.Now().UnixMilli()
	config.UpdatedBy = req.UpdatedBy
	bytes, err := json.Marshal(req.Value)
	if err != nil {
		logger.Error(ctx, "failed to marshal config: %v", err)
		return errors.ErrInvalidConfig
	}
	config.Value = string(bytes)
	err = manager.repository.SaveConfig(ctx, config)
	if err != nil {
		logger.Error(ctx, "failed to save config: %v", err)
		return errors.NewInternalServerErr("failed to save config, Err: %v", err)
	}
	err = parser.Unmarshal(ctx, config.Value, manager.registeredConfigs[req.Id])
	if err != nil {
		logger.Error(ctx, "failed to unmarshal config: %v", err)
		return errors.ErrInvalidConfig
	}
	logger.Info(ctx, "config updated successfully: %v", req.Id)
	return nil
}

func (manager *ConfigManagerImpl) getConfigByKey(ctx context.Context, key string) (*models.GetConfigResponse, error) {
	config, err := manager.repository.GetConfig(ctx, key)
	if err != nil {
		logger.Error(ctx, "failed to get config: %v", err)
		return nil, errors.NewInternalServerErr("failed to get config, Err: %v", err)
	}
	if config == nil {
		logger.Error(ctx, "config not found: %v", key)
		return nil, errors.ErrConfigNotFound
	}
	var obj []models.ConfigObject
	err = json.Unmarshal([]byte(config.Value), &obj)
	if err != nil {
		logger.Error(ctx, "failed to unmarshal config: %v", err)
		return nil, errors.ErrInvalidConfig
	}

	return &models.GetConfigResponse{Configs: obj}, nil
}

func (manager *ConfigManagerImpl) getConfigsMetadata(_ context.Context) (*models.ConfigMetadataResponse, error) {
	keys := make([]string, 0, len(manager.registeredConfigs))
	for k := range manager.registeredConfigs {
		keys = append(keys, k)
	}
	return &models.ConfigMetadataResponse{
		ServiceInfo: models.ServiceInfo{
			Name:        manager.config.ServiceName,
			Description: manager.config.ServiceDescription,
		},
		ConfigInfo: models.ConfigInfo{
			ConfigKeys: keys,
		},
	}, nil
}

func (manager *ConfigManagerImpl) addConfigToMap(_ context.Context, cfg config) {
	manager.registeredConfigs[cfg.Key()] = cfg
}
