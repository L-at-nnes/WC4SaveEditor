using WC4SaveEditor.Core.Models;

namespace WC4SaveEditor.Core.Operations;

public abstract class PendingChange
{
    public abstract bool IsStructural { get; }
    public abstract string Description { get; }
    public abstract void Apply(SaveDocument doc);
}

public sealed class MaxMoneyChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Max money";
    public override void Apply(SaveDocument doc) => StatOperations.SetPlayerMaxCurrency(doc, ConquestOperations.MainPlayer, 9999);
}

public sealed class MaxCityTechChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Max city tech";
    public override void Apply(SaveDocument doc) => StatOperations.SetPlayerMaxCityTech(doc, ConquestOperations.MainPlayer, 4);
}

public sealed class MaxCityLevelChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Max city level";
    public override void Apply(SaveDocument doc) => StatOperations.SetPlayerMaxCityLevel(doc, ConquestOperations.MainPlayer, 4);
}

public sealed class MaxMoraleChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Max morale";
    public override void Apply(SaveDocument doc) => StatOperations.SetPlayerMaxMorale(doc, ConquestOperations.MainPlayer);
}

public sealed class HealAlliesChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Heal allies";
    public override void Apply(SaveDocument doc) => StatOperations.RestoreAlliesHealth(doc, ConquestOperations.MainPlayer);
}

public sealed class WeakenEnemiesChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Weaken enemies";
    public override void Apply(SaveDocument doc) => StatOperations.WeakenEnemies(doc, ConquestOperations.MainPlayer);
}

public sealed class ConquerAllChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Conquer all";
    public override void Apply(SaveDocument doc) => ConquestOperations.ConvertAllPlayers(doc);
}

public sealed class UniteTeamChange : PendingChange
{
    public override bool IsStructural => false;
    public override string Description => "Unite team";
    public override void Apply(SaveDocument doc) => ConquestOperations.ConvertTeam(doc);
}

public sealed class ConquerPlayerChange(int oldPlayer, int newPlayer) : PendingChange
{
    public int OldPlayer { get; } = oldPlayer;
    public int NewPlayer { get; } = newPlayer;
    public override bool IsStructural => false;
    public override string Description => $"Conquer player {OldPlayer} -> {NewPlayer}";
    public override void Apply(SaveDocument doc) => ConquestOperations.ConvertPlayer(doc, OldPlayer, NewPlayer);
}

public sealed class ConquerProvinceChange(ushort provinceCode, int newPlayer) : PendingChange
{
    public ushort ProvinceCode { get; } = provinceCode;
    public int NewPlayer { get; } = newPlayer;
    public override bool IsStructural => false;
    public override string Description => $"Conquer province {ProvinceCode} -> {NewPlayer}";
    public override void Apply(SaveDocument doc) => ConquestOperations.ConvertProvince(doc, ProvinceCode, NewPlayer);
}

public sealed class ChangeTeamChange(int playerId, uint teamId) : PendingChange
{
    public int PlayerId { get; } = playerId;
    public uint TeamId { get; } = teamId;
    public override bool IsStructural => false;
    public override string Description => $"Player {PlayerId} -> team {TeamId}";
    public override void Apply(SaveDocument doc) => ConquestOperations.ChangePlayerTeam(doc, PlayerId, TeamId);
}

public sealed class ConvertTileChange(int x, int y, int newPlayer) : PendingChange
{
    public int X { get; } = x;
    public int Y { get; } = y;
    public int NewPlayer { get; } = newPlayer;
    public override bool IsStructural => false;
    public override string Description => $"Tile ({X},{Y}) -> {NewPlayer}";
    public override void Apply(SaveDocument doc) => ConquestOperations.ConvertTile(doc, X, Y, NewPlayer);
}

public sealed class SpawnUnitChange(ushort provinceCode, byte unitType, int level, int count) : PendingChange
{
    public ushort ProvinceCode { get; } = provinceCode;
    public byte UnitType { get; } = unitType;
    public int Level { get; } = level;
    public int Count { get; } = count;
    public override bool IsStructural => true;
    public override string Description => $"Spawn {Count}x type {UnitType} (lvl {Level}) in province {ProvinceCode}";
    public override void Apply(SaveDocument doc) => SpawnUnitsOperations.SpawnUnits(doc, ProvinceCode, UnitType, Level, Count);
}
