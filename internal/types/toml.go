package types

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

type TomlCodec struct{}

func (TomlCodec) Extension() string { return ".toml" }

func (TomlCodec) Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(v); err != nil {
		return nil, fmt.Errorf("toml marshal: %w", err)
	}
	return buf.Bytes(), nil
}

func (TomlCodec) Unmarshal(data []byte, v any) error {
	if err := toml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("toml unmarshal: %w", err)
	}
	return nil
}

type TomlBuffer struct {
	data []byte
}

func (b *TomlBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *TomlBuffer) Bytes() []byte {
	return b.data
}
