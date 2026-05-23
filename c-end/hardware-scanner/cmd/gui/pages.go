package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/ocmaster/hardware-scanner/scanner"
)

// Shared state — giu loop picks up changes every frame
var (
	hw        scanner.HardwareInfo
	scanned   bool
	scanning  bool
	scanErr   string
	apiURL    = "https://localhost/api/v1"
	language  = "zh-CN"
	autoScan  = true
	autoRan   bool
	langIdx   int32
	shareCode string
	uploading bool
	uploadErr string
)

func init() {
	if language == "en-US" {
		langIdx = 1
	}
}

func loop() {
	// Auto-scan on first frame
	if autoScan && !autoRan {
		autoRan = true
		go doScan()
	}

	g.SingleWindow().Layout(
		g.TabBar().ID("MainTabs").TabItems(
			g.TabItem("硬件扫描").Layout(scanTab()...),
			g.TabItem("上传分享").Layout(uploadTab()...),
			g.TabItem("设置").Layout(settingsTab()...),
			g.TabItem("关于").Layout(aboutTab()...),
		),
	)
}

func doScan() {
	logger.Println("scan started")
	scanning = true
	scanErr = ""
	result := scanner.ScanAll()
	logger.Printf("scan result: CPU=%s Motherboard=%s RAM=%s GPU=%s",
		result.CPU.Model, result.Motherboard.Model, result.RAM.TotalCapacity, result.GPU.Model)
	if result.CPU.Model == "" {
		logger.Println("WARNING: cpu.model empty — scan may have failed")
		scanErr = "扫描失败: 未检测到硬件信息"
	} else {
		hw = result
		scanned = true
		logger.Println("scan completed successfully")
	}
	scanning = false
}

// ====== Scan Tab ======

func scanTab() g.Layout {
	var layout g.Layout

	layout = append(layout,
		g.Dummy(0, 16),
		g.Row(
			g.Button("开始扫描").Size(120, 36).OnClick(func() {
				if !scanning {
					go doScan()
				}
			}),
			g.Button("导出 TXT").Size(100, 36).OnClick(func() {
				if scanned {
					exportTXT(hw)
					logger.Println("TXT exported")
				}
			}),
		),
	)

	if scanning {
		layout = append(layout,
			g.Dummy(0, 8),
			g.Label("扫描中..."),
			g.ProgressBar(-1).Size(300, 20),
		)
	}

	if scanErr != "" && !scanning {
		layout = append(layout,
			g.Dummy(0, 8),
			g.Label("错误: "+scanErr),
		)
	}

	if scanned && !scanning {
		layout = append(layout, g.Dummy(0, 12), g.Separator())
		layout = append(layout, hardwareTable()...)
	}

	return layout
}

func hardwareTable() g.Layout {
	return g.Layout{
		g.Dummy(0, 4),
		g.Table().
			Flags(g.TableFlagsBorders|g.TableFlagsResizable|g.TableFlagsScrollY).
			Size(0, 400).
			Columns(
				g.TableColumn("组件").InnerWidthOrWeight(100),
				g.TableColumn("属性").InnerWidthOrWeight(140),
				g.TableColumn("值").InnerWidthOrWeight(460),
			).Rows(
			g.TableRow(g.Label("CPU"), g.Label("型号"), g.Label(hw.CPU.Model)),
			g.TableRow(g.Label("CPU"), g.Label("核心/线程"), g.Label(fmt.Sprintf("%dC / %dT", hw.CPU.Cores, hw.CPU.Threads))),
			g.TableRow(g.Label("CPU"), g.Label("基频"), g.Label(hw.CPU.BaseFreq)),
			g.TableRow(g.Label("主板"), g.Label("品牌"), g.Label(hw.Motherboard.Brand)),
			g.TableRow(g.Label("主板"), g.Label("型号"), g.Label(hw.Motherboard.Model)),
			g.TableRow(g.Label("主板"), g.Label("BIOS"), g.Label(hw.Motherboard.BiosVersion)),
			g.TableRow(g.Label("内存"), g.Label("总容量"), g.Label(hw.RAM.TotalCapacity)),
			g.TableRow(g.Label("内存"), g.Label("条数"), g.Label(fmt.Sprintf("%d", hw.RAM.StickCount))),
			g.TableRow(g.Label("内存"), g.Label("通道"), g.Label(fmt.Sprintf("%d", hw.RAM.ChannelCount))),
			g.TableRow(g.Label("显卡"), g.Label("型号"), g.Label(hw.GPU.Model)),
			g.TableRow(g.Label("显卡"), g.Label("显存"), g.Label(hw.GPU.VRAM)),
		),
	}
}

// ====== Upload Tab ======

func uploadTab() g.Layout {
	var layout g.Layout

	layout = append(layout,
		g.Dummy(0, 16),
		g.Label("上传硬件信息并获取分享码"),
		g.Dummy(0, 8),
		g.Button("上传硬件信息").Size(140, 36).OnClick(func() {
			if !uploading {
				uploading = true
				uploadErr = ""
				go func() {
					logger.Println("upload started")
					if !scanned {
						hw = scanner.ScanAll()
						scanned = true
					}
					shareCode = "123456" // TODO: real HTTP POST
					logger.Println("upload done, shareCode=" + shareCode)
					uploading = false
				}()
			}
		}),
	)

	if uploading {
		layout = append(layout,
			g.Dummy(0, 8),
			g.Label("上传中..."),
			g.ProgressBar(-1).Size(300, 20),
		)
	}

	if uploadErr != "" {
		layout = append(layout, g.Label("错误: "+uploadErr))
	}

	if shareCode != "" && !uploading {
		layout = append(layout,
			g.Dummy(0, 12),
			g.Label("分享码 (7天有效):"),
			g.Style().SetFontSize(28).To(g.Label(shareCode)),
		)
	}

	return layout
}

// ====== Settings Tab ======

func settingsTab() g.Layout {
	return g.Layout{
		g.Dummy(0, 16),
		g.Label("API 地址"),
		g.InputText(&apiURL).Size(400),
		g.Dummy(0, 8),
		g.Label("语言"),
		g.Combo("", language, []string{"zh-CN", "en-US"}, &langIdx).Size(120),
		g.Dummy(0, 8),
		g.Checkbox("启动时自动扫描", &autoScan),
	}
}

// ====== About Tab ======

func aboutTab() g.Layout {
	ver := "0.0.1"
	if data, err := os.ReadFile(filepath.Join("..", "..", "VERSION")); err == nil {
		ver = strings.TrimSpace(string(data))
	}
	return g.Layout{
		g.Dummy(0, 40),
		g.Style().SetFontSize(24).To(g.Label("超频大师 OCMaster")),
		g.Dummy(0, 8),
		g.Label("版本: " + ver),
		g.Dummy(0, 20),
		g.Label("技术栈: Go + giu (Dear ImGui)"),
		g.Label("  · GUI — 跨平台 (Win/Mac/Linux)"),
		g.Label("  · CLI — 硬件扫描 (WMI/sysfs/IOKit)"),
		g.Label("  · Backend + 商家端 — Go + Vue 3"),
		g.Dummy(0, 8),
		g.Label("日志: " + os.TempDir() + "/OCMaster/"),
	}
}

// ====== Helpers ======

func exportTXT(hw scanner.HardwareInfo) {
	lines := []string{
		"OCMaster",
		"=======",
		fmt.Sprintf("CPU: %s | %dC/%dT | %s", hw.CPU.Model, hw.CPU.Cores, hw.CPU.Threads, hw.CPU.BaseFreq),
		fmt.Sprintf("Motherboard: %s %s | BIOS %s", hw.Motherboard.Brand, hw.Motherboard.Model, hw.Motherboard.BiosVersion),
		fmt.Sprintf("RAM: %s | %d sticks | %d channels", hw.RAM.TotalCapacity, hw.RAM.StickCount, hw.RAM.ChannelCount),
		fmt.Sprintf("GPU: %s | %s", hw.GPU.Model, hw.GPU.VRAM),
	}
	filename := fmt.Sprintf("ocmaster-hardware-%s.txt", time.Now().Format("20060102-150405"))
	os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	logger.Println("exported TXT to " + filename)
}
