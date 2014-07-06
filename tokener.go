package template

import (
	. "github.com/jzaikovs/t"
	"regexp"
	"strings"
)

type tokener interface {
	Render(*renderState, Map)
	Tokens() []tokener
	IsOpen() bool
	IsClose() bool
	IsPair(tokener) bool
	Name() string
	AddTokens([]tokener)
	Done()
}

func tokensParse(code string) []tokener {
	token := regexp.MustCompile(`<!--{#?((if|ifn|for) )?(@?[\w\._]+)#?}-->|{#?((if|ifn|for) )?(@?[\w\._]+)#?}`)
	tokens := make([]tokener, 0)
	var (
		typ, name string
		last      = 0
	)

	for _, v := range token.FindAllStringSubmatchIndex(code, -1) {
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

		tokens = append(tokens, tokensNew(tag, typ, name))
		last = match[0][1]
	}

	tokens = append(tokens, newEmptyToken(code[last:]))
	return tokens
}

func tokensCompile(tokens []tokener) []tokener {

	var last_open tokener = nil
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

func tokensNew(full, typ, name string) tokener {
	switch typ {
	case "if":
		return &tokenShow{newToken(name, full[len(full)-2] == '#', full[1] == '#'), false}
	case "ifn":
		return &tokenShow{newToken(name, full[len(full)-2] == '#', full[1] == '#'), true}
	case "for":
		return &tokenLoop{newToken(name, full[len(full)-2] == '#', full[1] == '#')}
	default:
		return newtokenValue(name)
	}
}

// function for recursive rendering
func tokensRender(buffer *renderState, tokens []tokener, binds Map) {
	for _, token := range tokens {
		token.Render(buffer, binds)
	}
}
