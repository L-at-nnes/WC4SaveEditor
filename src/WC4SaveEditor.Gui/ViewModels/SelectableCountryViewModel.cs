using CommunityToolkit.Mvvm.ComponentModel;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class SelectableCountryViewModel(CountryCardViewModel card) : ObservableObject
{
    public CountryCardViewModel Card { get; } = card;

    [ObservableProperty]
    private bool _isChecked;

    [ObservableProperty]
    private int _badgeCount;

    [ObservableProperty]
    private bool _isSelectable = true;
}
