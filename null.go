package null

import (
	"database/sql/driver"
)

type (
	// Type represents a nullable type. Default value is null.
	//
	// Keep in mind, that only following types T are compatible with database/sql and can be used as
	// parameters for queries or scanning destinations:
	// int64, float64, bool, []byte, string, time.Time.
	Type[T any] struct {
		value T
		ok    bool
	}
)

// WithValue returns not null value of T.
func WithValue[T any](value T) Type[T] {
	return Type[T]{
		ok:    true,
		value: value,
	}
}

// WithValueFromPtr returns null Type if valuePtr is nil, Type with actual value otherwise.
func WithValueFromPtr[T any](valuePtr *T) Type[T] {
	if valuePtr == nil {
		return Type[T]{}
	}

	return Type[T]{
		ok:    true,
		value: *valuePtr,
	}
}

// RawValue returns actual value if Type is not null, default value of T otherwise.
func (s *Type[T]) RawValue() T {
	if !s.ok {
		defaultValue[T]()
	}

	return s.value
}

// RawValuePtr returns pointer to actual value if Type is not null, nil otherwise.
func (s *Type[T]) RawValuePtr() *T {
	if !s.ok {
		return nil
	}

	return &s.value
}

// CheckedValue returns actual value and true if Type is not null, default value of T and false otherwise.
func (s *Type[T]) CheckedValue() (T, bool) {
	if !s.ok {
		return defaultValue[T](), false
	}

	return s.value, true
}

// CheckedValuePtr returns pointer to actual value and true if Type is not null, nil and false otherwise
func (s *Type[T]) CheckedValuePtr() (*T, bool) {
	if !s.ok {
		return nil, false
	}

	return &s.value, true
}

// IsNull returns true if value is null.
func (s *Type[T]) IsNull() bool {
	return !s.ok
}

// SetNull sets Type to null.
func (s *Type[T]) SetNull() {
	s.ok = false
	s.value = defaultValue[T]()
}

// SetValue sets Type to not null value v.
func (s *Type[T]) SetValue(v T) {
	s.ok = true
	s.value = v
}

// SetValueFromPtr sets Type to not null value if valuePtr is not nil, null value otherwise.
func (s *Type[T]) SetValueFromPtr(valuePtr *T) {
	if valuePtr == nil {
		s.ok = false
		s.value = defaultValue[T]()
	} else {
		s.ok = true
		s.value = *valuePtr
	}
}

// MarshalJSON implements json.Marshaler. Uses JsonMarshaler for marshaling not null values.
func (s Type[T]) MarshalJSON() ([]byte, error) {
	if !s.ok {
		return nullBytes, nil
	}

	return JsonMarshaler(s.value)
}

// UnmarshalJSON implements json.Unmarshaler. Uses JsonUnmarshaler for unmarshaling not null values.
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

// Value implements driver.Valuer for working with database.
// Returned error is nil only when T is one of the following types:
// int64, float64, bool, []byte, string, time.Time.
func (s Type[T]) Value() (driver.Value, error) {
	if !s.ok {
		return nil, nil
	}

	return s.value, nil
}

// Scan implements sql.Scanner for working with database.
// Returned err is nil only when src is nil or underlying src type and T are one of the following types:
// int64, float64, bool, []byte, string, time.Time.
func (s *Type[T]) Scan(src any) error {
	switch src.(type) {
	case nil:
		s.ok = false
		var v T
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
