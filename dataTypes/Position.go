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

	// X and Y sit at the top and bottom of the long, so shifting them into
	// place on a signed value sign extends them already.  Z is masked out of
	// the middle, so its 26th bit has to be carried down by hand.
	pos := Position{
		X: long >> 38,
		Y: long << 52 >> 52,
		Z: (long >> 12) & 0x3FFFFFF,
	}

	if pos.Z >= 1<<25 {
		pos.Z -= 1 << 26
	}

	return &pos, 8
}

// WritePosition packs a position into the single long the protocol sends it
// as, the way BlockPos.asLong does: 26 bits of X, then 12 of Y, then 26 of Z,
// each two's complement.  A long that is already packed is written through
// unchanged.
func WritePosition(position interface{}) []byte {
	var packed int64

	switch value := position.(type) {
	case *Position:
		packed = (value.X&0x3FFFFFF)<<38 | (value.Z&0x3FFFFFF)<<12 | value.Y&0xFFF
	case Position:
		packed = (value.X&0x3FFFFFF)<<38 | (value.Z&0x3FFFFFF)<<12 | value.Y&0xFFF
	case int64:
		packed = value
	default:
		fmt.Println("!!! Cannot write position of type", position)
	}

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(packed))
	return b
}
