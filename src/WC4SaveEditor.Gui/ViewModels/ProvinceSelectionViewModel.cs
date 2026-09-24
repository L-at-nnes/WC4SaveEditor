using CommunityToolkit.Mvvm.ComponentModel;
using WC4SaveEditor.Core.Operations;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class ProvinceSelectionViewModel(ProvinceInfo province, bool isChecked) : ObservableObject
{
    public ProvinceInfo Province { get; } = province;

    [ObservableProperty]
    private bool _isChecked = isChecked;
}
