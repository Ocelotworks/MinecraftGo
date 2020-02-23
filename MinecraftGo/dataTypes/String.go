package dataTypes

import (
	"../entity"
	"bytes"
	"encoding/json"
	"fmt"
)

func ReadString(buf []byte) (interface{}, int) {
	stringLength, stringStart := ReadVarInt(buf)

	output := string(bytes.Runes(buf[stringStart:]))[:stringLength.(int)]

	return output, stringStart + len([]byte(output))

}

func WriteString(input interface{}) []byte {
	var b []byte
	switch input.(type) {
	case entity.ChatMessageComponent:
		b, _ = json.Marshal(input)
		break
	case string:
		b = []byte(input.(string))
	}

	output := WriteVarInt(len(b))
	fmt.Println("String Length ", len(b))
	fmt.Println(output)
	output = append(output, b...)
	return output
}
