package structScanner

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"reflect"
)

type PacketStructScanner struct{}

// prefixedArrayTag marks a slice-of-struct field that is encoded as a varInt
// element count followed by that many consecutive elements. Each element is
// itself encoded field by field using the element struct's own proto tags.
const prefixedArrayTag = "prefixedArray"

func (pss *PacketStructScanner) StructScan(packet *packetType.Packet, buf []byte) {
	pss.scanStruct(reflect.ValueOf(*packet).Elem(), buf)
}

// scanStruct fills the struct v from buf and returns the number of bytes consumed.
func (pss *PacketStructScanner) scanStruct(v reflect.Value, buf []byte) int {
	t := v.Type()

	cursor := 0

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}

		if tag == prefixedArrayTag {
			if len(buf) < cursor {
				fmt.Println("Cursor overrun")
				continue
			}
			cursor += pss.scanPrefixedArray(v.Field(fieldIndex), buf[cursor:])
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

	return cursor
}

// scanPrefixedArray reads a varInt length followed by that many struct elements
// into the slice field, returning the number of bytes consumed.
func (pss *PacketStructScanner) scanPrefixedArray(slice reflect.Value, buf []byte) int {
	if slice.Kind() != reflect.Slice || slice.Type().Elem().Kind() != reflect.Struct {
		fmt.Println("!!! prefixedArray requires a slice of structs, got ", slice.Type())
		return 0
	}

	length, cursor := dataTypes.ReadVarInt(buf)
	count := length.(int)
	// Every element is at least one byte on the wire, so a count larger than the
	// remaining buffer is malformed (and must not drive a huge allocation).
	if count < 0 || count > len(buf)-cursor {
		fmt.Println("!!! Invalid prefixedArray length ", count)
		slice.Set(reflect.MakeSlice(slice.Type(), 0, 0))
		return cursor
	}

	slice.Set(reflect.MakeSlice(slice.Type(), count, count))
	for i := 0; i < count; i++ {
		cursor += pss.scanStruct(slice.Index(i), buf[cursor:])
	}

	return cursor
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

		if tag == prefixedArrayTag {
			segment = pss.marshalPrefixedArray(v.FieldByName(field.Name))
		} else if tag == "playerArray" { // TODO: why is this here?
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

			if tag == "nbt" {
				segment = dataWriteMap[tag](v.FieldByName(field.Name).Addr().Interface())
			} else {
				segment = dataWriteMap[tag](v.FieldByName(field.Name).Interface())
			}

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

// marshalPrefixedArray writes a varInt element count followed by each struct element.
func (pss *PacketStructScanner) marshalPrefixedArray(slice reflect.Value) []byte {
	if slice.Kind() != reflect.Slice || slice.Type().Elem().Kind() != reflect.Struct {
		fmt.Println("!!! prefixedArray requires a slice of structs, got ", slice.Type())
		return nil
	}

	segment := dataTypes.WriteVarInt(slice.Len())
	for i := 0; i < slice.Len(); i++ {
		segment = append(segment, pss.UnmarshalData(slice.Index(i).Addr().Interface())...)
	}

	return segment
}
