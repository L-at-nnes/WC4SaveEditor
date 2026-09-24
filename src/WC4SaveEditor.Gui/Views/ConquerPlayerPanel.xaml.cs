using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class ConquerPlayerPanel : UserControl
{
    public ConquerPlayerPanel()
    {
        InitializeComponent();
    }

    private void Annexer_Click(object sender, MouseButtonEventArgs e)
    {
        if (DataContext is not ConquerPlayerViewModel viewModel)
        {
            return;
        }
        if (sender is not FrameworkElement { DataContext: CountryCardViewModel card })
        {
            return;
        }
        viewModel.SelectedAnnexer = card;
    }

    private void Target_Click(object sender, MouseButtonEventArgs e)
    {
        if (DataContext is not ConquerPlayerViewModel viewModel)
        {
            return;
        }
        if (sender is not FrameworkElement { DataContext: SelectableCountryViewModel target })
        {
            return;
        }
        target.IsChecked = !target.IsChecked;
        viewModel.ToggleTarget(target);
    }
}
