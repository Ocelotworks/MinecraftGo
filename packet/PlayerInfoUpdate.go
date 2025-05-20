package packet

type PlayerInfoUpdate struct {
	ActionsUpper byte                     `proto:"unsignedByte"`
	ActionsLower byte                     `proto:"unsignedByte"` // TODO: This is an enumset?
	Players      []PlayerInfoUpdatePlayer `proto:"array"`
}

type PlayerInfoUpdatePlayer struct {
	UUID    []byte                       `proto:"uuid"`
	Actions PlayerInfoUpdatePlayerAction `proto:"array"` // TODO: this definitely isnt right
}

// TODO ???
type PlayerInfoUpdatePlayerAction any
