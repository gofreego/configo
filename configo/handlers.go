package configo

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/gofreego/configo/docs"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/response"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed static/*
var content embed.FS

func (c *configManagerImpl) handleSwagger(w http.ResponseWriter, r *http.Request) {
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
func (c *configManagerImpl) handleUI(w http.ResponseWriter, r *http.Request) {
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
	data, err := content.ReadFile("static" + filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if filePath == "/index.html" {
		data = []byte(strings.Replace(string(data), `<base href="/">`, fmt.Sprintf(`<base href="%s">`, endpoint), 1))
	}
	// Determine content type and serve the file
	w.Header().Set("Content-Type", getContentType(filePath))
	w.Write(data)
}

// Swagger doc
// @Summary Get config
// @Description Get config by key
// @Tags Config
// @Accept json
// @Produce json
// @Param key query string true "config key"
// @Success 200 {object} ConfigObject
// @Failure 400 {object} any
// @Router /configs/confio/{key} [get]
func (c *configManagerImpl) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("key is required in query params"))
		return
	}
	cfg, err := c.getConfigByKey(r.Context(), key)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, cfg)
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
func (c *configManagerImpl) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfgUpdateRequest UpdateConfigRequest
	err := json.NewDecoder(r.Body).Decode(&cfgUpdateRequest)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("failed to decode request body, Err: %s", err.Error()))
		return
	}

	if cfgUpdateRequest.Key == "" {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("config not registered or invalid config key"))
		return
	}

	err = c.updateConfig(r.Context(), &cfgUpdateRequest)
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
func (c *configManagerImpl) handleGetConfigMetadata(w http.ResponseWriter, r *http.Request) {
	metadata, err := c.getConfigsMetadata(r.Context())
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, metadata)
}
