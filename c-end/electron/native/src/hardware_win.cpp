#include "hardware_common.h"
#include <windows.h>
#include <comdef.h>
#include <Wbemidl.h>
#pragma comment(lib, "wbemuuid.lib")

// Windows hardware detection via WMI + SMBus + SPD
// Compile: node-gyp rebuild --target=... --arch=x64

static std::string WmiQuery(const char* wql, const char* property) {
    // TODO: Initialize COM, connect to WMI, execute query
    // IWbemLocator → ConnectServer → ExecQuery → enumerate
    return "";
}

HardwareInfo ScanHardware() {
    HardwareInfo hw;

    // CPU via Win32_Processor
    hw.cpu.model = "Unknown CPU";
    hw.cpu.cores = 0;
    hw.cpu.threads = 0;
    hw.cpu.baseFreq = "0 MHz";

    // Motherboard via Win32_BaseBoard / Win32_BIOS
    hw.motherboard.brand = "Unknown";
    hw.motherboard.model = "Unknown";
    hw.motherboard.chipset = "Unknown";

    // RAM via Win32_PhysicalMemory + SPD (SMBus I2C)
    hw.ram.totalCapacity = "0 GB";
    hw.ram.stickCount = 0;
    hw.ram.channelCount = 0;

    // GPU via Win32_VideoController
    hw.gpu.model = "Unknown GPU";
    hw.gpu.vram = "0 GB";

    // PSU: SMBus read PMBus registers 0x96-0x99
    hw.psu.ratedWattage = "";
    hw.psu.source = "manual";

    // Cooler: undetectable programmatically
    hw.cooler.type = "air";
    hw.cooler.source = "manual";

    return hw;
}
