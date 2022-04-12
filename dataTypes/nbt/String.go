package nbt

import "encoding/binary"

type String struct {
	Data string
}

func NewString(str string) String {
	return String{Data: str}
}

func (s *String) GetValue() interface{} {
	return s.Data
}

func (s *String) SetValue(i interface{}) {
	s.Data = i.(string)
}

func (s *String) GetType() int {
	return 8
}

func (s *String) Read(buf []byte) int {
	length := binary.BigEndian.Uint16(buf[:2])
	s.Data = string(buf[2 : 2+length])
	return 2 + int(length)
}

func (s *String) Write() []byte {
	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(len(s.Data)))
	return append(output, s.Data...)
}
