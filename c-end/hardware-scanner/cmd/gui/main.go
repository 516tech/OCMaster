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
	var logDir string
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		logDir = filepath.Join(home, "OCMaster", "logs")
	} else {
		logDir = filepath.Join(os.TempDir(), "OCMaster")
	}
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("ocmaster-%s.log", time.Now().Format("20060102-150405")))
	f, err := os.Create(logPath)
	if err == nil {
		logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
		fmt.Printf("Log: %s\n", logPath)
	} else {
		logger = log.New(os.Stderr, "OCMaster: ", log.LstdFlags)
		fmt.Fprintf(os.Stderr, "Cannot create log: %v\n", err)
	}
	logger.Printf("=== OCMaster v%s === %s/%s ===", version, runtime.GOOS, runtime.GOARCH)
	logger.Printf("log file: %s", logPath)
}

func main() {
	initLog()
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("FATAL main: %v", r)
			fmt.Fprintf(os.Stderr, "FATAL: %v\nLog: %s\n", r, logDir())
		}
	}()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "scan":
			runCLIScan()
			return
		case "version", "--version", "-v":
			fmt.Printf("ocmaster version %s\n", version)
			return
		case "help", "--help", "-h":
			fmt.Println("OCMaster  GUI | ocmaster scan  CLI | ocmaster scan --out json")
			return
		}
	}

	logger.Println("GUI mode starting")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	wnd := g.NewMasterWindow("超频大师 OCMaster", 1100, 780, 0)
	wnd.SetStyle(ImHexDark())
	wnd.SetTargetFPS(30)
	logger.Println("window created, entering loop")
	wnd.Run(loop)
}

func logDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "OCMaster", "logs")
	}
	return filepath.Join(os.TempDir(), "OCMaster")
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
