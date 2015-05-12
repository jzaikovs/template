package template

import (
	"github.com/jzaikovs/reflectutils"
	"github.com/jzaikovs/t"
)

// function for converting structures to map
func convertToMap(s interface{}) map[string]interface{} {
	switch v := s.(type) {
	case map[string]interface{}:
		return v
	case t.Map:
		return v
	default:
		return reflectutils.ToMap(s)
	}
}
