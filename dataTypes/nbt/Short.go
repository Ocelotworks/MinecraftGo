package nbt

import "encoding/binary"

type Short struct {
	Data int16
}

func NewShort(short int16) Short {
	return Short{Data: short}
}

func (s *Short) GetValue() interface{} {
	return s.Data
}

func (s *Short) SetValue(i interface{}) {
	s.Data = i.(int16)
}

func (s *Short) GetType() int {
	return 2
}

func (s *Short) Read(buf []byte) int {
	s.Data = int16(binary.BigEndian.Uint16(buf[:2]))
	return 2
}

func (s *Short) Write() []byte {
	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(s.Data))
	return output
}
