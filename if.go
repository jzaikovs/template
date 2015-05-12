package template

import (
	"reflect"

	"github.com/jzaikovs/t"
)

type tokenShowStruct struct {
	*tokenStruct
	invert bool
}

func (token *tokenShowStruct) Render(rendering *renderState, binds t.Map) {
	bind, ok := binds[token.name]
	if !ok {
		bind, ok = rendering.globals[token.name]
	}

	if !ok {
		if token.invert {
			// bind non boolean bind then render
			tokensRender(rendering, token.tokens, binds)
		}
		return // if bind doesn't exists then there is nothing to render
	}

	switch v := token.readStruc(bind).(type) {
	case bool:
		if v != token.invert {
			// if bind is boolean then render only if it is true
			tokensRender(rendering, token.tokens, binds)
		}
	default:
		val := reflect.ValueOf(v)
		isnil := false
		if val.Kind() == reflect.Ptr {
			isnil = val.IsNil()
		}
		if !isnil && !token.invert {
			tokensRender(rendering, token.tokens, binds)
		} else if isnil && token.invert {
			tokensRender(rendering, token.tokens, binds)
		}
	}
}
