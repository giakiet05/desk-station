package serial

import "encoding/json"

type FirmwareMessage struct {
	V       string          `json:"v"`
	Topic   FirmwareTopic   `json:"topic"`
	Type    FirmwareType    `json:"type"`
	Seq     uint64          `json:"seq"`
	Ts      uint64          `json:"ts"`
	Payload json.RawMessage `json:"payload"`
}
