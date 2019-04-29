package turing

type StateFunction func(*[]Layer) State

type State struct {
	Functions map[string] StateFunction
	Accepted bool
}

func (s State) NextFn(layers []Layer) (next StateFunction, exists bool) {
	state := ""
	for _, layer := range layers {
		state += string(layer.Current())
	}
	next, exists = s.Functions[state]
	return
}