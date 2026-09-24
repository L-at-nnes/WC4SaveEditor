using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class ConquerProvincePanel : UserControl
{
    public ConquerProvincePanel()
    {
        InitializeComponent();
    }

    private void Annexer_Click(object sender, MouseButtonEventArgs e)
    {
        if (DataContext is not ConquerProvinceViewModel viewModel)
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
        if (DataContext is not ConquerProvinceViewModel viewModel)
        {
            return;
        }
        if (sender is not FrameworkElement { DataContext: SelectableCountryViewModel target })
        {
            return;
        }
        if (viewModel.SelectedAnnexer is null)
        {
            MessageBox.Show("Pick an annexer on the left first.", "Conquer province", MessageBoxButton.OK, MessageBoxImage.Information);
            return;
        }

        var provinces = viewModel.GetProvincesFor(target);
        if (provinces.Count == 0)
        {
            MessageBox.Show($"{target.Card.Name} has no provinces to pick from.", "Conquer province", MessageBoxButton.OK, MessageBoxImage.Information);
            return;
        }

        var currentSelection = viewModel.GetSelection(target);
        var items = provinces.Select(p => new ProvinceSelectionViewModel(p, currentSelection.Contains(p.CoordinateCode))).ToList();

        var dialog = new ProvincePickerDialog(target.Card.Name, items) { Owner = Window.GetWindow(this) };
        if (dialog.ShowDialog() == true)
        {
            viewModel.ApplySelection(target, dialog.SelectedCodes);
        }
    }
}
