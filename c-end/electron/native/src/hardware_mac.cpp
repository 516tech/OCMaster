#include "hardware_common.h"
#include <sys/types.h>
#include <sys/sysctl.h>

// macOS hardware detection via IOKit + sysctl
// Compile: clang++ -std=c++17 -framework IOKit -framework Foundation

static std::string SysctlByName(const char* name) {
    char buf[256];
    size_t len = sizeof(buf);
    if (sysctlbyname(name, buf, &len, NULL, 0) == 0)
        return std::string(buf);
    return "";
}

HardwareInfo ScanHardware() {
    HardwareInfo hw;

    // CPU via sysctl machdep.cpu
    hw.cpu.model = SysctlByName("machdep.cpu.brand_string");
    hw.cpu.cores = 0;
    hw.cpu.threads = 0;

    // Motherboard: IOKit IOServiceGetMatchingService("IOPlatformExpertDevice")
    hw.motherboard.brand = "Apple";
    hw.motherboard.model = SysctlByName("hw.model");

    // RAM via sysctl hw.memsize
    hw.ram.totalCapacity = "0 GB";
    hw.ram.stickCount = 0;

    // GPU: IOKit matching "IOAccelerator" or "AGXAccelerator"
    hw.gpu.model = "Unknown GPU";
    hw.gpu.vram = "0 GB";

    hw.psu.ratedWattage = "";
    hw.psu.source = "manual";
    hw.cooler.type = "air";
    hw.cooler.source = "manual";

    return hw;
}
