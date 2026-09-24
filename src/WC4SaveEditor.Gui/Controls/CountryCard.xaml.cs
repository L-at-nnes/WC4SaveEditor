using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Controls;

public partial class CountryCard : UserControl
{
    public static readonly DependencyProperty IsDraggableProperty =
        DependencyProperty.Register(nameof(IsDraggable), typeof(bool), typeof(CountryCard), new PropertyMetadata(false));

    public static readonly DependencyProperty CardDataProperty =
        DependencyProperty.Register(nameof(CardData), typeof(CountryCardViewModel), typeof(CountryCard), new PropertyMetadata(null));

    public CountryCardViewModel? CardData
    {
        get => (CountryCardViewModel?)GetValue(CardDataProperty);
        set => SetValue(CardDataProperty, value);
    }

    public static readonly DependencyProperty IsCheckedProperty =
        DependencyProperty.Register(nameof(IsChecked), typeof(bool), typeof(CountryCard), new PropertyMetadata(false, OnHighlightChanged));

    public static readonly DependencyProperty IsSelectedProperty =
        DependencyProperty.Register(nameof(IsSelected), typeof(bool), typeof(CountryCard), new PropertyMetadata(false, OnHighlightChanged));

    public bool IsDraggable
    {
        get => (bool)GetValue(IsDraggableProperty);
        set => SetValue(IsDraggableProperty, value);
    }

    public bool IsChecked
    {
        get => (bool)GetValue(IsCheckedProperty);
        set => SetValue(IsCheckedProperty, value);
    }

    public bool IsSelected
    {
        get => (bool)GetValue(IsSelectedProperty);
        set => SetValue(IsSelectedProperty, value);
    }

    private static readonly Brush AccentBorder = (Brush)Application.Current.Resources["AccentBrush"];
    private static readonly Brush DefaultBorder = (Brush)Application.Current.Resources["BorderBrush0"];
    private static readonly Brush HighlightBackground = new SolidColorBrush(Color.FromRgb(0x2A, 0x23, 0x18));
    private static readonly Brush DefaultBackground = (Brush)Application.Current.Resources["SurfaceRaisedBrush"];

    private Point _dragStart;

    public CountryCard()
    {
        InitializeComponent();
    }

    private static void OnHighlightChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
    {
        if (d is not CountryCard card)
        {
            return;
        }
        var highlighted = card.IsChecked || card.IsSelected;
        card.Root.BorderBrush = highlighted ? AccentBorder : DefaultBorder;
        card.Root.Background = highlighted ? HighlightBackground : DefaultBackground;
    }

    private void Root_PreviewMouseLeftButtonDown(object sender, MouseButtonEventArgs e) => _dragStart = e.GetPosition(null);

    private void Root_MouseMove(object sender, MouseEventArgs e)
    {
        if (!IsDraggable || e.LeftButton != MouseButtonState.Pressed || CardData is null)
        {
            return;
        }

        var current = e.GetPosition(null);
        if (Math.Abs(current.X - _dragStart.X) < SystemParameters.MinimumHorizontalDragDistance &&
            Math.Abs(current.Y - _dragStart.Y) < SystemParameters.MinimumVerticalDragDistance)
        {
            return;
        }

        DragDrop.DoDragDrop(this, CardData, DragDropEffects.Move);
    }
}
