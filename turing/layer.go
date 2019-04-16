package turing

type Layer struct {
	position int
	layers   map[int] rune
}

func New(input string) Layer {
	l := Layer{}
	l.position = 0
	l.layers = make(map[int] rune, len(input))
	for i, symbol := range input {
		l.layers[i] = symbol
	}
	return l
}

func (l *Layer) Right(r rune) {
	l.layers[l.position] = r
	l.position = l.position + 1
}

func (l *Layer) Left(r rune) {
	l.layers[l.position] = r
	l.position = l.position - 1
}

func (l *Layer) Stand(r rune) {
	l.layers[l.position] = r
}

func (l *Layer) Current() rune {
	v, k := l.layers[l.position]
	if !k {
		v = rune("_"[0])
	}
	return v
}

func (l Layer) Position() int {
	return l.position
}

func (l Layer) ToString() (layer string, position string) {
	layer = ""
	position = ""
	for index := l.position - 15; index <= l.position + 15; index++ {
		v, k := l.layers[index]
		if k {
			layer += string(v)
		} else {
			layer += "_"
		}
	}
	position += "               ^"
	return
}

func (l Layer) Count(r rune) int {
	count := 0
	for _, v := range l.layers {
		if v == r {
			count++
		}
	}
	return count
}