package nbt

import (
	"encoding/binary"
	"math"
)

type Float struct {
	Data float32
}

func NewFloat(float float32) Float {
	return Float{Data: float}
}

func (f *Float) GetValue() interface{} {
	return f.Data
}

func (f *Float) SetValue(i interface{}) {
	f.Data = i.(float32)
}

func (f *Float) GetType() int {
	return 5
}

func (f *Float) Read(buf []byte) int {
	f.Data = math.Float32frombits(binary.BigEndian.Uint32(buf))
	return 4
}

func (f *Float) Write() []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, math.Float32bits(f.Data))
	return output
}
