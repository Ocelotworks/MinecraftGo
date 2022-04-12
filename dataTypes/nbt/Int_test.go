package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteInt(t *testing.T) {
	tests := []int32{
		1, 2, 3, 4, 1000, -1, -10000,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Float Test %d", i), func(t *testing.T) {
			input := NewInt(test)
			write := input.Write()
			output := Int{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
