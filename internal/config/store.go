package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

const DefaultConfigFile = ".iris.json"

type Store struct {
	path  string
	codec Codec
}

func NewStore(path string) (*Store, error) {
	if path == "" {
		path = DefaultConfigFile
	}
	codec, err := CodecForExtension(filepath.Ext(path))
	if err != nil {
		return nil, err
	}
	return &Store{path: path, codec: codec}, nil
}

func (s *Store) Path() string { return s.path }

func (s *Store) Load() (*types.IrisConfig, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("read %s: %w", s.path, ierrors.ErrConfigPermission)
		}
		return nil, fmt.Errorf("read %s: %w", s.path, err)
	}

	var cfg types.IrisConfig
	if err := s.codec.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", s.path, ierrors.ErrMalformedConfig)
	}
	if cfg.Servers == nil {
		cfg.Servers = make(map[string]types.MCPServer)
	}
	return &cfg, nil
}

func (s *Store) Save(cfg *types.IrisConfig) error {
	data, err := s.codec.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".iris-tmp-*")
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("create temp in %s: %w", dir, ierrors.ErrConfigPermission)
		}
		return fmt.Errorf("create temp in %s: %w", dir, err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpName, s.path); err != nil {
		_ = os.Remove(tmpName)
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("rename to %s: %w", s.path, ierrors.ErrConfigPermission)
		}
		return fmt.Errorf("rename to %s: %w", s.path, err)
	}
	return nil
}
