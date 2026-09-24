using System.Windows;
using WC4SaveEditor.Gui.Services;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class SettingsDialog : Window
{
    private readonly MainViewModel _mainViewModel;

    public SettingsDialog(MainViewModel mainViewModel)
    {
        InitializeComponent();
        _mainViewModel = mainViewModel;
    }

    private void Refresh_Click(object sender, RoutedEventArgs e)
    {
        _mainViewModel.RefreshSaveFileList();
        NotificationCenter.Info("Save file list refreshed.");
    }

    private void Close_Click(object sender, RoutedEventArgs e) => Close();
}
