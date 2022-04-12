package nbt

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type IntArray struct {
	Data []int32
}

func NewIntArray(arr []int32) IntArray {
	return IntArray{Data: arr}
}

func (i *IntArray) GetValue() interface{} {
	return i.Data
}

func (i *IntArray) SetValue(val interface{}) {
	i.Data = val.([]int32)
}

func (i *IntArray) GetType() int {
	return 11
}

func (i *IntArray) Read(buf []byte) int {
	length := int(int32(binary.BigEndian.Uint32(buf[:4])))
	i.Data = make([]int32, length)
	cursor := 4
	for x := 0; x < length; x++ {
		i.Data[x] = int32(binary.BigEndian.Uint32(buf[cursor : cursor+4]))
		cursor += 4
	}

	return cursor
}

func (i *IntArray) Write() []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(len(i.Data)))
	for _, value := range i.Data {
		valueBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(valueBytes, uint32(value))
		output = append(output, valueBytes...)
	}
	fmt.Println(hex.Dump(output))

	return output
}
