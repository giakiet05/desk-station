package ir

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os/exec"

	"github.com/bendahl/uinput"
)

type IrActionId string

const (
	IrActionVolume IrActionId = "volume"
	IrActionMute   IrActionId = "mute"

	IrActionMouseClick  IrActionId = "mouse_click"
	IrActionMouseMove   IrActionId = "mouse_move"
	IrActionMouseScroll IrActionId = "mouse_scroll"

	IrActionRunCommand IrActionId = "run_command"

	IrActionPlayPause     IrActionId = "play_pause"
	IrActionNextTrack     IrActionId = "next_track"
	IrActionPreviousTrack IrActionId = "previous_track"

	IrActionPrintScreen IrActionId = "print_screen"
	IrActionEnter       IrActionId = "enter"
	IrActionEscape      IrActionId = "escape"

	IrActionArrow IrActionId = "arrow"
)

type IrActionContext struct {
	keyboard uinput.Keyboard
	mouse    uinput.Mouse
	logger   *slog.Logger
}

func NewIrActionContext(keyboard uinput.Keyboard, mouse uinput.Mouse, logger *slog.Logger) *IrActionContext {
	if keyboard == nil || mouse == nil || logger == nil {
		panic("All parameters for IrActionContext must be provided")
	}

	return &IrActionContext{
		keyboard: keyboard,
		mouse:    mouse,
		logger:   logger,
	}
}

type IrAction func(ctx *IrActionContext, rawParams json.RawMessage) error

func Volume(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p VolumeParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for volume action")
	}

	switch p.Delta {
	case -1:
		return ctx.keyboard.KeyPress(uinput.KeyVolumedown)
	case 1:
		return ctx.keyboard.KeyPress(uinput.KeyVolumeup)
	default:
		return errors.New("invalid parameter value for volume action")
	}

}

func Mute(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyMute)
}

func MouseClick(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p MouseClickParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for mouse click action")
	}

	var clickFunc func() error

	switch p.Button {
	case -1:
		clickFunc = ctx.mouse.RightClick
	case 0:
		clickFunc = ctx.mouse.MiddleClick
	case 1:
		clickFunc = ctx.mouse.LeftClick
	default:
		return errors.New("invalid parameter value for mouse click action")
	}

	for i := 0; i < p.Count; i++ {
		if err := clickFunc(); err != nil {
			return err
		}
	}

	return nil
}

func MouseMove(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p MouseMoveParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for mouse move action")
	}

	return ctx.mouse.Move(int32(p.DX), int32(p.DY))
}

func MouseScroll(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p MouseScrollParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for mouse scroll action")
	}

	return ctx.mouse.Wheel(p.Direction, int32(p.Amount))
}

func PlayPause(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyPlaypause)
}

func NextTrack(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyNextsong)
}

func PreviousTrack(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyPrevioussong)
}

func PrintScreen(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyPrint)
}

func Enter(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyEnter)
}

func Escape(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.keyboard.KeyPress(uinput.KeyEsc)
}

func Arrow(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p ArrowKeyParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for arrow key action")
	}

	switch p.Direction {
	case "up":
		return ctx.keyboard.KeyPress(uinput.KeyUp)
	case "down":
		return ctx.keyboard.KeyPress(uinput.KeyDown)
	case "left":
		return ctx.keyboard.KeyPress(uinput.KeyLeft)
	case "right":
		return ctx.keyboard.KeyPress(uinput.KeyRight)
	default:
		return errors.New("invalid parameter value for arrow key action")
	}
}

func RunCommand(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p RunCommandParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for run command action")
	}

	if p.Command == "" {
		return errors.New("command cannot be empty for run command action")
	}

	cmd := exec.Command(p.Command, p.Args...)

	err := cmd.Start()
	if err != nil {
		return errors.New("failed to start command")
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			ctx.logger.Error("Command execution failed", "command", p.Command, "args", p.Args, "error", err)
		}
		ctx.logger.Debug("Command execution finished", "command", p.Command, "args", p.Args)

	}()

	return nil
}
