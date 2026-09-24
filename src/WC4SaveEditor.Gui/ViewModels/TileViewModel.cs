using System.Windows;
using System.Windows.Media;
using CommunityToolkit.Mvvm.ComponentModel;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class TileViewModel(int row, int col, PointCollection points) : ObservableObject
{
    public int Row { get; } = row;
    public int Col { get; } = col;
    public PointCollection Points { get; } = points;

    [ObservableProperty]
    private Brush _fill = Brushes.Gray;

    [ObservableProperty]
    private bool _hasPendingChange;

    /// <summary>Owner last painted, used to skip recoloring tiles whose owner hasn't changed.</summary>
    public byte? LastPaintedOwner { get; set; }
}
