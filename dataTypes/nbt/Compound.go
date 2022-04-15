package nbt

import (
	"encoding/binary"
)

type Compound struct {
	Data map[string]NBTValue
}

func NewCompound(compound map[string]NBTValue) Compound {
	return Compound{Data: compound}
}

func (c *Compound) GetValue() interface{} {
	return c.Data
}

func (c *Compound) SetValue(i interface{}) {
	c.Data = i.(map[string]NBTValue)
}

func (c *Compound) GetType() int {
	return 10
}

func (c *Compound) Read(buf []byte) int {
	c.Data = make(map[string]NBTValue)
	cursor := 0
	for {
		if cursor >= len(buf) {
			// fmt.Println("Buffer overrun ", cursor)
			break
		}
		itemType := buf[cursor]
		cursor++ // Advance 1 byte
		if itemType == 0 {
			// fmt.Printf("Stopping at cursor=%d because itemType=%d\n", cursor, itemType)
			break
		}
		nameLength := int(binary.BigEndian.Uint16(buf[cursor : cursor+2]))
		cursor += 2 // Advance 2 bytes for name length
		name := string(buf[cursor : cursor+nameLength])
		cursor += nameLength // Advance the length of the tag name
		item := NewValue(itemType)
		cursor += item.Read(buf[cursor:])
		// fmt.Printf("Loaded new value of type=%d name=%s cursor=%d\n", itemType, name, cursor)
		c.Data[name] = item
	}
	return cursor
}

func (c *Compound) Write() []byte {
	output := make([]byte, 0)
	for name, item := range c.Data {
		output = append(output, byte(item.GetType()))
		nameLength := make([]byte, 2)
		binary.BigEndian.PutUint16(nameLength, uint16(len(name)))
		output = append(output, nameLength...)
		output = append(output, name...)

		output = append(output, item.Write()...)
	}

	// Tag End
	output = append(output, 0x00)

	//fmt.Println(hex.Dump(output))

	return output
}
