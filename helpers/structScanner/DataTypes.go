package structScanner

import "github.com/Ocelotworks/MinecraftGo/dataTypes"

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":          dataTypes.ReadLong,
	"varInt":        dataTypes.ReadVarInt,
	"string":        dataTypes.ReadString,
	"raw":           dataTypes.ReadRaw,
	"short":         dataTypes.ReadShort,
	"unsignedShort": dataTypes.ReadUnsignedShort,
	"bool":          dataTypes.ReadBoolean,
	"unsignedByte":  dataTypes.ReadUnsignedByte,
	"int":           dataTypes.ReadInt,
	//"intArray":		 dataTypes.ReadIntArray,
	"float":           dataTypes.ReadFloat,
	"double":          dataTypes.ReadDouble,
	"uuid":            dataTypes.ReadUUID,
	"varIntByteArray": dataTypes.ReadVarIntByteArray,
}

var dataWriteMap = map[string]func(interface{}) []byte{
	"long":           dataTypes.WriteLong,
	"varInt":         dataTypes.WriteVarInt,
	"string":         dataTypes.WriteString,
	"raw":            dataTypes.WriteRaw,
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
}
