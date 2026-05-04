package ierrors

import "errors"

var (
	ErrServerNotFound     = errors.New("server not found")
	ErrMalformedConfig    = errors.New("malformed config")
	ErrProviderNotFound   = errors.New("provider not found")
	ErrConfigPermission   = errors.New("config file permission denied")
	ErrPathTraversal      = errors.New("path traversal detected")
	ErrUnsupportedVersion = errors.New("unsupported config version")
)
