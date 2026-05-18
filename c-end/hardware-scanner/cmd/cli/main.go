package main

import (
	"encoding/json"
	"fmt"
	"os"

	scanner "github.com/ocmaster/hardware-scanner/scanner"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "scan":
		outputFormat := "table"
		if len(os.Args) > 2 && os.Args[2] == "--out" && len(os.Args) > 3 {
			outputFormat = os.Args[3]
		}
		runScan(outputFormat)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", cmd)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`OCMaster CLI

  ocmaster scan               扫描硬件，输出文本表格
  ocmaster scan --out json    扫描硬件，输出 JSON
  ocmaster help               帮助信息`)
}

func runScan(format string) {
	info := scanner.ScanAll()
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(info)
	default:
		printTable(info)
	}
}

func printTable(hw scanner.HardwareInfo) {
	fmt.Printf(`
╔══════════════════════════════════════╗
║       超频大师 OCMaster 扫描结果     ║
╠══════════════════════════════════════╣
║ CPU:  %-30s ║
║       核心: %-4d  线程: %-4d         ║
║       基频: %-25s ║
╠══════════════════════════════════════╣
║ 主板: %-30s ║
║       型号: %-25s ║
╠══════════════════════════════════════╣
║ 内存: %-30s ║
║       条数: %-4d  通道: %-4d         ║
╠══════════════════════════════════════╣
║ GPU:  %-30s ║
╚══════════════════════════════════════╝
`,
		hw.CPU.Model, hw.CPU.Cores, hw.CPU.Threads, hw.CPU.BaseFreq,
		hw.Motherboard.Brand, hw.Motherboard.Model,
		hw.RAM.TotalCapacity, hw.RAM.StickCount, hw.RAM.ChannelCount,
		hw.GPU.Model)
}
