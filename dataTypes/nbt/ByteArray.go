package nbt

import "encoding/binary"

type ByteArray struct {
	Data []byte
}

func NewByteArray(b []byte) ByteArray {
	return ByteArray{Data: b}
}

func (b *ByteArray) GetValue() interface{} {
	return b.Data
}

func (b *ByteArray) SetValue(i interface{}) {
	b.Data = i.([]byte)
}

func (b *ByteArray) GetType() int {
	return 7
}

func (b *ByteArray) Read(buf []byte) int {
	length := int32(binary.BigEndian.Uint32(buf[:4]))
	b.Data = buf[4 : 4+length]
	return 4 + int(length)
}

func (b *ByteArray) Write() []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(len(b.Data)))
	return append(output, b.Data...)
}
