package template

import (
	"github.com/jzaikovs/t"
)

type tokenIncludeStruct struct {
	*tokenStruct
	template *Template
}

func newIncludeToken(name string) (token *tokenIncludeStruct) {
	token = new(tokenIncludeStruct)
	token.tokenStruct = newToken(name, false, false)
	return
}

func (token *tokenIncludeStruct) Render(rendering *renderState, binds t.Map) {
	if token.template == nil {
		if token.template = Get(token.Name()); token.template == nil {
			panic("include can't be found")
		}
	}

	token.writeString(rendering, token.template.Render(binds))
}
