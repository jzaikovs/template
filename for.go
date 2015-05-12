package template

import (
	"log"

	"github.com/jzaikovs/reflectutils"
	"github.com/jzaikovs/t"
)

type tokenLoopStruct struct {
	*tokenStruct
}

// Render renders loop structructure from data binds by
// trying to cast bind as collection
func (token *tokenLoopStruct) Render(rendering *renderState, binds t.Map) {
	bind, ok := binds[token.name]
	if !ok {
		return
	}

	switch v := token.readStruc(bind).(type) {
	case []map[string]interface{}:
		if v != nil {
			for _, row := range v {
				tokensRender(rendering, token.tokens, row)
			}
		}
	case map[string]interface{}:
		if v != nil {
			for key, value := range v {
				tokensRender(rendering, token.tokens, map[string]interface{}{"key": key, "value": value})
			}
		}
	case []interface{}:
		for _, elem := range v {
			tokensRender(rendering, token.tokens, convertToMap(elem))
		}
	case []t.Map:
		if v != nil {
			for _, row := range v {
				tokensRender(rendering, token.tokens, row)
			}
		}
	case t.Map:
		if v != nil {
			for key, value := range v {
				tokensRender(rendering, token.tokens, map[string]interface{}{"key": key, "value": value})
			}
		}
	case []string:
		if v != nil {
			for _, str := range v {
				tokensRender(rendering, token.tokens, map[string]interface{}{token.name: str})
			}
		}
	default:
		if reflectutils.IsSlice(v) {
			reflectutils.Foreach(v, func(i int, elem interface{}) bool {
				tokensRender(rendering, token.tokens, convertToMap(elem))
				return true
			})
		} else {
			log.Println("Loop: type not supported")
		}
	}
}
