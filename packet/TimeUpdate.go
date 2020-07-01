package packet

type TimeUpdate struct {
	WorldAge  int64 `proto:"long"`
	TimeOfDay int64 `proto:"long"`
}

func (stu *TimeUpdate) GetPacketId() int {
	return 0x4E
}
