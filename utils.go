package null

import "errors"

const (
	nullString = "null"
)

var (
	nullBytes = []byte(nullString)

	ErrScanningTypeMismatch = errors.New("scanning type mismatch")
	ErrTypeIsNotSupported   = errors.New("type is not supported by driver.Valuer")
)
