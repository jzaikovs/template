package template

import (
	"fmt"
	. "github.com/jzaikovs/t"
	"strings"
)

type t_token struct {
	name string
	// token is opening tag (for, if)
	opened bool
	// token is cosing tag (for, if)
	cosed bool
	// token reprezents structure
	structure bool

	// if token is part of (for, if), it my contain child tokens
	tokens []i_token

	// field path in structure
	fields []string

	// exp: this is for language tokens, idea is that language tokens will be specific {@...}, example: {@username}
	lang    bool
	langval string
}

func newEmptyToken(template string) *t_token {
	return &t_token{name: template}
}

func newToken(code string, opens, closes bool) (this *t_token) {
	this = &t_token{code, opens, closes, false, nil, nil, false, ""}

	// check if token is structure {struct.field.field}
	if i := strings.Index(code, "."); i >= 0 { // if contains dot then is structure
		this.structure = true
		this.name = code[:i]
		this.fields = strings.Split(code[i+1:], ".")
		return
	}

	this.name = code

	// language token
	if Config.HandleLang && this.name[0] == '@' {
		this.lang = true
		if lang, ok := Config.Dictionary[this.name[1:]]; ok {
			this.langval = lang
		} else {
			this.langval = this.name
		}
	}

	return this
}

func (this *t_token) Render(buffer *renderState, binds Map) {
	buffer.WriteString(this.name)
}

func (this *t_token) Tokens() []i_token {
	return this.tokens
}

func (this *t_token) IsOpen() bool {
	return this.opened
}

func (this *t_token) IsClose() bool {
	return this.cosed
}

func (this *t_token) IsPair(other i_token) bool {
	if other == nil {
		return false
	}
	return this.Name() == other.Name() && ((this.IsOpen() && other.IsClose()) || (other.IsOpen() && this.IsClose()))
}

func (this *t_token) Name() string {
	return this.name
}

func (this *t_token) AddTokens(tokens []i_token) {
	if this.tokens == nil {
		this.tokens = make([]i_token, 0)
	}

	this.tokens = append(this.tokens, tokens...)
}

func (this *t_token) String() string {
	return fmt.Sprintf("{o:%v|c:%v|n:%s|c:%v|s:%v|f:%v}", this.opened, this.cosed, this.name, this.tokens, this.structure, this.fields)
}

func (this *t_token) Done() {
	this.opened = false
	this.cosed = false
}

func (this *t_token) bindFields(s interface{}, render func(interface{})) {
	m := convert_to_map(s)
	if m == nil {
		return
	}

	n := len(this.fields) - 1 // fie need last
	var (
		val interface{}
		ok  bool
	)

	for i, field := range this.fields {
		if val, ok = m[field]; !ok {
			return
		}
		if i < n {
			if m = convert_to_map(val); m == nil {
				return
			}
		}
	}

	render(val)
}

func (this *t_token) readStruc(bind interface{}) (b interface{}) {
	if this.structure {
		this.bindFields(bind, func(val interface{}) {
			b = val
		})
		return
	}
	b = bind
	return
}

func (this *t_token) writeString(rendering *renderState, val interface{}) {
	switch v := val.(type) {
	case *int64:
		rendering.WriteString(fmt.Sprint(*v))
	case *string:
		rendering.WriteString(fmt.Sprint(*v))
	default:
		rendering.WriteString(fmt.Sprint(val))
	}
}
