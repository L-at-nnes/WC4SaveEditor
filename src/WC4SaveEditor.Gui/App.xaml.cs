using System.Windows;
using System.Windows.Threading;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui;

public partial class App : Application
{
    protected override void OnStartup(StartupEventArgs e)
    {
        base.OnStartup(e);
        DispatcherUnhandledException += OnDispatcherUnhandledException;
    }

    private static void OnDispatcherUnhandledException(object sender, DispatcherUnhandledExceptionEventArgs e)
    {
        NotificationCenter.Error($"Unexpected error: {e.Exception.Message}");
        e.Handled = true;
    }
}
