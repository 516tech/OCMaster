using System.Reflection;
using System.Runtime.InteropServices;
using System.Text.Json;
using OCMaster.App.Models;

namespace OCMaster.App.Services;

public static class NativeInterop
{
    private const string DllResourceName = "OCMaster.App.hardware_scanner.dll";
    private static string? _extractedDllPath;
    public static bool IsAvailable { get; private set; }

    static NativeInterop()
    {
        ExtractAndLoad();
    }

    /// <summary>从嵌入资源提取 Go DLL 到临时目录并加载</summary>
    private static void ExtractAndLoad()
    {
        try
        {
            // 检查是否已从嵌入资源提取
            var asm = Assembly.GetExecutingAssembly();
            using var stream = asm.GetManifestResourceStream(DllResourceName);

            if (stream == null)
            {
                // 未嵌入 — 尝试同目录下的 DLL
                var localDll = Path.Combine(AppContext.BaseDirectory, "hardware_scanner.dll");
                if (File.Exists(localDll))
                {
                    _extractedDllPath = localDll;
                }
                else
                {
                    IsAvailable = false;
                    return;
                }
            }
            else
            {
                // 提取到持久目录 (只需一次)
                var dir = Path.Combine(Environment.GetFolderPath(
                    Environment.SpecialFolder.LocalApplicationData), "OCMaster");
                Directory.CreateDirectory(dir);
                _extractedDllPath = Path.Combine(dir, "hardware_scanner.dll");

                // 版本检查：如果已存在且大小相同则跳过
                if (!File.Exists(_extractedDllPath) ||
                    new FileInfo(_extractedDllPath).Length != stream.Length)
                {
                    using var fs = new FileStream(_extractedDllPath, FileMode.Create, FileAccess.Write);
                    stream.CopyTo(fs);
                }
            }

            // 加载原生 DLL
            if (_extractedDllPath != null && File.Exists(_extractedDllPath))
            {
                NativeLibrary.Load(_extractedDllPath);
                IsAvailable = true;
            }
        }
        catch
        {
            IsAvailable = false;
        }
    }

    // ====== P/Invoke 声明 ======

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

    // ====== 托管 API ======

    private static readonly JsonSerializerOptions JsonOpts = new()
    {
        PropertyNameCaseInsensitive = true
    };

    public static HardwareInfo? ScanAllManaged()
    {
        if (!IsAvailable) return null;
        var ptr = ScanAll();
        try
        {
            var json = Marshal.PtrToStringUTF8(ptr);
            if (string.IsNullOrEmpty(json)) return null;
            return JsonSerializer.Deserialize<HardwareInfo>(json, JsonOpts);
        }
        finally { FreeString(ptr); }
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
