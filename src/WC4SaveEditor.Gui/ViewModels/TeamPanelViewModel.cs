using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class TeamPanelViewModel : ObservableObject
{
    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private SaveDocument? _document;

    [ObservableProperty]
    private ObservableCollection<TeamColumnViewModel> _columns = [];

    public TeamPanelViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
    }

    public void Reset(SaveDocument document)
    {
        _document = document;

        var cards = document.Players
            .Select((p, i) => new CountryCardViewModel(i, p))
            .ToList();

        var columns = new ObservableCollection<TeamColumnViewModel>();
        foreach (var teamId in cards.Select(c => c.OriginalTeamId).Distinct().OrderBy(t => t))
        {
            var column = new TeamColumnViewModel(teamId, $"Team {teamId}");
            foreach (var card in cards.Where(c => c.OriginalTeamId == teamId))
            {
                column.Members.Add(card);
            }
            columns.Add(column);
        }
        Columns = columns;
    }

    public void MoveCard(CountryCardViewModel card, uint targetTeamId)
    {
        if (card.CurrentTeamId == targetTeamId)
        {
            return;
        }

        foreach (var column in Columns)
        {
            column.Members.Remove(card);
        }
        card.CurrentTeamId = targetTeamId;

        var destination = Columns.FirstOrDefault(c => c.TeamId == targetTeamId);
        destination?.Members.Add(card);

        _queue.RemoveAll(c => c is ChangeTeamChange ct && ct.PlayerId == card.PlayerId);
        if (targetTeamId != card.OriginalTeamId)
        {
            _queue.Add(new ChangeTeamChange(card.PlayerId, targetTeamId));
            NotificationCenter.Info($"Will move {card.Name} to Team {targetTeamId} on Save.");
        }
        _notifyChanged();
    }

    [RelayCommand]
    private void UniteTeam()
    {
        if (_document is null || Columns.Count == 0)
        {
            return;
        }

        var mainTeamId = Columns.SelectMany(c => c.Members).First(c => c.PlayerId == ConquestOperations.MainPlayer).CurrentTeamId;

        _queue.RemoveAll<ChangeTeamChange>();
        _queue.RemoveAll<UniteTeamChange>();
        _queue.Add(new UniteTeamChange());

        var allCards = Columns.SelectMany(c => c.Members).ToList();
        foreach (var column in Columns)
        {
            column.Members.Clear();
        }
        var mainColumn = Columns.First(c => c.TeamId == mainTeamId);
        foreach (var card in allCards)
        {
            card.CurrentTeamId = mainTeamId;
            mainColumn.Members.Add(card);
        }

        NotificationCenter.Info("Will move every player to your team on Save.");
        _notifyChanged();
    }
}
