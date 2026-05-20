using Microsoft.UI.Xaml;
using OCMaster.App.Services;
using System.Runtime.InteropServices;
using System.Text.Json;

namespace OCMaster.App;

public partial class App : Application
{
    public static Window CurrentWindow { get; private set; } = null!;

    [DllImport("kernel32.dll")]
    private static extern bool AttachConsole(uint dwProcessId);
    private const uint ATTACH_PARENT_PROCESS = 0xFFFFFFFF;

    public App()
    {
        this.InitializeComponent();
    }

    protected override void OnLaunched(LaunchActivatedEventArgs args)
    {
        // CLI 模式: ocmaster.exe --cli scan [--out json] [--lang en-US]
        var cliArgs = Environment.GetCommandLineArgs();
        if (cliArgs.Length > 1 && cliArgs[1] == "--cli")
        {
            RunCLI(cliArgs);
            Environment.Exit(0);
            return;
        }

        // GUI 模式
        CurrentWindow = new MainWindow();
        CurrentWindow.Activate();
    }

    private static void RunCLI(string[] args)
    {
        AttachConsole(ATTACH_PARENT_PROCESS);

        if (args.Length < 3 || args[2] == "help" || args[2] == "--help")
        {
            PrintHelp();
            return;
        }

        var cmd = args[2];
        switch (cmd)
        {
            case "scan":
                RunScan(args);
                break;
            case "version":
                Console.WriteLine($"ocmaster version {GetVersion()}");
                break;
            default:
                Console.Error.WriteLine($"未知命令: {cmd}");
                Environment.Exit(1);
                break;
        }
    }

    private static void RunScan(string[] args)
    {
        if (!NativeInterop.IsAvailable)
        {
            Console.Error.WriteLine("扫描模块不可用 — hardware_scanner.dll 未找到");
            Environment.Exit(1);
            return;
        }

        var hw = NativeInterop.ScanAllManaged();
        if (hw == null)
        {
            Console.Error.WriteLine("扫描失败");
            Environment.Exit(1);
            return;
        }

        var outFormat = "table";
        for (int i = 0; i < args.Length; i++)
        {
            if (args[i] == "--out" && i + 1 < args.Length)
                outFormat = args[i + 1];
        }

        if (outFormat == "json")
        {
            Console.WriteLine(JsonSerializer.Serialize(hw, new JsonSerializerOptions
            { WriteIndented = true, PropertyNamingPolicy = JsonNamingPolicy.CamelCase }));
        }
        else
        {
            PrintTable(hw);
        }
    }

    private static void PrintTable(dynamic hw)
    {
        Console.WriteLine();
        Console.WriteLine("  OCMaster Scan Result");
        Console.WriteLine("  ====================");
        Console.WriteLine($"  CPU:  {hw.CPU.Model} ({hw.CPU.Cores}C/{hw.CPU.Threads}T) @ {hw.CPU.BaseFreq}");
        Console.WriteLine($"  MB:   {hw.Motherboard.Brand} {hw.Motherboard.Model}");
        Console.WriteLine($"  RAM:  {hw.RAM.TotalCapacity} ({hw.RAM.StickCount} sticks, {hw.RAM.ChannelCount} ch)");
        Console.WriteLine($"  GPU:  {hw.GPU.Model} ({hw.GPU.VRAM})");
        Console.WriteLine();
    }

    private static void PrintHelp()
    {
        Console.WriteLine("OCMaster CLI");
        Console.WriteLine();
        Console.WriteLine("  ocmaster --cli scan                扫描硬件，输出文本表格");
        Console.WriteLine("  ocmaster --cli scan --out json     扫描硬件，输出 JSON");
        Console.WriteLine("  ocmaster --cli help                帮助信息");
        Console.WriteLine("  ocmaster                           启动图形界面");
    }

    private static string GetVersion()
    {
        try { return typeof(App).Assembly.GetName().Version?.ToString() ?? "0.0.1"; }
        catch { return "0.0.1"; }
    }
}
