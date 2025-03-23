package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/configo/configo"
	_ "github.com/gofreego/configo/docs" // Import generated docs
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/cache/memory"
	"github.com/gofreego/goutils/logger"
)

// @title Config Manager APIs
// @version 1.0
// @description This API is for demonstration purposes only.
// @termsOfService http://github.com/gofreego/configo/readme.md

// @contact.name Developers
// @contact.url http://www.github.com/gofreego
// @contact.email pavanyewale1996@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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

func getRegistar(router gin.IRouter) configo.RouteRegistrar {
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
	configo, err := configo.NewConfigManager(ctx, &configo.ConfigManagerConfig{
		ServiceName:        "Test Service",
		ServiceDescription: "The Test Service is a scalable and modular system designed to facilitate automated and manual testing processes across different domains. It enables developers, QA engineers, and businesses to validate functionality, performance, security, and reliability of software applications, APIs, or user knowledge in an examination environment. The service provides a structured approach to defining, executing, and reporting tests, ensuring high quality and accuracy in results.",
	}, NewRepo())
	if err != nil {
		panic(err)
	}

	router := gin.New()
	group := router.Group("/myservice")
	err = configo.RegisterRoute(ctx, getRegistar(group))
	if err != nil {
		panic(err)
	}

	logger.Info(ctx, "Swagger UI served at http://localhost:8085/myservice/configo/swagger/index.html")
	logger.Info(ctx, "Starting server on port 8085")
	router.Run(":8085")
}
