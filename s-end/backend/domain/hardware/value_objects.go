package hardware

type CpuInfo struct {
	Model       string `json:"model"`
	Cores       int    `json:"cores"`
	Threads     int    `json:"threads"`
	BaseFreq    string `json:"base_freq"`
	BiosVersion string `json:"bios_version"`
}

type MotherboardInfo struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Chipset      string `json:"chipset"`
	BiosVersion  string `json:"bios_version"`
}

type RamStick struct {
	Capacity  string `json:"capacity"`
	Frequency string `json:"frequency"`
	Timings   string `json:"timings"`
	DieType   string `json:"die_type"`
}

type RamInfo struct {
	TotalCapacity string     `json:"total_capacity"`
	StickCount    int        `json:"stick_count"`
	Sticks        []RamStick `json:"sticks"`
	ChannelCount  int        `json:"channel_count"`
}

type GpuInfo struct {
	Model string `json:"model"`
	VRAM  string `json:"vram"`
}

type PsuInfo struct {
	RatedWattage string `json:"rated_wattage"`
	Source       string `json:"source"` // "smbus" or "manual"
}

type CoolerInfo struct {
	Type   string `json:"type"`   // "air", "aio", "custom_loop"
	Source string `json:"source"` // "auto" or "manual"
}

type HardwareInfo struct {
	CPU         *CpuInfo         `json:"cpu"`
	Motherboard *MotherboardInfo `json:"motherboard"`
	RAM         *RamInfo         `json:"ram"`
	GPU         *GpuInfo         `json:"gpu"`
	PSU         *PsuInfo         `json:"psu"`
	Cooler      *CoolerInfo      `json:"cooler"`
}
