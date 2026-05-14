package multiline

import "errors"

// ErrEmptyPattern is returned when an empty start pattern is supplied to New.
var ErrEmptyPattern = errors.New("multiline: start pattern must not be empty")
