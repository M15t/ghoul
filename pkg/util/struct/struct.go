package structutil

import (
	"github.com/imdatngo/mergo"
)

// ToMap converts a struct into a map, respect the json tag name
func ToMap(in interface{}) map[string]interface{} {
	// Note: json.Marshal is too heavy
	// jsonstr, _ := json.Marshal(in)
	// json.Unmarshal(jsonstr, &out)
	out := make(map[string]interface{})
	mergo.Map(&out, in, mergo.WithJSONTagLookup)
	return out
}
