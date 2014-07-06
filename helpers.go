package template

import (
	"github.com/jzaikovs/reflectutils"
	. "github.com/jzaikovs/t"
)

// function for converting structures to map
func convert_to_map(s interface{}) map[string]interface{} {
	switch v := s.(type) {
	case map[string]interface{}:
		return v
	case Map:
		return v
	default:
		return reflectutils.ToMap(s)
	}
}
