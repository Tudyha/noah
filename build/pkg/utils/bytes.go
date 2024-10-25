package utils

import "encoding/json"

func AnyToBytes(data any) (b []byte, err error) {
	switch d := data.(type) {
	case []byte:
		b = d
	case string:
		b = []byte(d)
	default:
		b, err = json.Marshal(data)
	}
	return b, nil
}
