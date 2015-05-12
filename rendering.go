package template

import (
	"bytes"

	"github.com/jzaikovs/t"
)

type renderState struct {
	*bytes.Buffer
	globals t.Map
}

func newRendering(binds t.Map) *renderState {
	state := &renderState{new(bytes.Buffer), binds}
	return state
}
