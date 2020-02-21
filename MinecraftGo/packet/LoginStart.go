package packet

import (
	"crypto/x509"
	"fmt"

	"../dataTypes"
)

type LoginStart struct {
	Username string `proto:"string"`
}

func (ls *LoginStart) GetPacketId() int {
	return 0x00
}

func (ls *LoginStart) Handle(packet []byte, connection *Connection) {
	fmt.Println("Username ", ls.Username)

	fmt.Println(connection.Key.PublicKey)

	publicKey := x509.MarshalPKCS1PublicKey(&connection.Key.PublicKey)

	//publicKey, _ := asn1.Marshal(connection.Key.PublicKey)

	encryptionPacket := Packet(&EncryptionRequest{
		ServerID:          "",
		PublicKeyLength:   len(publicKey),
		PublicKey:         publicKey,
		VerifyTokenLength: 4,
		VerifyToken:       []byte{0x01, 0x02, 0x03, 0x04},
	})

	connection.SendPacket(&encryptionPacket)

	return

	returnPacket := Packet(&LoginSuccess{
		UUID:     "5d8af060-129e-419c-b3ac-c0dad1c91181",
		Username: "UnacceptableUse",
	})
	connection.SendPacket(&returnPacket)

	connection.State = PLAY

	joinGame := Packet(&JoinGame{
		EntityID:            120,
		Gamemode:            0,
		Dimension:           0,
		HashedSeed:          71495747907944700,
		MaxPlayers:          128,
		LevelType:           "default",
		ViewDistance:        32,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
	})

	connection.SendPacket(&joinGame)

	pluginMessage := Packet(&PluginMessage{
		IsServer:   false,
		Identifier: "minecraft:brand",
		ByteArray:  dataTypes.WriteString("BigPMC"),
	})

	connection.SendPacket(&pluginMessage)

	difficulty := Packet(&ServerDifficulty{
		Difficulty:       0,
		DifficultyLocked: false,
	})

	connection.SendPacket(&difficulty)

	//TODO: basically everything after held item change before player position

	//TODO light shit

	//connection.State = PLAY

}
