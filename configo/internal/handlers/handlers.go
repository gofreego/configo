package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/configo/configo/internal/models"
	"github.com/gofreego/configo/configo/internal/service"
	"github.com/gofreego/configo/configo/internal/ui"
	"github.com/gofreego/configo/configo/internal/utils"
	"github.com/gofreego/configo/docs"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/goutils/response"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (c *Handler) Swagger(ctx *gin.Context) {
	logger.Debug(ctx, "Swagger UI")
	docs.SwaggerInfo.BasePath = strings.Split(ctx.Request.URL.Path, "/configo/v1/swagger")[0]
	httpSwagger.Handler()(ctx.Writer, ctx.Request)
}

// Swagger doc
// @Summary UI
// @Description UI
// @Tags Config
// @Accept json
// @Produce html
// @Success 200 {string} string "UI"
// @Failure 400 {object} any
// @Router /configo/web/ [get]
func (c *Handler) UI(ctx *gin.Context) {
	// Ensure the path is correct (handle root path and default file)
	path := strings.Split(ctx.Request.URL.Path, "/configo/v1/web")
	if len(path) < 2 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	endpoint := path[0] + "/configo/v1/web/"
	filePath := path[1]
	if filePath == "" || filePath == "/" {
		filePath = "/index.html"
	}
	logger.Debug(ctx, "UI file path: %s ", filePath)
	logger.Debug(ctx, "UI endpoint: %s", endpoint)
	// Open the requested file from embedded FS
	data, err := ui.GetStatic().ReadFile("static" + filePath)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	if filePath == "/index.html" {
		data = []byte(strings.Replace(string(data), `<base href="/">`, fmt.Sprintf(`<base href="%s">`, endpoint), 1))
	}
	// Determine content type and serve the file
	ctx.Writer.Header().Set("Content-Type", utils.GetContentType(filePath))
	ctx.Writer.Write(data)
}

// Swagger doc
// @Summary Get config
// @Description Get config by id
// @Tags Config
// @Accept json
// @Produce json
// @Param id query string true "config id"
// @Success 200 {object} GetConfigResponse
// @Failure 400 {object} any
// @Router /configo/config [get]
func (c *Handler) GetConfig(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		response.WriteError(ctx, customerrors.BAD_REQUEST_ERROR("id is required in query params"))
		return
	}
	res, err := c.service.GetConfigByKey(ctx, key)
	if err != nil {
		response.WriteError(ctx, err)
		return
	}
	response.WriteSuccess(ctx, res)
}

// Swagger doc
// @Summary Save config
// @Description Save config
// @Tags Config
// @Accept json
// @Produce json
// @Param config body UpdateConfigRequest true "config object"
// @Success 200 {string} string "config saved successfully"
// @Failure 400 {object} any
// @Router /configo/config [post]
func (c *Handler) SaveConfig(ctx *gin.Context) {
	var cfgUpdateRequest models.UpdateConfigRequest
	err := ctx.BindJSON(&cfgUpdateRequest)
	if err != nil {
		response.WriteError(ctx, customerrors.BAD_REQUEST_ERROR("failed to decode request body, Err: %s", err.Error()))
		return
	}

	err = cfgUpdateRequest.Validate()
	if err != nil {
		response.WriteError(ctx, err)
		return
	}

	err = c.service.UpdateConfig(ctx, &cfgUpdateRequest)
	if err != nil {
		response.WriteError(ctx, err)
		return
	}
	response.WriteSuccess(ctx, "config saved successfully")
}

// Swagger doc
// @Summary Get all config keys
// @Description Get all config keys
// @Tags Config
// @Accept json
// @Produce json
// @Success 200 {object} configMetadataResponse
// @Failure 400 {object} any
// @Router /configo/metadata [get]
func (c *Handler) GetConfigMetadata(ctx *gin.Context) {
	metadata, err := c.service.GetConfigsMetadata(ctx)
	if err != nil {
		response.WriteError(ctx, err)
		return
	}
	response.WriteSuccess(ctx, metadata)
}
