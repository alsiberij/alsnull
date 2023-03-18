package null

import "encoding/json"

type (
	// CustomJsonMarshaler is defined in order to allow you custom (even with not encoding/json) json marshaling
	//
	// Src is always value (not a pointer), type switch may take place.
	//
	// Default one is JsonMarshaler
	//
	// Example of implementing custom marshaler:
	//
	//	import "encoding/json"
	//	import "time"
	//
	//	func Marshaler(value any) ([]byte, error) {
	//		switch V := value.(type) {
	//		case time.Time:
	//			return []byte(V.Format("\"02.01.2006\"")), nil
	//		default:
	//			return json.Marshal(value)
	//		}
	//	}
	CustomJsonMarshaler func(src any) ([]byte, error)

	// CustomJsonUnmarshaler is defined in order to allow you custom (even with not encoding/json) json unmarshaling
	//
	// Dst is always pointer (not a value), type switch may take place.
	//
	// Default one is JsonUnmarshaler
	//
	// Example of implementing custom unmarshaler:
	//
	//	import "encoding/json"
	//	import "time"
	//
	//	func Unmarshaler(b []byte, dst any) error {
	//  	switch V := dst.(type) {
	//		case *time.Time:
	//			t, err := time.Parse("\"02.01.2006\"", string(b))
	//			if err != nil {
	//				return err
	//			}
	//			*V = t
	//			return nil
	//		default:
	//			return json.Unmarshal(b, dst)
	//		}
	//	}
	CustomJsonUnmarshaler func(b []byte, dst any) error
)

var (
	// JsonMarshaler is used for json marshaling. Default marshaler uses encoding/json package with no custom options.
	// Is called only when value is considered not as null
	JsonMarshaler CustomJsonMarshaler = func(value any) ([]byte, error) {
		return json.Marshal(value)
	}

	// JsonUnmarshaler is used for json unmarshaling. Default unmarshaler uses encoding/json package with no custom options.
	// Is called only when b not equals to `null`
	JsonUnmarshaler CustomJsonUnmarshaler = func(b []byte, dst any) error {
		return json.Unmarshal(b, dst)
	}
)
