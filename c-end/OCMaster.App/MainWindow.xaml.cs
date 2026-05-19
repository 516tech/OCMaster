using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using OCMaster.App.Pages;

namespace OCMaster.App;

public sealed partial class MainWindow : Window
{
    public MainWindow()
    {
        this.InitializeComponent();
        MainNav.SelectionChanged += OnNavSelectionChanged;
        ContentFrame.Navigate(typeof(ScanPage));
    }

    private void OnNavSelectionChanged(NavigationView sender, NavigationViewSelectionChangedEventArgs args)
    {
        if (args.SelectedItem is NavigationViewItem item)
        {
            switch (item.Tag?.ToString())
            {
                case "scan":     ContentFrame.Navigate(typeof(ScanPage)); break;
                case "upload":   ContentFrame.Navigate(typeof(UploadPage)); break;
                case "settings": ContentFrame.Navigate(typeof(SettingsPage)); break;
                case "about":    ContentFrame.Navigate(typeof(AboutPage)); break;
            }
        }
    }
}
