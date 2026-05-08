package service

import (
	"app/internal/bus"
	"app/internal/handler/ir"
	"app/internal/serial"
	"context"
	"errors"
	"log/slog"
	"sync"
)

type IrRouter struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	eventBus  *bus.EventBus
	logger    *slog.Logger
	started   bool
	eventChan bus.EventChan
	handler   *ir.IrEventHandler
}

func NewIrRouter(ctx context.Context, handler *ir.IrEventHandler, eventBus *bus.EventBus, logger *slog.Logger) *IrRouter {
	return &IrRouter{
		ctx:      ctx,
		handler:  handler,
		eventBus: eventBus,
		logger:   logger,
	}
}

func (r *IrRouter) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		return nil
	}

	if r.eventBus == nil {
		return errors.New("event bus cannot be nil")
	}

	parentCtx := r.ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	ch := r.eventBus.Subscribe(bus.BusTopicIr)
	if ch == nil {
		cancel()
		return errors.New("failed to subscribe to event bus")
	}

	r.cancel = cancel
	r.eventChan = ch
	r.started = true

	r.logger.Info("IR router started")
	go r.loop(ctx, ch)
	return nil
}

func (r *IrRouter) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return
	}

	cancel := r.cancel

	ch := r.eventChan

	r.started = false
	r.cancel = nil
	r.eventChan = nil

	if cancel != nil {
		cancel()
	}

	if ch != nil {
		r.eventBus.Unsubscribe(bus.BusTopicIr, ch)
	}

	r.logger.Info("IR router stopped")
}

func (r *IrRouter) loop(ctx context.Context, ch bus.EventChan) {
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("IR router loop exiting")
			return
		case event, ok := <-ch:
			if !ok {
				r.logger.Info("IR router event channel closed")
				return
			}
			r.handleEvent(event)
		}
	}

}

func (r *IrRouter) handleEvent(event bus.Event) {
	if event == nil {
		return
	}

	switch p := event.Payload().(type) {
	case *serial.IrPayload:
		if p != nil {
			r.logger.Debug("IR event received", "address", p.Address, "command", p.Command)
			if err := r.handler.Handle(p); err != nil {
				r.logger.Error("Error handling event", "err", err, "event")
			}

		}
	case serial.IrPayload:
		r.logger.Debug("IR event received", "address", p.Address, "command", p.Command)
		if err := r.handler.Handle(&p); err != nil {
			r.logger.Error("Error handling event", "err", err, "event", event)
		}

	default:
		r.logger.Error("Unhandled event type", slog.Any("event", event))
	}
}
