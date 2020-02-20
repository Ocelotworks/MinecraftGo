package packet

import (
	"fmt"
)

type LoginStart struct {
	Username string `proto:"string"`
}

func (ls *LoginStart) GetPacketId() int {
	return 0x00
}

func (ls *LoginStart) Handle(packet []byte, connection *Connection) {
	//Just send the pong right back
	fmt.Println("Username ", ls.Username)

	publicKey := connection.Key.N.Bytes()

	returnPacket := Packet(&EncryptionRequest{
		ServerID:          "",
		PublicKeyLength:   len(publicKey),
		PublicKey:         publicKey,
		VerifyTokenLength: 4,
		VerifyToken:       []byte{0x01, 0x02, 0x03, 0x04},
	})
	connection.SendPacket(&returnPacket)
}
