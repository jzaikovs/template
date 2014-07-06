package template

import (
	"bytes"
	. "github.com/jzaikovs/t"
)

type renderState struct {
	*bytes.Buffer
	globals Map
}

func newRendering(binds Map) *renderState {
	state := &renderState{new(bytes.Buffer), binds}
	return state
}
