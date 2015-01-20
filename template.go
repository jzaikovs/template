package template

import (
	. "github.com/jzaikovs/t"
	"io/ioutil"
	"log"
	"time"
)

type Template struct {
	name      string
	tokens    []i_token
	cache     string
	cachetime time.Time
	static    map[string]string
}

func New() *Template {
	tmp := new(Template)
	tmp.static = make(map[string]string)
	return tmp
}

func (this *Template) Compile(template string) {
	this.tokens = tokensCompile(tokensParse(template))
}

// Method for view to load file data
func (this *Template) CompileFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}
	this.Compile(string(data))
	return nil
}

// Method for rendering view
func (this *Template) Render(binds Map) string {
	buffer := newRendering(binds)
	tokensRender(buffer, this.tokens, binds)
	return buffer.String()
}

func (this *Template) RenderCache(t time.Time, binds Map) string {
	if this.cachetime.Before(t) {
		// cache is older than time
		log.Println("cache-miss on", this.name)
		this.cache = this.Render(binds)
		this.cachetime = time.Now()
	}
	return this.cache
}

func (this *Template) Static(key string, binds Map) string {
	static, ok := this.static[key]
	if !ok {
		static = this.Render(binds)
		this.static[key] = static
	}
	return static
}
