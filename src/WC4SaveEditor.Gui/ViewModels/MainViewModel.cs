using System.Collections.ObjectModel;
using System.IO;
using System.Windows;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Core.SaveFile;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class MainViewModel : ObservableObject
{
    public ChangeQueue Queue { get; } = new();

    [ObservableProperty]
    private SaveDocument? _document;

    [ObservableProperty]
    private ObservableCollection<string> _saveFiles = [];

    [ObservableProperty]
    private string? _selectedSaveFile;

    [ObservableProperty]
    private bool _isDocumentLoaded;

    public QuickTogglesViewModel QuickToggles { get; }
    public TeamPanelViewModel TeamPanel { get; }
    public ConquerPlayerViewModel ConquerPlayer { get; }
    public ConquerProvinceViewModel ConquerProvince { get; }
    public SpawnUnitsViewModel SpawnUnits { get; }
    public MapPanelViewModel MapPanel { get; }

    public int PendingCount => Queue.Count;

    private string? _lastLoadedSaveFile;
    private bool _isRevertingSelection;

    public MainViewModel()
    {
        QuickToggles = new QuickTogglesViewModel(Queue, NotifyQueueChanged);
        TeamPanel = new TeamPanelViewModel(Queue, NotifyQueueChanged);
        ConquerPlayer = new ConquerPlayerViewModel(Queue, NotifyQueueChanged);
        ConquerProvince = new ConquerProvinceViewModel(Queue, NotifyQueueChanged);
        SpawnUnits = new SpawnUnitsViewModel(Queue, NotifyQueueChanged);
        MapPanel = new MapPanelViewModel(Queue, NotifyQueueChanged);

        RefreshSaveFileList();
        var defaultFile = SaveLocator.FindDefaultSaveFile();
        if (defaultFile is not null)
        {
            SelectedSaveFile = defaultFile;
        }
    }

    public void RefreshSaveFileList()
    {
        var files = SaveLocator.FindSaveFiles();
        SaveFiles = new ObservableCollection<string>(files);
    }

    partial void OnSelectedSaveFileChanged(string? value)
    {
        if (string.IsNullOrEmpty(value) || _isRevertingSelection)
        {
            return;
        }

        if (PendingCount > 0)
        {
            var result = MessageBox.Show(
                $"You have {PendingCount} unsaved change(s) on {(_lastLoadedSaveFile is null ? "the current file" : Path.GetFileName(_lastLoadedSaveFile))}. Switching files will discard them. Continue?",
                "Unsaved changes",
                MessageBoxButton.YesNo,
                MessageBoxImage.Warning);

            if (result != MessageBoxResult.Yes)
            {
                _isRevertingSelection = true;
                SelectedSaveFile = _lastLoadedSaveFile;
                _isRevertingSelection = false;
                return;
            }
        }

        LoadDocument(value);
    }

    private void LoadDocument(string path)
    {
        try
        {
            Queue.Clear();
            Document = SaveFileReader.Read(path);
            _lastLoadedSaveFile = path;
            IsDocumentLoaded = true;

            QuickToggles.Reset(Document);
            TeamPanel.Reset(Document);
            ConquerPlayer.Reset(Document);
            ConquerProvince.Reset(Document);
            SpawnUnits.Reset(Document);
            MapPanel.Reset(Document);

            NotifyQueueChanged();
            NotificationCenter.Info($"Loaded {Path.GetFileName(path)} ({Document.Players.Count} players, {Document.Cities.Count} cities, {Document.Units.Count} units)");
        }
        catch (Exception ex)
        {
            IsDocumentLoaded = false;
            NotificationCenter.Error($"Failed to load '{Path.GetFileName(path)}': {ex.Message}");
        }
    }

    private void NotifyQueueChanged()
    {
        OnPropertyChanged(nameof(PendingCount));
        SaveCommand.NotifyCanExecuteChanged();

        try
        {
            MapPanel.RefreshPreview();
        }
        catch (Exception ex)
        {
            NotificationCenter.Error($"Couldn't preview that change: {ex.Message}");
        }
    }

    [RelayCommand(CanExecute = nameof(CanSave))]
    private void Save()
    {
        if (Document is null)
        {
            return;
        }

        try
        {
            Queue.Commit(Document);
            NotificationCenter.Success($"Saved {Path.GetFileName(Document.FilePath)} (backup kept as .bak).");
            LoadDocument(Document.FilePath);
        }
        catch (Exception ex)
        {
            NotificationCenter.Error($"Save failed: {ex.Message}");
        }
    }

    private bool CanSave() => IsDocumentLoaded && PendingCount > 0;
}
