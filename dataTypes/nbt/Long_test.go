package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteLong(t *testing.T) {
	tests := []int64{
		1, 2, 3, 4, -1, -10000,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Long Test %d", i), func(t *testing.T) {
			input := NewLong(test)
			write := input.Write()
			output := Long{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
