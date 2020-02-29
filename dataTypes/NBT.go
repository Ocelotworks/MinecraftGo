package dataTypes

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
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
	//b, _ := json.Marshal(compound)
	//fmt.Println(string(b))

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

func getNBTWriteFunction(index byte, name string) func(input interface{}) []byte {
	if name == "none" {
		switch index {
		case 1:
			return NBTWriteByte
		case 2:
			return NBTWriteSignedShort
		case 3:
			return NBTWriteSignedInteger
		case 4:
			return NBTWriteSignedLong
		case 5:
			return NBTWriteFloat
		case 6:
			return NBTWriteDouble
		case 7:
			return NBTWriteByteArray
		case 8:
			return NBTWriteString
		case 9:
			return NBTWriteList
		case 10:
			return NBTWriteCompound
		case 11:
			return NBTWriteIntArray
		case 12:
			return NBTWriteLongArray
		}
		fmt.Println("!!!!! Unknown Type ", index)
		return nil
	}
	fmt.Println("write function ", name)
	return NBTWriteNamed(getNBTWriteFunction(index, "none"), index, name)
}

type NBTNamed struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
	Type byte        `json:"type"`
}

func NBTReadNamed(function func(buf []byte) (interface{}, int), typeId byte) func(buf []byte) (interface{}, int) {
	return func(buf []byte) (interface{}, int) {
		//fmt.Printf("Reading named type of %d\n", typeId)
		named := NBTNamed{}
		name, cursor := NBTReadString(buf)

		named.Name = name.(string)
		fmt.Println("--- TAG START: ", named.Name)
		if function == nil {
			fmt.Println("??? Nil function reading Named NBT Tag")
			return named, cursor
		}
		value, length := function(buf[cursor:])
		cursor += length
		named.Data = value
		named.Type = typeId
		fmt.Println("--- TAG END: ", named.Name)

		return named, cursor
	}
}

func NBTWriteNamed(function func(input interface{}) []byte, index byte, name string) func(input interface{}) []byte {
	return func(input interface{}) []byte {
		fmt.Printf("Writing named tag '%s' of type %d\n", name, index)
		output := append(NBTWriteByte(index), NBTWriteString(name)...)
		output = append(output, function(input)...)
		//fmt.Println(hex.Dump(output))
		return output
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
	return NBTWriteUnsignedShort(uint16(short.(int16)))
}

func NBTReadSignedInteger(buf []byte) (interface{}, int) {
	return int32(binary.BigEndian.Uint32(buf[:4])), 4
}

func NBTWriteSignedInteger(num interface{}) []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(num.(int32)))
	return output
}

func NBTReadSignedLong(buf []byte) (interface{}, int) {
	return int64(binary.BigEndian.Uint64(buf)), 8
}

func NBTWriteSignedLong(long interface{}) []byte {
	output := make([]byte, 8)
	binary.BigEndian.PutUint64(output, uint64(long.(int64)))
	return output
}

func NBTReadFloat(buf []byte) (interface{}, int) {
	num := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(num), 4
}

func NBTWriteFloat(float interface{}) []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, math.Float32bits(float.(float32)))
	return output
}

func NBTReadDouble(buf []byte) (interface{}, int) {
	num := binary.BigEndian.Uint64(buf)
	return math.Float64frombits(num), 8
}

func NBTWriteDouble(float interface{}) []byte {
	output := make([]byte, 8)
	binary.BigEndian.PutUint64(output, math.Float64bits(float.(float64)))
	return output
}

func NBTReadString(buf []byte) (interface{}, int) {
	stringLength, cursor := NBTReadUnsignedShort(buf)
	//fmt.Println("Read string length ", stringLength)
	return string(buf[cursor : int(stringLength.(uint16))+cursor]), int(stringLength.(uint16)) + cursor
}

func NBTWriteString(str interface{}) []byte {
	output := []byte(str.(string))
	//fmt.Println("Writing string length ", uint16(len(output)))
	output = append(NBTWriteUnsignedShort(uint16(len(output))), output...)
	return output
}

func NBTReadByteArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	//fmt.Printf("Read Byte Array length %d cursor %d\n", length, cursor)
	output := make([]byte, length.(int32))
	for i := 0; int32(i) < (length.(int32)); i++ {
		//fmt.Printf("Reading byte %d cursor pos %d\n", i, cursor)
		uncast, length := NBTReadByte(buf[cursor:])
		output[i] = uncast.(byte)
		cursor += length
	}
	return output, cursor
}

func NBTWriteByteArray(arr interface{}) []byte {
	array := arr.([]byte)
	return append(NBTWriteSignedInteger(int32(len(array))), array...)
}

func NBTReadIntArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	//fmt.Println("Int array length ", length)
	output := make([]byte, length.(int32))
	for i := 0; int32(i) < (length.(int32)); i++ {
		//fmt.Println("Reading int num ", i)
		uncast, length := NBTReadSignedInteger(buf[i+cursor:])
		output[i] = byte(uncast.(int32))
		cursor += length
	}
	return output, cursor
}

func NBTWriteIntArray(arr interface{}) []byte {
	array := arr.([]int32)
	output := NBTWriteSignedInteger(len(array) * 4)
	for _, val := range array {
		output = append(output, NBTWriteSignedInteger(val)...)
	}
	return output
}

func NBTReadLongArray(buf []byte) (interface{}, int) {
	length, cursor := NBTReadSignedInteger(buf)
	output := make([]int64, length.(int32))
	//fmt.Println("Long Array length ", length)
	//fmt.Println(hex.Dump(buf))
	for i := 0; int32(i) < (length.(int32)); i++ {
		//fmt.Println("Reading long num", i, " at cursor pos ", cursor)
		uncast, length := NBTReadSignedLong(buf[cursor:])
		//fmt.Println("Received as ", uncast)
		output[i] = uncast.(int64)
		cursor += length
	}
	return output, cursor
}

func NBTWriteLongArray(arr interface{}) []byte {
	array := arr.([]int64)
	output := NBTWriteSignedInteger(int32(len(array)))
	for _, val := range array {
		output = append(output, NBTWriteSignedLong(val)...)
	}
	return output
}

func NBTReadCompound(buf []byte) (interface{}, int) {
	//fmt.Println("Reading Compound")
	return NBTRead(buf, 0)
}

func NBTRead(buf []byte, cursor int) ([]interface{}, int) {
	into := make([]interface{}, 0)
	//fmt.Println("Reading Compound at cursor ", cursor)
	for {
		if cursor >= len(buf) {
			//fmt.Println("NBT Read Cursor overrun ", cursor, len(buf))
			break
		}
		contentsType := buf[cursor]
		cursor++
		if contentsType == 0x00 {
			//fmt.Println("NBT Tag End")
			break
		}
		//fmt.Println("NBTRead Contents Type ", contentsType)
		readMode := getNBTReadFunction(contentsType, false)
		interf, length := readMode(buf[cursor:])
		cursor += length
		into = append(into, interf)
	}
	//fmt.Println("Finished reading compound at cursor", cursor)
	return into, cursor
}

func NBTWriteCompound(input interface{}) []byte {
	compound := input.([]interface{})
	output := make([]byte, 0)
	for _, element := range compound {
		output = append(output, NBTWrite(element)...)
	}
	return append(output, 0x00)
}

func NBTWrite(element interface{}) []byte {
	fmt.Printf("Element Type %T\n", element)
	output := make([]byte, 0)
	switch element.(type) {
	case []interface{}:
		output = append(output, NBTWriteCompound(element)...)
		break
	case NBTNamed:
		namedTag := element.(NBTNamed)
		fmt.Printf("Named tag '%s' type %d\n", namedTag.Name, namedTag.Type)
		output = append(output, getNBTWriteFunction(namedTag.Type, namedTag.Name)(namedTag.Data)...)
		break
	default:
		fmt.Printf("!!!!! Unknown format %T", element)
	}
	return output
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

	listLength, length := NBTReadSignedInteger(buf[cursor:])
	//fmt.Println("List Length", listLength)
	cursor += length

	if list.Type == 0 {
		fmt.Println("??? List has 0 type")
		return list, cursor
	}

	if listLength == 0 {
		fmt.Println("??? List has 0 length")
		return list, cursor
	}

	list.Values = make([]interface{}, listLength.(int32))

	//fmt.Println(hex.Dump(buf))

	readFunction := getNBTReadFunction(list.Type, true)

	for i := 0; i < int(listLength.(int32)); i++ {
		//fmt.Println("Reading tag for list type ", list.Type)
		tag, length := readFunction(buf[cursor:])
		cursor += length
		list.Values[i] = tag
	}

	return list, cursor
}

func NBTWriteList(input interface{}) []byte {
	list := input.(NBTList)
	fmt.Println("Input list ", list)
	output := append(NBTWriteByte(list.Type), NBTWriteSignedInteger(int32(len(list.Values)))...)
	writeFunc := getNBTWriteFunction(list.Type, "none")
	for _, item := range list.Values {
		output = append(output, writeFunc(item)...)
	}

	return output
}

var tagTypes = []string{
	"End", "Byte", "Short", "Int", "Long", "Float", "Double", "Byte_Array", "String", "List", "Compound", "Int_Array", "Long_Array",
}

func NBTToString(input interface{}, level int) string {
	output := ""
	spacing := ""
	for i := 0; i < level; i++ {
		spacing += " "
	}
	switch input.(type) {
	case []interface{}:
		output += fmt.Sprintf(" [Compound length %d] ", len(input.([]interface{})))
		for _, elem := range input.([]interface{}) {
			output += fmt.Sprintf("\n%s{ (%T) %s\n%s}\n", spacing, elem, NBTToString(elem, level+1), spacing)
		}
		break
	case NBTNamed:
		tag := input.(NBTNamed)
		output += fmt.Sprintf("TAG_%s('%s'):", tagTypes[tag.Type], tag.Name)
		output += NBTToString(tag.Data, level+1)
		break
	case NBTList:
		list := input.(NBTList)
		output += fmt.Sprintf("%d entries (Type %d)\n", len(list.Values), list.Type)
		for _, elem := range list.Values {
			output += spacing + NBTToString(elem, level+1) + "\n"
		}
		break
	case []uint8:
		array := input.([]uint8)
		if len(array) > 10 {
			output += fmt.Sprintf("%s[%d entires]\n", spacing, len(array))
			break
		}

		for _, elem := range array {
			output += fmt.Sprintf("%s%d\n", spacing, elem)
		}
		break
	default:
		output += fmt.Sprintf("%v\n", input)
		fmt.Printf("??? Unknown type %T\n", input)
		break
	}

	return output
}

func NBTAsMap(input []interface{}) interface{} {
	output := map[string]interface{}{}

	for i, elem := range input {
		switch elem.(type) {
		case NBTNamed:
			named := elem.(NBTNamed)
			name := named.Name
			if named.Name == "" {
				name = "Unnamed"
			}
			output[name] = NBTAsMap([]interface{}{named.Data})
			break
		case []interface{}:
			output[fmt.Sprintf("Compound_%d", i)] = NBTAsMap(elem.([]interface{}))
		case NBTList:
			output[fmt.Sprintf("List-%d", i)] = NBTAsMap((elem.(NBTList)).Values)
		default:
			//fmt.Printf("unknown type %T\n", elem)
			return elem
		}
	}
	return output

}
