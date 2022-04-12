package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteDouble(t *testing.T) {
	tests := []float64{
		1, 2, 3, 4, 1.1, 3.1415, -1, -10000,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Double Test %d", i), func(t *testing.T) {
			input := NewDouble(test)
			write := input.Write()
			output := Double{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
