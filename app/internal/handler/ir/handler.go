package ir

import (
	"app/internal/serial"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
)

type BoundAction struct {
	ActionId   IrActionId
	Repeatable bool
	Params     json.RawMessage
}

type IrEventHandler struct {
	mu           sync.RWMutex
	ctx          *IrActionContext
	keyActionMap map[string]BoundAction
	actionMap    map[IrActionId]IrAction
	logger       *slog.Logger
}

func NewIrEventHandler(ctx *IrActionContext, logger *slog.Logger) *IrEventHandler {
	if ctx == nil {
		panic("IrActionContext cannot be nil")

	}

	if ctx.keyboard == nil || ctx.mouse == nil {
		panic("IrActionContext must have both Keyboard and Mouse initialized")
	}

	if logger == nil {
		panic("Logger cannot be nil")
	}

	return &IrEventHandler{
		ctx:          ctx,
		keyActionMap: make(map[string]BoundAction),
		actionMap:    make(map[IrActionId]IrAction),
		logger:       logger,
	}
}

func (eh *IrEventHandler) RegisterAction(actionId IrActionId, action IrAction) {
	eh.actionMap[actionId] = action
	eh.logger.Debug("Action registered", "action_id", actionId)
}

func (eh *IrEventHandler) Handle(payload *serial.IrPayload) error {
	if payload == nil {
		return errors.New("payload is nil")
	}

	key := NewIrKey(payload)

	eh.mu.RLock()
	boundAction, exists := eh.keyActionMap[key.String()]
	eh.mu.RUnlock()

	if !exists {
		eh.logger.Debug("No action mapped for IR key", "key", key.String())
		return errors.New("no action mapped for this key")
	}

	if !boundAction.Repeatable && payload.IsRepeat {
		eh.logger.Debug("Button is not repeatable and this is a repeat signal, ignoring", "key", key.String())
		return nil
	}

	action, exists := eh.actionMap[boundAction.ActionId]

	if !exists {
		eh.logger.Error("Action not found", "action_id", boundAction.ActionId)
		return errors.New("no action found for action ID: " + string(boundAction.ActionId))
	}

	err := action(eh.ctx, boundAction.Params)
	if err != nil {
		eh.logger.Error("Failed to execute action", "action_id", boundAction.ActionId, "error", err)
		return err
	}

	eh.logger.Debug("Action executed successfully", "action_id", boundAction.ActionId)
	return nil
}

func (eh *IrEventHandler) ReloadKeyActions(newMap map[string]BoundAction) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	eh.keyActionMap = newMap
	eh.logger.Info("Key actions reloaded", "total_actions", len(newMap))
}
