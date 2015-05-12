package template

import (
	"testing"
	"time"

	. "github.com/jzaikovs/t"
)

func TestTokenShow(t *testing.T) {
	template := New()
	template.Compile(`F{if a#}O{#if a}O`)

	binds := make(Map)
	a := template.Render(binds)
	b := "FO"
	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}

	binds["a"] = true
	a = template.Render(binds)
	b = "FOO"
	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}
}

func TestTokenBinds(t *testing.T) {
	now := time.Now()
	template := New()
	template.Compile(`{time}`)
	data := map[string]interface{}{
		"time": now,
	}

	if template.Render(data) != now.String() {
		t.Error("time.Time won't render correctly")
	}
}

func TestTokenShowStruct(t *testing.T) {
	template := New()
	template.Compile(`{if error#}{error}{#if error}{ifn error#}test{#ifn error}{ifn show#}test2{#ifn show}`)
	binds := make(Map)

	binds["show"] = true
	a := template.Render(binds)
	// will show just 2nd tag becaust - 1. error is not set in if tag and show is true in if not tag
	if a != "test" {
		t.Errorf("\na:%s\nb:%s\n", a, "test")
	}

	binds["error"] = "Error"
	binds["show"] = false

	a = template.Render(binds)

	if a != "Errortest2" {
		t.Errorf("\na:%s\nb:%s\n", a, "Errortest2")
	}
}

func TestTokenShowStructAlt(t *testing.T) {
	template := New()
	template.Compile(`<!--{if show_it#}-->show_it is true<!--{#if show_it}-->`)
	binds := make(Map)

	a := template.Render(binds)
	if a != "" {
		t.Error("must be empty result, got:", a)
	}
}

func TestTokenLoop(t *testing.T) {
	//fmt.Println("TestTokenLoop======")

	template := New()
	template.Compile(`A{for f#}A{V}B{#for f}B`)
	binds := make(Map)

	a := template.Render(binds)
	b := "AB"
	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}

	arr := make([]map[string]interface{}, 5)

	for i := 0; i < 5; i++ {
		arr[i] = map[string]interface{}{"V": i}
	}

	binds["f"] = arr

	a = template.Render(binds)
	b = "AA0BA1BA2BA3BA4BB"

	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}
	//fmt.Println("======TestTokenLoop")
}

func TestTokenValue(t *testing.T) {
	//fmt.Println("TestTokenValue======")
	template := New()
	template.Compile(`A{if a#}{X}{#if a}{X}{for f#}A{V}B{#for f}B`)
	binds := make(Map)
	binds["X"] = "test"
	a := template.Render(binds)
	b := "AtestB"

	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}

	binds["a"] = true
	a = template.Render(binds)
	b = "AtesttestB"

	if a != b {
		t.Errorf("\na:%s\nb:%s\n", a, b)
	}
	//fmt.Println("======TestTokenValue")
}

// Test for struture tokens
func TestTokenStructure(t *testing.T) {
	var a struct {
		A string
		B []string
		C *int64
		D *int64
	}

	a.A = "foo"
	a.B = []string{"X", "Y"}

	var i int64 = 123
	a.D = &i

	template := New()
	template.Compile(`a.A:{a.A} {if a.C#}a.C:{a.C} {#if a.C}{if a.D#}a.D:{a.D} {#if a.D}A.a.A {for a.B#}_{a}_{#for a.B}`)
	have := template.Render(map[string]interface{}{"a": a})
	need := "a.A:foo a.D:123 A.a.A _X__Y_"
	if have != need {
		t.Log(have)
		t.Log(need)
		t.Error("Structures rendering incorrectly")
	}
}

func Tes2tTokenStructure2(t *testing.T) {
	type test struct {
		Foo int64
		Bar *int64
		Cat *int64
	}

	foo := test{}
	foo.Foo = 123
	foo.Bar = &foo.Foo

	template := New()
	template.Compile(`{if foo.Bar#}foo.Bar:{foo.Bar} {#if foo.Bar}{if foo.Cat#}foo.Cat:{foo.Cat} {#if foo.Cat}{ifn foo.Cat#}no cat{#ifn foo.Cat}`)
	have := template.Render(map[string]interface{}{"foo": foo})
	need := "foo.Bar:123 no cat"
	if have != need {
		t.Log(have)
		t.Log(need)
		t.Error("Structures rendering incorrectly")
	}
}

func TestIfAndIfnot(t *testing.T) {
	template := New()
	template.Compile(`<!--{if test#}-->x<!--{#if test}--><!--{ifn test#}-->y<!--{#ifn test}-->`)
	out := template.Render(map[string]interface{}{"test": true})
	if out != `x` {
		t.Error("`ifn` is not rendering")
	}
}
