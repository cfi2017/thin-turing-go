package turing

import (
	"bufio"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"time"
)

type TuringMachine struct {
	layers  []Layer
	states  []State
	debug   bool
	gifPath string
	gif bool
	fast bool
	maxIterations int
	stack bool
	delay time.Duration
}

func MakeMachine(layerCount int,
	states []State,
	debug bool,
	fast bool,
	delay int,
	maxIterations int,
	stack bool,
	gif bool,
	gifPath string) TuringMachine {
	machine := TuringMachine{}
	machine.layers = make([]Layer, layerCount)
	machine.states = states
	machine.debug = debug
	machine.gif = gif
	machine.fast = fast
	machine.stack = stack
	machine.maxIterations = maxIterations
	machine.delay = time.Duration(delay)
	machine.gifPath = gifPath
	for i := 0; i < layerCount; i++ {
		machine.layers[i] = New("")
	}
	return machine
}

func (machine TuringMachine) Run(input string) ([]Layer, State, int) {

	var images []*image.Paletted
	var delays []int
	var bounds image.Rectangle
	var col color.Color
	if machine.gif {
		bounds = image.Rect(0, 0, 240, len(machine.layers)*30 + 15)
		col = color.RGBA{200, 100, 0, 255}
		images = make([]*image.Paletted, 0)
		delays = make([]int, 0)
		defer generateGif(machine.gifPath, &images, &delays)
	}

	var reader *bufio.Reader
	if !machine.fast && machine.debug {
		reader = bufio.NewReader(os.Stdin)
	}

	machine.layers[0] = New(input)

	state := machine.states[0]
	next, nextExists := state.NextFn(machine.layers)
	iteration := 0
	for nextExists {
		if iteration == machine.maxIterations {
			break;
		}
		if !machine.fast {

			if machine.stack {
				fmt.Print("\033[H\033[2J")
			}
			fmt.Printf("Iteration #%v\n", iteration)
			for index, layer := range machine.layers {
				lString, pString := layer.ToString()
				fmt.Printf("layer [%v]:[%v]:\n", index, layer.Position())
				fmt.Printf("%s\n", lString)
				fmt.Printf("%s\n", pString)
			}

		}
		if machine.gif && iteration < 1000 {
			rgba := image.NewRGBA(bounds)

			for index, layer := range machine.layers {
				l, p := layer.ToString()
				x, y := 10, (index + 1) * 30
				point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
				d := &font.Drawer{
					Dst:  rgba,
					Src:  image.NewUniform(col),
					Face: basicfont.Face7x13,
					Dot:  point,
				}
				d.DrawString(l)
				y = ((index + 1) * 30) + 15
				point = fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
				d = &font.Drawer{
					Dst:  rgba,
					Src:  image.NewUniform(col),
					Face: basicfont.Face7x13,
					Dot:  point,
				}
				d.DrawString(p)
			}

			img := image.NewPaletted(bounds, palette.Plan9)
			draw.Draw(img, img.Rect, rgba, bounds.Min, draw.Over)
			images = append(images, img)
			delays = append(delays, 25)
		}
		if machine.debug && !machine.fast {
			reader.ReadByte()
		}
		state = next(&machine.layers)
		next, nextExists = state.NextFn(machine.layers)
		iteration++
		if !machine.fast && machine.delay != 0 {
			time.Sleep(time.Millisecond * machine.delay)
		}
	}
	if nextExists {
		fmt.Println("Exceeded maximum iteration count. Exiting.")
	}

	return machine.layers, state, iteration;
}

func generateGif(name string, images *[]*image.Paletted, delays *[]int) {
	f, _ := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0600)
	(*delays)[len(*delays)-1] = 150
	defer f.Close()
	err := gif.EncodeAll(f, &gif.GIF{
		Image: *images,
		Delay: *delays,
	})
	if err != nil {
		log.Fatalf("Error encoding GIF: %v", err)
	} else {
		fmt.Printf("Saved GIF to %s.\n", name)
	}
}