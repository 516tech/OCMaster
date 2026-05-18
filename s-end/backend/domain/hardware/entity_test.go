package hardware

import (
	"testing"
	"time"
)

func TestHardwareUpload_IsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		h := &HardwareUpload{ExpiresAt: time.Now().Add(1 * time.Hour)}
		if h.IsExpired() {
			t.Error("expected not expired")
		}
	})
	t.Run("expired", func(t *testing.T) {
		h := &HardwareUpload{ExpiresAt: time.Now().Add(-1 * time.Hour)}
		if !h.IsExpired() {
			t.Error("expected expired")
		}
	})
}

func TestHardwareInfo_ZeroValues(t *testing.T) {
	info := &HardwareInfo{}
	if info.CPU != nil {
		t.Error("expected nil CPU")
	}
}
