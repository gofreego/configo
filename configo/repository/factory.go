package repository

import (
	"github.com/gofreego/configo/configo/repository/consul"
	"github.com/gofreego/configo/configo/repository/database"
	"github.com/gofreego/configo/configo/repository/memory"
	"github.com/gofreego/configo/configo/repository/zookeeper"
	"github.com/gofreego/goutils/customerrors"
)

type Name string

const (
	// memeory can be used for testing purposes, we recommend using consul, zookeeper or any other repository for production
	Memory    Name = "memory"
	Consul    Name = "consul"
	Zookeeper Name = "zookeeper"
	Database  Name = "database"
)

type Config struct {
	Name      Name
	Consul    consul.Config
	Zookeeper zookeeper.Config
}

func NewRepository(cfg *Config) (Repository, error) {
	switch cfg.Name {
	case Memory:
		return memory.NewRepository()
	case Consul:
		return consul.NewRepository(&cfg.Consul)
	case Zookeeper:
		return zookeeper.NewRepository(&cfg.Zookeeper)
	case Database:
		return database.NewRepository()
	default:
		return nil, customerrors.BAD_REQUEST_ERROR("invalid repository name")
	}
}
