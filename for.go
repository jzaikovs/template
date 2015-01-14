package template

import (
	"github.com/jzaikovs/reflectutils"
	. "github.com/jzaikovs/t"
	"log"
)

type t_token_loop struct {
	*t_token
}

func (this *t_token_loop) Render(rendering *renderState, binds Map) {
	bind, ok := binds[this.name]
	if !ok {
		return
	}

	switch v := this.readStruc(bind).(type) {
	case []map[string]interface{}:
		if v != nil {
			for _, row := range v {
				tokensRender(rendering, this.tokens, row)
			}
		}
	case map[string]interface{}:
		if v != nil {
			for key, value := range v {
				tokensRender(rendering, this.tokens, map[string]interface{}{"key": key, "value": value})
			}
		}
	case []interface{}:
		for _, elem := range v {
			tokensRender(rendering, this.tokens, convert_to_map(elem))
		}
	case []Map:
		if v != nil {
			for _, row := range v {
				tokensRender(rendering, this.tokens, row)
			}
		}
	case Map:
		if v != nil {
			for key, value := range v {
				tokensRender(rendering, this.tokens, map[string]interface{}{"key": key, "value": value})
			}
		}
	case []string:
		if v != nil {
			for _, str := range v {
				tokensRender(rendering, this.tokens, map[string]interface{}{this.name: str})
			}
		}
	default:
		if reflectutils.IsSlice(v) {
			reflectutils.Foreach(v, func(i int, elem interface{}) bool {
				tokensRender(rendering, this.tokens, convert_to_map(elem))
				return true
			})
		} else {
			log.Println("Loop: type not supported")
		}
	}
}
