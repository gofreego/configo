package secretmanager

import (
	"context"

	"github.com/gofreego/configo/secretmanager/impls/aws"
	"github.com/gofreego/configo/secretmanager/impls/vault"
	"github.com/gofreego/configo/secretmanager/models"
	"github.com/gofreego/goutils/customerrors"
)

type Config struct {
	Name  Name
	AWS   aws.Config
	Vault vault.Config
}

type Name string

const (
	AWS   Name = "aws"
	Vault Name = "vault"
)

type Manager interface {
	GetSecret(ctx context.Context, key string) (*models.Secret, error)
}

func NewManager(cfg *Config) (Manager, error) {
	switch cfg.Name {
	case AWS:
		return aws.NewAWSSecretManager(&cfg.AWS)
	case Vault:
		return vault.NewVaultSecretManager(&cfg.Vault)
	default:
		return nil, customerrors.BAD_REQUEST_ERROR("invalid secret manager name")
	}
}
