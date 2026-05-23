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

func initLog() {
	logDir := filepath.Join(os.TempDir(), "OCMaster")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("ocmaster-%s.log", time.Now().Format("20060102-150405")))
	f, err := os.Create(logPath)
	if err == nil {
		logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stderr, "OCMaster: ", log.LstdFlags)
	}
	logger.Println("=== OCMaster started ===")
	logger.Println("log: " + logPath)
}

func main() {
	initLog()
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("FATAL panic: %v", r)
			fmt.Fprintf(os.Stderr, "FATAL: %v\nSee log: %s\n", r, os.TempDir()+"/OCMaster/")
			os.Exit(1)
		}
	}()

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

	logger.Println("starting GUI — " + runtime.GOOS + "/" + runtime.GOARCH)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	wnd := g.NewMasterWindow("超频大师 OCMaster", 1100, 780, 0)
	wnd.SetStyle(ImHexDark())
	wnd.SetTargetFPS(30)

	// Delayed auto-scan (1s after startup)
	go func() {
		time.Sleep(1 * time.Second)
		if autoScan {
			logger.Println("auto-scan triggered")
			doScan()
		}
	}()

	logger.Println("entering main loop")
	wnd.Run(loop)
}

func printHelp() {
	fmt.Println("OCMaster")
	fmt.Println("  ocmaster              GUI")
	fmt.Println("  ocmaster scan          CLI")
	fmt.Println("  ocmaster scan --out json")
	fmt.Println("  ocmaster version")
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
