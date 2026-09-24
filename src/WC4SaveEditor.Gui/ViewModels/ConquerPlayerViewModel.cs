using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class ConquerPlayerViewModel : ObservableObject
{
    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private bool _suppressHandling;

    [ObservableProperty]
    private ObservableCollection<CountryCardViewModel> _annexers = [];

    [ObservableProperty]
    private ObservableCollection<SelectableCountryViewModel> _targets = [];

    [ObservableProperty]
    private CountryCardViewModel? _selectedAnnexer;

    public int CheckedCount => Targets.Count(t => t.IsChecked);

    public ConquerPlayerViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
    }

    public void Reset(SaveDocument document)
    {
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

        _queue.RemoveAll<ConquerPlayerChange>();
        foreach (var target in Targets)
        {
            target.IsChecked = false;
        }
        OnPropertyChanged(nameof(CheckedCount));
        _notifyChanged();
    }

    public void ToggleTarget(SelectableCountryViewModel target)
    {
        if (SelectedAnnexer is null)
        {
            target.IsChecked = false;
            return;
        }

        if (!target.IsSelectable)
        {
            target.IsChecked = false;
            NotificationCenter.Error($"{target.Card.Name} is the annexer - pick a different target.");
            return;
        }

        _queue.RemoveAll(c => c is ConquerPlayerChange pc && pc.OldPlayer == target.Card.PlayerId && pc.NewPlayer == SelectedAnnexer.PlayerId);
        if (target.IsChecked)
        {
            _queue.Add(new ConquerPlayerChange(target.Card.PlayerId, SelectedAnnexer.PlayerId));
            NotificationCenter.Info($"Will annex {target.Card.Name} into {SelectedAnnexer.Name} on Save.");
        }

        OnPropertyChanged(nameof(CheckedCount));
        _notifyChanged();
    }
}
