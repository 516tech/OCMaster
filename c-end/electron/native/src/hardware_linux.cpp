#include "hardware_common.h"
#include <fstream>
#include <sstream>
#include <cstring>

// Linux hardware detection via sysfs + dmidecode + i2c-tools
// Requires: dmidecode installed, i2c-tools for SPD

static std::string ReadFile(const char* path) {
    std::ifstream f(path);
    if (!f.is_open()) return "";
    std::stringstream ss; ss << f.rdbuf(); return ss.str();
}

static std::string ExecCmd(const char* cmd) {
    // Popen shell command
    return "";
}

HardwareInfo ScanHardware() {
    HardwareInfo hw;

    // CPU via /proc/cpuinfo
    hw.cpu.model = "Unknown CPU";
    hw.cpu.cores = 0;
    hw.cpu.threads = 0;

    // Motherboard via dmidecode -t baseboard
    hw.motherboard.brand = "Unknown";
    hw.motherboard.model = "Unknown";

    // RAM via /sys/devices/system/edac/mc/ + i2c SPD
    hw.ram.totalCapacity = "0 GB";
    hw.ram.stickCount = 0;

    // GPU via /sys/class/drm/ or lspci
    hw.gpu.model = "Unknown GPU";
    hw.gpu.vram = "0 GB";

    hw.psu.ratedWattage = "";
    hw.psu.source = "manual";
    hw.cooler.type = "air";
    hw.cooler.source = "manual";

    return hw;
}
