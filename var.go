package template

import (
	. "github.com/jzaikovs/t"
)

type tokenValue struct {
	*token
}

func newtokenValue(name string) *tokenValue {
	return &tokenValue{newToken(name, false, false)}
}

func (this *tokenValue) Render(rendering *renderState, binds Map) {
	if this.lang {
		rendering.WriteString(this.langval)
		return
	}

	s, ok := binds[this.name]
	if !ok {
		s, ok = rendering.globals[this.name]
	}

	if ok && s != nil {
		this.writeString(rendering, this.readStruc(s))
	}
}
