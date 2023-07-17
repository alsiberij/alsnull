package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alsiberij/alsnull"
	"log"
	"time"
)

type (
	//Nullable is defined in order to implement custom JSON marshaling and unmarshaling, and sql/database compatability
	Nullable[T any] struct {
		null.Type[T]
	}
)

func (t *Nullable[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("null"), nil
	}

	vPtr := t.RawValuePtr()

	switch v := any(vPtr).(type) {
	case *time.Time: //Custom time.Time marshaling
		return []byte(v.Format("\"2006-01-02\"")), nil
	default:
		return json.Marshal(vPtr)
	}
}

func (t *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		t.SetNull()
		return nil
	}

	switch ptr := any(t.RawValuePtr()).(type) {
	case *time.Time: //Custom time.Time unmarshaling
		ts, err := time.Parse("\"2006-01-02\"", string(bytes))
		if err != nil {
			return err
		}

		t.SetValue(t.DefaultValue())
		*ptr = ts
	default:
		var v T
		err := json.Unmarshal(bytes, &v)
		if err != nil {
			return err
		}

		t.SetValue(v)
	}

	return nil
}

// Value Implements driver.Valuer
func (t *Nullable[T]) Value() (driver.Value, error) {
	if t.IsNull() {
		return nil, nil
	}

	return t.RawValue(), nil
}

// Scan implements sql.Scanner
func (t *Nullable[T]) Scan(src any) error {
	switch src.(type) {
	case nil:
		t.SetNull()
	default:
		v, ok := src.(T)
		if !ok {
			return errors.New("unsupported")
		}
		t.SetValue(v)
	}
	return nil
}

func main() {
	type (
		Item struct {
			Id          int                 `json:"id"`
			Code        string              `json:"code"`
			Description Nullable[string]    `json:"description"`
			Comment     Nullable[string]    `json:"comment"`
			Ca          Nullable[time.Time] `json:"ca"`
			Ua          Nullable[time.Time] `json:"ua"`
		}

		ItemAggregator struct {
			ItemRequired     Item           `json:"item"`
			ItemNotRequired1 Nullable[Item] `json:"itemNotRequired1"`
			ItemNotRequired2 Nullable[Item] `json:"itemNotRequired2"`
		}
	)

	marshalledValue := []byte(`{"item":{"id":1,"code":"CODE","description":null,"comment":"Comment","ca":"2023-07-17","ua":null},"itemNotRequired1":null,"itemNotRequired2":{"id":2,"code":"CODE","description":"Description","comment":null,"ca":null,"ua":"2023-07-17"}}`)

	buffer := bytes.NewBuffer(marshalledValue)

	var items ItemAggregator
	err := json.NewDecoder(buffer).Decode(&items)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(&items)
	if err != nil {
		log.Fatal(err)
	}

	if string(b) != string(marshalledValue) {
		log.Fatal("NOT EQUAL")
	}

	fmt.Println("OK")
}
