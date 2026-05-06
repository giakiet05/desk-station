package serial

type FirmwarePayload interface {
	isFirmwarePayload()
}

type IrPayload struct {
	RawCode  string `json:"raw_code"`
	Address  uint16 `json:"address"`
	Command  uint16 `json:"command"`
	IsRepeat bool   `json:"is_repeat"`
}

func (p *IrPayload) isFirmwarePayload() {}

type Dht11Payload struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Valid       bool    `json:"valid"`
}

func (p *Dht11Payload) isFirmwarePayload() {}

type ButtonPayload struct {
	Pressed bool `json:"pressed"`
}

func (p *ButtonPayload) isFirmwarePayload() {}

type DevicePayload struct {
	Status string `json:"status"`
}

func (p *DevicePayload) isFirmwarePayload() {}

type CommandAckPayload struct {
	CmdId  uint64 `json:"cmd_id"`
	OK     bool   `json:"ok"`
	Reason string `json:"reason"`
}

func (p *CommandAckPayload) isFirmwarePayload() {}

type LedCommandPayload struct {
	CmdId    uint64 `json:"cmd_id"`
	RedOn    bool   `json:"red_on"`
	YellowOn bool   `json:"yellow_on"`
	GreenOn  bool   `json:"green_on"`
}

func (p *LedCommandPayload) isFirmwarePayload() {}
