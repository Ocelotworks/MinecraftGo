package entity

type ServerListPingResponse struct {
	Version     ServerListPingVersion `json:"version"`
	Players     ServerListPingPlayers `json:"players"`
	Favicon     string                `json:"favicon"`
	Description ChatMessageComponent  `json:"description"`
}

type ServerListPingVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerListPingPlayers struct {
	Max    int                            `json:"max"`
	Online int                            `json:"online"`
	Sample []ServerListPingPlayerListItem `json:"sample"`
}

type ServerListPingPlayerListItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
