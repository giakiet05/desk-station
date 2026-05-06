package service

import (
	"app/internal/bus"
	"app/internal/serial"
)

func RegisterRouters() {
	RegisterRouter(serial.FirmwareTopicEventIr, bus.BusTopicIr)
	RegisterRouter(serial.FirmwareTopicEventDht11, bus.BusTopicDht11)
	RegisterRouter(serial.FirmwareTopicEventButton, bus.BusTopicButton)
	RegisterRouter(serial.FirmwareTopicEventDevice, bus.BusTopicDevice)

}
