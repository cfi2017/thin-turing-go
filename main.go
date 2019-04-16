package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"zhaw.ch/thin/turing/turing"
	"zhaw.ch/thin/turing/util"
)

func main() {

	defer timeTrack(time.Now(), "turing")

	inFile := flag.String("input", "", "input file")
	inputRef := flag.String("text", "11011", "text input")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	input := *inputRef

	if *inFile != "" {
		fmt.Println("Getting input from file.")
		dat, err := ioutil.ReadFile(*inFile)
		check(err)
		input = string(dat)
	}

	fmt.Printf("Debug: %v\n", *debug)

	states := make([]turing.State, 6)
	states[5] = turing.State{Accepted: true, Functions: make(map[string]turing.StateFunction, 0)}
	states[4] = turing.State{Accepted: false, Functions: make(map[string]turing.StateFunction, 2)}
	states[3] = turing.State{Accepted: false, Functions: make(map[string]turing.StateFunction, 2)}
	states[2] = turing.State{Accepted: false, Functions: make(map[string]turing.StateFunction, 2)}
	states[1] = turing.State{Accepted: false, Functions: make(map[string]turing.StateFunction, 2)}
	states[0] = turing.State{Accepted: false, Functions: make(map[string]turing.StateFunction, 2)}

	_state5 := &states[5]
	_state4 := &states[4]
	_state3 := &states[3]
	_state2 := &states[2]
	_state1 := &states[1]
	_state0 := &states[0]

	states[0].Functions["1_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("_", 0))
		return *_state1
	}
	states[0].Functions["0_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("_", 0))
		return *_state4
	}

	states[1].Functions["0_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("0", 0))
		return *_state2
	}
	states[1].Functions["1_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("1", 0))
		return *_state1
	}

	states[2].Functions["1_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("1", 0))
		(*_layers)[1].Right(util.StringToRune("1", 0))
		return *_state2
	}
	states[2].Functions["__"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Left(util.StringToRune("_", 0))
		return *_state3
	}

	states[3].Functions["0_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Left(util.StringToRune("0", 0))
		return *_state3
	}
	states[3].Functions["1_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Left(util.StringToRune("1", 0))
		return *_state3
	}
	states[3].Functions["__"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("_", 0))
		return *_state0
	}

	states[4].Functions["0_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("_", 0))
		return *_state4
	}
	states[4].Functions["1_"] = func(_layers *[]turing.Layer) turing.State {
		(*_layers)[0].Right(util.StringToRune("_", 0))
		return *_state4
	}
	states[4].Functions["__"] = func(_layers *[]turing.Layer) turing.State {
		return *_state5
	}

	machine := turing.MakeMachine(2, states, *debug)
	layers, state, iteration := machine.Run(input)
	fmt.Printf("Word was accepted: %v\n", state.Accepted)

	fmt.Printf("Total Iteration Count: %v\n", iteration)
	for index, layer := range layers {
		lString, pString := layer.ToString()
		fmt.Printf("layer [%v]:[%v]:\n", index, layer.Position())
		fmt.Printf("%s\n", lString)
		fmt.Printf("%s\n", pString)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}