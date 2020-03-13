package structScanner

import (
	"reflect"
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
)

type PacketStructScanner struct {}

func (pss *PacketStructScanner) StructScan(packet *packetType.Packet, buf []byte) {
	v := reflect.ValueOf(*packet).Elem()
	t := reflect.TypeOf(*packet).Elem()

	cursor := 0

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}
		if dataReadMap[tag] == nil {
			fmt.Println("!!! Unknown tag type ", tag)
			continue
		}

		if len(buf) < cursor {
			fmt.Println("Cursor overrun")
			continue
		}
		val, end := dataReadMap[tag](buf[cursor:])

		cursor += end
		//fmt.Printf("Reading tag %s into field %s value %v cursor %d\n", tag, field.Name, val, cursor)

		v.FieldByName(field.Name).Set(reflect.ValueOf(val))
	}
}

func (pss *PacketStructScanner) UnmarshalData(input interface{}) []byte {
	v := reflect.ValueOf(input).Elem()
	t := reflect.TypeOf(input).Elem()

	payload := make([]byte, 0)

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}

		val := v.FieldByName(field.Name).Interface()
		if val == nil {
			fmt.Println("nil value!!!", field.Name)
			continue
		}

		var segment []byte

		if tag == "playerArray" {
			playerData := val.([]entity.Player)
			segment = dataTypes.WriteVarInt(len(playerData))
			//fmt.Println("Player Data Array Length ", len(playerData))
			for _, player := range playerData {
				segment = append(segment, pss.UnmarshalData(&player)...)
			}
			//fmt.Println(hex.Dump(segment))
		} else if tag == "playerPropertiesArray" {
			playerProperties := val.([]entity.PlayerProperty)
			//fmt.Println("Player Property Array Length ", len(playerProperties))
			segment = dataTypes.WriteVarInt(len(playerProperties))
			if len(playerProperties) > 0 {
				for _, playerProperty := range playerProperties {
					if playerProperty.Signature != "" {
						playerProperty.Signed = true
					}
					segment = append(segment, pss.UnmarshalData(&playerProperty)...)
				}
			}
		} else {
			if dataWriteMap[tag] == nil {
				fmt.Println("!!!! Unknown tag type ", tag)
				continue
			}

			segment = dataWriteMap[tag](v.FieldByName(field.Name).Interface())
		}
		//if len(segment) < 100 {
		//	fmt.Printf("Field %s type %s coded as value %v (Between  %d - %d)\n", field.Name, tag, val, len(payload), len(payload)+len(segment))
		//}else{
		//	fmt.Printf("Field %s type %s coded as value [Big value] (Between  %d - %d)\n", field.Name, tag, len(payload), len(payload)+len(segment))
		//}

		payload = append(payload, segment...)
	}

	return payload
}


