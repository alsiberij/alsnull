# Description
Package provides generic nullable types.
<br>
Provided types are compatible with `encoding/json` package (`json.Marshaler` and `json.Unmarshaler`)
and partially with `sql/database` package (`sql.Scanner` and `driver.Valuer`).<br>
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

type (
	SimpleStructure struct {
		I  int               `json:"i"`
		S  null.Type[string] `json:"s"`
		NB null.Type[bool]   `json:"nb"`
	}

	ComplicatedStructure struct {
		F  null.Type[float64]         `json:"f"`
		O  null.Type[SimpleStructure] `json:"o"`
		NO null.Type[SimpleStructure] `json:"no"`
	}
)

func main() {
	// Default value is null,
	// so `null` will be printed
	var nullableInt64 null.Type[int64]
	_ = json.NewEncoder(os.Stdout).Encode(nullableInt64)

	// In this case we explicitly set not null value `STRING`,
	// so `"STRING"` will be printed
	var nullableString = null.WithValue("STRING")
	_ = json.NewEncoder(os.Stdout).Encode(nullableString)

	// This huge complicated structure with nullable and not nullable primitives and embedded structures
	// will be successfully encoded as expected:
	// {
	//   "f":3.14,
	//   "o":{
	//      "i":314,
	//      "s":"3,14",
	//      "nb":null
	//   },
	//   "no":null
	// }
	object := ComplicatedStructure{
		F: null.WithValue(3.14),
		O: null.WithValue(SimpleStructure{
			I:  314,
			S:  null.WithValue("3,14"),
			NB: null.Type[bool]{}, // Default is always null
		}),
		NO: null.Type[SimpleStructure]{}, // Default is always null
	}
	_ = json.NewEncoder(os.Stdout).Encode(object)

	// Unmarshalling works the same way
	var target ComplicatedStructure
	_ = json.Unmarshal([]byte(`{"f":3.14,"o":{"i":314,"s":"3,14","nb":null},"no":null}`), &target)
	fmt.Println(target)
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
Compatible with `pgx`.<br>
Keep in mind that only `int64`, `float64`, `bool`, `[]byte`, `string`, `time.Time` can be used for such purposes.

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
		NullableTime: null.WithValue(time.Unix(1672531200, 0)),
	}

	// Json marshalling was overridden, so the result is `{"t":"01.01.2023"}`
	// Not RFC3339 which is default behavior of marshaling time.Time
	b, _ := json.Marshal(&someStruct)
	fmt.Println(string(b))

	// Json unmarshalling was also overridden so not null initial value 1672531200 will be printed
	var nextStruct TestStruct
	_ = json.Unmarshal(b, &nextStruct)
	fmt.Println(nextStruct.NullableTime.RawValue().Unix())
}

```