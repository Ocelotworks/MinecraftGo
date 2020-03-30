package dataTypes

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNBTReadWriteSignedInteger(t *testing.T) {
	for i := -1024; i < 1024; i++ {
		write := NBTWriteSignedInteger(int32(i))
		read, _ := NBTReadSignedInteger(write)
		//fmt.Printf("Write %v Read %v\n", i, read)
		assert.Equal(t, read, int32(i))
	}
}

func TestNBTReadWriteFloat(t *testing.T) {
	for i := 100; i < 1024; i++ {
		write := NBTWriteFloat(float32(i) / 10)
		read, _ := NBTReadFloat(write)
		//fmt.Printf("Write %v Read %v\n", float32(i)/10, read)
		assert.Equal(t, read, float32(i)/10)
	}
}

func TestNBTReadWriteDouble(t *testing.T) {
	for i := 100; i < 1024; i++ {
		write := NBTWriteDouble(float64(i) / 10)
		read, _ := NBTReadDouble(write)
		//fmt.Printf("Write %v Read %v\n", float64(i)/10, read)
		assert.Equal(t, read, float64(i)/10)
	}
}

func TestNBTReadWriteUnsignedShort(t *testing.T) {
	for i := 0; i < 65536; i++ {
		write := NBTWriteUnsignedShort(uint16(i))
		read, _ := NBTReadUnsignedShort(write)
		//fmt.Printf("Write %v Read %v\n", float64(i)/10, read)
		assert.Equal(t, read, uint16(i))
	}
}

func TestNBTReadWriteString(t *testing.T) {
	tests := []string{
		"hello", "this is a long string with spaces", "a", "", "weeeeeeeeeeeee", "1234567890987654321",
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("String Test %s", test), func(t *testing.T) {
			write := NBTWriteString(test)
			read, _ := NBTReadString(write)
			assert.Equal(t, read, test)
		})
	}
}

func TestNBTReadWriteList(t *testing.T) {
	tests := []NBTList{
		{
			Type:   1, //Byte
			Values: []interface{}{uint8(0x01), uint8(0x02), uint8(0x03), uint8(0x04)},
		},
		{
			Type:   2, //Short
			Values: []interface{}{int16(1), int16(2), int16(3), int16(4), int16(5)},
		},
		{
			Type:   3, //int
			Values: []interface{}{int32(1), int32(2), int32(3), int32(4), int32(5)},
		},
		{
			Type:   4, //long
			Values: []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5)},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("NBT List Test %d", i), func(t *testing.T) {
			write := NBTWriteList(test)
			read, _ := NBTReadList(write)
			assert.Equal(t, test, read)
		})
	}
}

func TestNBTReadWriteTest(t *testing.T) {
	inData, exception := ioutil.ReadFile("../nbt-test/hello_world.nbt")

	assert.Nil(t, exception)

	fmt.Println("===READING===")
	read, _ := ReadNBT(inData)
	fmt.Println("====WRITING====")
	write := NBTWriteCompound(read)
	fmt.Println("Write Finished:")
	fmt.Println(hex.Dump(write))
	fmt.Println("====READING AGAIN===")
	read2, _ := ReadNBT(write)

	assert.Equal(t, read2, read)
}

//func TestNBTReadWriteBigTest(t *testing.T) {
//	buf, exception := ioutil.ReadFile("../nbt-test/bigtest.nbt")
//
//	compressed := bytes.NewReader(buf)
//	zr, exception := gzip.NewReader(compressed)
//	assert.Nil(t, exception)
//
//	uncompressed, exception := ioutil.ReadAll(zr)
//	assert.Nil(t, exception)
//
//	inData := uncompressed
//
//	assert.Nil(t, exception)
//
//	fmt.Println("===READING===")
//	read, _ := ReadNBT(inData)
//	fmt.Println("====WRITING====")
//	write := NBTWriteCompound(read)
//	fmt.Println("====READING AGAIN===")
//	read2, _ := ReadNBT(write)
//
//	assert.Equal(t, len(read2.([]interface{})), len(read.([]interface{})))
//}

func TestNBTReadWriteLongArray(t *testing.T) {
	tests := [][]int64{
		{1, 2, 3, 4, 5, 6, 7},
		{10, 9, 8, 7, 6, 5, 4},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("NBT Long Array Test %d", i), func(t *testing.T) {
			write := NBTWriteLongArray(test)
			read, _ := NBTReadLongArray(write)
			fmt.Println(hex.Dump(write))
			fmt.Println(read)
			assert.Equal(t, test, read)
		})
	}
}
