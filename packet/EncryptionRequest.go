package packet

type EncryptionRequest struct {
	ServerID          string `proto:"string"`
	PublicKeyLength   int    `proto:"varInt"`
	PublicKey         []byte `proto:"raw"`
	VerifyTokenLength int    `proto:"varInt"`
	VerifyToken       []byte `proto:"raw"`
}

func (er *EncryptionRequest) GetPacketId() int {
	return 0x01
}

func (er *EncryptionRequest) Handle(packet []byte, connection *Connection) {
	//No op
}
