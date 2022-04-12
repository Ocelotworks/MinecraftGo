package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteList(t *testing.T) {

	tests := [][]NBTValue{}
	for i, test := range tests {
		t.Run(fmt.Sprintf("List Test %d", i), func(t *testing.T) {
			input := NewList(test)
			write := input.Write()
			output := List{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
