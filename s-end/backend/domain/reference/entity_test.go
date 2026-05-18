package reference

import "testing"

func TestCategory_Values(t *testing.T) {
	if CategoryCPU != "cpu" { t.Error("expected cpu") }
	if CategoryRAM != "ram" { t.Error("expected ram") }
	if CategoryCooler != "cooler" { t.Error("expected cooler") }
}

func TestReferenceData_Struct(t *testing.T) {
	r := ReferenceData{ID: 1, Category: CategoryCPU, ModelKey: "i7-13700K", Params: `{"freq":"5.5GHz"}`}
	if r.ID != 1 { t.Error("unexpected ID") }
	if r.ModelKey != "i7-13700K" { t.Error("unexpected ModelKey") }
	if r.Params == "" { t.Error("expected non-empty params") }
}
