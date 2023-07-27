package null

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type (
	//Nullable is defined in order to implement default JSON marshaling and unmarshaling, and sql/database compatability
	Nullable[T any] struct {
		Type[T]
	}
)

// NullableValue returns not null Nullable with value.
func NullableValue[T any](value T) Nullable[T] {
	return Nullable[T]{
		Type: Type[T]{
			value: value,
			ok:    true,
		},
	}
}

// NullableValueFromPtr returns null Nullable if valuePtr is nil, Nullable with actual value otherwise.
func NullableValueFromPtr[T any](valuePtr *T) Nullable[T] {
	if valuePtr == nil {
		return Nullable[T]{
			Type: Type[T]{},
		}
	}

	return Nullable[T]{
		Type: Type[T]{
			value: *valuePtr,
			ok:    true,
		},
	}
}

func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	if !t.ok {
		return nullBytes, nil
	}

	return json.Marshal(&t.value)
}

func (t *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == nullString {
		t.value = t.DefaultValue()
		t.ok = false
		return nil
	}

	var v T
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}

	t.value = v
	t.ok = true

	return nil
}

// Value Implements driver.Valuer. Supported T are:
// int, int32, int64, uint, uint32, uint64, float32, float64, bool, []byte, string, time.Time
func (t Nullable[T]) Value() (driver.Value, error) {
	var value driver.Value

	switch v := interface{}(t.value).(type) {
	case int:
		value = int64(v)
	case int32:
		value = int64(v)
	case int64:
		value = v
	case uint:
		value = int64(v)
	case uint32:
		value = int64(v)
	case uint64:
		value = int64(v)
	case float32:
		value = float64(v)
	case float64:
		value = v
	case bool:
		value = v
	case []byte:
		value = v
	case string:
		value = v
	case time.Time:
		value = v
	default:
		return nil, ErrTypeIsNotSupported
	}

	if !t.ok {
		return nil, nil
	}

	return value, nil
}

// Scan implements sql.Scanner. Supported T are:
// int, int32, int64, uint, uint32, uint64, float32, float64, bool, []byte, string, time.Time
func (t *Nullable[T]) Scan(src any) error {
	switch ptr := interface{}(&t.value).(type) {
	case *int:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = int(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *int32:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = int32(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *int64:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = value
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *uint:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = uint(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *uint32:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = uint32(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *uint64:
		if src != nil {
			value, ok := src.(int64)
			if ok {
				*ptr = uint64(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *float32:
		if src != nil {
			value, ok := src.(float64)
			if ok {
				*ptr = float32(value)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *float64:
		if src != nil {
			value, ok := src.(float64)
			if ok {
				*ptr = value
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *bool:
		if src != nil {
			value, ok := src.(bool)
			if ok {
				*ptr = value
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *string:
		if src != nil {
			str, ok := src.(string)
			if ok {
				*ptr = str
				t.ok = true
				break
			}
			b, ok := src.([]byte)
			if ok {
				*ptr = string(b)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *[]byte:
		if src != nil {
			b, ok := src.([]byte)
			if ok {
				*ptr = b
				t.ok = true
				break
			}
			str, ok := src.(string)
			if ok {
				*ptr = []byte(str)
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	case *time.Time:
		if src != nil {
			value, ok := src.(time.Time)
			if ok {
				*ptr = value
				t.ok = true
			}
		} else {
			t.value = t.DefaultValue()
			t.ok = false
		}
	default:
		return ErrScanningTypeMismatch
	}

	return nil
}
