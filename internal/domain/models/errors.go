package models

import "errors"

var (
	ErrCitizenReportNotFound = errors.New("citizen report not found")
	ErrInvalidCitizenReport  = errors.New("invalid citizen report data")
)
