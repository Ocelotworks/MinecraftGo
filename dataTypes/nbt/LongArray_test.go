package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteLongArray(t *testing.T) {
	tests := [][]int64{
		{1, 2, 3, 4},
		{-1, 0, 1},
		{-1, -10000},
		{},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("LongArray Test %d", i), func(t *testing.T) {
			input := NewLongArray(test)
			write := input.Write()
			output := LongArray{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
