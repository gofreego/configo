package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/configo/configo"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/cache/memory"
)

type Repo struct {
	cache cache.Cache
}

func NewRepo() *Repo {
	return &Repo{
		cache: memory.NewCache(),
	}
}

func (r *Repo) SaveConfig(ctx context.Context, cfg *configo.Config) error {
	return r.cache.Set(ctx, cfg.Key, cfg)
}

func (r *Repo) GetConfig(ctx context.Context, key string) (*configo.Config, error) {
	var cfg configo.Config
	err := r.cache.GetV(ctx, key, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Key == "" {
		return nil, nil
	}
	return &cfg, nil
}

func getRegistar(router *gin.Engine) configo.RouteRegistrar {
	return func(method, path string, handler http.HandlerFunc) error {
		ginHandler := func(c *gin.Context) {
			handler(c.Writer, c.Request)
		}
		router.Handle(method, path, ginHandler)
		return nil
	}
}

func main() {
	ctx := context.Background()
	configo, err := configo.NewConfigManager(ctx, &configo.ConfigManagerConfig{}, NewRepo())
	if err != nil {
		panic(err)
	}

	router := gin.New()

	err = configo.RegisterRoute(ctx, getRegistar(router))
	if err != nil {
		panic(err)
	}
	router.Run(":8085")
}
