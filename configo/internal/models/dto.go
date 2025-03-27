package models

import "github.com/gofreego/goutils/customerrors"

type UpdateConfigRequest struct {
	Id        string         `json:"id"`
	Value     []ConfigObject `json:"configs"`
	UpdatedBy string         `json:"-"`
}

func (req *UpdateConfigRequest) Validate() error {
	if req.Id == "" {
		return customerrors.BAD_REQUEST_ERROR("id is required")
	}
	if len(req.Value) == 0 {
		return customerrors.BAD_REQUEST_ERROR("value is required")
	}
	for _, config := range req.Value {
		err := config.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type GetConfigResponse struct {
	Configs []ConfigObject `json:"configs"`
}

type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ConfigInfo struct {
	ConfigKeys []string `json:"configKeys"`
}

type ConfigMetadataResponse struct {
	ServiceInfo ServiceInfo `json:"serviceInfo"`
	ConfigInfo  ConfigInfo  `json:"configInfo"`
}
