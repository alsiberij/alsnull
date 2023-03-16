package null

import (
	"database/sql/driver"
	"time"
)

type (
	// SupportedTypes provides types that can be used for generic Type
	SupportedTypes interface {
		int64 | float64 | bool | string | time.Time
	}

	// Type represents a nullable type that can used either for json marshaling/unmarshaling or sql/database interaction
	Type[T SupportedTypes] struct {
		value T
		ok    bool
	}
)

// NewType provides typed Type with given value if ok is true, otherwise value is considered as null
//
// Deprecated: second parameter is useless in current semantics. Function will be removed in future releases.
// Try NewTypeWithValue instead
func NewType[T SupportedTypes](value T, ok bool) Type[T] {
	e := Type[T]{
		ok: ok,
	}

	if ok {
		e.value = value
	}

	return e
}

// NewTypeWithValue provides not null value of type T
func NewTypeWithValue[T SupportedTypes](value T) Type[T] {
	return Type[T]{
		ok:    true,
		value: value,
	}
}

// ValueOrZero returns either actual value or default value for type T depending on whether Type is null
func (s *Type[T]) ValueOrZero() T {
	return s.value
}

// RawValue returns actual value and true if Type is not null, otherwise default value for T and false
func (s *Type[T]) RawValue() (T, bool) {
	return s.value, s.ok
}

// IsNull returns true if value is null
func (s *Type[T]) IsNull() bool {
	return !s.ok
}

// SetNull sets Type to null. Existing value will be overwritten with default one
func (s *Type[T]) SetNull() {
	s.ok = false

	var v T
	s.value = v
}

// SetValue sets Type to not null value
func (s *Type[T]) SetValue(v T) {
	s.ok = true
	s.value = v
}

// Equal checks if v is equal to current Type. Returns true if both of Type are null or have equal not null values
func (s *Type[T]) Equal(v Type[T]) bool {
	return s.ok && v.ok && (!s.ok || s.value == v.value)
}

// MarshalJSON implements json.Marshaler. Uses JsonMarshaler for marshaling not null values
func (s Type[T]) MarshalJSON() ([]byte, error) {
	if !s.ok {
		return nullBytes, nil
	}

	return JsonMarshaler(s.value)
}

// UnmarshalJSON implements json.Unmarshaler. Uses JsonUnmarshaler for unmarshaling not null values
func (s *Type[T]) UnmarshalJSON(b []byte) error {
	if string(b) == string(nullBytes) {
		s.ok = false

		var v T
		s.value = v
	} else {
		var v T
		err := JsonUnmarshaler(b, &v)
		if err != nil {
			return err
		}

		s.ok = true
		s.value = v
	}

	return nil
}

// Value implements driver.Valuer for database interacting
func (s Type[T]) Value() (driver.Value, error) {
	if !s.ok {
		return nil, nil
	}

	return s.value, nil
}

// Scan implements sql.Scanner for database interacting
func (s *Type[T]) Scan(src any) error {
	switch V := src.(type) {
	case nil:
		s.ok = false
		var v T
		s.value = v
	case []byte:
		src = string(V)
		v, ok := src.(T)
		if !ok {
			return errUnsupportedValue
		}
		s.ok = true
		s.value = v
	default:
		v, ok := src.(T)
		if !ok {
			return errUnsupportedValue
		}
		s.ok = true
		s.value = v
	}
	return nil
}
