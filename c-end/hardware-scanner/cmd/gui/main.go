package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	g "github.com/AllenDang/giu"
	"github.com/ocmaster/hardware-scanner/scanner"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "scan":
			runCLIScan()
			return
		case "version", "--version", "-v":
			fmt.Printf("ocmaster version %s\n", version)
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	wnd := g.NewMasterWindow("超频大师 OCMaster", 1000, 720, 0)
	wnd.SetStyle(giuDarkTheme())
	wnd.SetTargetFPS(30)
	wnd.Run(loop)
}

func printHelp() {
	fmt.Println("OCMaster — 超频大师")
	fmt.Println("  ocmaster              启动 GUI")
	fmt.Println("  ocmaster scan          命令行扫描")
	fmt.Println("  ocmaster scan --out json  JSON 输出")
	fmt.Println("  ocmaster version       版本信息")
}

func runCLIScan() {
	outFormat := "table"
	for i, a := range os.Args {
		if a == "--out" && i+1 < len(os.Args) {
			outFormat = os.Args[i+1]
		}
	}
	info := scanner.ScanAll()
	if outFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(info)
	} else {
		fmt.Printf("CPU: %s (%dC/%dT) @ %s\n", info.CPU.Model, info.CPU.Cores, info.CPU.Threads, info.CPU.BaseFreq)
		fmt.Printf("Motherboard: %s %s (BIOS %s)\n", info.Motherboard.Brand, info.Motherboard.Model, info.Motherboard.BiosVersion)
		fmt.Printf("RAM: %s | %d sticks | %d channels\n", info.RAM.TotalCapacity, info.RAM.StickCount, info.RAM.ChannelCount)
		fmt.Printf("GPU: %s (%s)\n", info.GPU.Model, info.GPU.VRAM)
	}
}
