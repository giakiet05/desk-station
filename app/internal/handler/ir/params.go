package ir

type VolumeParams struct {
	Delta int `json:"delta"` // -1 for volume down, 1 for volume up
}

type MouseMoveParams struct {
	DX int `json:"dx"`
	DY int `json:"dy"`
}

type MouseClickParams struct {
	Button int `json:"button"` // -1: right, 0: middle, 1: left
	Count  int `json:"count"`
}

type MouseScrollParams struct {
	Direction bool `json:"direction"` // true: horizontal, false: vertical
	Amount    int  `json:"amount"`    //positive for scroll up/right, negative for scroll down/left
}

type RunCommandParams struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type ArrowKeyParams struct {
	Direction string `json:"direction"` // "up", "down", "left", "right"
}
