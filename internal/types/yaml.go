package types

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YamlCodec struct{}

func (YamlCodec) Extension() string { return ".yaml" }

func (YamlCodec) Marshal(v any) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("yaml marshal: %w", err)
	}
	return data, nil
}

func (YamlCodec) Unmarshal(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("yaml unmarshal: %w", err)
	}
	return nil
}
