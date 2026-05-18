package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	scanner "github.com/ocmaster/hardware-scanner/scanner"
)

var version = "dev"

var lang = "zh-CN"

var msgs = map[string]map[string]string{
	"help_title":   {"zh-CN": "OCMaster CLI 超频大师", "en-US": "OCMaster CLI"},
	"help_scan":    {"zh-CN": "扫描硬件，输出文本表格", "en-US": "Scan hardware, output as table"},
	"help_json":    {"zh-CN": "扫描硬件，输出 JSON", "en-US": "Scan hardware, output as JSON"},
	"help_help":    {"zh-CN": "帮助信息", "en-US": "Help information"},
	"unknown_cmd":  {"zh-CN": "未知命令", "en-US": "Unknown command"},
	"scan_result":  {"zh-CN": "超频大师 OCMaster 扫描结果", "en-US": "OCMaster Scan Result"},
	"cpu":          {"zh-CN": "CPU", "en-US": "CPU"},
	"cores":        {"zh-CN": "核心", "en-US": "Cores"},
	"threads":      {"zh-CN": "线程", "en-US": "Threads"},
	"base_freq":    {"zh-CN": "基频", "en-US": "Base Freq"},
	"motherboard":  {"zh-CN": "主板", "en-US": "Motherboard"},
	"model":        {"zh-CN": "型号", "en-US": "Model"},
	"ram":          {"zh-CN": "内存", "en-US": "RAM"},
	"sticks":       {"zh-CN": "条数", "en-US": "Sticks"},
	"channels":     {"zh-CN": "通道", "en-US": "Channels"},
	"gpu":          {"zh-CN": "GPU", "en-US": "GPU"},
	"scan_failed":  {"zh-CN": "扫描失败", "en-US": "Scan failed"},
}

func t(key string) string {
	if m, ok := msgs[key]; ok {
		if v, ok := m[lang]; ok {
			return v
		}
	}
	return key
}

// parseFlag 解析 --flag value 形式的参数
func parseFlag(args []string, flag string) (string, bool) {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

func main() {
	// 全局 --lang 参数
	if v, ok := parseFlag(os.Args, "--lang"); ok {
		lang = v
	}
	// 移除 --lang 及其值，简化后续解析
	var cleanArgs []string
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--lang" {
			i++ // skip value
			continue
		}
		cleanArgs = append(cleanArgs, os.Args[i])
	}

	if len(cleanArgs) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := cleanArgs[1]
	switch cmd {
	case "scan":
		outputFormat := "table"
		if v, ok := parseFlag(cleanArgs, "--out"); ok {
			outputFormat = v
		}
		runScan(outputFormat)
	case "version", "--version", "-v":
		fmt.Printf("ocmaster version %s\n", version)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "%s: %s\n", t("unknown_cmd"), cmd)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`%s

  ocmaster scan                  %s
  ocmaster scan --out json       %s
  ocmaster scan --lang en-US     %s
  ocmaster version
  ocmaster help                  %s
`, t("help_title"), t("help_scan"), t("help_json"),
		strings.Replace(t("help_scan"), "文本表格", "English", 1),
		t("help_help"))
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
	if lang == "zh-CN" {
		fmt.Printf(`
╔══════════════════════════════════════╗
║       %-30s ║
╠══════════════════════════════════════╣
║ %s: %-32s ║
║      %s: %-4d  %s: %-4d          ║
║      %s: %-26s ║
╠══════════════════════════════════════╣
║ %s: %-32s ║
║      %s: %-24s ║
╠══════════════════════════════════════╣
║ %s: %-32s ║
║      %s: %-4d  %s: %-4d          ║
╠══════════════════════════════════════╣
║ %s: %-32s ║
╚══════════════════════════════════════╝
`,
			t("scan_result"),
			t("cpu"), hw.CPU.Model,
			t("cores"), hw.CPU.Cores, t("threads"), hw.CPU.Threads,
			t("base_freq"), hw.CPU.BaseFreq,
			t("motherboard"), hw.Motherboard.Brand,
			t("model"), hw.Motherboard.Model,
			t("ram"), hw.RAM.TotalCapacity,
			t("sticks"), hw.RAM.StickCount, t("channels"), hw.RAM.ChannelCount,
			t("gpu"), hw.GPU.Model)
	} else {
		fmt.Printf(`
+======================================+
|       %-30s |
+======================================+
| %s: %-32s |
|      %s: %-4d  %s: %-4d          |
|      %s: %-26s |
+--------------------------------------+
| %s: %-32s |
|      %s: %-24s |
+--------------------------------------+
| %s: %-32s |
|      %s: %-4d  %s: %-4d          |
+--------------------------------------+
| %s: %-32s |
+======================================+
`,
			t("scan_result"),
			t("cpu"), hw.CPU.Model,
			t("cores"), hw.CPU.Cores, t("threads"), hw.CPU.Threads,
			t("base_freq"), hw.CPU.BaseFreq,
			t("motherboard"), hw.Motherboard.Brand,
			t("model"), hw.Motherboard.Model,
			t("ram"), hw.RAM.TotalCapacity,
			t("sticks"), hw.RAM.StickCount, t("channels"), hw.RAM.ChannelCount,
			t("gpu"), hw.GPU.Model)
	}
}
