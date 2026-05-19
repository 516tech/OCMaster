using Microsoft.UI.Xaml;

namespace OCMaster.App;

public partial class App : Application
{
    public static Window CurrentWindow { get; private set; } = null!;

    public App()
    {
        this.InitializeComponent();
    }

    protected override void OnLaunched(LaunchActivatedEventArgs args)
    {
        CurrentWindow = new MainWindow();
        CurrentWindow.Activate();
    }
}
