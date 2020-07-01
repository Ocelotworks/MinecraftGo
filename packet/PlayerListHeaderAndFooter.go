package packet

type PlayerListHeaderAndFooter struct {
	Header string `proto:"string"`
	Footer string `proto:"string"`
}

func (plhaf *PlayerListHeaderAndFooter) GetPacketId() int {
	return 0x53
}
