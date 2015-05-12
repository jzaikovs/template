package template

import (
	"fmt"
	"strings"

	"github.com/jzaikovs/t"
)

type tokenStruct struct {
	name string
	// token is opening tag (for, if)
	opened bool
	// token is cosing tag (for, if)
	cosed bool
	// token reprezents structure
	structure bool

	// if token is part of (for, if), it my contain child tokens
	tokens []Token

	// field path in structure
	fields []string

	// exp: token is for language tokens, idea is that language tokens will be specific {@...}, example: {@username}
	lang    bool
	langval string
}

func newEmptyToken(template string) *tokenStruct {
	return &tokenStruct{name: template}
}

func newToken(code string, opens, closes bool) (token *tokenStruct) {
	token = &tokenStruct{code, opens, closes, false, nil, nil, false, ""}

	// check if token is structure {struct.field.field}
	if i := strings.Index(code, "."); i >= 0 { // if contains dot then is structure
		token.structure = true
		token.name = code[:i]
		token.fields = strings.Split(code[i+1:], ".")
		return
	}

	token.name = code

	// language token
	if Config.HandleLang && token.name[0] == '@' {
		token.lang = true
		if lang, ok := Config.Dictionary[token.name[1:]]; ok {
			token.langval = lang
		} else {
			token.langval = token.name
		}
	}

	return token
}

func (token *tokenStruct) Render(buffer *renderState, binds t.Map) {
	buffer.WriteString(token.name)
}

func (token *tokenStruct) Tokens() []Token {
	return token.tokens
}

func (token *tokenStruct) IsOpen() bool {
	return token.opened
}

func (token *tokenStruct) IsClose() bool {
	return token.cosed
}

func (token *tokenStruct) IsPair(other Token) bool {
	if other == nil {
		return false
	}
	return token.Name() == other.Name() && ((token.IsOpen() && other.IsClose()) || (other.IsOpen() && token.IsClose()))
}

func (token *tokenStruct) Name() string {
	return token.name
}

func (token *tokenStruct) AddTokens(tokens []Token) {
	if token.tokens == nil {
		token.tokens = make([]Token, 0)
	}

	token.tokens = append(token.tokens, tokens...)
}

func (token *tokenStruct) String() string {
	return fmt.Sprintf("{o:%v|c:%v|n:%s|c:%v|s:%v|f:%v}", token.opened, token.cosed, token.name, token.tokens, token.structure, token.fields)
}

func (token *tokenStruct) Done() {
	token.opened = false
	token.cosed = false
}

func (token *tokenStruct) bindFields(s interface{}, render func(interface{})) {
	m := convertToMap(s)
	if m == nil {
		return
	}

	n := len(token.fields) - 1 // fie need last
	var (
		val interface{}
		ok  bool
	)

	for i, field := range token.fields {
		if val, ok = m[field]; !ok {
			return
		}
		if i < n {
			if m = convertToMap(val); m == nil {
				return
			}
		}
	}

	render(val)
}

func (token *tokenStruct) readStruc(bind interface{}) (b interface{}) {
	if token.structure {
		token.bindFields(bind, func(val interface{}) {
			b = val
		})
		return
	}
	b = bind
	return
}

func (token *tokenStruct) writeString(rendering *renderState, val interface{}) {
	switch v := val.(type) {
	case *int64:
		rendering.WriteString(fmt.Sprint(*v))
	case *string:
		rendering.WriteString(fmt.Sprint(*v))
	default:
		rendering.WriteString(fmt.Sprint(val))
	}
}
