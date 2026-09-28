using System.IO;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using CommunityToolkit.Mvvm.ComponentModel;
using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class CountryCardViewModel : ObservableObject
{
    public int PlayerId { get; }
    public string Name { get; }
    public ImageSource? FlagSource { get; }
    public Brush ColorBrush { get; }
    public bool IsAlive { get; }

    [ObservableProperty]
    private uint _currentTeamId;

    [ObservableProperty]
    private bool _isSelected;

    public uint OriginalTeamId { get; }

    public CountryCardViewModel(int playerId, CountryData data, SaveDocument document)
    {
        PlayerId = playerId;
        Name = Countries.TryGetName((byte)data.CountryId, out var name) ? name : $"Country {data.CountryId}";
        FlagSource = TryLoadFlag(data.CountryId);
        var hex = Countries.ColorBytesToHex(data.PrimaryColor);
        ColorBrush = new SolidColorBrush((Color)ColorConverter.ConvertFromString("#" + hex));
        OriginalTeamId = data.TeamId;
        CurrentTeamId = data.TeamId;
        IsAlive = playerId == ConquestOperations.MainPlayer || ConquestOperations.PlayerHasTerritory(document, playerId);
    }

    private static ImageSource? TryLoadFlag(uint countryId)
    {
        try
        {
            var uri = new Uri($"pack://application:,,,/Assets/Flags/{countryId}.png", UriKind.Absolute);
            return new BitmapImage(uri);
        }
        catch (IOException)
        {
            return null;
        }
    }
}
