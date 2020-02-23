package dataTypes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat(t *testing.T) {
	tests := []float32{
		0, 1, 2, 3, 4, 5, 0.1, 0.2, 0.3, 1.5, 10.5, 11.5, 300.5, 1.123,
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Float test %v", test), func(t *testing.T) {
			write := WriteFloat(test)
			read, _ := ReadFloat(write)
			fmt.Printf("Write %v Read %v\n", write, read)
			assert.Equal(t, test, read)
		})
	}
}

func TestDouble(t *testing.T) {
	tests := []float64{
		0, 1, 2, 3, 4, 5, 0.1, 0.2, 0.3, 1.5, 10.5, 11.5, 300.5, 1.123, 5, 256, 36.458,
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Double test %v", test), func(t *testing.T) {
			write := WriteDouble(test)
			read, _ := ReadDouble(write)
			fmt.Printf("Write %v Read %v\n", write, read)
			assert.Equal(t, test, read)
		})
	}
}

func TestDoubleMultiple(t *testing.T) {
	array := []float64{5, 256, 36.458}
	output := WriteDouble(array[0])
	output = append(output, WriteDouble(array[1])...)
	output = append(output, WriteDouble(array[2])...)
	read1, cursor := ReadDouble(output)
	read2, length := ReadDouble(output[cursor:])
	cursor += length
	read3, length := ReadDouble(output[cursor:])
	cursor += length
	reRead := []float64{read1.(float64), read2.(float64), read3.(float64)}
	assert.Equal(t, reRead, array)
}
