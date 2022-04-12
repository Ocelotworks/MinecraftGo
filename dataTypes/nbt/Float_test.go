package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteFloat(t *testing.T) {
	tests := []float32{
		1, 2, 3, 4, 1.1, 3.1415, -1, -10000,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Float Test %d", i), func(t *testing.T) {
			input := NewFloat(test)
			write := input.Write()
			output := Float{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
