package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/pantheon-org/iris/internal/ierrors"
)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	Extension() string
}

type jsonCodec struct{}

func (jsonCodec) Extension() string { return ".json" }

func (jsonCodec) Marshal(v any) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	return data, nil
}

func (jsonCodec) Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
}

type yamlCodec struct{}

func (yamlCodec) Extension() string { return ".yaml" }

func (yamlCodec) Marshal(v any) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("yaml marshal: %w", err)
	}
	return data, nil
}

func (yamlCodec) Unmarshal(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("yaml unmarshal: %w", err)
	}
	return nil
}

type tomlCodec struct{}

func (tomlCodec) Extension() string { return ".toml" }

func (tomlCodec) Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(v); err != nil {
		return nil, fmt.Errorf("toml marshal: %w", err)
	}
	return buf.Bytes(), nil
}

func (tomlCodec) Unmarshal(data []byte, v any) error {
	if err := toml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("toml unmarshal: %w", err)
	}
	return nil
}

func CodecForExtension(ext string) (Codec, error) {
	switch ext {
	case ".json":
		return jsonCodec{}, nil
	case ".yaml", ".yml":
		return yamlCodec{}, nil
	case ".toml":
		return tomlCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported extension %q: %w", ext, ierrors.ErrMalformedConfig)
	}
}
