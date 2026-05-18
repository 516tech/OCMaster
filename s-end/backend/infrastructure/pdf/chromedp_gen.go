package pdf

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const reportHTML = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>超频建议报告</title>
<style>body{font-family:Arial,sans-serif;padding:20px} h2{color:#333} table{border-collapse:collapse;width:100%} th,td{border:1px solid #ccc;padding:8px} th{background:#f5f5f5}</style></head>
<body>
<h2>超频建议报告</h2>
<h3>CPU 超频建议</h3><p>频率: {{.CpuFreq}} | 电压: {{.CpuVoltage}}</p><p>{{.CpuNotes}}</p>
<h3>内存超频建议</h3><p>频率: {{.RamFreq}} | 时序: {{.RamTimings}} | 电压: {{.RamVoltage}}</p><p>{{.RamNotes}}</p>
<h3>稳定性测试</h3><p>工具: {{.TestTool}} | 时长: {{.TestDuration}}</p>
<h3>风险提示</h3><p>{{.RiskWarning}}</p>
</body></html>`

type ReportData struct {
	CpuFreq, CpuVoltage, CpuNotes                string
	RamFreq, RamTimings, RamVoltage, RamNotes    string
	TestTool, TestDuration, RiskWarning          string
}

func findChromePath() string {
	paths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// fallback: $PATH lookup
	if p, err := exec.LookPath("google-chrome"); err == nil {
		return p
	}
	if p, err := exec.LookPath("chromium"); err == nil {
		return p
	}
	return ""
}

func GenerateReport(data *ReportData) ([]byte, error) {
	tmpl, _ := template.New("report").Parse(reportHTML)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(findChromePath()),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pdfBuf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, buf.String()).Do(ctx)
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	return pdfBuf, err
}
