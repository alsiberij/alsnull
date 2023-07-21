package null

type (
	// Type represents a nullable type. Default value is null.
	Type[T any] struct {
		value T
		ok    bool
	}
)

// TypeValue returns not null value of T.
func TypeValue[T any](value T) Type[T] {
	return Type[T]{
		ok:    true,
		value: value,
	}
}

// TypeValueFromPtr returns null Type if valuePtr is nil, Type with actual value otherwise.
func TypeValueFromPtr[T any](valuePtr *T) Type[T] {
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
		return s.DefaultValue()
	}

	return s.value
}

// RawValuePtr returns pointer to actual value. Is useful only when Type is not null
func (s *Type[T]) RawValuePtr() *T {
	return &s.value
}

// CheckedValue returns actual value and true if Type is not null, default value of T and false otherwise.
func (s *Type[T]) CheckedValue() (T, bool) {
	if !s.ok {
		return s.DefaultValue(), false
	}

	return s.value, true
}

// CheckedValuePtr returns pointer to actual value and true if Type is not null, false otherwise
func (s *Type[T]) CheckedValuePtr() (*T, bool) {
	return &s.value, s.ok
}

// IsNull returns true if value is null.
func (s *Type[T]) IsNull() bool {
	return !s.ok
}

// SetNull sets Type to null.
func (s *Type[T]) SetNull() {
	s.ok = false
	s.value = s.DefaultValue()
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
		s.value = s.DefaultValue()
	} else {
		s.ok = true
		s.value = *valuePtr
	}
}

// DefaultValue returns default value of T
func (s *Type[T]) DefaultValue() T {
	var v T
	return v
}
