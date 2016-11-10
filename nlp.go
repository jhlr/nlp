package nlp

import (
	"errors"

	"github.com/rakyll/portmidi"
)

const (
	Blank  int64 = 0
	Lime   int64 = 49
	Yellow int64 = 51
	Green  int64 = 48
	Orange int64 = 19
	Red    int64 = 3
	Ignore int64 = -1
)

type Board struct {
	out     *portmidi.Stream
	in      <-chan portmidi.Event
	main    [9][9]int64
	channel chan Event
	on      map[Event]func()
}

type Event struct {
	X, Y    uint8
	K, Done bool
}

func NewDefault() (*Board, error) {
	max := portmidi.DeviceID(portmidi.CountDevices())
	in, out := 0, 0
	ok := [2]bool{false, false}
	for d := portmidi.DeviceID(0); d < max; d++ {
		info := portmidi.Info(d)
		if !ok[0] && info.IsInputAvailable {
			in = int(d)
			ok[0] = true
		}
		if !ok[1] && info.IsOutputAvailable {
			out = int(d)
			ok[1] = true
		}
	}
	if !ok[0] || !ok[1] {
		return nil, errors.New("portmidi: missing devices")
	}
	return New(in, out)
}

func New(inId, outId int) (*Board, error) {
	in, err0 := portmidi.NewInputStream(portmidi.DeviceID(inId), 1024)
	out, err1 := portmidi.NewOutputStream(portmidi.DeviceID(outId), 1024, 0)
	if err0 != nil {
		return nil, err0
	}
	if err1 != nil {
		return nil, err1
	}
	b := &Board{
		in:      in.Listen(),
		out:     out,
		channel: make(chan Event, 1),
		on:      make(map[Event]func()),
	}
	go func() {
		for {
			ev := <-b.in
			var e Event
			e.K = ev.Data2 != 0
			e.Done = false
			if ev.Status == 176 {
				e.Y, e.X = 0, uint8(ev.Data1-104)
			} else {
				e.Y, e.X = uint8(ev.Data1/16)+1, uint8(ev.Data1%16)
			}
			if foo := b.on[e]; foo != nil {
				foo()
				e.Done = true
			}
			b.channel <- e
		}
	}()
	return b, nil
}

func (b *Board) Set(x, y uint8, color int64) {
	if color == Ignore {
		return
	}
	if x > 8 || (x == 8 && y == 0) {
		return
	}
	var key int64
	if y == 0 {
		key = 128 + int64(x)
	} else {
		key = int64(y-1)*16 + int64(x)
	}
	b.out.WriteShort(0x90, key, color)
	b.main[y][x] = color
}

func (b *Board) Get(x, y uint8) int64 {
	return b.main[y][x]
}
