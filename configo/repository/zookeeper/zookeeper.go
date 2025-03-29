package zookeeper

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/gofreego/configo/configo/models"
)

type Config struct {
	Host     string
	Prefix   string
	Username string
	Password string
}

type Repository struct {
	cfg  *Config
	conn *zk.Conn
}

// NewRepository creates a new ZooKeeper repository
func NewRepository(cfg *Config) (*Repository, error) {
	// Parse servers
	var servers []string
	if cfg.Host != "" {
		servers = strings.Split(cfg.Host, ",")
	} else {
		servers = []string{"127.0.0.1:2181"} // Default to localhost if nothing specified
	}

	// Ensure prefix starts with /
	if !strings.HasPrefix(cfg.Prefix, "/") {
		cfg.Prefix = "/" + cfg.Prefix
	}

	// Connect to ZooKeeper with default timeout
	conn, _, err := zk.Connect(servers, 10*time.Second)
	if err != nil {
		// Log error or handle it according to your app's error handling strategy
		return nil, err
	}

	// Set up authentication if credentials provided
	if cfg.Username != "" && cfg.Password != "" {
		auth := fmt.Sprintf("%s:%s", cfg.Username, cfg.Password)
		conn.AddAuth("digest", []byte(auth))
	}

	repo := &Repository{
		cfg:  cfg,
		conn: conn,
	}

	// Ensure base path exists
	repo.ensurePath(cfg.Prefix)

	return repo, nil
}

// ensurePath ensures that a path exists, creating it if necessary
func (r *Repository) ensurePath(path string) error {
	if r.conn == nil || path == "" || path == "/" {
		return nil
	}

	// Check if path exists
	exists, _, err := r.conn.Exists(path)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Create parent path first
	parentIdx := strings.LastIndex(path, "/")
	if parentIdx > 0 {
		parent := path[:parentIdx]
		if err := r.ensurePath(parent); err != nil {
			return err
		}
	}

	// Create the path
	_, err = r.conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	return err
}

// GetConfig retrieves a config from ZooKeeper by key
func (r *Repository) GetConfig(ctx context.Context, key string) (*models.Config, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("zookeeper connection not established")
	}

	// Normalize key path
	nodePath := path.Join(r.cfg.Prefix, key)

	// Get data from ZooKeeper
	data, _, err := r.conn.Get(nodePath)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, fmt.Errorf("config not found for key: %s", key)
		}
		return nil, fmt.Errorf("failed to get config from ZooKeeper: %w", err)
	}

	// Unmarshal data
	var cfg models.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config for key %s: %w", key, err)
	}

	// Set the key
	cfg.Key = key
	return &cfg, nil
}

// SaveConfig saves a config to ZooKeeper
func (r *Repository) SaveConfig(ctx context.Context, cfg *models.Config) error {
	if r.conn == nil {
		return fmt.Errorf("zookeeper connection not established")
	}

	// Marshal config to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Normalize key path
	nodePath := path.Join(r.cfg.Prefix, cfg.Key)

	// Ensure parent path exists
	parentPath := nodePath[:strings.LastIndex(nodePath, "/")]
	if err := r.ensurePath(parentPath); err != nil {
		return fmt.Errorf("failed to ensure parent path: %w", err)
	}

	// Check if node exists
	exists, stat, err := r.conn.Exists(nodePath)
	if err != nil {
		return fmt.Errorf("failed to check if config exists: %w", err)
	}

	if exists {
		// Update existing node
		_, err = r.conn.Set(nodePath, data, stat.Version)
	} else {
		// Create new node
		_, err = r.conn.Create(nodePath, data, 0, zk.WorldACL(zk.PermAll))
	}

	if err != nil {
		return fmt.Errorf("failed to save config to ZooKeeper: %w", err)
	}

	return nil
}

// Close closes the ZooKeeper connection
func (r *Repository) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}
