package suggestion

import "testing"

func TestSuggestion_Struct(t *testing.T) {
	s := Suggestion{
		ID: 1, MerchantID: 10, HardwareID: 100,
		CpuSuggestion: &CpuSuggestion{Frequency: "5.5GHz", Voltage: "1.35V", Notes: "逐步提升"},
		RamSuggestion: &RamSuggestion{Frequency: "6400MHz", Timings: "32-40-40-84", Voltage: "1.40V"},
		StabilityTest: &StabilityTest{Tool: "Prime95", Duration: "2小时"},
		RiskWarning: "超频有风险",
	}
	if s.MerchantID != 10 { t.Error("unexpected MerchantID") }
	if s.CpuSuggestion == nil { t.Error("expected CpuSuggestion") }
	if s.CpuSuggestion.Frequency != "5.5GHz" { t.Error("unexpected frequency") }
}

func TestCpuSuggestion_Defaults(t *testing.T) {
	c := CpuSuggestion{}
	if c.Frequency != "" { t.Error("expected empty frequency") }
}

func TestRamSuggestion_Defaults(t *testing.T) {
	r := RamSuggestion{}
	if r.Timings != "" { t.Error("expected empty timings") }
}

func TestStabilityTest_Defaults(t *testing.T) {
	s := StabilityTest{}
	if s.Tool != "" { t.Error("expected empty tool") }
}
