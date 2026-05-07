package service

import (
	"app/internal/bus"
	"app/internal/serial"
	"context"
	"log/slog"
	"sync"
)

type bridgeEvent struct {
	topic   bus.BusTopic
	payload any
}

func (e *bridgeEvent) Topic() bus.BusTopic {
	return e.topic
}

func (e *bridgeEvent) Payload() any {
	return e.payload
}

var topicRouter = make(map[serial.FirmwareTopic]bus.BusTopic)

func RegisterRouter(serialTopic serial.FirmwareTopic, busTopic bus.BusTopic) {
	topicRouter[serialTopic] = busTopic
}

type SerialBridge struct {
	mu      sync.Mutex
	started bool
	ctx     context.Context
	cancel  context.CancelFunc

	eventBus *bus.EventBus
	byIdPath string
	baudRate int
	logger   *slog.Logger
}

func NewSerialBridge(ctx context.Context, eventBus *bus.EventBus, byIdPath string, baudRate int, logger *slog.Logger) *SerialBridge {
	return &SerialBridge{
		ctx:      ctx,
		eventBus: eventBus,
		byIdPath: byIdPath,
		baudRate: baudRate,
		logger:   logger,
	}
}

func (s *SerialBridge) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return nil
	}

	parentCtx := s.ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)
	s.cancel = cancel

	s.logger.Info("Serial bridge started", "port", s.byIdPath, "baudrate", s.baudRate)
	serial.ListenFirmware(
		ctx, s.byIdPath, s.baudRate, s.onMessage, s.onError,
	)
	s.started = true
	return nil
}

func (s *SerialBridge) onMessage(msg serial.FirmwareMessage, payload serial.FirmwarePayload) {
	busTopic, exists := topicRouter[msg.Topic]
	if !exists {
		s.logger.Warn("Unregistered firmware topic", "topic", msg.Topic)
		return
	}

	s.logger.Debug("Firmware message received", "topic", msg.Topic, "payload", payload)
	s.eventBus.Publish(&bridgeEvent{
		topic:   busTopic,
		payload: payload,
	})
}

func (s *SerialBridge) onError(err error) {
	s.logger.Error("Serial communication error", "error", err)
	s.eventBus.Publish(&bridgeEvent{
		topic:   bus.BusTopicSerialError,
		payload: err,
	})
}

func (s *SerialBridge) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return
	}

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.started = false
	s.logger.Info("Serial bridge stopped")
}
