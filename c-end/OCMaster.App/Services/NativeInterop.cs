using System.Runtime.InteropServices;
using System.Text.Json;
using OCMaster.App.Models;

namespace OCMaster.App.Services;

public static class NativeInterop
{
    // Go DLL 导出函数声明
    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern IntPtr ScanAll();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern IntPtr ScanCPU();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern IntPtr ScanRAM();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern IntPtr ScanGPU();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern IntPtr ScanMotherboard();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    private static extern void FreeString(IntPtr ptr);

    private static readonly JsonSerializerOptions JsonOpts = new()
    {
        PropertyNameCaseInsensitive = true
    };

    /// <summary>调用 Go DLL 扫描全部硬件，返回 HardwareInfo</summary>
    public static HardwareInfo? ScanAllManaged()
    {
        var ptr = ScanAll();
        try
        {
            var json = Marshal.PtrToStringUTF8(ptr);
            if (string.IsNullOrEmpty(json)) return null;
            return JsonSerializer.Deserialize<HardwareInfo>(json, JsonOpts);
        }
        finally
        {
            FreeString(ptr);
        }
    }

    private static string CallAndFree(IntPtr ptr)
    {
        try { return Marshal.PtrToStringUTF8(ptr) ?? ""; }
        finally { FreeString(ptr); }
    }

    public static string ScanCPUJson() => CallAndFree(ScanCPU());
    public static string ScanRAMJson() => CallAndFree(ScanRAM());
    public static string ScanGPUJson() => CallAndFree(ScanGPU());
    public static string ScanMotherboardJson() => CallAndFree(ScanMotherboard());
}
