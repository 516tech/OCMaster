package scanner

import (
	"encoding/json"
	"runtime"
	"testing"
)

func TestScanAll_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin only test")
	}
	hw := ScanAll()

	t.Logf("CPU: %s (%dC/%dT) @ %s", hw.CPU.Model, hw.CPU.Cores, hw.CPU.Threads, hw.CPU.BaseFreq)
	t.Logf("Motherboard: %s %s", hw.Motherboard.Brand, hw.Motherboard.Model)
	t.Logf("RAM: %s, %d sticks, %d channels", hw.RAM.TotalCapacity, hw.RAM.StickCount, hw.RAM.ChannelCount)
	t.Logf("GPU: %s (%s)", hw.GPU.Model, hw.GPU.VRAM)

	if hw.CPU.Model == "" {
		t.Error("CPU model should not be empty")
	}
	if hw.CPU.Cores == 0 {
		t.Error("CPU cores should be > 0")
	}
	if hw.Motherboard.Model == "" {
		t.Error("motherboard model should not be empty")
	}
	if hw.RAM.TotalCapacity == "" || hw.RAM.TotalCapacity == "0 GB" {
		t.Error("RAM total capacity should not be empty")
	}
}

func TestScanAll_JSONOutput(t *testing.T) {
	hw := ScanAll()
	b, err := json.MarshalIndent(hw, "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	t.Logf("\n%s", string(b))

	var parsed HardwareInfo
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
}

func TestScanCPU_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin only test")
	}
	cpu := scanCPU()
	if cpu.Model == "" {
		t.Error("cpu model empty")
	}
	t.Logf("CPU model=%s cores=%d threads=%d freq=%s", cpu.Model, cpu.Cores, cpu.Threads, cpu.BaseFreq)
}

func TestScanGPU_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin only test")
	}
	gpu := scanGPU()
	t.Logf("GPU model=%s vram=%s", gpu.Model, gpu.VRAM)
}
