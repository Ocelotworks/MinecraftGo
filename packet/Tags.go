package packet

// TODO this has changed
type Tags struct {
	TagCount int    `proto:"varInt"`
	Tags     []byte `proto:"raw"`
}

func (t *Tags) GetPacketId() int {
	return 98
}
