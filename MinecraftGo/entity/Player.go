package entity

type Player struct {
	UUID           []byte               `proto:"uuid"`
	Username       string               `proto:"string"`
	Properties     []PlayerProperty     `proto:"playerPropertiesArray"`
	Gamemode       int                  `proto:"varInt"`
	Ping           int                  `proto:"varInt"`
	HasDisplayName bool                 `proto:"bool"`
	DisplayName    ChatMessageComponent `proto:"string"`
	EntityID       int
	X              float64
	Y              float64
	Z              float64
	Yaw            float32
	Pitch          float32
	Settings       PlayerSettings
}

type PlayerSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           int
	ChatColours        bool
	DisplayedSkinParts byte
	MainHand           int
}

type PlayerProperty struct {
	Name      string `proto:"string" json:"name"`
	Value     string `proto:"string" json:"value"`
	Signed    bool   `proto:"bool"`
	Signature string `proto:"string" json:"signature"`
}
