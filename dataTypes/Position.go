package dataTypes

import (
	"encoding/binary"
	"fmt"
)

type Position struct {
	X int64
	Y int64
	Z int64
}

func (p *Position) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}

func ReadPosition(buf []byte) (interface{}, int) {
	long := int64(binary.BigEndian.Uint64(buf[:8]))

	pos := Position{
		X: long >> 38,
		Y: long << 52 >> 52,
		Z: (long >> 12) & 0x3FFFFFF,
	}

	if pos.X >= 2<<25 {
		pos.X -= 2 << 26
	}
	if pos.Y >= 2<<11 {
		pos.Y -= 2 << 12
	}
	if pos.Z >= 2<<25 {
		pos.Z -= 2 << 26
	}

	return &pos, 8
}

func WritePosition(long interface{}) []byte {
	// TODO
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(long.(int64)))
	return b
}
