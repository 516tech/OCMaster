using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Models;
using OCMaster.App.Services;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using WinRT;

namespace OCMaster.App.Pages;

public sealed partial class ScanPage : Page, INotifyPropertyChanged
{
    private HardwareInfo _hardware = new();
    private bool _isScanning;
    private string _errorMessage = "";

    public HardwareInfo Hardware { get => _hardware; set { _hardware = value; OnPropertyChanged(); OnPropertyChanged(nameof(HasHardware)); } }
    public bool IsScanning { get => _isScanning; set { _isScanning = value; OnPropertyChanged(); OnPropertyChanged(nameof(IsNotScanning)); OnPropertyChanged(nameof(ShowHint)); } }
    public bool IsNotScanning => !_isScanning;
    public bool ShowHint => !_isScanning && !HasHardware;
    public string ErrorMessage { get => _errorMessage; set { _errorMessage = value; OnPropertyChanged(); OnPropertyChanged(nameof(HasError)); } }
    public bool HasError => !string.IsNullOrEmpty(_errorMessage);
    public bool HasHardware => !string.IsNullOrEmpty(Hardware.CPU?.Model);
    public bool DllMissing => !NativeInterop.IsAvailable;
    public string CpuCoresText => $"{Hardware.CPU?.Cores ?? 0}C / {Hardware.CPU?.Threads ?? 0}T";

    public ScanPage()
    {
        this.InitializeComponent();
    }

    private async void ScanBtn_Click(object sender, RoutedEventArgs e)
    {
        if (!NativeInterop.IsAvailable)
        {
            ErrorMessage = "扫描模块未安装 — 请确认 hardware_scanner.dll 已编译";
            return;
        }

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
                    throw new Exception("扫描无结果 — 请确认当前系统支持硬件检测");
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

    private async void ExportBtn_Click(object sender, RoutedEventArgs e)
    {
        var h = Hardware;
        var txt = string.Join("\n",
            $"CPU: {h.CPU?.Model} | {CpuCoresText} | {h.CPU?.BaseFreq}",
            $"Motherboard: {h.Motherboard?.Brand} {h.Motherboard?.Model} | BIOS {h.Motherboard?.BiosVersion}",
            $"RAM: {h.RAM?.TotalCapacity} | {h.RAM?.StickCount} sticks | {h.RAM?.ChannelCount} channels",
            $"GPU: {h.GPU?.Model} | {h.GPU?.VRAM}"
        );

        var savePicker = new Windows.Storage.Pickers.FileSavePicker
        {
            SuggestedStartLocation = Windows.Storage.Pickers.PickerLocationId.Desktop,
            SuggestedFileName = $"ocmaster-hardware-{DateTime.Now:yyyyMMdd-HHmmss}"
        };
        savePicker.FileTypeChoices.Add("Text", new[] { ".txt" });

        var hwnd = WindowNative.GetWindowHandle(App.CurrentWindow);
        Interop.InitializeWithWindow.Initialize(savePicker, hwnd);

        var file = await savePicker.PickSaveFileAsync();
        if (file != null)
            await Windows.Storage.FileIO.WriteTextAsync(file, txt);
    }

    public event PropertyChangedEventHandler? PropertyChanged;
    private void OnPropertyChanged([CallerMemberName] string? name = null)
        => PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(name));
}
