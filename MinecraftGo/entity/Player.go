package entity

type Player struct {
	Username string
	UUID     []byte
	EntityID int
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	Settings PlayerSettings
}

type PlayerSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           int
	ChatColours        bool
	DisplayedSkinParts byte
	MainHand           int
}
