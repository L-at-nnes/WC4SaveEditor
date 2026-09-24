using System.Runtime.InteropServices;
using System.Windows;
using System.Windows.Interop;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui;

public partial class MainWindow : Window
{
    [DllImport("dwmapi.dll", PreserveSig = true)]
    private static extern int DwmSetWindowAttribute(IntPtr hwnd, int attribute, ref int value, int valueSize);

    private const int DwmwaUseImmersiveDarkModeBefore20H1 = 19;
    private const int DwmwaUseImmersiveDarkMode = 20;

    public MainWindow()
    {
        InitializeComponent();
        DataContext = new MainViewModel();
        SourceInitialized += (_, _) => EnableDarkTitleBar();
    }

    private void EnableDarkTitleBar()
    {
        var hwnd = new WindowInteropHelper(this).Handle;
        var useDarkMode = 1;
        if (DwmSetWindowAttribute(hwnd, DwmwaUseImmersiveDarkMode, ref useDarkMode, sizeof(int)) != 0)
        {
            DwmSetWindowAttribute(hwnd, DwmwaUseImmersiveDarkModeBefore20H1, ref useDarkMode, sizeof(int));
        }
    }

    private void ScrollLeft_Click(object sender, RoutedEventArgs e) => SaveFileScroller.ScrollToHorizontalOffset(SaveFileScroller.HorizontalOffset - 250);

    private void ScrollRight_Click(object sender, RoutedEventArgs e) => SaveFileScroller.ScrollToHorizontalOffset(SaveFileScroller.HorizontalOffset + 250);

    private void Settings_Click(object sender, RoutedEventArgs e)
    {
        var dialog = new Views.SettingsDialog((ViewModels.MainViewModel)DataContext) { Owner = this };
        dialog.ShowDialog();
    }
}
