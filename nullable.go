package null

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

func (t *Nullable[T]) MarshalJSON() ([]byte, error) {
	if !t.ok {
		return []byte("null"), nil
	}

	return json.Marshal(&t.value)
}

func (t *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
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

// Value Implements driver.Valuer
func (t *Nullable[T]) Value() (driver.Value, error) {
	if !t.ok {
		return nil, nil
	}

	return t.value, nil
}

// Scan implements sql.Scanner
func (t *Nullable[T]) Scan(src any) error {
	switch src.(type) {
	case nil:
		t.value = t.DefaultValue()
		t.ok = false
	default:
		v, ok := src.(T)
		if !ok {
			return errors.New("unsupported")
		}
		t.value = v
		t.ok = true
	}
	return nil
}
