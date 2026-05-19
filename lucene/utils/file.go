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

func WriteJsonFile(filename string, v any) error {
	f, err := os.OpenFile(filename, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(v); err != nil {
		return err
	}

	return nil
}
