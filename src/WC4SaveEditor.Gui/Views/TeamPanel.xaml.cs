using System.Windows;
using System.Windows.Controls;
using WC4SaveEditor.Gui.ViewModels;

namespace WC4SaveEditor.Gui.Views;

public partial class TeamPanel : UserControl
{
    public TeamPanel()
    {
        InitializeComponent();
    }

    private void Column_DragOver(object sender, DragEventArgs e)
    {
        e.Effects = e.Data.GetDataPresent(typeof(CountryCardViewModel)) ? DragDropEffects.Move : DragDropEffects.None;
        e.Handled = true;
    }

    private void Column_Drop(object sender, DragEventArgs e)
    {
        if (DataContext is not TeamPanelViewModel viewModel)
        {
            return;
        }
        if (sender is not FrameworkElement { DataContext: TeamColumnViewModel column })
        {
            return;
        }
        if (e.Data.GetData(typeof(CountryCardViewModel)) is not CountryCardViewModel card)
        {
            return;
        }

        viewModel.MoveCard(card, column.TeamId);
    }
}
