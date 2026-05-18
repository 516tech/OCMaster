using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Services;

namespace OCMaster.App.Pages;

public sealed partial class SettingsPage : Page
{
    private readonly AppConfig _config;

    public SettingsPage()
    {
        this.InitializeComponent();
        _config = ConfigService.Load();

        ApiUrlInput.Text = _config.ApiUrl;
        for (int i = 0; i < LanguageCombo.Items.Count; i++)
        {
            if ((LanguageCombo.Items[i] as ComboBoxItem)?.Tag?.ToString() == _config.Language)
            {
                LanguageCombo.SelectedIndex = i;
                break;
            }
        }
        AutoScanCheck.IsChecked = _config.AutoScan;
    }

    private void SaveBtn_Click(object sender, RoutedEventArgs e)
    {
        _config.ApiUrl = ApiUrlInput.Text.Trim();
        _config.Language = (LanguageCombo.SelectedItem as ComboBoxItem)?.Tag?.ToString() ?? "zh-CN";
        _config.AutoScan = AutoScanCheck.IsChecked ?? true;
        ConfigService.Save(_config);

        if (Frame.CanGoBack) Frame.GoBack();
    }

    private void CancelBtn_Click(object sender, RoutedEventArgs e)
    {
        if (Frame.CanGoBack) Frame.GoBack();
    }
}
