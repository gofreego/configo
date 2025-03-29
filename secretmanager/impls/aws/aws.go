package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gofreego/configo/secretmanager/models"
)

type Config struct {
	Region    string `json:"region" default:"us-east-1"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Token     string `json:"token"`
	Endpoint  string `json:"endpoint"`
}

type AWSSecretManager struct {
	client *secretsmanager.SecretsManager
	config *Config
}

// NewAWSSecretManager creates a new AWS Secret Manager client
func NewAWSSecretManager(cfg *Config) (*AWSSecretManager, error) {
	if cfg == nil {
		cfg = &Config{
			Region: "us-east-1",
		}
	}

	// Create AWS configuration
	awsConfig := &aws.Config{
		Region: aws.String(cfg.Region),
	}

	// Set custom endpoint if specified (useful for testing or using localstack)
	if cfg.Endpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.Endpoint)
	}

	// Set credentials if provided
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, cfg.Token))
	}

	// Create a new session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create Secrets Manager client
	client := secretsmanager.New(sess)

	return &AWSSecretManager{
		client: client,
		config: cfg,
	}, nil
}

// GetSecret retrieves a secret by key from AWS Secrets Manager
func (m *AWSSecretManager) GetSecret(ctx context.Context, key string) (*models.Secret, error) {
	// Construct input for GetSecretValue operation
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	// Call the AWS Secrets Manager API
	result, err := m.client.GetSecretValue(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", key, err)
	}

	// Get the secret string
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else if result.SecretBinary != nil {
		// Handle binary secrets if needed
		return nil, fmt.Errorf("binary secrets are not supported")
	} else {
		return nil, fmt.Errorf("secret is empty")
	}

	// Parse the secret JSON
	var secret models.Secret
	if err := json.Unmarshal([]byte(secretString), &secret); err != nil {
		// Try to handle case where the secret might be just a token string
		token := secretString
		secret.Token = &token
		return &secret, nil
	}

	return &secret, nil
}
