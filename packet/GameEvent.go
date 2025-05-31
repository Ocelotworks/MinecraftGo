package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EventType byte

const (
	EventTypeNoRespawnBlockAvailable = 0
	EventTypeBeginRaining            = 1
	EventTypeEndRaining              = 2
	EventTypeChangeGameMode          = 3
	EventTypeWinGame                 = 4
	EventTypeDemoEvent               = 5
	EventTypeArrowHitPlayer          = 6
	EventTypeRainLevelChange         = 7
	EventTypePufferfishSting         = 9
	EventTypeElderGuardian           = 10
	EventTypeEnableRespawnScreen     = 11
	EventTypeLimitedCrafting         = 12
	EventTypeWaitingForChunks        = 13
)

type GameEvent struct {
	Event byte    `proto:"unsignedByte"`
	Value float32 `proto:"float"`
}

func (ka *GameEvent) GetPacketId() int {
	return constants.CBGameEvent //Client
}

/**
func (ka *KeepAlive) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("KeepAlive", ka)
}
*/
