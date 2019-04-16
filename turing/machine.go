package turing

import (
	"fmt"
)

type TuringMachine struct {
	layers []Layer
	states []State
	debug bool
}

func MakeMachine(layerCount int, states []State, debug bool) TuringMachine {
	machine := TuringMachine{}
	machine.layers = make([]Layer, layerCount)
	machine.states = states
	machine.debug = debug
	for i := 0; i < layerCount; i++ {
		machine.layers[i] = New("")
	}
	return machine
}

func (machine TuringMachine) Run(input string) ([]Layer, State, int) {
	machine.layers[0] = New(input)

	state := machine.states[0]
	next, nextExists := state.NextFn(machine.layers)
	iteration := 0
	for nextExists {
		if iteration == 200000000 {
			break;
		}
		if machine.debug {

			fmt.Printf("Iteration #%v\n", iteration)
			for index, layer := range machine.layers {
				lString, pString := layer.ToString()
				fmt.Printf("layer [%v]:[%v]:\n", index, layer.Position())
				fmt.Printf("%s\n", lString)
				fmt.Printf("%s\n", pString)
			}

		}
		state = next(&machine.layers)
		next, nextExists = state.NextFn(machine.layers)
		iteration++
	}
	if !nextExists {
		fmt.Println("reached a dead end. exiting.")
	} else {
		fmt.Println("exceeded maximum iteration count. exiting.")
	}

	return machine.layers, state, iteration;
}
