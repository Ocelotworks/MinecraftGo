package nbt

import "encoding/binary"

type Int struct {
	Data int32
}

func NewInt(int int32) Int {
	return Int{Data: int}
}

func (i *Int) GetValue() interface{} {
	return i.Data
}

func (i *Int) SetValue(val interface{}) {
	i.Data = val.(int32)
}

func (i *Int) GetType() int {
	return 3
}

func (i *Int) Read(buf []byte) int {
	i.Data = int32(binary.BigEndian.Uint32(buf[:4]))
	return 4
}

func (i *Int) Write() []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(i.Data))
	return output
}
