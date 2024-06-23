package err

import "errors"

var (
	NoRecord             = errors.New("no records found")
	InternalServiceError = errors.New("internal service error")
	DuplicateEntry       = errors.New("duplicate entry")
)
