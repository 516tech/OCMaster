package suggestion

import "time"

type CpuSuggestion struct {
	Frequency string `json:"frequency"`
	Voltage   string `json:"voltage"`
	Notes     string `json:"notes"`
}

type RamSuggestion struct {
	Frequency string `json:"frequency"`
	Timings   string `json:"timings"`
	Voltage   string `json:"voltage"`
	Notes     string `json:"notes"`
}

type StabilityTest struct {
	Tool     string `json:"tool"`
	Duration string `json:"duration"`
}

type Suggestion struct {
	ID             uint64
	MerchantID     uint64
	HardwareID     uint64
	CpuSuggestion  *CpuSuggestion
	RamSuggestion  *RamSuggestion
	StabilityTest  *StabilityTest
	RiskWarning    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
