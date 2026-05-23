package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/ocmaster/hardware-scanner/scanner"
)

var version = "dev"
var logger *log.Logger

func init() {
	logDir := filepath.Join(os.TempDir(), "OCMaster")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("ocmaster-%s.log", time.Now().Format("20060102-150405")))
	f, err := os.Create(logPath)
	if err == nil {
		logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
		logger.Println("OCMaster started — " + logPath)
	} else {
		logger = log.New(os.Stderr, "OCMaster: ", log.LstdFlags)
	}
}

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

	logger.Println("starting GUI mode")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	wnd := g.NewMasterWindow("超频大师 OCMaster", 1000, 720, 0)
	wnd.SetStyle(giuDarkTheme())
	wnd.SetTargetFPS(30)
	logger.Println("window created, entering main loop")
	wnd.Run(loop)
}

func printHelp() {
	fmt.Println("OCMaster")
	fmt.Println("  ocmaster              GUI")
	fmt.Println("  ocmaster scan          CLI scan")
	fmt.Println("  ocmaster scan --out json  JSON")
	fmt.Println("  ocmaster version       version")
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
