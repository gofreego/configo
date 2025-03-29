package consul

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofreego/configo/configo/models"
	"github.com/hashicorp/consul/api"
)

type Config struct {
	Host       string
	Datacenter string
	Token      string
	Prefix     string
}

type Repository struct {
	cfg    *Config
	client *api.Client
}

func NewRepository(cfg *Config) (*Repository, error) {
	// Configure Consul client
	config := api.DefaultConfig()
	config.Address = cfg.Host
	config.Datacenter = cfg.Datacenter

	if cfg.Token != "" {
		config.Token = cfg.Token
	}

	// Create Consul client
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &Repository{
		cfg:    cfg,
		client: client,
	}, nil
}

func (r *Repository) GetConfig(ctx context.Context, key string) (*models.Config, error) {

	// Get KV client
	kv := r.client.KV()

	// Construct the full key with prefix
	fullKey := fmt.Sprintf("%s%s", r.cfg.Prefix, key)

	// Get the key
	pair, _, err := kv.Get(fullKey, nil)

	if err != nil {
		return nil, fmt.Errorf("consul get failed for key %s: %w", key, err)
	}

	// Check if key exists
	if pair == nil {
		return nil, fmt.Errorf("config not found for key: %s", key)
	}

	// Unmarshal the value
	var cfg models.Config
	if err := json.Unmarshal(pair.Value, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config for key %s: %w", key, err)
	}

	// Set the key from the request
	cfg.Key = key

	return &cfg, nil
}

func (r *Repository) SaveConfig(ctx context.Context, cfg *models.Config) error {

	// Get KV client
	kv := r.client.KV()

	// Construct the full key with prefix
	fullKey := fmt.Sprintf("%s%s", r.cfg.Prefix, cfg.Key)

	// Marshal the config to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create KV pair
	pair := &api.KVPair{
		Key:   fullKey,
		Value: data,
	}

	// Put the key-value pair
	_, err = kv.Put(pair, nil)

	if err != nil {
		return fmt.Errorf("consul save failed for key %s: %w", cfg.Key, err)
	}

	return nil
}
