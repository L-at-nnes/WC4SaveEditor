namespace WC4SaveEditor.Core.Data;

public static class GameModes
{
    public const uint Campaign = 1;
    public const uint Conquest = 2;
    public const uint Frontier = 6;
    public const uint ArmyGroup = 10;

    public static string GetName(uint gameMode) => gameMode switch
    {
        Campaign => "Campaign",
        Conquest => "Conquest",
        Frontier => "Frontier",
        ArmyGroup => "ArmyGroup",
        _ => "Unknown"
    };
}
