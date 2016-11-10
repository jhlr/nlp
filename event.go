package nlp

func (b *Board) Event() <-chan Event {
	return b.channel
}

func (b *Board) On(x, y uint8, k bool, foo func()) {
	b.on[Event{
		K: k, Y: y, X: x,
		Done: false,
	}] = foo
}

func (b *Board) Press(x, y uint8) {
	foo := b.on[Event{
		K: true,
		X: x, Y: y,
		Done: false,
	}]
	if foo != nil {
		foo()
	}
}

func (b *Board) Release(x, y uint8) {
	foo := b.on[Event{
		K: false,
		X: x, Y: y,
		Done: false,
	}]
	if foo != nil {
		foo()
	}
}
