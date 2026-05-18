using System.Text.Json.Serialization;

namespace OCMaster.App.Models;

public class CpuInfo
{
    [JsonPropertyName("model")] public string Model { get; set; } = "";
    [JsonPropertyName("cores")] public int Cores { get; set; }
    [JsonPropertyName("threads")] public int Threads { get; set; }
    [JsonPropertyName("baseFreq")] public string BaseFreq { get; set; } = "";
    [JsonPropertyName("biosVersion")] public string BiosVersion { get; set; } = "";
}

public class MotherboardInfo
{
    [JsonPropertyName("brand")] public string Brand { get; set; } = "";
    [JsonPropertyName("model")] public string Model { get; set; } = "";
    [JsonPropertyName("chipset")] public string Chipset { get; set; } = "";
    [JsonPropertyName("biosVersion")] public string BiosVersion { get; set; } = "";
}

public class RamStick
{
    [JsonPropertyName("capacity")] public string Capacity { get; set; } = "";
    [JsonPropertyName("frequency")] public string Frequency { get; set; } = "";
    [JsonPropertyName("timings")] public string Timings { get; set; } = "";
    [JsonPropertyName("dieType")] public string DieType { get; set; } = "";
}

public class RamInfo
{
    [JsonPropertyName("totalCapacity")] public string TotalCapacity { get; set; } = "";
    [JsonPropertyName("stickCount")] public int StickCount { get; set; }
    [JsonPropertyName("channelCount")] public int ChannelCount { get; set; }
    [JsonPropertyName("sticks")] public List<RamStick> Sticks { get; set; } = new();
}

public class GpuInfo
{
    [JsonPropertyName("model")] public string Model { get; set; } = "";
    [JsonPropertyName("vram")] public string VRAM { get; set; } = "";
}

public class PsuInfo
{
    [JsonPropertyName("ratedWattage")] public string RatedWattage { get; set; } = "";
    [JsonPropertyName("source")] public string Source { get; set; } = "manual";
}

public class CoolerInfo
{
    [JsonPropertyName("type")] public string Type { get; set; } = "air";
    [JsonPropertyName("source")] public string Source { get; set; } = "manual";
}

public class HardwareInfo
{
    [JsonPropertyName("cpu")] public CpuInfo CPU { get; set; } = new();
    [JsonPropertyName("motherboard")] public MotherboardInfo Motherboard { get; set; } = new();
    [JsonPropertyName("ram")] public RamInfo RAM { get; set; } = new();
    [JsonPropertyName("gpu")] public GpuInfo GPU { get; set; } = new();
    [JsonPropertyName("psu")] public PsuInfo PSU { get; set; } = new();
    [JsonPropertyName("cooler")] public CoolerInfo Cooler { get; set; } = new();
}
