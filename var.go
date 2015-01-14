package template

import (
	. "github.com/jzaikovs/t"
)

type t_token_var struct {
	*t_token
}

func new_t_token_var(name string) *t_token_var {
	return &t_token_var{newToken(name, false, false)}
}

func (this *t_token_var) Render(rendering *renderState, binds Map) {
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
