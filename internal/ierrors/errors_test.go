package ierrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
)

func TestErrServerNotFound_directMatch_isTrue(t *testing.T) {
	if !errors.Is(ierrors.ErrServerNotFound, ierrors.ErrServerNotFound) {
		t.Fatal("expected ErrServerNotFound to match itself")
	}
}

func TestErrServerNotFound_wrapped_isTrue(t *testing.T) {
	wrapped := fmt.Errorf("ctx: %w", ierrors.ErrServerNotFound)
	if !errors.Is(wrapped, ierrors.ErrServerNotFound) {
		t.Fatal("expected wrapped ErrServerNotFound to match via errors.Is")
	}
}

func TestErrServerNotFound_crossMatch_isFalse(t *testing.T) {
	if errors.Is(ierrors.ErrServerNotFound, ierrors.ErrMalformedConfig) {
		t.Fatal("expected ErrServerNotFound to not match ErrMalformedConfig")
	}
}

func TestErrMalformedConfig_directMatch_isTrue(t *testing.T) {
	if !errors.Is(ierrors.ErrMalformedConfig, ierrors.ErrMalformedConfig) {
		t.Fatal("expected ErrMalformedConfig to match itself")
	}
}

func TestErrMalformedConfig_wrapped_isTrue(t *testing.T) {
	wrapped := fmt.Errorf("ctx: %w", ierrors.ErrMalformedConfig)
	if !errors.Is(wrapped, ierrors.ErrMalformedConfig) {
		t.Fatal("expected wrapped ErrMalformedConfig to match via errors.Is")
	}
}

func TestErrProviderNotFound_directMatch_isTrue(t *testing.T) {
	if !errors.Is(ierrors.ErrProviderNotFound, ierrors.ErrProviderNotFound) {
		t.Fatal("expected ErrProviderNotFound to match itself")
	}
}

func TestErrProviderNotFound_wrapped_isTrue(t *testing.T) {
	wrapped := fmt.Errorf("ctx: %w", ierrors.ErrProviderNotFound)
	if !errors.Is(wrapped, ierrors.ErrProviderNotFound) {
		t.Fatal("expected wrapped ErrProviderNotFound to match via errors.Is")
	}
}

func TestErrConfigPermission_directMatch_isTrue(t *testing.T) {
	if !errors.Is(ierrors.ErrConfigPermission, ierrors.ErrConfigPermission) {
		t.Fatal("expected ErrConfigPermission to match itself")
	}
}

func TestErrConfigPermission_wrapped_isTrue(t *testing.T) {
	wrapped := fmt.Errorf("ctx: %w", ierrors.ErrConfigPermission)
	if !errors.Is(wrapped, ierrors.ErrConfigPermission) {
		t.Fatal("expected wrapped ErrConfigPermission to match via errors.Is")
	}
}

func TestErrConfigPermission_crossMatch_isFalse(t *testing.T) {
	if errors.Is(ierrors.ErrConfigPermission, ierrors.ErrProviderNotFound) {
		t.Fatal("expected ErrConfigPermission to not match ErrProviderNotFound")
	}
}
