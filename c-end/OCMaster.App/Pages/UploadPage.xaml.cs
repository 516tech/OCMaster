using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Models;
using OCMaster.App.Services;

namespace OCMaster.App.Pages;

public sealed partial class UploadPage : Page
{
    private string? _shareCode;

    public UploadPage()
    {
        this.InitializeComponent();
    }

    private async void UploadBtn_Click(object sender, RoutedEventArgs e)
    {
        UploadBtn.IsEnabled = false;
        UploadProgress.Visibility = Visibility.Visible;
        ErrorBar.IsOpen = false;

        try
        {
            // 1. 扫描硬件
            var hw = await Task.Run(() => NativeInterop.ScanAllManaged());
            if (hw == null || string.IsNullOrEmpty(hw.CPU?.Model))
            {
                ErrorBar.Message = "扫描硬件失败，请先到「硬件扫描」页确认扫描正常";
                ErrorBar.IsOpen = true;
                return;
            }

            // 2. 上传
            var config = ConfigService.Load();
            var client = new ApiClient(config.ApiUrl);
            var code = await client.UploadHardware(hw);

            if (!string.IsNullOrEmpty(code))
            {
                _shareCode = code;
                CodeText.Text = code;
                CodePanel.Visibility = Visibility.Visible;
            }
            else
            {
                ErrorBar.Message = "上传失败，请检查网络连接和 API 地址";
                ErrorBar.IsOpen = true;
            }
        }
        catch (Exception ex)
        {
            ErrorBar.Message = $"上传失败: {ex.Message}";
            ErrorBar.IsOpen = true;
        }
        finally
        {
            UploadBtn.IsEnabled = true;
            UploadProgress.Visibility = Visibility.Collapsed;
        }
    }

    private void CopyBtn_Click(object sender, RoutedEventArgs e)
    {
        if (_shareCode != null)
        {
            var pkg = new Windows.ApplicationModel.DataTransfer.DataPackage();
            pkg.SetText(_shareCode);
            Windows.ApplicationModel.DataTransfer.Clipboard.SetContent(pkg);
        }
    }

    private async void DeleteBtn_Click(object sender, RoutedEventArgs e)
    {
        if (_shareCode == null) return;

        var config = ConfigService.Load();
        var client = new ApiClient(config.ApiUrl);

        try
        {
            await client.DeleteHardware(_shareCode);
            _shareCode = null;
            CodePanel.Visibility = Visibility.Collapsed;
        }
        catch (Exception ex)
        {
            ErrorBar.Message = $"删除失败: {ex.Message}";
            ErrorBar.IsOpen = true;
        }
    }
}
