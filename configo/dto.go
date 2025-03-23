package configo

type UpdateConfigRequest struct {
	Key       string         `json:"key"`
	Value     []ConfigObject `json:"value"`
	UpdatedBy string         `json:"-"`
}
