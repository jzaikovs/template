package template

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestBig(t *testing.T) {
	template := Get("????")
	if template != nil {
		t.Error("got template when not intended")
		return
	}

	template = Get("test")

	type t_dummy_base struct {
		BaseID int
	}

	type t_dummy struct {
		*t_dummy_base
		Id   int
		Str  string
		Time time.Time
		Next *t_dummy
	}

	var struct_is_nil *t_dummy
	dummy_time, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"show_check":    true,
		"arr":           []int{0, 1, 2, 3, 4},
		"struct_is_nil": struct_is_nil,
		"GLOBAL":        "X",
		"map":           map[string]interface{}{"a": "A", "b": 15},
		"dummy":         &t_dummy{&t_dummy_base{321}, 123, "dummy", dummy_time, &t_dummy{nil, 987, "next", dummy_time, nil}},
	}

	out := strings.Split(template.Render(data), "\r\n")

	b, _ := ioutil.ReadFile("test/expected.txt")
	arr := strings.Split(string(b), "\r\n")

	if len(arr) != len(out) {
		t.Log(len(out))
		t.Log(len(arr))
		t.Error("output not what expected")
	}

	for i, line := range arr {
		arr[i] = strings.Trim(line, "\r\n")
		x := strings.Trim(out[i], "\r\n")
		if x != arr[i] {
			t.Log(arr[i])
			t.Log(x)
			t.Error("output line not what expected")
		}
	}

}
