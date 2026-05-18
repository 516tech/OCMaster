package scanner

// Windows WMI + Win32 API 实现
// 编译: go build -buildmode=c-shared -o hardware_scanner.dll

func ScanAll() HardwareInfo {
	hw := HardwareInfo{}
	hw.CPU = scanCPUWindows()
	hw.Motherboard = scanMotherboardWindows()
	hw.RAM = scanRAMWindows()
	hw.GPU = scanGPUWindows()
	hw.PSU = PsuInfo{RatedWattage: "", Source: "manual"}
	hw.Cooler = CoolerInfo{Type: "air", Source: "manual"}
	return hw
}

func scanCPUWindows() CpuInfo {
	// TODO: Win32_Processor via WMI
	// SELECT Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed FROM Win32_Processor
	return CpuInfo{Model: "Unknown CPU", Cores: 0, Threads: 0, BaseFreq: "0 MHz"}
}

func scanMotherboardWindows() MotherboardInfo {
	// TODO: Win32_BaseBoard + Win32_BIOS
	return MotherboardInfo{Brand: "Unknown", Model: "Unknown", Chipset: "Unknown"}
}

func scanRAMWindows() RamInfo {
	// TODO: Win32_PhysicalMemory → Capacity, Speed, PartNumber, Manufacturer
	// SPD via SMBus I2C (需要 Ring0 驱动或 inpout32.dll)
	return RamInfo{TotalCapacity: "0 GB", StickCount: 0, ChannelCount: 0}
}

func scanGPUWindows() GpuInfo {
	// TODO: Win32_VideoController → Name, AdapterRAM
	return GpuInfo{Model: "Unknown GPU", VRAM: "0 GB"}
}
