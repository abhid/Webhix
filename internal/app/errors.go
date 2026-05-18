package app

import "errors"

var (
	ErrOpenDatabase          = errors.New("open database")
	ErrMigrateDatabase       = errors.New("migrate database")
	ErrCloseDatabase         = errors.New("close database")
	ErrAuthSetup             = errors.New("auth setup")
	ErrAuthRequired          = errors.New("auth is required")
	ErrInvalidTrustedProxies = errors.New("invalid trusted proxies")
)
