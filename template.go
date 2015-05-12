package template

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/jzaikovs/t"
)

// Template represents single compiled template
type Template struct {
	name      string
	tokens    []Token
	cache     string
	cachetime time.Time
	static    map[string]string
}

func New() *Template {
	tmp := new(Template)
	tmp.static = make(map[string]string)
	return tmp
}

func (template *Template) Compile(code string) {
	template.tokens = tokensCompile(tokensParse(code))
}

// Method for view to load file data
func (template *Template) CompileFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}
	template.Compile(string(data))
	return nil
}

// Render renders template using passed data
func (template *Template) Render(binds t.Map) string {
	buffer := newRendering(binds)
	tokensRender(buffer, template.tokens, binds)
	return buffer.String()
}

func (template *Template) RenderCache(t time.Time, binds t.Map) string {
	if template.cachetime.Before(t) {
		// cache is older than time
		log.Println("cache-miss on", template.name)
		template.cache = template.Render(binds)
		template.cachetime = time.Now()
	}
	return template.cache
}

func (template *Template) Static(key string, binds t.Map) string {
	static, ok := template.static[key]
	if !ok {
		static = template.Render(binds)
		template.static[key] = static
	}
	return static
}
