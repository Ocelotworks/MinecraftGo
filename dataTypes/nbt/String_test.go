package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteString(t *testing.T) {
	tests := []string{
		"Hello World", "", "NBT Test", "Hello NBT",
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("String Test %d", i), func(t *testing.T) {
			input := NewString(test)
			write := input.Write()
			output := String{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
