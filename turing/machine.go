package turing

import (
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

type TuringMachine struct {
	layers  []Layer
	states  []State
	debug   bool
	gifPath string
	gif bool
}

func MakeMachine(layerCount int, states []State, debug bool, gif bool, gifPath string) TuringMachine {
	machine := TuringMachine{}
	machine.layers = make([]Layer, layerCount)
	machine.states = states
	machine.debug = debug
	machine.gif = gif
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
		bounds = image.Rect(0, 0, 300, len(machine.layers)*30 + 15)
		col = color.RGBA{200, 100, 0, 255}
		images = make([]*image.Paletted, 0)
		delays = make([]int, 0)
		defer generateGif(machine.gifPath, &images, &delays)
	}

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
		if machine.gif && iteration < 1000 {
			rgba := image.NewRGBA(bounds)

			for index, layer := range machine.layers {
				l, p := layer.ToString()
				x, y := 0, (index + 1) * 30
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
		state = next(&machine.layers)
		next, nextExists = state.NextFn(machine.layers)
		iteration++
	}
	if nextExists {
		fmt.Println("exceeded maximum iteration count. exiting.")
	}


	return machine.layers, state, iteration;
}

func generateGif(name string, images *[]*image.Paletted, delays *[]int) {
	f, _ := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0600)
	(*delays)[len(*delays)-1] = 150
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: *images,
		Delay: *delays,
	})
	fmt.Printf("Saved GIF to %s.\n", name)
}