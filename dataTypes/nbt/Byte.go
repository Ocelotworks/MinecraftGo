package nbt

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
	b.Data = i.(byte)
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
