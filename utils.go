package null

import (
	"errors"
)

var (
	errUnsupportedValue = errors.New("unsupported value")
	nullBytes           = []byte("null")
)
