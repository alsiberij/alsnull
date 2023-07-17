# Description
Package provides generic nullable types.
## Example
Base type
```go
type (
	//Nullable is defined in order to implement custom JSON marshaling and unmarshaling and database sql compatability
	Nullable[T any] struct {
		null.Type[T]
	}
)
```
Implementing `json.Marshaler` and `json.Unmarshaler`
```go
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
```
You can implement `sql.Scanner` and `driver.Valuer` for compatability with `sql/database` so provided nullable types can be used for either
scanning as destinations, or performing query with parameters. Example:

```go
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
```

```go
// ...
var nullableInt64 Nullable[int64]
// ...
_ = rows.Scan(&nullableInt64)
// ...
_, _ = conn.Exec(context.Background(), 'DELETE FROM table WHERE id = $1', &nullableInt64)
// ...
```
Compatible with `pgx`.<br>
Keep in mind that only `int64`, `float64`, `bool`, `[]byte`, `string`, `time.Time` can be used for such purposes.
