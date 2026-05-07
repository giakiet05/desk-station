package ir

import (
	"encoding/json"
	"os"
)

type PresetFile struct {
	Address uint16       `json:"address"`
	Items   []PresetItem `json:"items"`
}

type PresetItem struct {
	Command    uint16          `json:"command"`
	Action     IrActionId      `json:"action"`
	Repeatable bool            `json:"repeatable"`
	Params     json.RawMessage `json:"params"`
}

func LoadPreset(filePath string, eh *IrEventHandler) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var fileData PresetFile
	if err := json.Unmarshal(data, &fileData); err != nil {
		return err
	}

	address := fileData.Address

	for _, item := range fileData.Items {
		key := IrKey{
			Address: address,
			Command: item.Command,
		}

		eh.RegisterKeyAction(key, item.Action, item.Params, item.Repeatable)

	}
	return nil
}
