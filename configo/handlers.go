package configo

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/gofreego/configo/docs"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/response"
	httpSwagger "github.com/swaggo/http-swagger"
)

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
// @Router /configs/ui [get]
func (c *configManagerImpl) handleUI(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Form</title>
</head>
<body>
    <h2>Sample Form</h2>
    <form action="/submit" method="post">
        <label for="name">Name:</label>
        <input type="text" id="name" name="name" required>
        <br><br>
        <label for="email">Email:</label>
        <input type="email" id="email" name="email" required>
        <br><br>
        <input type="submit" value="Submit">
    </form>
</body>
</html>

	`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
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
// @Router /configs/config/{key} [get]
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
// @Router /configs/config [post]
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

type configMetadataResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keys        []string `json:"keys"`
}

// Swagger doc
// @Summary Get all config keys
// @Description Get all config keys
// @Tags Config
// @Accept json
// @Produce json
// @Success 200 {object} configMetadata
// @Failure 400 {object} any
// @Router /configs/metadata [get]
func (c *configManagerImpl) handleGetConfigMetadata(w http.ResponseWriter, r *http.Request) {
	metadata, err := c.getConfigsMetadata(r.Context())
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, metadata)
}
