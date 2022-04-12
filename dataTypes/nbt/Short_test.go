package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteShort(t *testing.T) {
	tests := []int16{
		1, 2, 3, 4, -1, -100,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Long Test %d", i), func(t *testing.T) {
			input := NewShort(test)
			write := input.Write()
			output := Short{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
