using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class ConquerProvinceViewModel : ObservableObject
{
    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private bool _suppressHandling;
    private SaveDocument? _document;
    private readonly Dictionary<int, HashSet<ushort>> _selectedProvincesByTarget = [];

    [ObservableProperty]
    private ObservableCollection<CountryCardViewModel> _annexers = [];

    [ObservableProperty]
    private ObservableCollection<SelectableCountryViewModel> _targets = [];

    [ObservableProperty]
    private CountryCardViewModel? _selectedAnnexer;

    public ConquerProvinceViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
    }

    public void Reset(SaveDocument document)
    {
        _document = document;
        _selectedProvincesByTarget.Clear();
        _suppressHandling = true;
        var cards = document.Players.Select((p, i) => new CountryCardViewModel(i, p)).ToList();
        Annexers = new ObservableCollection<CountryCardViewModel>(cards);
        Targets = new ObservableCollection<SelectableCountryViewModel>(cards.Select(c => new SelectableCountryViewModel(c)));
        SelectedAnnexer = null;
        _suppressHandling = false;
    }

    partial void OnSelectedAnnexerChanged(CountryCardViewModel? value)
    {
        foreach (var annexer in Annexers)
        {
            annexer.IsSelected = ReferenceEquals(annexer, value);
        }
        foreach (var target in Targets)
        {
            target.IsSelectable = value is null || target.Card.PlayerId != value.PlayerId;
        }

        if (_suppressHandling)
        {
            return;
        }
        _queue.RemoveAll<ConquerProvinceChange>();
        _selectedProvincesByTarget.Clear();
        foreach (var target in Targets)
        {
            target.BadgeCount = 0;
        }
        _notifyChanged();
    }

    public List<ProvinceInfo> GetProvincesFor(SelectableCountryViewModel target) =>
        _document is null ? [] : ProvinceLookup.GetPlayerProvinces(_document, target.Card.PlayerId);

    public IReadOnlySet<ushort> GetSelection(SelectableCountryViewModel target) =>
        _selectedProvincesByTarget.TryGetValue(target.Card.PlayerId, out var set) ? set : new HashSet<ushort>();

    public void ApplySelection(SelectableCountryViewModel target, IReadOnlyCollection<ushort> selectedProvinceCodes)
    {
        if (SelectedAnnexer is null || !target.IsSelectable)
        {
            return;
        }

        _queue.RemoveAll(c => c is ConquerProvinceChange pc && pc.NewPlayer == SelectedAnnexer.PlayerId
            && _selectedProvincesByTarget.TryGetValue(target.Card.PlayerId, out var old) && old.Contains(pc.ProvinceCode));

        _selectedProvincesByTarget[target.Card.PlayerId] = new HashSet<ushort>(selectedProvinceCodes);
        target.BadgeCount = selectedProvinceCodes.Count;

        foreach (var code in selectedProvinceCodes)
        {
            _queue.Add(new ConquerProvinceChange(code, SelectedAnnexer.PlayerId));
        }

        if (selectedProvinceCodes.Count > 0)
        {
            Services.NotificationCenter.Info($"Will annex {selectedProvinceCodes.Count} province(s) from {target.Card.Name} into {SelectedAnnexer.Name} on Save.");
        }
        _notifyChanged();
    }
}
