package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWriteCompound(t *testing.T) {
	byteValue := NewByte(0x00)
	byteArrayValue := NewByteArray([]byte{0x00, 0x01, 0x02, 0x03})
	compoundValue := NewCompound(map[string]NBTValue{
		"test1": &byteValue,
		"test2": &byteArrayValue,
		"test3": &byteArrayValue,
	})
	doubleValue := NewDouble(1.2345)
	floatValue := NewFloat(1.2345)
	intValue := NewInt(12345)
	listValue := NewList([]NBTValue{&byteValue, &byteValue, &byteValue, &byteValue})
	longValue := NewLong(12345)
	shortValue := NewShort(1234)
	stringValue := NewString("Hello NBT")

	tests := []map[string]NBTValue{
		{"test1": &byteValue},
		{"test2": &byteArrayValue},
		{"test1": &byteValue, "test2": &byteArrayValue},
		{"test1": &byteValue, "test2": &byteArrayValue, "test3": &byteArrayValue},
		{"nestedCompound": &compoundValue, "intValue": &intValue},
		{"doubleValue": &doubleValue, "floatValue": &floatValue},
		{"listValue": &listValue, "longValue": &longValue},
		{"shortValue": &shortValue, "stringValue": &stringValue},
		{},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Compound Test %d", i), func(t *testing.T) {
			input := NewCompound(test)
			write := input.Write()
			output := Compound{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
