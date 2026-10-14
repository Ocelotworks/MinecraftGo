package dataTypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadPosition(t *testing.T) {

	//27487791108035 100,-61,103

	// 1245123 0 4035 303

	// 4035 0 -61 0

	// 4035 dec - (0, -61, 0)
	output, cursor := ReadPosition([]byte{0, 0, 0, 0, 0, 0, 15, 195})
	pos := output.(*Position)
	assert.Equal(t, 8, cursor)
	assert.Equal(t, int64(0), pos.X)
	assert.Equal(t, int64(-61), pos.Y)
	assert.Equal(t, int64(0), pos.Z)

}

func TestWritePositionRoundTrips(t *testing.T) {
	for _, pos := range []Position{
		{X: 0, Y: 0, Z: 0},
		{X: 0, Y: -61, Z: 0},
		{X: 100, Y: -61, Z: 103},
		{X: -1, Y: -1, Z: -1},
		{X: 33554431, Y: 2047, Z: -33554432},
	} {
		written := WritePosition(&pos)
		assert.Equal(t, 8, len(written))

		output, cursor := ReadPosition(written)
		assert.Equal(t, 8, cursor)
		assert.Equal(t, pos, *(output.(*Position)))
	}
}

func TestWritePositionPassesThroughPackedLong(t *testing.T) {
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 15, 195}, WritePosition(int64(4035)))
}
