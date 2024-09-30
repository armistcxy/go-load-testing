package attacker

import (
	"encoding/json"
	"fmt"
	"os"
)

type Figure struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Fields map[string]string `json:"fields"`
}

func RetrieveFigure(path string) (*Figure, error) {
	fig := Figure{}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fields from json file, error: %s", err)
	}

	if err := json.NewDecoder(file).Decode(&fig); err != nil {
		return nil, err
	}

	return &fig, nil
}
