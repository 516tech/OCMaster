package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ScanAll() HardwareInfo {
	hw := HardwareInfo{}
	hw.CPU = scanCPULinux()
	hw.Motherboard = scanMotherboardLinux()
	hw.RAM = scanRAMLinux()
	hw.GPU = scanGPULinux()
	hw.PSU = PsuInfo{RatedWattage: "", Source: "manual"}
	hw.Cooler = CoolerInfo{Type: "air", Source: "manual"}
	return hw
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func scanCPULinux() CpuInfo {
	cpu := CpuInfo{}

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return CpuInfo{Model: "Unknown CPU", Cores: 0, Threads: 0, BaseFreq: "0 MHz"}
	}

	lines := strings.Split(string(data), "\n")
	coreSet := make(map[string]bool)
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "model name":
			if cpu.Model == "" {
				cpu.Model = val
			}
		case "cpu cores":
			v, _ := strconv.Atoi(val)
			if v > cpu.Cores {
				cpu.Cores = v
			}
		case "cpu MHz":
			if cpu.BaseFreq == "" {
				if mhz, err := strconv.ParseFloat(val, 64); err == nil {
					cpu.BaseFreq = fmt.Sprintf("%.0f MHz", mhz)
				}
			}
		case "physical id":
			coreSet[val] = true
		case "processor":
			cpu.Threads++
		}
	}

	if cpu.Threads == 0 {
		cpu.Threads = cpu.Cores * 2
	}
	if cpu.BaseFreq == "" {
		cpu.BaseFreq = "Unknown"
	}

	return cpu
}

func scanMotherboardLinux() MotherboardInfo {
	mb := MotherboardInfo{}

	mb.Brand = readFile("/sys/class/dmi/id/board_vendor")
	mb.Model = readFile("/sys/class/dmi/id/board_name")
	mb.BiosVersion = readFile("/sys/class/dmi/id/bios_version")

	if mb.Brand == "" {
		// Fallback: dmidecode
		out, err := exec.Command("dmidecode", "-s", "system-manufacturer").Output()
		if err == nil {
			mb.Brand = strings.TrimSpace(string(out))
		}
	}
	if mb.Model == "" {
		out, err := exec.Command("dmidecode", "-s", "system-product-name").Output()
		if err == nil {
			mb.Model = strings.TrimSpace(string(out))
		}
	}

	if mb.Brand == "" {
		mb.Brand = "Unknown"
	}
	if mb.Model == "" {
		mb.Model = "Unknown"
	}

	return mb
}

func scanRAMLinux() RamInfo {
	ram := RamInfo{Sticks: []RamStick{}}

	// 总内存: /proc/meminfo MemTotal
	memInfo, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		for _, line := range strings.Split(string(memInfo), "\n") {
			if strings.HasPrefix(line, "MemTotal:") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					kb, _ := strconv.ParseInt(parts[1], 10, 64)
					ram.TotalCapacity = fmt.Sprintf("%.0f GB", float64(kb)/1_048_576)
				}
			}
		}
	}

	// dmidecode -t 17: 每条内存槽详情
	out, err := exec.Command("dmidecode", "-t", "17").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		var stick *RamStick
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "Size:") {
				if stick != nil && stick.Capacity != "" && stick.Capacity != "No Module Installed" {
					ram.Sticks = append(ram.Sticks, *stick)
				}
				stick = &RamStick{}
				val := strings.TrimSpace(strings.TrimPrefix(trimmed, "Size:"))
				stick.Capacity = val
			}
			if stick != nil {
				switch {
				case strings.HasPrefix(trimmed, "Speed:"):
					stick.Frequency = strings.TrimSpace(strings.TrimPrefix(trimmed, "Speed:"))
				case strings.HasPrefix(trimmed, "Type:"):
					stick.DieType = strings.TrimSpace(strings.TrimPrefix(trimmed, "Type:"))
				case strings.HasPrefix(trimmed, "Configured Memory Speed:"):
					if stick.Frequency == "" || stick.Frequency == "Unknown" {
						stick.Frequency = strings.TrimSpace(strings.TrimPrefix(trimmed, "Configured Memory Speed:"))
					}
				}
			}
		}
		if stick != nil && stick.Capacity != "" && stick.Capacity != "No Module Installed" {
			ram.Sticks = append(ram.Sticks, *stick)
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
		ram.TotalCapacity = "Unknown"
	}

	return ram
}

func scanGPULinux() GpuInfo {
	gpu := GpuInfo{Model: "Unknown GPU", VRAM: "Unknown"}

	// lspci: 查找 VGA/3D controller
	out, err := exec.Command("lspci").Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "VGA") || strings.Contains(line, "3D") || strings.Contains(line, "Display") {
				parts := strings.SplitN(line, ": ", 2)
				if len(parts) >= 2 {
					gpu.Model = strings.TrimSpace(parts[1])
					// 截断过长名称
					if len(gpu.Model) > 80 {
						gpu.Model = gpu.Model[:77] + "..."
					}
				}
				break
			}
		}
	}

	// /sys/class/drm 寻找显存信息 (card*/device/mem_info_vram_used 或类似)
	entries, _ := os.ReadDir("/sys/class/drm")
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "card") {
			vramPath := fmt.Sprintf("/sys/class/drm/%s/device/mem_info_vram_total", e.Name())
			if v, err := os.ReadFile(vramPath); err == nil {
				val := strings.TrimSpace(string(v))
				if vramBytes, err := strconv.ParseInt(val, 10, 64); err == nil && vramBytes > 0 {
					gpu.VRAM = fmt.Sprintf("%.0f GB", float64(vramBytes)/1_073_741_824)
				}
			}
		}
	}

	return gpu
}
