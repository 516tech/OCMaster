package pdf

import (
	"testing"
)

func TestGenerateReport(t *testing.T) {
	chromePath := findChromePath()
	if chromePath == "" {
		t.Skip("chrome not found, skipping PDF test")
	}

	data := &ReportData{
		CpuFreq:      "5.5GHz",
		CpuVoltage:   "1.35V",
		CpuNotes:     "建议逐步提升频率",
		RamFreq:      "6400MHz",
		RamTimings:   "32-40-40-84",
		RamVoltage:   "1.40V",
		RamNotes:     "注意散热",
		TestTool:     "Prime95",
		TestDuration: "2小时",
		RiskWarning:  "超频有风险，请在专业人员指导下操作",
	}

	pdf, err := GenerateReport(data)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}
	if len(pdf) == 0 {
		t.Error("expected non-empty PDF")
	}
	// PDF magic bytes
	if len(pdf) < 4 || pdf[0] != '%' || pdf[1] != 'P' || pdf[2] != 'D' || pdf[3] != 'F' {
		t.Error("output does not start with %PDF magic")
	}
	t.Logf("PDF generated: %d bytes", len(pdf))
}

func TestFindChromePath(t *testing.T) {
	path := findChromePath()
	t.Logf("chrome path: %s", path)
}

func TestReportData_Empty(t *testing.T) {
	chromePath := findChromePath()
	if chromePath == "" {
		t.Skip("chrome not found")
	}
	pdf, err := GenerateReport(&ReportData{RiskWarning: "test"})
	if err != nil {
		t.Fatalf("empty report failed: %v", err)
	}
	if len(pdf) < 100 {
		t.Error("expected meaningful PDF content")
	}
}
