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
	hw        scanner.HardwareInfo
	scanned   bool
	scanning  bool
	scanErr   string
	apiURL    = "https://localhost/api/v1"
	language  = "zh-CN"
	autoScan  = true
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
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("PANIC in loop: %v", r)
		}
	}()

	g.SingleWindow().Layout(
		g.TabBar().ID("MainNav").Flags(g.TabBarFlagsReorderable).TabItems(
			g.TabItem("硬件扫描").Layout(scanTab()...),
			g.TabItem("上传分享").Layout(uploadTab()...),
			g.TabItem("设置").Layout(settingsTab()...),
			g.TabItem("关于").Layout(aboutTab()...),
		),
	)
}

func doScan() {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("PANIC in doScan: %v", r)
			scanErr = fmt.Sprintf("扫描崩溃: %v", r)
			scanning = false
		}
	}()

	logger.Println("scan: starting")
	scanning = true
	scanErr = ""

	result := scanner.ScanAll()
	logger.Printf("scan: CPU=%s cores=%d", result.CPU.Model, result.CPU.Cores)
	logger.Printf("scan: MB=%s %s", result.Motherboard.Brand, result.Motherboard.Model)
	logger.Printf("scan: RAM=%s sticks=%d", result.RAM.TotalCapacity, result.RAM.StickCount)
	logger.Printf("scan: GPU=%s vram=%s", result.GPU.Model, result.GPU.VRAM)

	if result.CPU.Cores == 0 && result.CPU.Model == "" {
		logger.Println("scan: WARNING — empty result")
		scanErr = "扫描失败: 未检测到硬件信息"
	} else {
		hw = result
		scanned = true
		logger.Println("scan: OK")
	}
	scanning = false
}

// ====== Scan Tab ======

func scanTab() g.Layout {
	var w g.Layout

	w = append(w, g.Dummy(0, 12))
	w = append(w, g.Row(
		g.Button("开始扫描").Size(130, 34).OnClick(func() {
			if !scanning {
				go doScan()
			}
		}),
		g.Button("导出 TXT").Size(100, 34).OnClick(func() {
			if scanned {
				exportTXT(hw)
			}
		}),
	))

	if scanning {
		w = append(w, g.Dummy(0, 10))
		w = append(w, g.Label("扫描中..."))
		w = append(w, g.ProgressBar(-1).Size(400, 18))
	}

	if scanErr != "" && !scanning {
		w = append(w, g.Dummy(0, 10))
		w = append(w, g.Label("错误: "+scanErr))
	}

	if scanned && !scanning {
		w = append(w, g.Dummy(0, 14))
		w = append(w, g.Separator())
		w = append(w, g.Dummy(0, 8))
		w = append(w, g.Label("硬件信息"))
		w = append(w, g.Dummy(0, 6))
		w = append(w, g.Table().
			Flags(g.TableFlagsBorders|g.TableFlagsResizable|g.TableFlagsRowBg|g.TableFlagsScrollY).
			Size(0, 420).
			Columns(
				g.TableColumn("组件").InnerWidthOrWeight(80),
				g.TableColumn("属性").InnerWidthOrWeight(130),
				g.TableColumn("值").InnerWidthOrWeight(490),
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
	w = append(w, g.Dummy(0, 12))
	w = append(w, g.Label("上传硬件信息并获取分享码"))
	w = append(w, g.Dummy(0, 8))
	w = append(w, g.Button("上传硬件信息").Size(140, 34).OnClick(func() {
		if !uploading {
			uploading = true
			uploadErr = ""
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Printf("PANIC in upload: %v", r)
						uploadErr = fmt.Sprintf("崩溃: %v", r)
						uploading = false
					}
				}()
				logger.Println("upload: starting")
				if !scanned {
					hw = scanner.ScanAll()
					scanned = true
				}
				shareCode = "123456" // TODO: real HTTP POST
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
		w = append(w, g.Dummy(0, 8))
		w = append(w, g.Label("错误: "+uploadErr))
	}

	if shareCode != "" && !uploading {
		w = append(w, g.Dummy(0, 14))
		w = append(w, g.Label("分享码 (7天有效):"))
		w = append(w, g.Dummy(0, 6))
		w = append(w, g.Style().SetFontSize(28).To(g.Label(shareCode)))
	}

	return w
}

// ====== Settings Tab ======

func settingsTab() g.Layout {
	return g.Layout{
		g.Dummy(0, 12),
		g.Label("API 地址"),
		g.Dummy(0, 2),
		g.InputText(&apiURL).Size(420),
		g.Dummy(0, 10),
		g.Label("语言"),
		g.Dummy(0, 2),
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
		g.Label("技术栈"),
		g.Label("  Go + giu (Dear ImGui)  —  跟 ImHex 同款 UI 框架"),
		g.Label("  WMI / sysfs / IOKit  —  三平台硬件扫描"),
		g.Label("  Go Backend + Vue 3  —  商家端管理平台"),
		g.Dummy(0, 10),
		g.Label("调试日志"),
		g.Label("  " + filepath.Join(os.TempDir(), "OCMaster")),
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
