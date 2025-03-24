package configo

type UpdateConfigRequest struct {
	Key       string         `json:"id"`
	Value     []ConfigObject `json:"configs"`
	UpdatedBy string         `json:"-"`
}

type GetConfigResponse struct {
	Configs []ConfigObject `json:"configs"`
}
