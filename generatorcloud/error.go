package generatorcloud

import "errors"

var (
	ErrUnsupportedDatabase              = errors.New("unsupported database dialect")
	ErrInvalidRouterPath                = errors.New("router relative model path is empty")
	ErrInvalidAuthorityConfig           = errors.New("authority router configuration is incomplete")
	ErrModelPackageRequired             = errors.New("model package import path is required")
	ErrUnsupportedIDType                = errors.New("unsupported model ID type")
	ErrTimeRangeConfigRequired          = errors.New("default and allowed time range fields are required")
	ErrDefaultTimeRangeFieldNotAllowed = errors.New("default time range field is not allowed")
)
