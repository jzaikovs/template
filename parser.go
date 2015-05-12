package template

import (
	"regexp"
	"strings"

	"github.com/jzaikovs/t"
)

const token_core_exp = `{#?((if|ifn|for|include) )?(@?[\w\._]+)#?}`

var rexp_token = regexp.MustCompile(`<!--` + token_core_exp + `-->|` + token_core_exp)

// Token is interface for token strucktures in template
type Token interface {
	Render(*renderState, t.Map)
	Tokens() []Token
	IsOpen() bool
	IsClose() bool
	IsPair(Token) bool
	Name() string
	AddTokens([]Token)
	Done()
}

func tokensParse(code string) []Token {
	var (
		tokens    []Token
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

		tokens = append(tokens, createToken(tag, typ, name))
		last = match[0][1]
	}

	tokens = append(tokens, newEmptyToken(code[last:]))
	return tokens
}

func tokensCompile(tokens []Token) []Token {
	var (
		lastOpenToken Token
		lastIdx       = -1
	)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].IsOpen() {
			lastOpenToken = tokens[i]
			lastIdx = i
		}

		if tokens[i].IsClose() {
			if lastOpenToken == nil {
				//todo : found closing tag without opening tag
				panic("found closing tag without opening tag")
			}

			if tokens[i].IsPair(lastOpenToken) {
				// found closing tag to lastOpenToken
				tokens[lastIdx].AddTokens(tokens[lastIdx+1 : i]) //todo -1 ?
				tokens[lastIdx].Done()
				tokens = append(tokens[:lastIdx+1], tokens[i+1:]...)

				i = 0 // start from new
				lastIdx = -1
				lastOpenToken = nil
			}
		}
	}

	return tokens
}

func createToken(full, typ, name string) Token {
	switch typ {
	case "if":
		return &tokenShowStruct{newToken(name, full[len(full)-2] == '#', full[1] == '#'), false}
	case "ifn":
		return &tokenShowStruct{newToken(name, full[len(full)-2] == '#', full[1] == '#'), true}
	case "for":
		return &tokenLoopStruct{newToken(name, full[len(full)-2] == '#', full[1] == '#')}
	case "include":
		return newIncludeToken(name)
	default:
		return newVarToken(name)
	}
}

// function for recursive rendering
func tokensRender(buffer *renderState, tokens []Token, binds t.Map) {
	for _, token := range tokens {
		token.Render(buffer, binds)
	}
}
