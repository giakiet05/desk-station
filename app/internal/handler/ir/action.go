package ir

import (
	"encoding/json"
	"errors"

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
	Keyboard uinput.Keyboard
	Mouse    uinput.Mouse
}

func NewIrActionContext(keyboard uinput.Keyboard, mouse uinput.Mouse) *IrActionContext {
	if keyboard == nil || mouse == nil {
		panic("Both keyboard and mouse must be provided")
	}

	return &IrActionContext{
		Keyboard: keyboard,
		Mouse:    mouse,
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
		return ctx.Keyboard.KeyPress(uinput.KeyVolumedown)
	case 1:
		return ctx.Keyboard.KeyPress(uinput.KeyVolumeup)
	default:
		return errors.New("invalid parameter value for volume action")
	}

}

func Mute(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyMute)
}

func MouseClick(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p MouseClickParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for mouse click action")
	}

	var clickFunc func() error

	switch p.Button {
	case -1:
		clickFunc = ctx.Mouse.RightClick
	case 0:
		clickFunc = ctx.Mouse.MiddleClick
	case 1:
		clickFunc = ctx.Mouse.LeftClick
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

	return ctx.Mouse.Move(int32(p.DX), int32(p.DY))
}

func MouseScroll(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p MouseScrollParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for mouse scroll action")
	}

	return ctx.Mouse.Wheel(p.Direction, int32(p.Amount))
}

func RunCommand(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p RunCommandParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for run command action")
	}

	// Here you would implement the logic to run the command with the given arguments.
	// This is a placeholder and should be replaced with actual command execution code.
	return errors.New("run command action is not implemented yet")
}

func PlayPause(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyPlaypause)
}

func NextTrack(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyNextsong)
}

func PreviousTrack(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyPrevioussong)
}

func PrintScreen(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyPrint)
}

func Enter(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyEnter)
}

func Escape(ctx *IrActionContext, _ json.RawMessage) error {
	return ctx.Keyboard.KeyPress(uinput.KeyEsc)
}

func Arrow(ctx *IrActionContext, rawParams json.RawMessage) error {
	var p ArrowKeyParams

	if err := json.Unmarshal(rawParams, &p); err != nil {
		return errors.New("invalid parameters for arrow key action")
	}

	switch p.Direction {
	case "up":
		return ctx.Keyboard.KeyPress(uinput.KeyUp)
	case "down":
		return ctx.Keyboard.KeyPress(uinput.KeyDown)
	case "left":
		return ctx.Keyboard.KeyPress(uinput.KeyLeft)
	case "right":
		return ctx.Keyboard.KeyPress(uinput.KeyRight)
	default:
		return errors.New("invalid parameter value for arrow key action")
	}
}
