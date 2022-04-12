package nbt

import "encoding/binary"

type Long struct {
	Data int64
}

func NewLong(long int64) Long {
	return Long{Data: long}
}

func (i *Long) GetValue() interface{} {
	return i.Data
}

func (i *Long) SetValue(val interface{}) {
	i.Data = val.(int64)
}

func (i *Long) GetType() int {
	return 4
}

func (i *Long) Read(buf []byte) int {
	i.Data = int64(binary.BigEndian.Uint64(buf[:8]))
	return 8
}

func (i *Long) Write() []byte {
	output := make([]byte, 8)
	binary.BigEndian.PutUint64(output, uint64(i.Data))
	return output
}
