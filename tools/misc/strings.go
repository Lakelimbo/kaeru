package misc

import (
	"encoding/json"
)

// MapToJSONString simply converts an arbitrary struct to a string
// formatted as a JSON.
//
// (TO-DO): move it somewhere else?
func MapToJSONString[T any](m T) (string, error) {
	marshalled, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(marshalled), nil
}
