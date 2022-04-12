package nbt

import "encoding/binary"

type LongArray struct {
	Data []int64
}

func NewLongArray(arr []int64) LongArray {
	return LongArray{Data: arr}
}

func (i *LongArray) GetValue() interface{} {
	return i.Data
}

func (i *LongArray) SetValue(val interface{}) {
	i.Data = val.([]int64)
}

func (i *LongArray) GetType() int {
	return 12
}

func (i *LongArray) Read(buf []byte) int {
	length := int(int32(binary.BigEndian.Uint32(buf[:4])))
	i.Data = make([]int64, length)
	cursor := 4
	for x := 0; x < length; x++ {
		i.Data[x] = int64(binary.BigEndian.Uint64(buf[cursor : cursor+8]))
		cursor += 8
	}

	return cursor
}

func (i *LongArray) Write() []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(len(i.Data)))
	for _, value := range i.Data {
		valueBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(valueBytes, uint64(value))
		output = append(output, valueBytes...)
	}
	return output
}
