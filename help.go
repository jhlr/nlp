package nlp

import (
	"fmt"

	"github.com/rakyll/portmidi"
)

func init() {
	portmidi.Initialize()
}

func (b *Board) Count(x, y uint8, j, i int8) uint8 {
	if i == 0 && j == 0 {
		return 1
	} else if !IsBoard(x, y) {
		return 0
	}
	color := b.Get(x, y)
	count := uint8(0)
	for IsBoard(x, y) {
		c := b.Get(x, y)
		if c == color {
			count++
		} else {
			break
		}
		y += uint8(i)
		x += uint8(j)
	}
	return count
}

func IsMenu(x, y uint8) bool {
	return (x == 8) != (y == 0)
}

func IsBoard(x, y uint8) bool {
	return y <= 8 && y >= 1 && x <= 7 && x >= 0
}

func (b *Board) FillBoard(color int64) {
	for y := uint8(1); y < 9; y++ {
		for x := uint8(0); x < 8; x++ {
			b.Set(x, y, color)
		}
	}
}

func (b *Board) FillMenu(c0, c1 int64) {
	for i := uint8(0); i < 8; i++ {
		b.Set(i, 0, c0)
		b.Set(8, i+1, c1)
	}
}

func Color(r, g uint8) int64 {
	r %= 4
	g %= 4
	return int64(g*16 + r)
}

func ColorRG(c int64) (uint8, uint8) {
	return uint8(c) / 16, uint8(c) % 4
}

func PrintDeviceList() {
	max := portmidi.DeviceID(portmidi.CountDevices())
	for d := portmidi.DeviceID(0); d < max; d++ {
		info := portmidi.Info(d)
		fmt.Printf("[%d] %s (%s) :", d, info.Name, info.Interface)
		if info.IsInputAvailable {
			fmt.Printf(" input")
		}
		if info.IsOutputAvailable {
			fmt.Printf(" output")
		}
		if info.IsOpened {
			fmt.Printf(" open")
		}
		fmt.Println()
	}
}
