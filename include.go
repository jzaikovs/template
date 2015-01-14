package template

import (
	. "github.com/jzaikovs/t"
)

type t_token_include struct {
	*t_token
	template *Template
}

func new_t_token_include(name string) (this *t_token_include) {
	this = new(t_token_include)
	this.t_token = newToken(name, false, false)
	return
}

func (this *t_token_include) Render(rendering *renderState, binds Map) {
	if this.template == nil {
		if this.template = Get(this.Name()); this.template == nil {
			panic("include can't be found")
		}
	}

	this.writeString(rendering, this.template.Render(binds))
}
