namespace WC4SaveEditor.Core.Operations;

public sealed class ConversionStats
{
    public int TotalChanged { get; private set; }
    public Dictionary<byte, int> ChangesByPlayer { get; } = [];

    public void AddChange(byte oldPlayer, byte newPlayer, int x, int y)
    {
        TotalChanged++;
        ChangesByPlayer[oldPlayer] = ChangesByPlayer.GetValueOrDefault(oldPlayer) + 1;
    }
}
