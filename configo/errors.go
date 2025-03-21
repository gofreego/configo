package configo

import (
	"net/http"

	"github.com/gofreego/goutils/customerrors"
)

var (
	ErrConfigNotFound = customerrors.BAD_REQUEST_ERROR("config not found / not registered")
	ErrInvalidConfig  = customerrors.BAD_REQUEST_ERROR("invalid config")
)

func NewInternalServerErr(message string, params ...any) error {
	return customerrors.New(http.StatusInternalServerError, message, params...)
}
