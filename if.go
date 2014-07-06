package template

import (
	. "github.com/jzaikovs/t"
	"reflect"
)

type tokenShow struct {
	*token
	invert bool
}

func (this *tokenShow) Render(rendering *renderState, binds Map) {
	bind, ok := binds[this.name]
	if !ok {
		bind, ok = rendering.globals[this.name]
	}

	if !ok {
		if this.invert {
			// bind non boolean bind then render
			tokensRender(rendering, this.tokens, binds)
		}
		return // if bind doesn't exists then there is nothing to render
	}

	switch v := this.readStruc(bind).(type) {
	case bool:
		if v != this.invert {
			// if bind is boolean then render only if it is true
			tokensRender(rendering, this.tokens, binds)
		}
	default:
		val := reflect.ValueOf(v)
		isnil := false
		if val.Kind() == reflect.Ptr {
			isnil = val.IsNil()
		}
		if !isnil && !this.invert {
			tokensRender(rendering, this.tokens, binds)
		} else if isnil && this.invert {
			tokensRender(rendering, this.tokens, binds)
		}
	}
}
