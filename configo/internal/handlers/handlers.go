package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofreego/configo/configo/internal/models"
	"github.com/gofreego/configo/configo/internal/service"
	"github.com/gofreego/configo/configo/internal/ui"
	"github.com/gofreego/configo/configo/internal/utils"
	"github.com/gofreego/configo/docs"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/response"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (c *Handler) Swagger(w http.ResponseWriter, r *http.Request) {
	docs.SwaggerInfo.BasePath = strings.Split(r.URL.Path, "/configo/swagger")[0]
	httpSwagger.Handler()(w, r)
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
func (c *Handler) UI(w http.ResponseWriter, r *http.Request) {
	// Ensure the path is correct (handle root path and default file)
	path := strings.Split(r.URL.Path, "/configo/web")
	if len(path) < 2 {
		http.NotFound(w, r)
		return
	}
	endpoint := path[0] + "/configo/web/"
	filePath := path[1]
	if filePath == "" || filePath == "/" {
		filePath = "/index.html"
	}

	// Open the requested file from embedded FS
	data, err := ui.GetStatic().ReadFile("static" + filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if filePath == "/index.html" {
		data = []byte(strings.Replace(string(data), `<base href="/">`, fmt.Sprintf(`<base href="%s">`, endpoint), 1))
	}
	// Determine content type and serve the file
	w.Header().Set("Content-Type", utils.GetContentType(filePath))
	w.Write(data)
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
func (c *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("id is required in query params"))
		return
	}
	res, err := c.service.GetConfigByKey(r.Context(), key)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, res)
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
func (c *Handler) SaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfgUpdateRequest models.UpdateConfigRequest
	err := json.NewDecoder(r.Body).Decode(&cfgUpdateRequest)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("failed to decode request body, Err: %s", err.Error()))
		return
	}

	err = cfgUpdateRequest.Validate()
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}

	err = c.service.UpdateConfig(r.Context(), &cfgUpdateRequest)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, "config saved successfully")
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
func (c *Handler) GetConfigMetadata(w http.ResponseWriter, r *http.Request) {
	metadata, err := c.service.GetConfigsMetadata(r.Context())
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, metadata)
}
