package idex

import "encoding/json"

// UnmarshalErrorOnType is true when the error was a json.UnmarshalTypeError on type t
func UnmarshalErrorOnType(err error, t string) bool {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		if e.Value == t {
			return true
		}
	}
	return false
}
