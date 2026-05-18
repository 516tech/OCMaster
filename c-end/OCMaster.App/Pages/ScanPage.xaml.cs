using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Models;
using OCMaster.App.Services;
using System.ComponentModel;
using System.Runtime.CompilerServices;

namespace OCMaster.App.Pages;

public sealed partial class ScanPage : Page, INotifyPropertyChanged
{
    private HardwareInfo _hardware = new();
    private bool _isScanning;
    private string _errorMessage = "";
    private string _shareCode = "";

    public HardwareInfo Hardware { get => _hardware; set { _hardware = value; OnPropertyChanged(); } }
    public bool IsScanning { get => _isScanning; set { _isScanning = value; OnPropertyChanged(); OnPropertyChanged(nameof(IsNotScanning)); OnPropertyChanged(nameof(ShowHint)); } }
    public bool IsNotScanning => !_isScanning;
    public bool ShowHint => !_isScanning && !HasHardware;
    public string ErrorMessage { get => _errorMessage; set { _errorMessage = value; OnPropertyChanged(); OnPropertyChanged(nameof(HasError)); } }
    public bool HasError => !string.IsNullOrEmpty(_errorMessage);
    public bool HasHardware => !string.IsNullOrEmpty(Hardware.CPU?.Model);
    public bool HasShareCode => !string.IsNullOrEmpty(_shareCode);
    public string ShareCodeText => $"分享码: {_shareCode} (7天有效)";
    public string CpuCoresText => $"{Hardware.CPU?.Cores ?? 0}C / {Hardware.CPU?.Threads ?? 0}T";
    public Button CopyCodeBtn { get; } = new() { Content = "复制" };

    public ScanPage()
    {
        this.InitializeComponent();
        CopyCodeBtn.Click += async (_, _) =>
        {
            Windows.ApplicationModel.DataTransfer.DataPackage pkg = new();
            pkg.SetText(_shareCode);
            Windows.ApplicationModel.DataTransfer.Clipboard.SetContent(pkg);
        };
    }

    private async void ScanBtn_Click(object sender, RoutedEventArgs e)
    {
        IsScanning = true;
        ErrorMessage = "";

        try
        {
            await Task.Run(() =>
            {
                var result = NativeInterop.ScanAllManaged();
                if (result != null && !string.IsNullOrEmpty(result.CPU?.Model))
                {
                    DispatcherQueue.TryEnqueue(() => Hardware = result);
                }
                else
                {
                    throw new Exception("扫描模块不可用 — 请确认 hardware_scanner.dll 已编译");
                }
            });
        }
        catch (Exception ex)
        {
            ErrorMessage = ex.Message;
        }
        finally
        {
            IsScanning = false;
        }
    }

    private async void UploadBtn_Click(object sender, RoutedEventArgs e)
    {
        var config = ConfigService.Load();
        var client = new ApiClient(config.ApiUrl);

        try
        {
            var code = await client.UploadHardware(Hardware);
            if (!string.IsNullOrEmpty(code))
            {
                _shareCode = code;
                OnPropertyChanged(nameof(HasShareCode));
                OnPropertyChanged(nameof(ShareCodeText));
            }
            else
            {
                ErrorMessage = "上传失败，请检查网络连接和 API 地址";
            }
        }
        catch (Exception ex)
        {
            ErrorMessage = $"上传失败: {ex.Message}";
        }
    }

    private async void ExportBtn_Click(object sender, RoutedEventArgs e)
    {
        var h = Hardware;
        var txt = string.Join("\n",
            $"CPU: {h.CPU?.Model} | {CpuCoresText} | {h.CPU?.BaseFreq}",
            $"Motherboard: {h.Motherboard?.Brand} {h.Motherboard?.Model} | {h.Motherboard?.Chipset} | BIOS {h.Motherboard?.BiosVersion}",
            $"RAM: {h.RAM?.TotalCapacity} | {h.RAM?.StickCount} sticks | {h.RAM?.ChannelCount} channels",
            $"GPU: {h.GPU?.Model} | {h.GPU?.VRAM}"
        );

        var savePicker = new Windows.Storage.Pickers.FileSavePicker();
        savePicker.SuggestedStartLocation = Windows.Storage.Pickers.PickerLocationId.Desktop;
        savePicker.FileTypeChoices.Add("Text", new[] { ".txt" });
        savePicker.SuggestedFileName = $"ocmaster-hardware-{DateTime.Now:yyyyMMdd-HHmmss}";

        // WinUI 3 中 FileSavePicker 需要窗口句柄
        var hwnd = WinRT.Interop.WindowNative.GetWindowHandle(
            (Application.Current as App)?.GetType().GetProperty("_mainWindow",
                System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance)
                ?.GetValue(Application.Current) as Window ?? this);
        WinRT.Interop.InitializeWithWindow.Initialize(savePicker, hwnd);

        var file = await savePicker.PickSaveFileAsync();
        if (file != null)
        {
            await Windows.Storage.FileIO.WriteTextAsync(file, txt);
        }
    }

    public event PropertyChangedEventHandler? PropertyChanged;
    private void OnPropertyChanged([CallerMemberName] string? name = null)
        => PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(name));
}
