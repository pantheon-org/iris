package types

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/ierrors"
)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	Extension() string
}

func CodecForExtension(ext string) (Codec, error) {
	switch ext {
	case ".json":
		return JsonCodec{}, nil
	case ".yaml", ".yml":
		return YamlCodec{}, nil
	case ".toml":
		return TomlCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported extension %q: %w", ext, ierrors.ErrUnsupportedFormat)
	}
}
