package entity

type BlockData struct {
	States     []BlockState        `json:"states"`
	Properties map[string][]string `json:"properties"`
}

type BlockState struct {
	ID         int               `json:"id"`
	Default    bool              `json:"default"`
	Properties map[string]string `json:"properties"`
}
