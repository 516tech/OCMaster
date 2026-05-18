#ifndef HARDWARE_COMMON_H
#define HARDWARE_COMMON_H

#include <string>
#include <vector>

struct CpuInfo {
    std::string model;
    int cores = 0;
    int threads = 0;
    std::string baseFreq;
    std::string biosVersion;
};

struct MotherboardInfo {
    std::string brand;
    std::string model;
    std::string chipset;
    std::string biosVersion;
};

struct RamStick {
    std::string capacity;
    std::string frequency;
    std::string timings;
    std::string dieType;
};

struct RamInfo {
    std::string totalCapacity;
    int stickCount = 0;
    int channelCount = 0;
    std::vector<RamStick> sticks;
};

struct GpuInfo {
    std::string model;
    std::string vram;
};

struct PsuInfo {
    std::string ratedWattage;
    std::string source; // "smbus" or "manual"
};

struct CoolerInfo {
    std::string type; // "air", "aio", "custom_loop"
    std::string source; // "auto" or "manual"
};

struct HardwareInfo {
    CpuInfo cpu;
    MotherboardInfo motherboard;
    RamInfo ram;
    GpuInfo gpu;
    PsuInfo psu;
    CoolerInfo cooler;
};

HardwareInfo ScanHardware();

#endif
