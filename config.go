package template

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	templates = make(map[string]*Template)
	lock      sync.RWMutex
	Config    = new(configStruct)
)

type configStruct struct {
	Templates map[string]string `json:"templates"`
	// TODO: handle multiple languages
	HandleLang bool              `json:"handle_lang"`
	Dictionary map[string]string `json:"dictionary"`
}

// Get returns template with specific name
func Get(name string) *Template {
	lock.RLock()
	defer lock.RUnlock()

	if v, ok := templates[name]; ok {
		return v
	}
	log.Println("view not found!")
	return nil
}

func addTemplate(name string, template *Template) {
	lock.Lock()
	defer lock.Unlock()

	templates[name] = template
}

func init() {
	bytes, err := ioutil.ReadFile("template.json")
	if err != nil {
		log.Println(err)
		return
	}
	if err = json.Unmarshal(bytes, Config); err != nil {
		log.Println(err)
		return
	}

	for name, path := range Config.Templates {
		template := New()
		err := template.CompileFromFile(path)
		if err != nil {
			panic(err)
		}
		log.Printf("template %s (%s) created", name, path)
		addTemplate(name, template)
	}
}
