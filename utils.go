package null

import (
	"errors"
)

var (
	errUnsupportedValue = errors.New("unsupported value")
	nullBytes           = []byte("null")
)

func defaultValue[T any]() T {
	var v T
	return v
}
