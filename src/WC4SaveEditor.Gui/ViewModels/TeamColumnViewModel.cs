using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class TeamColumnViewModel(uint teamId, string title) : ObservableObject
{
    public uint TeamId { get; } = teamId;

    [ObservableProperty]
    private string _title = title;

    public ObservableCollection<CountryCardViewModel> Members { get; } = [];
}
