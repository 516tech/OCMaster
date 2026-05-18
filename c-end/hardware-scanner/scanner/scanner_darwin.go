package scanner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// runCmd 执行命令，带 3s 超时
func runCmd(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &out
	cmd.Run() // 忽略错误，超时时返回空
	return out.String()
}

func ScanAll() HardwareInfo {
	hw := HardwareInfo{}
	hw.CPU = scanCPU()
	hw.Motherboard = scanMotherboard()
	hw.RAM = scanRAM()
	hw.GPU = scanGPU()
	hw.PSU = PsuInfo{Source: "manual"}
	hw.Cooler = CoolerInfo{Type: "air", Source: "manual"}
	return hw
}

func sysctlByName(name string) string {
	s, err := unix.Sysctl(name)
	if err != nil {
		return ""
	}
	return s
}

func sysctlUint32(name string) uint32 {
	v, err := unix.SysctlUint32(name)
	if err != nil {
		return 0
	}
	return v
}

func sysctlUint64(name string) uint64 {
	v, err := unix.SysctlUint64(name)
	if err != nil {
		return 0
	}
	return v
}

func scanCPU() CpuInfo {
	cpu := CpuInfo{
		Model:       sysctlByName("machdep.cpu.brand_string"),
		BiosVersion: sysctlByName("machdep.cpu.brand_string"),
	}

	cpu.Cores = int(sysctlUint32("hw.physicalcpu"))
	cpu.Threads = int(sysctlUint32("hw.logicalcpu"))

	// Apple Silicon 无 cpufrequency，用品牌字符串中的频率信息
	// Intel Mac: hw.cpufrequency 返回 Hz
	if cpu.BaseFreq == "" {
		cpu.BaseFreq = "Apple Silicon"
	}

	return cpu
}

func scanMotherboard() MotherboardInfo {
	mb := MotherboardInfo{
		Brand: "Apple",
		Model: sysctlByName("hw.model"),
	}

	// 尝试获取 BIOS/固件版本
	bootArgs := runCmd("nvram", "boot-args")
	if len(bootArgs) > 0 {
		mb.BiosVersion = strings.TrimSpace(bootArgs)
	}
	if mb.BiosVersion == "" {
		mb.BiosVersion = "Apple UEFI"
	}

	return mb
}

func scanRAM() RamInfo {
	ram := RamInfo{}

	memSize := sysctlUint64("hw.memsize")
	if memSize > 0 {
		ram.TotalCapacity = fmt.Sprintf("%.0f GB", float64(memSize)/1_073_741_824)
	}

	// macOS 上通过 system_profiler 获取内存条详情
	out := runCmd("system_profiler", "SPMemoryDataType")
	if out != "" {
		ram = parseMacOSMemory(string(out), ram)
	}

	if ram.StickCount == 0 {
		ram.StickCount = 1
		ram.ChannelCount = 1
		ram.Sticks = []RamStick{{
			Capacity:  ram.TotalCapacity,
			Frequency: "Unknown",
			Timings:   "Unknown",
			DieType:   "Unknown",
		}}
	}

	// Apple Silicon 统一内存通常是双通道
	if runtime.GOARCH == "arm64" && ram.ChannelCount == 1 {
		ram.ChannelCount = 2
	}

	return ram
}

func parseMacOSMemory(output string, ram RamInfo) RamInfo {
	lines := strings.Split(output, "\n")
	var currentStick *RamStick
	stickCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "BANK") || strings.HasPrefix(trimmed, "DIMM") {
			if currentStick != nil {
				ram.Sticks = append(ram.Sticks, *currentStick)
			}
			currentStick = &RamStick{}
			stickCount++
		}

		if currentStick != nil {
			switch {
			case strings.Contains(trimmed, "Size:"):
				currentStick.Capacity = strings.TrimSpace(strings.TrimPrefix(trimmed, "Size:"))
			case strings.Contains(trimmed, "Type:"):
				currentStick.DieType = strings.TrimSpace(strings.TrimPrefix(trimmed, "Type:"))
			case strings.Contains(trimmed, "Speed:"):
				currentStick.Frequency = strings.TrimSpace(strings.TrimPrefix(trimmed, "Speed:"))
			}
		}
	}

	if currentStick != nil {
		ram.Sticks = append(ram.Sticks, *currentStick)
	}

	ram.StickCount = stickCount
	if stickCount >= 2 {
		ram.ChannelCount = 2
	} else if stickCount == 1 {
		ram.ChannelCount = 1
	}

	return ram
}

func scanGPU() GpuInfo {
	gpu := GpuInfo{Model: "Unknown GPU", VRAM: "Unknown"}

	out := runCmd("system_profiler", "SPDisplaysDataType")
	if out == "" {
		return gpu
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.Contains(trimmed, "Chipset Model:"):
			gpu.Model = strings.TrimSpace(strings.TrimPrefix(trimmed, "Chipset Model:"))
		case strings.Contains(trimmed, "VRAM"):
			gpu.VRAM = trimmed
		}
	}

	return gpu
}
