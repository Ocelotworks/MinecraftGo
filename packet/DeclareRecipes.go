package packet

// TODO: this has changed
type DeclareRecipes struct {
	NumRecipes int `proto:"varInt"`
	Recipes    []Recipe
}

func (dr *DeclareRecipes) GetPacketId() int {
	return 99
}

/**
func (dr *DeclareRecipes) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/

type Recipe struct {
	ID   string `proto:"string"`
	Type string `proto:"string"`
	Data []byte `proto:"raw"`
}
