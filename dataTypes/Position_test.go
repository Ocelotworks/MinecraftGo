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
