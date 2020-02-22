package dataTypes

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

func ReadNBT(buf []byte) (interface{}, int) {
	data := buf
	//gzip header
	if buf[0] == 0x1f && buf[1] == 0x8b {
		fmt.Println("NBT is compressed")
		var uncompressed []byte
		compressed := bytes.NewReader(buf)
		zr, exception := gzip.NewReader(compressed)
		if exception != nil {
			fmt.Println("Gzip error", exception)
			return nil, 0
		}
		uncompressed, exception = ioutil.ReadAll(zr)
		if exception != nil {
			fmt.Println("Exception reading gzip", exception)
			return nil, 0
		}
		data = uncompressed
	}

	compound, _ := NBTReadCompound(data)
	b, _ := json.Marshal(compound)
	fmt.Println(string(b))

	return compound, 0
}

func getNBTReadFunction(index byte, list bool) func(buf []byte) (interface{}, int) {
	if list {
		switch index {
		case 1:
			return NBTReadByte
		case 2:
			return NBTReadSignedShort
		case 3:
			return NBTReadSignedInteger
		case 4:
			return NBTReadSignedLong
		case 5:
			return NBTReadFloat
		case 6:
			return NBTReadDouble
		case 7:
			return NBTReadByteArray
		case 8:
			return NBTReadString
		case 9:
			return NBTReadList
		case 10:
			return NBTReadCompound
		case 11:
			return NBTReadIntArray
		case 12:
			return NBTReadLongArray

		}
		fmt.Println("!!!!! Unknown Type ", index)
		return nil
	}
	return NBTReadNamed(getNBTReadFunction(index, true), index)
}

func getNBTWriteFunction(index byte, list bool) func(input interface{}) []byte {
	if list {
		switch index {
		case 1:
			return WriteUnsignedByte
		case 2:
			return WriteUnsigned

		}
	}
}

type NBTNamed struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
	Type byte        `json:"type"`
}

func NBTReadNamed(function func(buf []byte) (interface{}, int), typeId byte) func(buf []byte) (interface{}, int) {
	return func(buf []byte) (interface{}, int) {
		named := NBTNamed{}
		name, cursor := NBTReadString(buf)

		named.Name = name.(string)
		if function == nil {
			return named, cursor
		}
		value, length := function(buf[cursor:])
		cursor += length
		named.Data = value
		named.Type = typeId

		return named, cursor
	}
}

func NBTReadByte(buf []byte) (interface{}, int) {
	return buf[0], 1
}

func NBTWriteByte(input interface{}) []byte {
	return []byte{input.(byte)}
}

func NBTReadUnsignedShort(buf []byte) (interface{}, int) {
	return binary.BigEndian.Uint16(buf[:2]), 2
}

func NBTWriteUnsignedShort(short interface{}) []byte {
	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, short.(uint16))
	return output
}

func NBTReadSignedShort(buf []byte) (interface{}, int) {
	val, cursor := NBTReadUnsignedShort(buf)
	return int16(val.(uint16)), cursor
}

func NBTWriteSignedShort(short interface{}) []byte {

}

func NBTReadSignedInteger(buf []byte) (interface{}, int) {
	return int32(binary.BigEndian.Uint32(buf[:4])), 4
}

func NBTReadSignedLong(buf []byte) (interface{}, int) {
	return int64(binary.BigEndian.Uint64(buf)), 8
}

func NBTReadFloat(buf []byte) (interface{}, int) {
	num := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(num), 4
}

func NBTReadDouble(buf []byte) (interface{}, int) {
	num := binary.BigEndian.Uint64(buf)
	return math.Float64frombits(num), 8
}

func NBTReadString(buf []byte) (interface{}, int) {
	stringLength, cursor := NBTReadUnsignedShort(buf)
	return string(buf[cursor : int(stringLength.(uint16))+cursor]), int(stringLength.(uint16)) + cursor
}

func NBTReadByteArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	output := make([]byte, length.(int32))
	for i := 0; int32(i) < (length.(int32)); i++ {
		uncast, length := NBTReadByte(buf[i+cursor:])
		output[i] = uncast.(byte)
		cursor += length
	}
	return output, cursor
}

func NBTReadIntArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	output := make([]byte, length.(int32))
	for i := 0; int32(i) < (length.(int32)); i++ {
		uncast, length := NBTReadSignedInteger(buf[i+cursor:])
		output[i] = byte(uncast.(int32))
		cursor += length
	}
	return output, cursor
}

func NBTReadLongArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	output := make([]byte, length.(int32))
	for i := 0; int32(i) <= (length.(int32)); i++ {
		uncast, length := NBTReadSignedLong(buf[i+cursor:])
		output[i] = uncast.(byte)
		cursor += length
	}
	return output, cursor
}

func NBTReadCompound(buf []byte) (interface{}, int) {
	return NBTRead(buf, 0)
}

func NBTRead(buf []byte, cursor int) ([]interface{}, int) {
	into := make([]interface{}, 0)
	for {
		if cursor > len(buf)-1 {
			fmt.Println("Cursor overrun")
			break
		}
		contentsType := buf[cursor]
		cursor++
		if contentsType == 0x00 {
			fmt.Println("NBT Tag End")
			break
		}
		fmt.Println("Contents Type ", contentsType)
		readMode := getNBTReadFunction(contentsType, false)
		interf, length := readMode(buf[cursor:])
		fmt.Println(interf.(NBTNamed))
		cursor += length
		into = append(into, interf)
	}
	return into, cursor
}

type NBTList struct {
	Type   byte          `json:"type"`
	Values []interface{} `json:"values"`
}

func NBTReadList(buf []byte) (interface{}, int) {
	list := NBTList{}

	list.Type = buf[0]
	//fmt.Println("List type ", list.Type)
	cursor := 1

	if list.Type == 0 {
		fmt.Println("List has 0 type")
		return list, cursor
	}

	listLength, length := NBTReadSignedInteger(buf[cursor:])
	//fmt.Println("List Length", listLength)
	cursor += length

	if listLength == 0 {
		fmt.Println("List has 0 length")
		return list, cursor
	}
	output := make([]interface{}, listLength.(int32))

	//fmt.Println(hex.Dump(buf))

	readFunction := getNBTReadFunction(list.Type, true)

	for i := 0; i < int(listLength.(int32)); i++ {
		tag, length := readFunction(buf[cursor:])
		cursor += length
		output[i] = tag
	}

	return list, cursor
}
