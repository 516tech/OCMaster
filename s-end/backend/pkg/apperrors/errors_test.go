package apperrors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	errs := []error{ErrNotFound, ErrUnauthorized, ErrForbidden, ErrInvalidInput, ErrConflict, ErrInternal, ErrShareCodeExpired}
	for _, e := range errs {
		if e.Error() == "" {
			t.Errorf("expected non-empty error message for %v", e)
		}
	}
}

func TestErrors_Wrapping(t *testing.T) {
	wrapped := errors.Join(ErrNotFound, errors.New("detail"))
	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("expected ErrNotFound in wrapped error")
	}
}
