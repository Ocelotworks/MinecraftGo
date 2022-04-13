package nbt

import "encoding/binary"

type List struct {
	Data []NBTValue
	Type byte
}

func NewList(values []NBTValue) List {
	if len(values) == 0 {
		return List{Data: values, Type: 10} // Compound maybe?
	}
	return List{Data: values, Type: byte(values[0].GetType())}
}

func (l *List) GetValue() interface{} {
	return l.Data
}

func (l *List) SetValue(i interface{}) {
	l.Data = i.([]NBTValue)
}

func (l *List) GetType() int {
	return 9
}

func (l *List) Read(buf []byte) int {
	listType := buf[0]
	l.Type = listType
	listLength := int(binary.BigEndian.Uint32(buf[1:5]))
	l.Data = make([]NBTValue, listLength)
	// Cursor starts at 5 because of listType + listLength
	cursor := 5
	for i := 0; i < listLength; i++ {
		item := NewValue(listType)
		cursor += item.Read(buf[cursor:])
		l.Data[i] = item
	}
	return cursor
}

func (l *List) Write() []byte {
	output := []byte{l.Type}
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(l.Data)))
	output = append(output, length...)
	for _, item := range l.Data {
		output = append(output, item.Write()...)
	}

	return output
}
