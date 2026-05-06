package bus

import (
	"app/internal/logging"
	"log/slog"
	"sync"
)

type Event interface {
	Topic() BusTopic
	Payload() any
}

type EventChan chan Event

type EventBus struct {
	subscribers      map[BusTopic][]EventChan
	mu               sync.RWMutex
	closed           bool
	subscriberBuffer int
	logger           *slog.Logger
}

// Close implements [EventBus].
func (e *EventBus) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return
	}
	e.closed = true

	for _, subscribers := range e.subscribers {
		for _, s := range subscribers {
			close(s)
		}
	}

	e.subscribers = make(map[BusTopic][]EventChan)
	e.logger.Info("Event bus closed")
}

// Publish implements [EventBus].
func (e *EventBus) Publish(event Event) {

	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.closed || event == nil {
		return
	}

	topic := event.Topic()
	if topic == "" {
		return
	}

	if subscribers, ok := e.subscribers[topic]; ok {
		for _, s := range subscribers {
			// Non-blocking send
			// Bus will not be block if there is 1 slow subscriber
			// However, the event will be dropped if channel is full
			select {
			case s <- event:
			default:
				e.logger.Warn("Event dropped - subscriber channel full", "topic", topic)
			}
		}
	}
}

// Subscribe implements [EventBus].
func (e *EventBus) Subscribe(topic BusTopic) EventChan {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed || topic == "" {
		return nil
	}

	ch := make(EventChan, e.subscriberBuffer)

	e.subscribers[topic] = append(e.subscribers[topic], ch)
	totalSubscribers := len(e.subscribers[topic])
	e.logger.Debug("Subscriber registered", "topic", topic, "total_subscribers", totalSubscribers)

	return ch
}

// Unsubscribe implements [EventBus].
func (e *EventBus) Unsubscribe(topic BusTopic, ch EventChan) {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed || topic == "" || ch == nil {
		return
	}

	subscribers, ok := e.subscribers[topic]
	if !ok || len(subscribers) == 0 {
		return
	}

	filtered := subscribers[:0]
	for _, sub := range subscribers {
		if sub == ch {
			close(ch)
		} else {
			filtered = append(filtered, sub)
		}
	}

	for i := len(filtered); i < len(subscribers); i++ {
		subscribers[i] = nil
	}

	if len(filtered) == 0 {
		delete(e.subscribers, topic)
		e.logger.Debug("Subscriber unregistered", "topic", topic, "total_subscribers", 0)
		return
	}

	e.subscribers[topic] = filtered
	e.logger.Debug("Subscriber unregistered", "topic", topic, "total_subscribers", len(filtered))
}

func NewEventBus(subscriberBuffer int) *EventBus {
	if subscriberBuffer <= 0 {
		subscriberBuffer = 16
	}

	return &EventBus{
		subscribers:      make(map[BusTopic][]EventChan),
		subscriberBuffer: subscriberBuffer,
		logger:           logging.NewLogger(),
	}
}
