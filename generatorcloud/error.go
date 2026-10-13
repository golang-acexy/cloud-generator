package generatorcloud

import "errors"

var (
	// ErrUnsupportedDatabase indicates a nil or unsupported database dialect.
	ErrUnsupportedDatabase = errors.New("unsupported database dialect")
	// ErrInvalidRouterPath indicates an empty router path component list.
	ErrInvalidRouterPath = errors.New("router relative model path is empty")
	// ErrInvalidAuthorityConfig indicates a missing authority expression, struct field, or column.
	ErrInvalidAuthorityConfig = errors.New("authority router configuration is incomplete")
	// ErrModelPackageRequired indicates an empty model import path.
	ErrModelPackageRequired = errors.New("model package import path is required")
	// ErrUnsupportedIDType indicates unsupported ID metadata needed by services or routers.
	ErrUnsupportedIDType = errors.New("unsupported model ID type")
	// ErrTimeRangeConfigRequired indicates missing default or allowed time fields.
	ErrTimeRangeConfigRequired = errors.New("default and allowed time range fields are required")
	// ErrDefaultTimeRangeFieldNotAllowed indicates the normalized default is absent from the allowlist.
	ErrDefaultTimeRangeFieldNotAllowed = errors.New("default time range field is not allowed")
)
