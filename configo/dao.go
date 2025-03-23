package configo

import (
	"reflect"

	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/utils"
)

type ConfigObject struct {
	Name        string         `json:"name"`
	Type        ConfigType     `json:"type"`
	Description string         `json:"description"`
	Required    bool           `json:"required"`
	Choices     []string       `json:"choices,omitempty"`
	Value       any            `json:"value"`
	Childrens   []ConfigObject `json:"children,omitempty"`
}

func (co ConfigObject) Validate() error {

	if co.Required && co.Value == nil {
		return customerrors.BAD_REQUEST_ERROR("config %s is required, please pass the value", co.Name)
	}

	if co.Value != nil {
		switch co.Type {
		case CONFIG_TYPE_STRING, CONFIG_TYPE_BIG_TEXT:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
		case CONFIG_TYPE_NUMBER:
			typ := reflect.TypeOf(co.Value).Kind()
			if utils.NotIn[reflect.Kind](typ, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64) {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: number", co.Name, co.Value)
			}
		case CONFIG_TYPE_BOOLEAN:
			if _, ok := co.Value.(bool); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: boolean", co.Name, co.Value)
			}
		case CONFIG_TYPE_JSON:
			if _, ok := co.Value.(map[string]any); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: json", co.Name, co.Value)
			}
		case CONFIG_TYPE_CHOICE:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
		default:
			return customerrors.BAD_REQUEST_ERROR("config %s has invalid type %s", co.Name, co.Type)
		}
	}
	for _, child := range co.Childrens {
		if err := child.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ConfigInfo struct {
	ConfigKeys []string `json:"configKeys"`
}

type configMetadataResponse struct {
	ServiceInfo ServiceInfo `json:"serviceInfo"`
	ConfigInfo  ConfigInfo  `json:"configInfo"`
}
