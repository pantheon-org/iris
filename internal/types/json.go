package types

import (
	"encoding/json"
	"fmt"
)

type JsonCodec struct{}

func (JsonCodec) Extension() string { return ".json" }

func (JsonCodec) Marshal(v any) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	return data, nil
}

func (JsonCodec) Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
}
