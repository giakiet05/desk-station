package serial

func RegisterPayloads() {
	RegisterPayload(FirmwareTopicEventIr, func() FirmwarePayload { return &IrPayload{} })
	RegisterPayload(FirmwareTopicEventDht11, func() FirmwarePayload { return &Dht11Payload{} })
	RegisterPayload(FirmwareTopicEventButton, func() FirmwarePayload { return &ButtonPayload{} })
	RegisterPayload(FirmwareTopicEventDevice, func() FirmwarePayload { return &DevicePayload{} })
	RegisterPayload(FirmwareTopicAckCommand, func() FirmwarePayload { return &CommandAckPayload{} })
	RegisterPayload(FirmwareTopicCommandLed, func() FirmwarePayload { return &LedCommandPayload{} })
}
