package hardware

import "time"

type HardwareUpload struct {
	ID           uint64
	ShareCode    string
	HardwareData HardwareInfo
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func (h *HardwareUpload) IsExpired() bool {
	return time.Now().After(h.ExpiresAt)
}
