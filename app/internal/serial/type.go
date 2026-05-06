package serial

type FirmwareType string

const (
	FirmwareTypeEvent   FirmwareType = "event"
	FirmwareTypeCommand FirmwareType = "command"
	FirmwareTypeAck     FirmwareType = "ack"
)
