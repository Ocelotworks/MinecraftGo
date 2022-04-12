package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteIntArray(t *testing.T) {
	tests := [][]int32{
		{1, 2, 3, 4},
		{-1, 0, 1},
		{-1, -10000},
		{},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("IntArray Test %d", i), func(t *testing.T) {
			input := NewIntArray(test)
			write := input.Write()
			output := IntArray{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
