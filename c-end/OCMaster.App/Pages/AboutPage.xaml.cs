using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Services;

namespace OCMaster.App.Pages;

public sealed partial class AboutPage : Page
{
    public AboutPage()
    {
        this.InitializeComponent();
        try
        {
            var ver = typeof(NativeInterop).Assembly.GetName().Version;
            VersionText.Text = $"版本: {ver?.ToString() ?? "0.0.1"}";
        }
        catch
        {
            VersionText.Text = "版本: 0.0.1";
        }
    }
}
