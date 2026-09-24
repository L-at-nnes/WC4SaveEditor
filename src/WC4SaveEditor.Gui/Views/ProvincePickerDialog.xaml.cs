using System.Windows;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class ProvincePickerDialog : Window
{
    private readonly List<ProvinceSelectionViewModel> _items;

    public IReadOnlyList<ushort> SelectedCodes { get; private set; } = [];

    public ProvincePickerDialog(string targetName, List<ProvinceSelectionViewModel> items)
    {
        InitializeComponent();
        _items = items;
        HeaderText.Text = $"Provinces owned by {targetName}";
        ProvinceList.ItemsSource = items;
    }

    private void Apply_Click(object sender, RoutedEventArgs e)
    {
        SelectedCodes = _items.Where(i => i.IsChecked).Select(i => i.Province.CoordinateCode).ToList();
        DialogResult = true;
    }

    private void Cancel_Click(object sender, RoutedEventArgs e) => DialogResult = false;
}
