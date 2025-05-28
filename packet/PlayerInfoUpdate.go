package packet

type PlayerInfoUpdate struct {
	ActionsUpper byte `proto:"unsignedByte"`
	ActionsLower byte `proto:"unsignedByte"` // TODO: This is an enumset?
	Players      []PlayerInfoUpdatePlayer
}

type PlayerInfoUpdatePlayer struct {
	UUID    []byte                       `proto:"uuid"`
	Actions PlayerInfoUpdatePlayerAction // TODO: this definitely isnt right
}

// TODO ???
type PlayerInfoUpdatePlayerAction any
