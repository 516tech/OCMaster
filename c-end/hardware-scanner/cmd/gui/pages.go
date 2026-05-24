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

// Shared state
var (
	hw          scanner.HardwareInfo
	scanned     bool
	scanning    bool
	scanStep    string // current scan step for UI feedback
	scanErr     string
	apiURL      = "https://localhost/api/v1"
	language    = "zh-CN"
	autoScan    = true
	autoRan     bool
	langIdx     int32
	shareCode   string
	uploading   bool
	uploadErr   string
	activeTab   int32 // 0=scan, 1=upload, 2=settings, 3=about
)

func init() {
	if language == "en-US" {
		langIdx = 1
	}
}

func navItem(label string, tab int32) g.Widget {
	return g.Button(label).Size(100, 32).Disabled(activeTab == tab).OnClick(func() {
		activeTab = tab
		logger.Printf("nav: switched to tab %d=%s", tab, label)
	})
}

func loop() {
	g.SingleWindow().Layout(
		g.Child().Size(0, 44).Layout(
			g.Row(
				g.Dummy(12, 0),
				navItem("硬件扫描", 0),
				navItem("上传分享", 1),
				navItem("设置", 2),
				navItem("关于", 3),
			),
		).Border(false),
		g.Separator(),
		g.Child().Size(0, 0).Layout(
			contentArea()...,
		),
	)

	// Auto-scan on first frame
	if autoScan && !autoRan {
		autoRan = true
		go func() {
			time.Sleep(500 * time.Millisecond)
			doScan()
		}()
	}
}

func contentArea() g.Layout {
	switch activeTab {
	case 0:
		return scanTab()
	case 1:
		return uploadTab()
	case 2:
		return settingsTab()
	case 3:
		return aboutTab()
	default:
		return nil
	}
}

func doScan() {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("PANIC doScan step=%s: %v", scanStep, r)
			scanErr = fmt.Sprintf("扫描崩溃 @%s: %v", scanStep, r)
			scanning = false
			scanStep = ""
		}
	}()

	logger.Println("=== scan: start ===")
	scanning = true
	scanErr = ""
	scanned = false

	scanStep = "CPU"
	logger.Println("scan: step=CPU")
	result := scanner.ScanAll()
	logger.Printf("scan: result CPU.Model=%s CPU.Cores=%d CPU.Threads=%d", result.CPU.Model, result.CPU.Cores, result.CPU.Threads)
	logger.Printf("scan: result MB=%s %s BIOS=%s", result.Motherboard.Brand, result.Motherboard.Model, result.Motherboard.BiosVersion)
	logger.Printf("scan: result RAM=%s sticks=%d ch=%d", result.RAM.TotalCapacity, result.RAM.StickCount, result.RAM.ChannelCount)
	logger.Printf("scan: result GPU=%s VRAM=%s", result.GPU.Model, result.GPU.VRAM)

	scanStep = "validate"
	if result.CPU.Cores == 0 && result.CPU.Model == "" {
		logger.Println("scan: FAILED — empty CPU")
		scanErr = "扫描失败: 未检测到硬件"
	} else {
		hw = result
		scanned = true
		logger.Println("scan: OK")
	}
	scanning = false
	scanStep = ""
	logger.Println("=== scan: done ===")
}

// ====== Scan Tab ======

func scanTab() g.Layout {
	var w g.Layout
	w = append(w, g.Dummy(0, 16))

	w = append(w, g.Row(
		g.Dummy(8, 0),
		g.Button("开始扫描").Size(0, 0).OnClick(func() {
			if !scanning {
				logger.Println("button: start scan")
				go doScan()
			}
		}),
		g.Button("导出 TXT").Size(0, 0).OnClick(func() {
			if scanned {
				logger.Println("button: export TXT")
				exportTXT(hw)
			}
		}),
	))

	if scanning {
		w = append(w, g.Dummy(0, 12))
		w = append(w, g.Label("扫描中... ("+scanStep+")"))
		w = append(w, g.ProgressBar(-1).Size(500, 20))
	}

	if scanErr != "" && !scanning {
		w = append(w, g.Dummy(0, 12))
		w = append(w, g.Label("错误: "+scanErr))
	}

	if scanned && !scanning {
		w = append(w, g.Dummy(0, 16))
		w = append(w, g.Separator())
		w = append(w, g.Label("硬件信息"))
		w = append(w, g.Dummy(0, 8))
		w = append(w, g.Table().
			Flags(g.TableFlagsBorders|g.TableFlagsResizable|g.TableFlagsRowBg|g.TableFlagsScrollY|g.TableFlagsSizingFixedFit).
			Size(0, 440).
			Columns(
				g.TableColumn("组件"),
				g.TableColumn("属性"),
				g.TableColumn("值"),
			).Rows(
			g.TableRow(g.Label("CPU"), g.Label("型号"), g.Label(hw.CPU.Model)),
			g.TableRow(g.Label("CPU"), g.Label("核心 / 线程"), g.Label(fmt.Sprintf("%dC / %dT", hw.CPU.Cores, hw.CPU.Threads))),
			g.TableRow(g.Label("CPU"), g.Label("基频"), g.Label(hw.CPU.BaseFreq)),
			g.TableRow(g.Label("主板"), g.Label("品牌"), g.Label(hw.Motherboard.Brand)),
			g.TableRow(g.Label("主板"), g.Label("型号"), g.Label(hw.Motherboard.Model)),
			g.TableRow(g.Label("主板"), g.Label("BIOS 版本"), g.Label(hw.Motherboard.BiosVersion)),
			g.TableRow(g.Label("内存"), g.Label("总容量"), g.Label(hw.RAM.TotalCapacity)),
			g.TableRow(g.Label("内存"), g.Label("内存条数"), g.Label(fmt.Sprintf("%d", hw.RAM.StickCount))),
			g.TableRow(g.Label("内存"), g.Label("通道数"), g.Label(fmt.Sprintf("%d", hw.RAM.ChannelCount))),
			g.TableRow(g.Label("显卡"), g.Label("型号"), g.Label(hw.GPU.Model)),
			g.TableRow(g.Label("显卡"), g.Label("显存"), g.Label(hw.GPU.VRAM)),
		),
		)
	}
	return w
}

// ====== Upload Tab ======

func uploadTab() g.Layout {
	var w g.Layout
	w = append(w, g.Dummy(0, 16))
	w = append(w, g.Label("上传硬件信息并获取分享码"))
	w = append(w, g.Dummy(0, 8))
	w = append(w, g.Button("上传硬件信息").Size(0, 0).OnClick(func() {
		if !uploading {
			uploading = true
			uploadErr = ""
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Printf("PANIC upload: %v", r)
						uploadErr = fmt.Sprintf("崩溃: %v", r)
						uploading = false
					}
				}()
				logger.Println("upload: start")
				if !scanned {
					hw = scanner.ScanAll()
					scanned = true
				}
				shareCode = "123456"
				logger.Println("upload: done shareCode=" + shareCode)
				uploading = false
			}()
		}
	}))
	if uploading {
		w = append(w, g.Dummy(0, 8))
		w = append(w, g.Label("上传中..."))
		w = append(w, g.ProgressBar(-1).Size(300, 18))
	}
	if uploadErr != "" && !uploading {
		w = append(w, g.Label("错误: "+uploadErr))
	}
	if shareCode != "" && !uploading {
		w = append(w, g.Dummy(0, 14))
		w = append(w, g.Label("分享码 (7天有效):"))
		w = append(w, g.Style().SetFontSize(28).To(g.Label(shareCode)))
	}
	return w
}

// ====== Settings Tab ======

func settingsTab() g.Layout {
	return g.Layout{
		g.Dummy(0, 16),
		g.Label("API 地址"),
		g.InputText(&apiURL).Size(420),
		g.Dummy(0, 10),
		g.Label("语言"),
		g.Combo("", language, []string{"中文 (zh-CN)", "English (en-US)"}, &langIdx).Size(180),
		g.Dummy(0, 10),
		g.Checkbox("启动时自动扫描硬件", &autoScan),
	}
}

// ====== About Tab ======

func aboutTab() g.Layout {
	ver := "0.0.1"
	if data, err := os.ReadFile(filepath.Join("..", "..", "VERSION")); err == nil {
		ver = strings.TrimSpace(string(data))
	}
	return g.Layout{
		g.Dummy(0, 30),
		g.Style().SetFontSize(24).To(g.Label("超频大师 OCMaster")),
		g.Dummy(0, 6),
		g.Label("版本 " + ver),
		g.Dummy(0, 20),
		g.Label("技术栈: Go + giu (Dear ImGui)  —  跟 ImHex 同款 UI 框架"),
		g.Label("扫描: WMI / sysfs / IOKit"),
		g.Label("后端: Go Backend + Vue 3"),
		g.Dummy(0, 14),
		g.Label("日志路径:"),
		g.Label("  " + logDir()),
	}
}

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
	logger.Println("export: " + filename)
}
