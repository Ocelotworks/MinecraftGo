package packet

type PlayerInfoAddPlayer struct {
	Action  int      `proto:"varInt"`
	Players []Player `proto:"playerArray"`
}

func (piap *PlayerInfoAddPlayer) GetPacketId() int {
	return 0x34
}

func (piap *PlayerInfoAddPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}

type Player struct {
	UUID           []byte           `proto:"uuid"`
	Username       string           `proto:"string"`
	Properties     []PlayerProperty `proto:"playerPropertiesArray"`
	Gamemode       int              `proto:"varInt"`
	Ping           int              `proto:"varInt"`
	HasDisplayname bool             `proto:"bool"`
	DisplayName    string           `proto:"string"`
}

type PlayerProperty struct {
	Name   string `proto:"string"`
	Value  string `proto:"string"`
	Signed bool   `proto:"bool"`
	//Signature string `proto:"string"`
}
