package entity

type BlockData struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	DisplayName  string          `json:"displayName"`
	Hardness     float64         `json:"hardness"`
	Resistance   float64         `json:"resistance"`
	HarvestTools map[string]bool `json:"harvestTools"`
	StackSize    int             `json:"stackSize"`
	Diggable     bool            `json:"diggable"`
	Material     string          `json:"material"`
	Transparent  bool            `json:"transparent"`
	EmitLight    int             `json:"emitLight"`
	FilterLight  int             `json:"filterLight"`
	DefaultState int             `json:"defaultState"`
	MinStateId   int             `json:"minStateId"`
	MaxStateId   int             `json:"maxStateId"`
	Drops        []int           `json:"drops"`
	BoundingBox  string          `json:"boundingBox"`
	States       []BlockState    `json:"states"`
}

type BlockState struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	NumValues int      `json:"num_values"`
	Values    []string `json:"values"`
}

type Block struct {
	BlockName  string
	Properties map[string]string
	X          int
	Y          int
	Z          int
	// TODO: Entities
}
