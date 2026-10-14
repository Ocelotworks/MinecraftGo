package packet

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
)

type JoinGame struct {
	EntityID            int      `proto:"int"`
	IsHardcore          bool     `proto:"bool"`
	DimensionNames      []string `proto:"stringArray"` // TODO: this is identifier array
	MaxPlayers          int      `proto:"varInt"`
	ViewDistance        int      `proto:"varInt"`
	SimulationDistance  int      `proto:"varInt"`
	ReducedDebugInfo    bool     `proto:"bool"`
	EnableRespawnScreen bool     `proto:"bool"`
	DoLimitedCrafting   bool     `proto:"bool"`
	DimensionType       int      `proto:"varInt"`
	DimensionName       string   `proto:"string"`
	HashedSeed          int64    `proto:"long"`
	Gamemode            int      `proto:"varInt"`
	PreviousGamemode    int      `proto:"varInt"` // optional: 0 for none, otherwise the gamemode id + 1
	IsDebug             bool     `proto:"bool"`
	IsFlat              bool     `proto:"bool"`
	HasDeathLocation    bool     `proto:"bool"` // TODO: must be false
	//DeathDimensionName  string   `proto:"string"`
	//DeathLocation int // TODO: this is a position
	PortalCooldown     int  `proto:"varInt"`
	SeaLevel           int  `proto:"varInt"`
	OnlineMode         bool `proto:"bool"`
	EnforcesSecureChat bool `proto:"bool"`
}

func (ls *JoinGame) GetPacketId() int {
	return constants.CBLogin
}

/**
func (ls *JoinGame) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
