package utils

import (
	"encoding/json"
	"io"
)

// FromJSON Deserialize JSON to the interface
func FromJSON(i interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(i)
}

// ToJSON Serialize the interface to JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}
