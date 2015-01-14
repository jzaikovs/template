package template

import (
	. "github.com/jzaikovs/t"
	"regexp"
	"strings"
)

const token_core_exp = `{#?((if|ifn|for|include) )?(@?[\w\._]+)#?}`

var rexp_token = regexp.MustCompile(`<!--` + token_core_exp + `-->|` + token_core_exp)

type i_token interface {
	Render(*renderState, Map)
	Tokens() []i_token
	IsOpen() bool
	IsClose() bool
	IsPair(i_token) bool
	Name() string
	AddTokens([]i_token)
	Done()
}

func tokensParse(code string) []i_token {
	tokens := make([]i_token, 0)
	var (
		typ, name string
		last      = 0
	)

	for _, v := range rexp_token.FindAllStringSubmatchIndex(code, -1) {
		match := make([][]int, 0, len(v)/2)
		// filtering and moving into groups
		for i := 0; i < len(v); i += 2 {
			if v[i] >= 0 {
				match = append(match, []int{v[i], v[i+1]})
			}
		}
		typ, name = "", ""
		tokens = append(tokens, newEmptyToken(code[last:match[0][0]]))

		if len(match) > 2 {
			// structure
			name = code[match[3][0]:match[3][1]]
			typ = code[match[2][0]:match[2][1]]

		} else {
			// just tag
			name = code[match[1][0]:match[1][1]]
		}

		tag := code[match[0][0]:match[0][1]]

		tag = strings.TrimPrefix(tag, "<!--")
		tag = strings.TrimSuffix(tag, "-->")

		tokens = append(tokens, create_token(tag, typ, name))
		last = match[0][1]
	}

	tokens = append(tokens, newEmptyToken(code[last:]))
	return tokens
}

func tokensCompile(tokens []i_token) []i_token {

	var last_open i_token = nil
	var last_idx = -1

	for i := 0; i < len(tokens); i++ {
		if tokens[i].IsOpen() {
			last_open = tokens[i]
			last_idx = i
		}

		if tokens[i].IsClose() {
			if last_open == nil {
				//todo : found closing tag without opening tag
				panic("found closing tag without opening tag")
			}

			if tokens[i].IsPair(last_open) {
				// found closing tag to last_open
				tokens[last_idx].AddTokens(tokens[last_idx+1 : i]) //todo -1 ?
				tokens[last_idx].Done()
				tokens = append(tokens[:last_idx+1], tokens[i+1:]...)

				i = 0 // start from new
				last_idx = -1
				last_open = nil
			}
		}
	}

	return tokens
}

func create_token(full, typ, name string) i_token {
	switch typ {
	case "if":
		return &t_token_show{newToken(name, full[len(full)-2] == '#', full[1] == '#'), false}
	case "ifn":
		return &t_token_show{newToken(name, full[len(full)-2] == '#', full[1] == '#'), true}
	case "for":
		return &t_token_loop{newToken(name, full[len(full)-2] == '#', full[1] == '#')}
	case "include":
		return new_t_token_include(name)
	default:
		return new_t_token_var(name)
	}
}

// function for recursive rendering
func tokensRender(buffer *renderState, tokens []i_token, binds Map) {
	for _, token := range tokens {
		token.Render(buffer, binds)
	}
}
