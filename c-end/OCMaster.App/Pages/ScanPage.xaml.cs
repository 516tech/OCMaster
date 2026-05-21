using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Models;
using OCMaster.App.Services;
using WinRT.Interop;

namespace OCMaster.App.Pages;

public sealed partial class ScanPage : Page
{
    private HardwareInfo? _hardware;

    public ScanPage()
    {
        this.InitializeComponent();
        UpdateUI();
    }

    private async void ScanBtn_Click(object sender, RoutedEventArgs e)
    {
        if (!NativeInterop.IsAvailable)
        {
            ErrorBar.Message = "扫描模块未安装 — hardware_scanner.dll 未找到";
            ErrorBar.IsOpen = true;
            return;
        }

        ScanBtn.IsEnabled = false;
        ScanProgress.Visibility = Visibility.Visible;
        ErrorBar.IsOpen = false;

        try
        {
            await Task.Run(() =>
            {
                var result = NativeInterop.ScanAllManaged();
                if (result != null && !string.IsNullOrEmpty(result.CPU?.Model))
                    DispatcherQueue.TryEnqueue(() => _hardware = result);
                else
                    throw new Exception("扫描无结果");
            });

            UpdateUI();
        }
        catch (Exception ex)
        {
            ErrorBar.Message = ex.Message;
            ErrorBar.IsOpen = true;
        }
        finally
        {
            ScanBtn.IsEnabled = true;
            ScanProgress.Visibility = Visibility.Collapsed;
        }
    }

    private void UpdateUI()
    {
        var hw = _hardware;
        bool hasData = hw != null && !string.IsNullOrEmpty(hw.CPU?.Model);

        HintText.Visibility = hasData ? Visibility.Collapsed : Visibility.Visible;
        HardwarePanel.Visibility = hasData ? Visibility.Visible : Visibility.Collapsed;
        DllWarning.Visibility = NativeInterop.IsAvailable ? Visibility.Collapsed : Visibility.Visible;

        if (!hasData) return;

        CpuModelText.Text = $"型号: {hw!.CPU.Model}";
        CpuCoresText.Text = $"核心/线程: {hw.CPU.Cores}C / {hw.CPU.Threads}T";
        CpuFreqText.Text = $"基频: {hw.CPU.BaseFreq}";

        MbBrandText.Text = $"品牌: {hw.Motherboard.Brand}";
        MbModelText.Text = $"型号: {hw.Motherboard.Model}";
        MbBiosText.Text = $"BIOS: {hw.Motherboard.BiosVersion}";

        RamCapacityText.Text = $"总容量: {hw.RAM.TotalCapacity}";
        RamSticksText.Text = $"条数: {hw.RAM.StickCount}";
        RamChannelsText.Text = $"通道: {hw.RAM.ChannelCount}";

        GpuModelText.Text = $"型号: {hw.GPU.Model}";
        GpuVRamText.Text = $"显存: {hw.GPU.VRAM}";
    }

    private async void ExportBtn_Click(object sender, RoutedEventArgs e)
    {
        var hw = _hardware;
        if (hw == null) return;

        var txt = string.Join("\n",
            $"CPU: {hw.CPU.Model} | {hw.CPU.Cores}C/{hw.CPU.Threads}T | {hw.CPU.BaseFreq}",
            $"Motherboard: {hw.Motherboard.Brand} {hw.Motherboard.Model} | BIOS {hw.Motherboard.BiosVersion}",
            $"RAM: {hw.RAM.TotalCapacity} | {hw.RAM.StickCount} sticks | {hw.RAM.ChannelCount} channels",
            $"GPU: {hw.GPU.Model} | {hw.GPU.VRAM}"
        );

        var savePicker = new Windows.Storage.Pickers.FileSavePicker
        {
            SuggestedStartLocation = Windows.Storage.Pickers.PickerLocationId.Desktop,
            SuggestedFileName = $"ocmaster-hardware-{DateTime.Now:yyyyMMdd-HHmmss}"
        };
        savePicker.FileTypeChoices.Add("Text", new[] { ".txt" });

        var hwnd = WindowNative.GetWindowHandle(App.CurrentWindow);
        InitializeWithWindow.Initialize(savePicker, hwnd);

        var file = await savePicker.PickSaveFileAsync();
        if (file != null)
            await Windows.Storage.FileIO.WriteTextAsync(file, txt);
    }
}
