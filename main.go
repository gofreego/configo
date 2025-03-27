package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/configo/configo"
	"github.com/gofreego/configo/configo/configs"
	"github.com/gofreego/configo/configo/models"
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

func (r *Repo) SaveConfig(ctx context.Context, cfg *models.Config) error {
	return r.cache.Set(ctx, cfg.Key, cfg)
}

func (r *Repo) GetConfig(ctx context.Context, key string) (*models.Config, error) {
	var cfg models.Config
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

type States struct {
	States map[string]string `json:"states"`
}

type RepositoryConfig struct {
	Name     string `name:"name" type:"choice" description:"Name of the repository" required:"true" choices:"memory,redis"`
	IsActive bool   `name:"isActive" type:"boolean" description:"Is the repository active" required:"false"`
	States   States `name:"statesMap" type:"json" description:"Map of states" required:"false"`
}

type KYCConfig struct {
	IsKYCEnabled bool `name:"isKycEnabled" type:"boolean" description:"Is KYC enabled" required:"true"`
}

type InvoiceConfig struct {
	IsInvoiceEnabled bool `name:"isInvoiceEnabled" type:"boolean" description:"Is Invoice enabled" required:"true"`
}

type ServiceConfig struct {
	KYCConfig     KYCConfig     `name:"KycConfig" type:"parent" description:"KYC Configuration"`
	InvoiceConfig InvoiceConfig `name:"InvoiceConfig" type:"parent" description:"Invoice Configuration"`
}

// Key implements configo.config.
func (r *RepositoryConfig) Key() string {
	return "Repository Config"
}

func main() {
	ctx := context.Background()
	configo, err := configo.NewConfigManager(ctx, &configs.ConfigManagerConfig{
		ServiceName:        "Test Service",
		ServiceDescription: "The Test Service is a scalable and modular system designed to facilitate automated and manual testing processes across different domains. It enables developers, QA engineers, and businesses to validate functionality, performance, security, and reliability of software applications, APIs, or user knowledge in an examination environment. The service provides a structured approach to defining, executing, and reporting tests, ensuring high quality and accuracy in results.",
	}, NewRepo())
	if err != nil {
		panic(err)
	}
	var repoConfig RepositoryConfig = RepositoryConfig{
		Name: "memory",
	}
	err = configo.RegisterConfig(ctx, &repoConfig)
	if err != nil {
		logger.Panic(ctx, "%v", err)
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			bytes, _ := json.Marshal(repoConfig.States)
			logger.Info(ctx, "States: %s", string(bytes))
		}
	}()

	var serviceConfig ServiceConfig = ServiceConfig{
		KYCConfig: KYCConfig{
			IsKYCEnabled: true,
		},
		InvoiceConfig: InvoiceConfig{
			IsInvoiceEnabled: true,
		},
	}
	err = configo.RegisterConfig(ctx, &serviceConfig)
	if err != nil {
		logger.Panic(ctx, "%v", err)
	}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	group := router.Group("/myservice")
	err = configo.RegisterRoute(ctx, getRegistar(group))
	if err != nil {
		panic(err)
	}

	logger.Info(ctx, "Swagger UI served at http://localhost:8085/myservice/configo/swagger/index.html")
	logger.Info(ctx, "Starting server on port 8085")
	router.Run(":8085")
}
