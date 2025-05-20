package structScanner

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
)

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":            dataTypes.ReadLong,
	"varInt":          dataTypes.ReadVarInt,
	"string":          dataTypes.ReadString,
	"raw":             dataTypes.ReadRaw,
	"byte":            dataTypes.ReadByte,
	"short":           dataTypes.ReadShort,
	"unsignedShort":   dataTypes.ReadUnsignedShort,
	"bool":            dataTypes.ReadBoolean,
	"unsignedByte":    dataTypes.ReadUnsignedByte,
	"int":             dataTypes.ReadInt,
	"intArray":        dataTypes.ReadIntArray,
	"float":           dataTypes.ReadFloat,
	"double":          dataTypes.ReadDouble,
	"uuid":            dataTypes.ReadUUID,
	"varIntByteArray": dataTypes.ReadVarIntByteArray,
	"bitset":          dataTypes.ReadBitSet,
	"nbt":             ReadNBT,
	"position":        dataTypes.ReadPosition,
}

var dataWriteMap = map[string]func(any) []byte{
	"long":           dataTypes.WriteLong,
	"varInt":         dataTypes.WriteVarInt,
	"string":         dataTypes.WriteString,
	"raw":            dataTypes.WriteRaw,
	"byte":           dataTypes.WriteByte,
	"short":          dataTypes.WriteShort,
	"unsignedShort":  dataTypes.WriteUnsignedShort,
	"bool":           dataTypes.WriteBoolean,
	"unsignedByte":   dataTypes.WriteUnsignedByte,
	"int":            dataTypes.WriteInt,
	"intArray":       dataTypes.WriteIntArray,
	"float":          dataTypes.WriteFloat,
	"double":         dataTypes.WriteDouble,
	"uuid":           dataTypes.WriteUUID,
	"entityMetadata": dataTypes.WriteEntityMetadata,
	"varIntArray":    dataTypes.WriteVarIntArray,
	"stringArray":    dataTypes.WriteStringArray,
	"bitset":         dataTypes.WriteBitSet,
	"nbt":            WriteNBT,
	"position":       dataTypes.WritePosition,
}

// Import cycle moment

func ReadNBT(buf []byte) (interface{}, int) {
	return nbt.ReadNBT(buf)
}

func WriteNBT(compoundStruct interface{}) []byte {
	output := nbt.NBTMarshal(compoundStruct)
	return output[:len(output)-1]
}
