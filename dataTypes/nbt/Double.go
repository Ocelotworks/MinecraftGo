package nbt

import (
	"encoding/binary"
	"math"
)

type Double struct {
	Data float64
}

func NewDouble(double float64) Double {
	return Double{Data: double}
}

func (f *Double) GetValue() interface{} {
	return f.Data
}

func (f *Double) SetValue(i interface{}) {
	f.Data = i.(float64)
}

func (f *Double) GetType() int {
	return 6
}

func (f *Double) Read(buf []byte) int {
	f.Data = math.Float64frombits(binary.BigEndian.Uint64(buf))
	return 8
}

func (f *Double) Write() []byte {
	output := make([]byte, 8)
	binary.BigEndian.PutUint64(output, math.Float64bits(f.Data))
	return output
}
