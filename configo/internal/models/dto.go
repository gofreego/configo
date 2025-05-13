package models

import "github.com/gofreego/goutils/customerrors"

type UpdateConfigRequest struct {
	Key       string         `json:"key"`
	Configs   []ConfigObject `json:"configs"`
	UpdatedBy string         `json:"-"`
}

func (req *UpdateConfigRequest) Validate() error {
	if req.Key == "" {
		return customerrors.BAD_REQUEST_ERROR("key is required")
	}
	if len(req.Configs) == 0 {
		return customerrors.BAD_REQUEST_ERROR("configs are required")
	}
	for _, config := range req.Configs {
		err := config.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type GetConfigResponse struct {
	Key       string         `json:"key"`
	Configs   []ConfigObject `json:"configs"`
	UpdatedBy string         `json:"updatedBy"`
	UpdatedAt int64          `json:"updatedAt"`
	CreatedAt int64          `json:"createdAt"`
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
