package scanner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var scanLog *log.Logger

func init() {
	logDir := filepath.Join(os.TempDir(), "OCMaster")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("ocmaster-scan-%s.log", time.Now().Format("20060102-150405")))
	f, err := os.Create(logPath)
	if err == nil {
		scanLog = log.New(f, "", log.LstdFlags)
	} else {
		scanLog = log.New(os.Stderr, "SCAN: ", log.LstdFlags)
	}
}

func ScanAll() HardwareInfo {
	scanLog.Println("=== ScanAll: START ===")
	hw := HardwareInfo{}
	scanLog.Println("ScanAll: CPU...")
	hw.CPU = scanCPUWindows()
	scanLog.Printf("ScanAll: CPU done — model=%s cores=%d", hw.CPU.Model, hw.CPU.Cores)
	scanLog.Println("ScanAll: Motherboard...")
	hw.Motherboard = scanMotherboardWindows()
	scanLog.Printf("ScanAll: MB done — brand=%s model=%s", hw.Motherboard.Brand, hw.Motherboard.Model)
	scanLog.Println("ScanAll: RAM...")
	hw.RAM = scanRAMWindows()
	scanLog.Printf("ScanAll: RAM done — total=%s sticks=%d", hw.RAM.TotalCapacity, hw.RAM.StickCount)
	scanLog.Println("ScanAll: GPU...")
	hw.GPU = scanGPUWindows()
	scanLog.Printf("ScanAll: GPU done — model=%s vram=%s", hw.GPU.Model, hw.GPU.VRAM)
	hw.PSU = PsuInfo{RatedWattage: "", Source: "manual"}
	hw.Cooler = CoolerInfo{Type: "air", Source: "manual"}
	scanLog.Println("=== ScanAll: DONE ===")
	return hw
}

// wmiQuery uses PowerShell Get-CimInstance (wmic is deprecated on Win10/11)
func wmiQuery(class string, props ...string) []string {
	scanLog.Printf("wmi: Win32_%s get %s", class, strings.Join(props, ","))
	psCmd := fmt.Sprintf(
		"Get-CimInstance Win32_%s | Select-Object %s | ConvertTo-Csv -NoTypeInformation",
		class, strings.Join(props, ","))
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		scanLog.Printf("wmi: Win32_%s FAILED — %v", class, err)
		return nil
	}
	scanLog.Printf("wmi: Win32_%s OK — %d bytes", class, len(out))
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Skip header line (CSV column names)
	var result []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || i == 0 {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}

// csvVal 返回 CSV 行中指定索引的值
func csvVal(line string, idx int) string {
	parts := strings.Split(line, ",")
	if idx < len(parts) {
		return strings.TrimSpace(parts[idx])
	}
	return ""
}

func scanCPUWindows() CpuInfo {
	cpu := CpuInfo{Model: "Unknown CPU", Cores: 0, Threads: 0, BaseFreq: "0 MHz"}

	lines := wmiQuery("Processor", "Name", "NumberOfCores", "NumberOfLogicalProcessors", "MaxClockSpeed")
	if len(lines) == 0 {
		return cpu
	}

	// 格式: HOSTNAME,Name,NumberOfCores,NumberOfLogicalProcessors,MaxClockSpeed
	line := lines[0]
	cpu.Model = csvVal(line, 1)
	if v, err := strconv.Atoi(csvVal(line, 2)); err == nil {
		cpu.Cores = v
	}
	if v, err := strconv.Atoi(csvVal(line, 3)); err == nil {
		cpu.Threads = v
	}
	if freq := csvVal(line, 4); freq != "" {
		if mhz, err := strconv.Atoi(freq); err == nil && mhz > 0 {
			cpu.BaseFreq = fmt.Sprintf("%.0f MHz", float64(mhz))
		}
	}

	if cpu.Model == "" {
		cpu.Model = "Unknown CPU"
	}
	return cpu
}

func scanMotherboardWindows() MotherboardInfo {
	mb := MotherboardInfo{Brand: "Unknown", Model: "Unknown"}

	// Win32_BaseBoard
	bbLines := wmiQuery("BaseBoard", "Manufacturer", "Product")
	if len(bbLines) > 0 {
		mb.Brand = csvVal(bbLines[0], 1)
		mb.Model = csvVal(bbLines[0], 2)
	}

	// BIOS version
	biosLines := wmiQuery("BIOS", "SMBIOSBIOSVersion")
	if len(biosLines) > 0 {
		mb.BiosVersion = csvVal(biosLines[0], 1)
	}

	if mb.Brand == "" {
		mb.Brand = "Unknown"
	}
	if mb.Model == "" {
		mb.Model = "Unknown"
	}
	return mb
}

func scanRAMWindows() RamInfo {
	ram := RamInfo{Sticks: []RamStick{}}

	// Win32_PhysicalMemory
	lines := wmiQuery("PhysicalMemory", "Capacity", "Speed", "ConfiguredClockSpeed", "SMBIOSMemoryType")
	totalBytes := uint64(0)
	for _, line := range lines {
		stick := RamStick{}
		capStr := csvVal(line, 1)
		if capBytes, err := strconv.ParseUint(capStr, 10, 64); err == nil && capBytes > 0 {
			totalBytes += capBytes
			stick.Capacity = fmt.Sprintf("%.0f GB", float64(capBytes)/1_073_741_824)
		}
		speed := csvVal(line, 3) // ConfiguredClockSpeed 优先
		if speed == "" {
			speed = csvVal(line, 2) // Speed fallback
		}
		if speed != "" {
			stick.Frequency = speed + " MHz"
		}
		dieType := csvVal(line, 4)
		stick.DieType = memTypeName(dieType)

		if stick.Capacity != "" {
			ram.Sticks = append(ram.Sticks, stick)
		}
	}

	if totalBytes > 0 {
		ram.TotalCapacity = fmt.Sprintf("%.0f GB", float64(totalBytes)/1_073_741_824)
	} else {
		// Fallback: Win32_OperatingSystem TotalVisibleMemorySize
		osLines := wmiQuery("OperatingSystem", "TotalVisibleMemorySize")
		if len(osLines) > 0 {
			if kb, err := strconv.ParseUint(csvVal(osLines[0], 1), 10, 64); err == nil && kb > 0 {
				ram.TotalCapacity = fmt.Sprintf("%.0f GB", float64(kb)/1_048_576)
			}
		}
	}

	ram.StickCount = len(ram.Sticks)
	if ram.StickCount >= 4 {
		ram.ChannelCount = 4
	} else if ram.StickCount >= 2 {
		ram.ChannelCount = 2
	} else if ram.StickCount == 1 {
		ram.ChannelCount = 1
	}
	if ram.TotalCapacity == "" {
		ram.TotalCapacity = "0 GB"
	}
	return ram
}

func memTypeName(code string) string {
	switch code {
	case "20": return "DDR"
	case "21": return "DDR2"
	case "22": return "DDR2 FB-DIMM"
	case "24": return "DDR3"
	case "26": return "DDR4"
	case "27": return "DDR4"
	case "28": return "DDR4"
	case "34": return "DDR5"
	default: return code
	}
}

func scanGPUWindows() GpuInfo {
	gpu := GpuInfo{Model: "Unknown GPU", VRAM: "Unknown"}

	lines := wmiQuery("VideoController", "Name", "AdapterRAM")
	if len(lines) == 0 {
		return gpu
	}

	line := lines[0]
	gpu.Model = csvVal(line, 1)
	if ramBytes, err := strconv.ParseUint(csvVal(line, 2), 10, 64); err == nil && ramBytes > 0 {
		gpu.VRAM = fmt.Sprintf("%.0f GB", float64(ramBytes)/1_073_741_824)
	}

	if gpu.Model == "" {
		gpu.Model = "Unknown GPU"
	}
	return gpu
}
