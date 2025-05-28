package structScanner

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
)

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":                      dataTypes.ReadLong,
	"varInt":                    dataTypes.ReadVarInt,
	"string":                    dataTypes.ReadString,
	"raw":                       dataTypes.ReadRaw,
	"byte":                      dataTypes.ReadByte,
	"short":                     dataTypes.ReadShort,
	"unsignedShort":             dataTypes.ReadUnsignedShort,
	"bool":                      dataTypes.ReadBoolean,
	"unsignedByte":              dataTypes.ReadUnsignedByte,
	"int":                       dataTypes.ReadInt,
	"intArray":                  dataTypes.ReadIntArray,
	"float":                     dataTypes.ReadFloat,
	"double":                    dataTypes.ReadDouble,
	"uuid":                      dataTypes.ReadUUID,
	"prefixedOptionalUuid":      dataTypes.ReadPrefixedOptionalUUID,
	"varIntByteArray":           dataTypes.ReadVarIntByteArray,
	"prefixedOptionalByteArray": dataTypes.ReadPrefixedOptionalByteArray,
	"bitset":                    dataTypes.ReadBitSet,
	"nbt":                       ReadNBT,
	"prefixedOptionalNbt":       ReadPrefixedOptionalNBT,
	"position":                  dataTypes.ReadPosition,
}

var dataWriteMap = map[string]func(any) []byte{
	"long":                      dataTypes.WriteLong,
	"varInt":                    dataTypes.WriteVarInt,
	"string":                    dataTypes.WriteString,
	"raw":                       dataTypes.WriteRaw,
	"byte":                      dataTypes.WriteByte,
	"short":                     dataTypes.WriteShort,
	"unsignedShort":             dataTypes.WriteUnsignedShort,
	"bool":                      dataTypes.WriteBoolean,
	"unsignedByte":              dataTypes.WriteUnsignedByte,
	"int":                       dataTypes.WriteInt,
	"intArray":                  dataTypes.WriteIntArray,
	"float":                     dataTypes.WriteFloat,
	"double":                    dataTypes.WriteDouble,
	"uuid":                      dataTypes.WriteUUID,
	"prefixedOptionalUuid":      dataTypes.WritePrefixedOptionalUUID,
	"entityMetadata":            dataTypes.WriteEntityMetadata,
	"varIntArray":               dataTypes.WriteVarIntArray,
	"stringArray":               dataTypes.WriteStringArray,
	"varIntByteArray":           dataTypes.WriteVarIntByteArray,
	"prefixedOptionalByteArray": dataTypes.WritePrefixedOptionalByteArray,
	"bitset":                    dataTypes.WriteBitSet,
	"nbt":                       WriteNBT,
	//"prefixedNbtArray":                  WritePrefixedNBTarray,
	"prefixedOptionalNbt": WritePrefixedOptionalNBT,
	"position":            dataTypes.WritePosition,
}

// Import cycle moment

func ReadNBT(buf []byte) (interface{}, int) {
	return nbt.ReadNBT(buf)
}

func WriteNBT(compoundStruct interface{}) []byte {
	output := nbt.NBTMarshal(compoundStruct)
	return output[:len(output)-1]
}

//func WritePrefixedNBTarray (compoundStruct interface{}) []byte {
//
//}

func ReadPrefixedOptionalNBT(buf []byte) (interface{}, int) {
	present, n := dataTypes.ReadBoolean(buf)
	if !(present.(bool)) {
		return nil, n
	}

	return ReadNBT(buf[n:])
}

func WritePrefixedOptionalNBT(compoundStruct interface{}) []byte {
	b := dataTypes.WriteBoolean(true) // TODO: support false
	b = append(b, nbt.NBTMarshal(compoundStruct)...)

	return b[:len(b)-1] // TODO: why is this off by one?
}
