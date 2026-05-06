package bus

type BusTopic string

const (
	BusTopicSerialError BusTopic = "serial.error"
	BusTopicIr          BusTopic = "fw.ir"
	BusTopicDht11       BusTopic = "fw.dht11"
	BusTopicButton      BusTopic = "fw.button"
	BusTopicDevice      BusTopic = "fw.device"
)
