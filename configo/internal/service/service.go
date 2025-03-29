package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/internal/errors"
	"github.com/gofreego/configo/configo/internal/models"
	"github.com/gofreego/configo/configo/internal/parser"
	"github.com/gofreego/configo/configo/repository"
	"github.com/gofreego/goutils/logger"
)

const (
	// DefaultConfigRefreshInSecs is the default time in seconds after which the config manager will refresh the configs.
	DefaultConfigRefreshInSecs = 10
)

type Service struct {
	repository        repository.Repository
	registeredConfigs registeredConfigsMap
	config            *configs.ConfigManagerConfig
}

func NewService(ctx context.Context, cfg *configs.ConfigManagerConfig, repo repository.Repository) (*Service, error) {
	if cfg == nil {
		cfg = &configs.ConfigManagerConfig{}
	}
	if repo == nil {
		return nil, errors.NewInternalServerErr("repository is required, got nil")
	}
	s := &Service{repository: repo, config: cfg, registeredConfigs: make(registeredConfigsMap)}
	go s.refreshConfigs(ctx)
	return s, nil
}

func (manager *Service) UpdateConfig(ctx context.Context, req *models.UpdateConfigRequest) error {

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
	bytes, err := json.Marshal(req.Configs)
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

func (manager *Service) GetConfigByKey(ctx context.Context, key string) (*models.GetConfigResponse, error) {
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

	return &models.GetConfigResponse{
		Key:       config.Key,
		Configs:   obj,
		UpdatedBy: config.UpdatedBy,
		UpdatedAt: config.UpdatedAt,
		CreatedAt: config.CreatedAt,
	}, nil
}

func (manager *Service) GetConfigsMetadata(_ context.Context) (*models.ConfigMetadataResponse, error) {
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

func (manager *Service) AddConfigToMap(_ context.Context, configDetails *ConfigDetails) {
	manager.registeredConfigs[configDetails.Key] = configDetails
}
