package packet

import (
	"bytes"
	"fmt"

	"github.com/Tnze/go-mc/nbt"
)

type ClientSettings struct {
	Locale             string `proto:"string"`
	ViewDistance       byte   `proto:"unsignedByte"`
	ChatMode           int    `proto:"varInt"`
	ChatColours        bool   `proto:"bool"`
	DisplayedSkinParts byte   `proto:"unsignedByte"`
	MainHand           int    `proto:"varInt"`
}

func (cs *ClientSettings) GetPacketId() int {
	return 0x05
}

func (cs *ClientSettings) Handle(packet []byte, connection *Connection) {
	fmt.Println("Got Client Settings:")
	fmt.Println("Locale: ", cs.Locale)
	fmt.Println("View Distance: ", cs.ViewDistance)
	fmt.Println("Chat Mode: ", cs.ChatMode)
	fmt.Println("Chat Colours: ", cs.ChatColours)
	fmt.Println("Displayed skin: ", cs.DisplayedSkinParts)
	fmt.Println("Main Hand: ", cs.MainHand)

	var buf bytes.Buffer

	encoder := nbt.NewEncoder(&buf)

	type BitTestStruct struct {
		MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	}

	heightMap := BitTestStruct{
		MotionBlocking: make([]int64, 0),
	}
	for i := 0; i < 36; i++ {
		heightMap.MotionBlocking = append(heightMap.MotionBlocking, int64(i))
	}

	exception := encoder.Encode(heightMap)

	fmt.Println(exception)

	chunkData := Packet(&ChunkData{
		X:                0,
		Z:                0,
		FullChunk:        true,
		PrimaryBitMask:   0b11111111111111111111111111111111,
		HeightMap:        buf.Bytes(),
		Biomes:           nil,
		DataSize:         0,
		Data:             make([]byte, 0),
		BlockEntityCount: 0,
		BlockEntities:    make([]byte, 0),
	})

	connection.SendPacket(&chunkData)

	//playerPos := Packet(&PlayerPositionAndLook{
	//	X:          12345,
	//	Y:          12345,
	//	Z:          12345,
	//	Yaw:        12345,
	//	Pitch:      12345,
	//	Flags:      0,
	//	TeleportID: 12345,
	//})
	//
	//connection.SendPacket(&playerPos)

	//TODO Player info

	//viewPos := Packet(&UpdateViewPosition{
	//	ChunkX: 1,
	//	ChunkZ: 2,
	//})
	//
	//connection.SendPacket(&viewPos)
}
