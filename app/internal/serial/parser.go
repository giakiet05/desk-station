package serial

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	se "github.com/tarm/serial"
)

type PayloadFactory func() FirmwarePayload

var payloadRegistry = make(map[FirmwareTopic]PayloadFactory)

func RegisterPayload(topic FirmwareTopic, factory PayloadFactory) {
	payloadRegistry[topic] = factory
}

func ParseFirmwareMessage(line []byte) (FirmwareMessage, FirmwarePayload, error) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return FirmwareMessage{}, nil, errors.New("empty message")
	}

	msg := FirmwareMessage{}

	if err := json.Unmarshal(trimmed, &msg); err != nil {
		return FirmwareMessage{}, nil, fmt.Errorf("failed to parse message: %w", err)
	}

	factory, exists := payloadRegistry[msg.Topic]
	if !exists {
		return msg, nil, fmt.Errorf("unknown topic: %s", msg.Topic)
	}

	payload := factory()

	// var payload any
	// switch msg.Topic {
	// case FirmwareTopicEventIR:
	// 	payload = &IrPayload{}
	// case FirmwareTopicEventDHT11:
	// 	payload = &Dht11Payload{}
	// case FirmwareTopicEventButton:
	// 	payload = &ButtonPayload{}
	// case FirmwareTopicEventDevice:
	// 	payload = &DevicePayload{}
	// case FirmwareTopicAckCommand:
	// 	payload = &CommandAckPayload{}
	// case FirmwareTopicCommandLED:
	// 	payload = &LedCommandPayload{}
	// case FirmwareTopicCommandPing:
	// 	return msg, msg.Payload, nil
	// default:
	// 	return msg, nil, fmt.Errorf("unknown topic: %s", msg.Topic)
	// }

	if len(msg.Payload) == 0 || string(msg.Payload) == "null" {
		return msg, nil, errors.New("empty payload")
	}

	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		return msg, nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	return msg, payload, nil
}

func ListenFirmware(ctx context.Context, byIdPath string, baudRate int, onMessage func(msg FirmwareMessage, payload any), onError func(err error)) {
	go func() {
		config := &se.Config{
			Name: byIdPath,
			Baud: baudRate,
		}
		port, err := se.OpenPort(config)
		if err != nil {
			if onError != nil {
				onError(fmt.Errorf("failed to open serial port: %w", err))
			}
			return
		}
		defer port.Close()

		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				_ = port.Close()
			case <-done:
			}
		}()
		defer close(done)

		reader := bufio.NewReader(port)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				if onError != nil {
					onError(err)
				}
				return
			}

			msg, payload, err := ParseFirmwareMessage(line)
			if err != nil {
				if onError != nil {
					onError(err)
				}
				continue
			}

			if onMessage != nil {
				onMessage(msg, payload)
			}
		}
	}()
}
