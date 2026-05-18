package utils

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(filename string, v any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(f).Decode(v); err != nil {
		return err
	}

	return nil
}
