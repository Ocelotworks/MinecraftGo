package dataTypes

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
)

const (
	nbtTagEnd    = 0
	nbtByte      = 1
	nbtShort     = 2
	nbtInt       = 3
	nbtLong      = 4
	nbtFloat     = 5
	nbtDouble    = 6
	nbtByteArray = 7
	nbtString    = 8
	nbtList      = 9
	nbtCompound  = 10
	nbtIntArray  = 11
	nbtLongArray = 12
)

func ReadNBT(buf []byte) ([]*NBTNamed, int) {
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

	return compound.([]*NBTNamed), 0
}

func getNBTReadFunction(index byte) func(buf []byte) (interface{}, int) {
	switch index {
	case nbtByte:
		return NBTReadByte
	case nbtShort:
		return NBTReadSignedShort
	case nbtInt:
		return NBTReadSignedInteger
	case nbtLong:
		return NBTReadSignedLong
	case nbtFloat:
		return NBTReadFloat
	case nbtDouble:
		return NBTReadDouble
	case nbtByteArray:
		return NBTReadByteArray
	case nbtString:
		return NBTReadString
	case nbtList:
		return NBTReadList
	case nbtCompound:
		return NBTReadCompound
	case nbtIntArray:
		return NBTReadIntArray
	case nbtLongArray:
		return NBTReadLongArray
	}
	fmt.Println("!!!!! Unknown Type ", index)
	return nil
}

func getNBTWriteFunction(index byte, name string) func(input interface{}) []byte {
	if strings.HasPrefix(name, "List_") {
		switch index {
		case nbtByte:
			return NBTWriteByte
		case nbtShort:
			return NBTWriteSignedShort
		case nbtInt:
			return NBTWriteSignedInteger
		case nbtLong:
			return NBTWriteSignedLong
		case nbtFloat:
			return NBTWriteFloat
		case nbtDouble:
			return NBTWriteDouble
		case nbtByteArray:
			return NBTWriteByteArray
		case nbtString:
			return NBTWriteString
		case nbtList:
			return NBTWriteList
		case nbtCompound:
			return NBTWriteCompound
		case nbtIntArray:
			return NBTWriteIntArray
		case nbtLongArray:
			return NBTWriteLongArray
		}
		fmt.Println("!!!!! Unknown Type ", index)
		return nil
	}
	fmt.Println("write function ", name)
	return NBTWriteNamed(getNBTWriteFunction(index, "List_"), index, name)
}

func getNBTType(value interface{}) byte {
	switch value.(type) {
	case byte:
	case bool:
	case int8:
		return nbtByte
	case int16:
		return nbtShort
	case int32:
		return nbtInt
	case int:
	case int64:
		return nbtLong
	case float32:
		return nbtFloat
	case float64:
		return nbtDouble
	case []byte:
	case []bool:
	case []int8:
		return nbtByteArray
	case string:
		return nbtString
	case []int32:
		return nbtIntArray
	case []int64:
	case []int:
		return nbtLongArray
	}
	reflectedValue := reflect.TypeOf(value)

	if isCustom(reflectedValue) {
		return nbtCompound
	}

	if reflectedValue.Kind().String() == "slice" {
		return nbtList
	}

	fmt.Printf("!!! Unknown compound type %T\n", value)
	return 0x00
}

func getNBTTypeFromString(name string) byte {
	switch name {
	case "end":
		return nbtTagEnd
	case "byte":
		return nbtByte
	case "short":
		return nbtShort
	case "int":
		return nbtInt
	case "long":
		return nbtLong
	case "float":
		return nbtFloat
	case "double":
		return nbtDouble
	case "byteArray":
		return nbtByteArray
	case "string":
		return nbtString
	case "list":
		return nbtList
	case "compound":
		return nbtCompound
	case "intArray":
		return nbtIntArray
	case "longArray":
		return nbtLongArray
	}
	fmt.Println("!!! Unknown compound type ", name)
	return 0x00
}

type NBTNamed struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
	Type byte        `json:"type"`
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

func NBTRead(buf []byte, cursor int) ([]*NBTNamed, int) {
	into := make([]*NBTNamed, 0)
	//fmt.Println("Reading Compound at cursor ", cursor)
	for {
		if cursor >= len(buf) {
			//fmt.Println("NBT Read Cursor overrun ", cursor, len(buf))
			break
		}
		contentsType := buf[cursor]
		cursor++
		if contentsType == nbtTagEnd {
			//fmt.Println("NBT Tag End")
			break
		}
		//fmt.Println("NBTRead Contents Type ", contentsType)
		named := NBTNamed{}
		name, length := NBTReadString(buf[cursor:])
		cursor += length
		named.Name = name.(string)
		named.Type = contentsType
		readMode := getNBTReadFunction(contentsType)
		interf, length := readMode(buf[cursor:])
		cursor += length
		named.Data = interf
		into = append(into, &named)
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
	return append(output, nbtTagEnd)
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

	named := NBTNamed{}
	named.Name = fmt.Sprintf("List_%d", cursor)
	named.Type = list.Type
	readFunction := getNBTReadFunction(list.Type)

	for i := 0; i < int(listLength.(int32)); i++ {
		//fmt.Println("Reading tag for list type ", list.Type)
		tag, length := readFunction(buf[cursor:])
		cursor += length
		list.Values[i] = tag
	}

	named.Data = list

	return named, cursor
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

func isCustom(t reflect.Type) bool {
	if t.PkgPath() != "" {
		return true
	}

	if k := t.Kind(); k == reflect.Array || k == reflect.Chan || k == reflect.Map ||
		k == reflect.Ptr || k == reflect.Slice {
		return isCustom(t.Elem()) || k == reflect.Map && isCustom(t.Key())
	} else if k == reflect.Struct {
		for i := t.NumField() - 1; i >= 0; i-- {
			if isCustom(t.Field(i).Type) {
				return true
			}
		}
	}

	return false
}

func NBTWriteStruct(input interface{}) []*NBTNamed {
	output := make([]*NBTNamed, 0)
	v := reflect.ValueOf(input).Elem()
	t := reflect.TypeOf(input).Elem()
	fmt.Println("Struct scanning Fields: ", t.NumField())
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("nbt")
		if !exists {
			fmt.Println("no NBT tag present on", field.Name)
			continue
		}
		tags := strings.SplitN(tag, ",", 2)

		tagName := tags[0]
		fmt.Println("Tag Name:", tags[0])

		if strings.HasSuffix(tagName, "*") {
			fmt.Println("can't convert wildcards to nbt")
			continue
		}

		namedStruct := NBTNamed{
			Name: tagName,
		}
		fieldValue := v.FieldByName(field.Name)
		if isCustom(field.Type) {
			namedStruct.Type = nbtCompound
			namedStruct.Data = NBTWriteStruct(fieldValue.Addr())
		}

		if len(tags) > 1 {
			fmt.Println("Tag has type override: ", tags[1])
			namedStruct.Type = getNBTTypeFromString(tags[1])
		} else {
			namedStruct.Type = getNBTType(fieldValue.Interface())
		}

		namedStruct.Data = fieldValue.Interface()

		output = append(output, &namedStruct)
	}
	return output
}

func NBTStructScan(nbt []*NBTNamed, output interface{}) {
	v := reflect.ValueOf(output).Elem()
	t := reflect.TypeOf(output).Elem()
	fmt.Println("Struct scanning Fields: ", t.NumField())
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("nbt")
		if !exists {
			fmt.Println("no NBT tag present on", field.Name)
			continue
		}
		tags := strings.SplitN(tag, ",", 2)
		tagName := tags[0]
		fmt.Println("Tag Name:", tags[0])

		isWildcard := false

		if strings.HasSuffix(tagName, "*") {
			isWildcard = true
			tagName = strings.TrimRight(tagName, "*")
		}
		var targetTag *NBTNamed

		if len(nbt) == 1 && nbt[0].Name == "" {
			nbt = nbt[0].Data.([]*NBTNamed)
		}

		for _, tag := range nbt {
			//fmt.Println(tag.Name)
			if (isWildcard && strings.HasPrefix(tag.Name, tagName)) || tag.Name == tagName {
				targetTag = tag
				break
			}
		}

		if targetTag == nil {
			fmt.Println("Could not find tag ", tagName)
			continue
		}
		fieldValue := v.FieldByName(field.Name)
		//Not a builtin
		if isCustom(field.Type) {
			if targetTag.Type == nbtCompound {
				NBTStructScan(targetTag.Data.([]*NBTNamed), fieldValue.Addr().Interface())
			} else {
				fmt.Printf("Target is %s not builtin... but target tag is %d not compound\n", field.Type.Name(), targetTag.Type)
			}
			continue
		}

		if field.Type.Kind().String() == "slice" {
			listTag, ok := targetTag.Data.(NBTNamed)
			values := targetTag.Data
			if ok {
				values = listTag.Data.(NBTList).Values
				if targetTag.Type != nbtList {
					fmt.Println("Target is a slice... but target tag is ", targetTag.Type, " not list")
					continue
				}
			}
			valuesValue := reflect.ValueOf(values)
			slice := reflect.MakeSlice(fieldValue.Type(), valuesValue.Len(), valuesValue.Len())
			for i := 0; i < valuesValue.Len(); i++ {
				elem := valuesValue.Index(i)
				slice.Index(i).Set(elem)
			}
			fieldValue.Set(slice)
			continue
		}

		fieldValue.Set(reflect.ValueOf(targetTag.Data))
	}
}
