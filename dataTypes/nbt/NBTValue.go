package nbt

import (
	"fmt"
	"strings"
)

type NBTValue interface {
	GetValue() interface{}
	SetValue(interface{})
	GetType() int
	Read(buf []byte) int
	Write() []byte
}

func IDFromName(name string) byte {
	switch name {
	case "end":
		return 0
	case "byte":
		return 1
	case "short":
		return 2
	case "int":
		return 3
	case "long":
		return 4
	case "float":
		return 5
	case "double":
		return 6
	case "byteArray":
		return 7
	case "string":
		return 8
	case "list":
		return 9
	case "compound":
		return 10
	case "intArray":
		return 11
	case "longArray":
		return 12
	}
	fmt.Println("Invalid name", name)
	return 0
}

func IDFromType(typeName string) byte {
	switch typeName {
	case "uint8":
		fallthrough
	case "byte":
		return 1
	case "int16":
		return 2
	case "int32":
		return 3
	case "int64":
		return 4
	case "float32":
		return 5
	case "float64":
		return 6
	case "[]byte":
		return 7
	case "string":
		return 8
	// How to lists :think:
	//case "list":
	//    return 9
	case "struct":
		return 10
	case "[]int32":
		return 11
	case "[]int64":
		return 12
	}

	if strings.HasPrefix(typeName, "[]") {
		return 9
	}

	fmt.Println("Unknown type", typeName)
	return 10
}

func NewValue(id byte) NBTValue {
	switch id {
	case 1:
		return &Byte{}
	case 2:
		return &Short{}
	case 3:
		return &Int{}
	case 4:
		return &Long{}
	case 5:
		return &Float{}
	case 6:
		return &Double{}
	case 7:
		return &ByteArray{}
	case 8:
		return &String{}
	case 9:
		return &List{}
	case 10:
		return &Compound{}
	case 11:
		return &IntArray{}
	case 12:
		return &LongArray{}
	}

	fmt.Println("Unknown type", id)

	return nil
}

func ReadNBT(data []byte) (Compound, int) {
	outerCompound := Compound{}
	return outerCompound, outerCompound.Read(data)
}
