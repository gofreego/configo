package configo

import (
	"github.com/gofreego/goutils/customerrors"
)

var (
	ErrConfigNotFound = customerrors.BAD_REQUEST_ERROR("config not found")
	ErrInvalidConfig  = customerrors.BAD_REQUEST_ERROR("invalid config")
)
