package nbt

import "fmt"

type Byte struct {
	Data byte
}

func NewByte(b byte) Byte {
	return Byte{Data: b}
}

func (b *Byte) GetValue() interface{} {
	return b.Data
}

func (b *Byte) SetValue(i interface{}) {
	asByte, ok := i.(byte)
	if ok {
		b.Data = asByte
		return
	}

	asBool, ok := i.(bool)
	if ok {
		if asBool {
			b.Data = 1
		} else {
			b.Data = 0
		}
		return
	}

	panic(fmt.Sprintf("byte is neither byte nor bool: %T", i))
}

func (b *Byte) GetType() int {
	return 1
}

func (b *Byte) Read(buf []byte) int {
	b.Data = buf[0]
	return 1
}

func (b *Byte) Write() []byte {
	return []byte{b.Data}
}
