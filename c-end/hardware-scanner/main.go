package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/ocmaster/hardware-scanner/scanner"
)

//export ScanAll
func ScanAll() *C.char {
	info := scanner.ScanAll()
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return C.CString(`{"error":"marshal failed"}`)
	}
	return C.CString(string(jsonBytes))
}

//export ScanCPU
func ScanCPU() *C.char {
	info := scanner.ScanAll()
	jsonBytes, _ := json.Marshal(info.CPU)
	return C.CString(string(jsonBytes))
}

//export ScanRAM
func ScanRAM() *C.char {
	info := scanner.ScanAll()
	jsonBytes, _ := json.Marshal(info.RAM)
	return C.CString(string(jsonBytes))
}

//export ScanGPU
func ScanGPU() *C.char {
	info := scanner.ScanAll()
	jsonBytes, _ := json.Marshal(info.GPU)
	return C.CString(string(jsonBytes))
}

//export ScanMotherboard
func ScanMotherboard() *C.char {
	info := scanner.ScanAll()
	jsonBytes, _ := json.Marshal(info.Motherboard)
	return C.CString(string(jsonBytes))
}

//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {} // c-shared 需要 main
