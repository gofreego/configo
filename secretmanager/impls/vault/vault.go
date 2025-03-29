package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofreego/configo/secretmanager/models"
	"github.com/hashicorp/vault/api"
)

type Config struct {
	Address string `json:"address" default:"http://127.0.0.1:8200"`
	Token   string `json:"token"`
	Path    string `json:"path" default:"secret"`
}

type VaultSecretManager struct {
	client *api.Client
	config *Config
}

// NewVaultSecretManager creates a new HashiCorp Vault secret manager client
func NewVaultSecretManager(cfg *Config) (*VaultSecretManager, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Set default values if not provided
	if cfg.Address == "" {
		cfg.Address = "http://127.0.0.1:8200"
	}
	if cfg.Path == "" {
		cfg.Path = "secret"
	}

	// Create Vault configuration
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.Address

	// Create Vault client
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	if cfg.Token == "" {
		return nil, errors.New("no authentication token provided")
	}

	client.SetToken(cfg.Token)

	return &VaultSecretManager{
		client: client,
		config: cfg,
	}, nil
}

// GetSecret retrieves a secret by key from HashiCorp Vault
func (m *VaultSecretManager) GetSecret(ctx context.Context, key string) (*models.Secret, error) {

	// Construct the path
	path := fmt.Sprintf("%s/data/%s", m.config.Path, key)

	// For KV v1, use a different path format
	if !m.isKVv2() {
		path = fmt.Sprintf("%s/%s", m.config.Path, key)
	}

	// Get the secret
	secret, err := m.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", key, err)
	}

	if secret == nil {
		return nil, fmt.Errorf("secret not found: %s", key)
	}

	// Parse the secret data
	var data map[string]interface{}
	if m.isKVv2() {
		// KV v2 store has a nested data structure
		dataRaw, ok := secret.Data["data"]
		if !ok {
			return nil, fmt.Errorf("secret data not found in KV v2 store")
		}
		data, ok = dataRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid secret data format")
		}
	} else {
		// KV v1 store has a flat data structure
		data = secret.Data
	}

	// Create secret model
	secretModel := &models.Secret{}

	// Extract token if present
	if tokenVal, ok := data["token"]; ok {
		if tokenStr, ok := tokenVal.(string); ok {
			secretModel.Token = &tokenStr
		}
	}

	// Extract username if present
	if usernameVal, ok := data["username"]; ok {
		if usernameStr, ok := usernameVal.(string); ok {
			secretModel.Username = &usernameStr
		}
	}

	// Extract password if present
	if passwordVal, ok := data["password"]; ok {
		if passwordStr, ok := passwordVal.(string); ok {
			secretModel.Password = &passwordStr
		}
	}

	return secretModel, nil
}

// isKVv2 checks if the configured path is using KV v2
func (m *VaultSecretManager) isKVv2() bool {
	// A simple check to determine if we're using KV v2
	// In a real implementation, you might want to check this more thoroughly
	mountPath := m.config.Path

	// Check if the mount exists and get its type
	mounts, err := m.client.Sys().ListMounts()
	if err != nil {
		// Default to true for safety
		return true
	}

	// Check if the mount exists
	mount, exists := mounts[mountPath+"/"]
	if exists {
		// Check if it's KV v2
		if mount.Type == "kv" && mount.Options != nil {
			if version, ok := mount.Options["version"]; ok && version == "2" {
				return true
			}
		}
	}

	return false
}
