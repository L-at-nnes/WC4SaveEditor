using System.Collections.ObjectModel;
using System.Windows;
using System.Windows.Media;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Core.Rendering;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class MapPanelViewModel : ObservableObject
{
    private const double HexSize = 11.0;

    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private SaveDocument? _document;
    private TileViewModel[][] _tileGrid = [];

    [ObservableProperty]
    private ObservableCollection<TileViewModel> _tiles = [];

    [ObservableProperty]
    private double _canvasWidth;

    [ObservableProperty]
    private double _canvasHeight;

    [ObservableProperty]
    private TileViewModel? _selectedTile;

    [ObservableProperty]
    private ObservableCollection<CountryCardViewModel> _players = [];

    [ObservableProperty]
    private CountryCardViewModel? _pickedOwner;

    [ObservableProperty]
    private bool _isFlyoutOpen;

    public MapPanelViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
    }

    public void Reset(SaveDocument document)
    {
        _document = document;
        IsFlyoutOpen = false;
        SelectedTile = null;
        Players = new ObservableCollection<CountryCardViewModel>(document.Players.Select((p, i) => new CountryCardViewModel(i, p)));

        var height = document.UnitOwnerData.Length;
        var width = height == 0 ? 0 : document.UnitOwnerData[0].Length;

        var tiles = new ObservableCollection<TileViewModel>();
        var grid = new TileViewModel[height][];
        double maxX = 0, maxY = 0;

        for (var row = 0; row < height; row++)
        {
            grid[row] = new TileViewModel[width];
            for (var col = 0; col < width; col++)
            {
                var (cx, cy) = HexLayout.GetHexCenter(col, row, HexSize);
                var hexPoints = HexLayout.GetHexPoints(cx, cy, HexSize);
                var points = new PointCollection(hexPoints.Select(p => new Point(p.X, p.Y)));

                var tile = new TileViewModel(row, col, points);
                tiles.Add(tile);
                grid[row][col] = tile;

                maxX = Math.Max(maxX, cx + HexSize);
                maxY = Math.Max(maxY, cy + HexSize);
            }
        }

        Tiles = tiles;
        _tileGrid = grid;
        CanvasWidth = maxX + HexSize;
        CanvasHeight = maxY + HexSize;

        RefreshPreview();
    }

    /// <summary>
    /// Recolors every tile from a fresh copy of the save file with all currently queued
    /// changes (from every panel, not just this one) applied to it, so the map always shows
    /// what Save would actually produce, without touching the document being edited.
    /// </summary>
    private static readonly Dictionary<string, SolidColorBrush> BrushCache = [];

    public void RefreshPreview()
    {
        if (_document is null || _tileGrid.Length == 0)
        {
            return;
        }

        var hasStructuralChanges = _queue.Pending.Any(c => c.IsStructural);
        var preview = _queue.Count == 0 ? _document : BuildPreviewDocument();

        var playerColors = TileColoring.BuildPlayerColors(preview);
        var cityOwnership = TileColoring.DetermineCityOwnership(preview);

        for (var row = 0; row < preview.UnitOwnerData.Length; row++)
        {
            for (var col = 0; col < preview.UnitOwnerData[row].Length; col++)
            {
                var owner = preview.UnitOwnerData[row][col];
                var tile = _tileGrid[row][col];

                // Spawning can claim a previously-neutral tile without changing its byte value
                // relative to last paint in some edge cases, so only fully skip the cheap path
                // when nothing structural is pending.
                if (!hasStructuralChanges && tile.LastPaintedOwner == owner)
                {
                    continue;
                }
                tile.LastPaintedOwner = owner;

                var colorHex = TileColoring.GetTileColor(preview, row, col, cityOwnership, playerColors);
                if (!BrushCache.TryGetValue(colorHex, out var brush))
                {
                    brush = new SolidColorBrush((Color)ColorConverter.ConvertFromString(colorHex));
                    brush.Freeze();
                    BrushCache[colorHex] = brush;
                }

                tile.Fill = brush;
                tile.HasPendingChange = owner != _document.UnitOwnerData[row][col];
            }
        }
    }

    private SaveDocument BuildPreviewDocument()
    {
        var preview = _document!.Clone();
        _queue.ApplyToPreview(preview);
        return preview;
    }

    public void SelectTile(TileViewModel tile)
    {
        if (_document is null)
        {
            return;
        }

        var currentOwner = _document.UnitOwnerData[tile.Row][tile.Col];
        if (currentOwner == ConquestOperations.TileUnowned)
        {
            NotificationCenter.Error("This tile has no owner - nothing to reassign.");
            return;
        }

        SelectedTile = tile;
        PickedOwner = Players.FirstOrDefault(p => p.PlayerId == currentOwner);
        IsFlyoutOpen = true;
    }

    [RelayCommand]
    private void ConfirmReassign()
    {
        if (SelectedTile is null || PickedOwner is null)
        {
            IsFlyoutOpen = false;
            return;
        }

        _queue.RemoveAll(c => c is ConvertTileChange tc && tc.X == SelectedTile.Col && tc.Y == SelectedTile.Row);

        var currentOwner = _document!.UnitOwnerData[SelectedTile.Row][SelectedTile.Col];
        if (PickedOwner.PlayerId != currentOwner)
        {
            _queue.Add(new ConvertTileChange(SelectedTile.Col, SelectedTile.Row, PickedOwner.PlayerId));
            NotificationCenter.Info($"Will reassign tile ({SelectedTile.Row},{SelectedTile.Col}) to {PickedOwner.Name} on Save.");
        }

        IsFlyoutOpen = false;
        _notifyChanged();
    }

    [RelayCommand]
    private void CancelReassign() => IsFlyoutOpen = false;
}
