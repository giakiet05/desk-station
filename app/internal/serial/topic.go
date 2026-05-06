package serial

type FirmwareTopic string

const (
	FirmwareTopicEventIr     FirmwareTopic = "evt.ir"
	FirmwareTopicEventDht11  FirmwareTopic = "evt.dht11"
	FirmwareTopicEventButton FirmwareTopic = "evt.button"
	FirmwareTopicEventDevice FirmwareTopic = "evt.device"
	FirmwareTopicCommandLed  FirmwareTopic = "cmd.led"
	FirmwareTopicCommandPing FirmwareTopic = "cmd.ping"
	FirmwareTopicAckCommand  FirmwareTopic = "ack.command"
)
