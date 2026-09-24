using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class MapPanel : UserControl
{
    private Point _dragStart;
    private Point _translateStart;
    private bool _isDragging;
    private bool _didPan;

    public MapPanel()
    {
        InitializeComponent();
    }

    private void MapHost_MouseWheel(object sender, MouseWheelEventArgs e)
    {
        var cursor = e.GetPosition(MapHost);
        var oldScale = MapScale.ScaleX;
        var newScale = Math.Clamp(oldScale * (e.Delta > 0 ? 1.15 : 1 / 1.15), 0.3, 6.0);

        var canvasPoint = new Point((cursor.X - MapTranslate.X) / oldScale, (cursor.Y - MapTranslate.Y) / oldScale);
        MapScale.ScaleX = newScale;
        MapScale.ScaleY = newScale;
        MapTranslate.X = cursor.X - canvasPoint.X * newScale;
        MapTranslate.Y = cursor.Y - canvasPoint.Y * newScale;

        e.Handled = true;
    }

    private void MapHost_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
    {
        _dragStart = e.GetPosition(MapHost);
        _translateStart = new Point(MapTranslate.X, MapTranslate.Y);
        _isDragging = true;
        _didPan = false;
        MapHost.CaptureMouse();
    }

    private void MapHost_MouseMove(object sender, MouseEventArgs e)
    {
        if (!_isDragging || e.LeftButton != MouseButtonState.Pressed)
        {
            return;
        }

        var current = e.GetPosition(MapHost);
        var dx = current.X - _dragStart.X;
        var dy = current.Y - _dragStart.Y;

        if (!_didPan && Math.Abs(dx) < 4 && Math.Abs(dy) < 4)
        {
            return;
        }

        _didPan = true;
        MapTranslate.X = _translateStart.X + dx;
        MapTranslate.Y = _translateStart.Y + dy;
    }

    private void MapHost_MouseLeftButtonUp(object sender, MouseButtonEventArgs e)
    {
        _isDragging = false;
        MapHost.ReleaseMouseCapture();

        if (_didPan)
        {
            return;
        }
        if (DataContext is not MapPanelViewModel viewModel)
        {
            return;
        }

        var hit = VisualTreeHelper.HitTest(MapCanvas, e.GetPosition(MapCanvas));
        if (hit?.VisualHit is FrameworkElement { DataContext: TileViewModel tile })
        {
            viewModel.SelectTile(tile);
        }
    }
}
