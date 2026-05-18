#include <napi.h>
#include "hardware_common.h"

Napi::Object HardwareInfoToNapi(const Napi::Env& env, const HardwareInfo& info) {
    Napi::Object obj = Napi::Object::New(env);

    Napi::Object cpu = Napi::Object::New(env);
    cpu.Set("model", info.cpu.model);
    cpu.Set("cores", info.cpu.cores);
    cpu.Set("threads", info.cpu.threads);
    cpu.Set("baseFreq", info.cpu.baseFreq);
    cpu.Set("biosVersion", info.cpu.biosVersion);
    obj.Set("cpu", cpu);

    Napi::Object mb = Napi::Object::New(env);
    mb.Set("brand", info.motherboard.brand);
    mb.Set("model", info.motherboard.model);
    mb.Set("chipset", info.motherboard.chipset);
    mb.Set("biosVersion", info.motherboard.biosVersion);
    obj.Set("motherboard", mb);

    Napi::Object ram = Napi::Object::New(env);
    ram.Set("totalCapacity", info.ram.totalCapacity);
    ram.Set("stickCount", info.ram.stickCount);
    ram.Set("channelCount", info.ram.channelCount);
    Napi::Array sticks = Napi::Array::New(env, info.ram.stickCount);
    for (int i = 0; i < info.ram.stickCount; i++) {
        Napi::Object stick = Napi::Object::New(env);
        stick.Set("capacity", info.ram.sticks[i].capacity);
        stick.Set("frequency", info.ram.sticks[i].frequency);
        stick.Set("timings", info.ram.sticks[i].timings);
        stick.Set("dieType", info.ram.sticks[i].dieType);
        sticks.Set(i, stick);
    }
    ram.Set("sticks", sticks);
    obj.Set("ram", ram);

    Napi::Object gpu = Napi::Object::New(env);
    gpu.Set("model", info.gpu.model);
    gpu.Set("vram", info.gpu.vram);
    obj.Set("gpu", gpu);

    Napi::Object psu = Napi::Object::New(env);
    psu.Set("ratedWattage", info.psu.ratedWattage);
    psu.Set("source", info.psu.source);
    obj.Set("psu", psu);

    Napi::Object cooler = Napi::Object::New(env);
    cooler.Set("type", info.cooler.type);
    cooler.Set("source", info.cooler.source);
    obj.Set("cooler", cooler);

    return obj;
}

Napi::Value Scan(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();
    HardwareInfo hw = ScanHardware();
    return HardwareInfoToNapi(env, hw);
}

Napi::Object Init(Napi::Env env, Napi::Object exports) {
    exports.Set("scan", Napi::Function::New(env, Scan));
    return exports;
}

NODE_API_MODULE(hardware_scanner, Init)
