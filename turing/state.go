package turing

type StateFunction func(*[]Layer) State

type State struct {
	Functions map[string] StateFunction
	Accepted bool
}

func (s State) NextFn(layers []Layer) (next StateFunction, exists bool) {
	symbolA, symbolB := layers[0].Current(), layers[1].Current()
	currentSymbols := string(symbolA) + string(symbolB)
	next, exists = s.Functions[currentSymbols]
	return
}