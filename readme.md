# Description
Package provides nullable types based on generics.
Type implements `json.Marshaler` and `json.Unmarshaler` in order to serialize/deserialize JSON,
`driver.Valuer` and `sql.Scanner` in order to be compatible with `sql/database` package.
## Examples
### Base usage
```go
package main

import (
	"encoding/json"
	"fmt"
	null "github.com/alsiberij/alsnull"
	"os"
)

func main() {
	// Default value is null,
	// so `null` will be printed
	var nullableInt64 null.Type[int64]
	_ = json.NewEncoder(os.Stdout).Encode(nullableInt64)

	// In this case we explicitly set not null value `STRING`,
	// so "STRING" will be printed
	var nullableString = null.NewTypeWithValue("STRING")
	_ = json.NewEncoder(os.Stdout).Encode(nullableString)

	// Unmarshalling nullable values not differs from unmarshaling primitives,
	// but in case of having null bytes we will get null value,
	// so the output will be {{3.14 true} {false false}}
	// indicating that float64 has not null value of 3.14 and bool is null
	var someStruct struct {
		NullableFloat64 null.Type[float64] `json:"f"`
		NullableBool64  null.Type[bool]    `json:"b"`
	}
	_ = json.Unmarshal([]byte(`{"f":3.14, "b":null}`), &someStruct)
	fmt.Println(someStruct)
}
```
Package is compatible with `sql/database` so provided nullable types can be used for either
scanning as destinations, or performing query with parameters. Example:

```go
var nullableInt64 null.Type[int64]
// ...
_ = rows.Scan(&nullableInt64)
// ...
_, _ = conn.ExecContext(context.Background(), 'DELETE FROM table WHERE id = $1', nullableInt64)
```

Also works with pgx.

### Advanced usage
If you want to use custom package or override type marshaling/unmarshaling you must define 
your own `CustomJsonMarshaler` or `CustomJsonUnmarshaler`. The following example shows how to override
json marshaling for `time.Time`:
```go
package main

import (
	"encoding/json"
	"fmt"
	null "github.com/alsiberij/alsnull"
	"time"
)

func init() {
	// The following assignments will override json marshaling/unmarshaling for time.Time type
	// Package encoding/json is used here, but you can choose third-party ones if you need

	null.JsonMarshaler = func(value any) ([]byte, error) {
		switch V := value.(type) {
		case time.Time:
			return []byte(V.Format("\"02.01.2006\"")), nil
		default:
			return json.Marshal(value)
		}
	}

	null.JsonUnmarshaler = func(b []byte, dst any) error {
		switch V := dst.(type) {
		case *time.Time:
			t, err := time.Parse("\"02.01.2006\"", string(b))
			if err != nil {
				return err
			}
			*V = t
			return nil
		default:
			return json.Unmarshal(b, dst)
		}
	}
}

type TestStruct struct {
	NullableTime null.Type[time.Time] `json:"t"`
}

func main() {
	someStruct := TestStruct{
		NullableTime: null.NewTypeWithValue(time.Unix(1672531200, 0)),
	}

	// Json marshalling was overridden, so the result is `{"t":"01.01.2023"}`
	// Not RFC3339 which is default behavior of marshaling time.Time
	b, _ := json.Marshal(&someStruct)
	fmt.Println(string(b))

	// Json unmarshalling was also overridden so not null initial value 1672531200 will be printed
	var nextStruct TestStruct
	_ = json.Unmarshal(b, &nextStruct)
	fmt.Println(nextStruct.NullableTime.ValueOrZero().Unix())
}

```