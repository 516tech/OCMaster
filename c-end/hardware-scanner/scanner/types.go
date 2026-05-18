package scanner

// CpuInfo CPU 信息
type CpuInfo struct {
	Model       string `json:"model"`
	Cores       int    `json:"cores"`
	Threads     int    `json:"threads"`
	BaseFreq    string `json:"baseFreq"`
	BiosVersion string `json:"biosVersion"`
}

// MotherboardInfo 主板信息
type MotherboardInfo struct {
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Chipset     string `json:"chipset"`
	BiosVersion string `json:"biosVersion"`
}

// RamStick 单条内存信息
type RamStick struct {
	Capacity  string `json:"capacity"`
	Frequency string `json:"frequency"`
	Timings   string `json:"timings"`
	DieType   string `json:"dieType"`
}

// RamInfo 内存汇总信息
type RamInfo struct {
	TotalCapacity string     `json:"totalCapacity"`
	StickCount    int        `json:"stickCount"`
	ChannelCount  int        `json:"channelCount"`
	Sticks        []RamStick `json:"sticks"`
}

// GpuInfo 显卡信息
type GpuInfo struct {
	Model string `json:"model"`
	VRAM  string `json:"vram"`
}

// PsuInfo 电源信息
type PsuInfo struct {
	RatedWattage string `json:"ratedWattage"`
	Source       string `json:"source"` // "smbus" | "manual"
}

// CoolerInfo 散热信息
type CoolerInfo struct {
	Type   string `json:"type"`   // "air" | "aio" | "custom_loop"
	Source string `json:"source"` // "auto" | "manual"
}

// HardwareInfo 全部硬件信息
type HardwareInfo struct {
	CPU         CpuInfo         `json:"cpu"`
	Motherboard MotherboardInfo `json:"motherboard"`
	RAM         RamInfo         `json:"ram"`
	GPU         GpuInfo         `json:"gpu"`
	PSU         PsuInfo         `json:"psu"`
	Cooler      CoolerInfo      `json:"cooler"`
}
