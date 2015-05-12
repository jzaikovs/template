package template

import (
	"github.com/jzaikovs/t"
)

type varTokenStruct struct {
	*tokenStruct
}

func newVarToken(name string) *varTokenStruct {
	return &varTokenStruct{newToken(name, false, false)}
}

func (token *varTokenStruct) Render(rendering *renderState, binds t.Map) {
	if token.lang {
		rendering.WriteString(token.langval)
		return
	}

	s, ok := binds[token.name]
	if !ok {
		s, ok = rendering.globals[token.name]
	}

	if ok && s != nil {
		token.writeString(rendering, token.readStruc(s))
	}
}
