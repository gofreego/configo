package models

import (
	"reflect"

	"github.com/gofreego/configo/configo/internal/constants"
	myutils "github.com/gofreego/configo/configo/internal/utils"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/utils"
)

/*
ConfigObject is a struct that represents a configuration object. it is used for marshalling and unmarshalling configuration objects.
also used for rendering UI for configuration management.
*/
type ConfigObject struct {
	Name        string               `json:"name"`
	Type        constants.ConfigType `json:"type"`
	Description string               `json:"description"`
	Required    bool                 `json:"required"`
	Choices     []string             `json:"choices,omitempty"`
	Value       any                  `json:"value"`
	Childrens   []ConfigObject       `json:"children,omitempty"`
}

func (co ConfigObject) Validate() error {

	if co.Required && co.Value == nil {
		return customerrors.BAD_REQUEST_ERROR("config %s is required, please pass the value", co.Name)
	}

	if co.Value != nil {
		switch co.Type {
		case constants.CONFIG_TYPE_STRING, constants.CONFIG_TYPE_BIG_TEXT:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
			if co.Required && co.Value == "" {
				return customerrors.BAD_REQUEST_ERROR("config %s is required, please pass the value of type string", co.Name)
			}
		case constants.CONFIG_TYPE_NUMBER:
			typ := reflect.TypeOf(co.Value).Kind()
			if utils.NotIn[reflect.Kind](typ, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64) {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: number", co.Name, co.Value)
			}

		case constants.CONFIG_TYPE_BOOLEAN:
			if _, ok := co.Value.(bool); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: boolean", co.Name, co.Value)
			}
		case constants.CONFIG_TYPE_JSON:
			if !myutils.IsValidJsonString(co.Value) {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: json string", co.Name, co.Value)
			}

		case constants.CONFIG_TYPE_CHOICE:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
			if len(co.Choices) == 0 {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid choices %v, Expect: non empty", co.Name, co.Choices)
			}

			if co.Required {
				if co.Value == "" {
					return customerrors.BAD_REQUEST_ERROR("config %s is required, please pass the value of type string", co.Name)
				}
				if utils.NotIn[string](co.Value.(string), co.Choices...) {
					return customerrors.BAD_REQUEST_ERROR("config %s has invalid value %s, Expect: %v", co.Name, co.Value, co.Choices)
				}
			}

		default:
			return customerrors.BAD_REQUEST_ERROR("config %s has invalid type %s, Expected : string,json,boolean,number,choice,parent,bigText", co.Name, co.Type)
		}
	}

	for _, child := range co.Childrens {
		if err := child.Validate(); err != nil {
			return err
		}
	}

	return nil
}
