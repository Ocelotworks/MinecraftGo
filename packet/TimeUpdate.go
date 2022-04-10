package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type TimeUpdate struct {
	WorldAge  int64 `proto:"long"`
	TimeOfDay int64 `proto:"long"`
}

func (stu *TimeUpdate) GetPacketId() int {
	return constants.CBTimeUpdate
}
