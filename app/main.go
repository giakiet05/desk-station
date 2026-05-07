package main

import (
	"app/internal/bus"
	"app/internal/config"
	"app/internal/handler/ir"
	"app/internal/logging"
	"app/internal/serial"
	"app/internal/service"
	"context"
	"embed"

	"github.com/bendahl/uinput"
	"github.com/joho/godotenv"
)

// go:embed all:frontend/dist
var assets embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		println("No .env file found, using environment variables")
	}

	byIdPath := "/dev/serial/by-id/usb-Arduino_RaspberryPi_Pico_C52864E6334D9F1F-if00"
	baudRate := 115200

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("desk-station"))
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("desk-station-mouse"))
	if err != nil {
		panic(err)
	}
	defer mouse.Close()

	bus := bus.NewEventBus(20)
	logger := logging.NewLogger()

	irActionContext := ir.NewIrActionContext(keyboard, mouse, logger)
	irHandler := ir.NewIrEventHandler(irActionContext, logger)
	irRouter := service.NewIrRouter(ctx, irHandler, bus, logger)

	serial.RegisterPayloads()
	ir.RegisterDefaultActions(irHandler)
	service.RegisterRouters()

	presetDir := "presets"

	newKeyActionMap, err := ir.ParsePresetDir(presetDir)
	if err != nil {
		panic(err)
	}

	irHandler.ReloadKeyActions(newKeyActionMap)

	config.WatchPresetDir(presetDir, irHandler, logger)

	serialBridge := service.NewSerialBridge(ctx, bus, byIdPath, baudRate, logger)

	if err := irRouter.Start(); err != nil {
		panic(err)
	}

	if err := serialBridge.Start(); err != nil {
		panic(err)
	}

	select {}
}

// // Create an instance of the app structure
// app := NewApp()

// // Create application with options
// err := wails.Run(&options.App{
// 	Title:  "app",
// 	Width:  1024,
// 	Height: 768,
// 	AssetServer: &assetserver.Options{
// 		Assets: assets,
// 	},
// 	BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
// 	OnStartup:        app.startup,
// 	Bind: []interface{}{
// 		app,
// 	},
// })

// if err != nil {
// 	println("Error:", err.Error())
// }
