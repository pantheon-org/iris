package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

const DefaultConfigFile = ".iris.json"

type Store struct {
	mu    sync.Mutex
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
	if cfg.Version != 0 && cfg.Version != 1 {
		return nil, fmt.Errorf("parse %s: version %d: %w", s.path, cfg.Version, ierrors.ErrUnsupportedVersion)
	}
	if cfg.Servers == nil {
		cfg.Servers = make(map[string]types.MCPServer)
	}
	return &cfg, nil
}

func (s *Store) Save(cfg *types.IrisConfig) (retErr error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	// Guarantee temp-file removal on any failure path.
	defer func() {
		if retErr != nil {
			if removeErr := os.Remove(tmpName); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				retErr = fmt.Errorf("%w; also failed to remove temp file %s: %v", retErr, tmpName, removeErr)
			}
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		// Close before the deferred Remove runs.
		if closeErr := tmp.Close(); closeErr != nil {
			return fmt.Errorf("write temp file: %w; close: %v", err, closeErr)
		}
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpName, s.path); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("rename to %s: %w", s.path, ierrors.ErrConfigPermission)
		}
		return fmt.Errorf("rename to %s: %w", s.path, err)
	}
	return nil
}
