package packet

type KnownPack struct {
	Namespace string `proto:"string"`
	ID        string `proto:"string"`
	Version   string `proto:"string"`
}
