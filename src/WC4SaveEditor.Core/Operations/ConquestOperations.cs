using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Operations;

public static class ConquestOperations
{
    public const byte TileUnowned = 255;
    public const byte MainPlayer = 0;
    public const int LandmineOwnerOffset = 2;

    public static ConversionStats ConvertPlayer(SaveDocument doc, int oldPlayer, int newPlayer)
    {
        if (oldPlayer < 0 || newPlayer < 0)
        {
            throw new ArgumentException("Player IDs cannot be negative");
        }
        if (oldPlayer >= doc.Players.Count)
        {
            throw new ArgumentException($"Old player ID {oldPlayer} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }
        if (newPlayer >= doc.Players.Count)
        {
            throw new ArgumentException($"New player ID {newPlayer} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }
        if (oldPlayer == newPlayer)
        {
            throw new ArgumentException("Old player and new player cannot be the same");
        }

        var stats = new ConversionStats();

        for (var row = 0; row < doc.UnitOwnerData.Length; row++)
        {
            for (var col = 0; col < doc.UnitOwnerData[row].Length; col++)
            {
                var owner = doc.UnitOwnerData[row][col];
                if (owner == TileUnowned || owner == MainPlayer || owner != (byte)oldPlayer)
                {
                    continue;
                }

                doc.UnitOwnerData[row][col] = (byte)newPlayer;
                stats.AddChange(owner, (byte)newPlayer, col, row);
            }
        }

        WriteAllUnitOwners(doc);

        var destinationPositions = CollectPlayerPositions(doc, newPlayer);
        ConvertUnitOwnershipMetadata(doc, destinationPositions);
        ConvertLandmineOwners(doc, oldPlayer, newPlayer);

        return stats;
    }

    public static (ConversionStats Stats, int MinesConverted) ConvertProvince(SaveDocument doc, ushort provinceCode, int newPlayer)
    {
        if (newPlayer < 0 || newPlayer >= doc.Players.Count)
        {
            throw new ArgumentException($"New player ID {newPlayer} does not exist");
        }

        var provinceExists = false;
        var positions = new HashSet<ushort>();
        var stats = new ConversionStats();
        var gameMode = doc.Header.GameMode;

        var attachedPorts = new List<CityData>();
        foreach (var city in doc.Cities)
        {
            if (!ProvinceLookup.IsPort(city.BuildingType))
            {
                continue;
            }
            var (portProvince, attached) = ProvinceLookup.ProvinceCodeForPort(doc, city);
            if (attached && portProvince == provinceCode)
            {
                attachedPorts.Add(city);
            }
        }

        for (var row = 0; row < doc.CityTiles.Length; row++)
        {
            for (var col = 0; col < doc.CityTiles[row].Length; col++)
            {
                if (doc.CityTiles[row][col] != provinceCode)
                {
                    continue;
                }
                provinceExists = true;
                var oldOwner = doc.UnitOwnerData[row][col];
                if (oldOwner == TileUnowned || oldOwner == (byte)newPlayer)
                {
                    continue;
                }
                doc.UnitOwnerData[row][col] = (byte)newPlayer;
                var unitCode = CoordinateMath.ConvertToCoordinateCode(row, col, (int)doc.Header.MapWidth, gameMode);
                positions.Add((ushort)unitCode);
                stats.AddChange(oldOwner, (byte)newPlayer, col, row);
            }
        }

        if (!provinceExists)
        {
            throw new ArgumentException($"Province coordinate code {provinceCode} does not exist");
        }

        foreach (var city in attachedPorts)
        {
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, gameMode);
            if (!valid)
            {
                continue;
            }
            var oldOwner = doc.UnitOwnerData[row][col];
            if (oldOwner == TileUnowned || oldOwner == (byte)newPlayer)
            {
                continue;
            }
            doc.UnitOwnerData[row][col] = (byte)newPlayer;
            positions.Add(city.CoordinateCode);
            stats.AddChange(oldOwner, (byte)newPlayer, col, row);
        }

        WriteAllUnitOwners(doc);
        ConvertUnitOwnershipMetadata(doc, positions);

        var minesConverted = 0;
        for (var i = 0; i < doc.Landmines.Count; i++)
        {
            var mine = doc.Landmines[i];
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(mine.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, gameMode);
            if (!valid || doc.CityTiles[row][col] != provinceCode)
            {
                continue;
            }
            var newMineOwner = (ushort)(newPlayer + 1);
            if (mine.Owner == newMineOwner)
            {
                continue;
            }
            var offset = doc.Offsets.LandmineOffsets[i];
            SaveFileWriter.WriteUInt16At(doc, offset + LandmineOwnerOffset, newMineOwner);
            mine.Owner = newMineOwner;
            minesConverted++;
        }

        return (stats, minesConverted);
    }

    public static void ConvertTile(SaveDocument doc, int x, int y, int newPlayer)
    {
        if (x < 0 || y < 0)
        {
            throw new ArgumentException($"Coordinates cannot be negative: ({x}, {y})");
        }
        if (y >= doc.UnitOwnerData.Length || x >= doc.UnitOwnerData[0].Length)
        {
            throw new ArgumentException($"Coordinates out of bounds: ({x}, {y}). Map size: {doc.UnitOwnerData[0].Length}x{doc.UnitOwnerData.Length}");
        }
        if (newPlayer < 0 || newPlayer >= doc.Players.Count)
        {
            throw new ArgumentException($"New player ID {newPlayer} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        var oldPlayer = doc.UnitOwnerData[y][x];
        if (oldPlayer == TileUnowned)
        {
            throw new InvalidOperationException($"Can't convert tile at ({y}, {x}) - tile has no owner");
        }

        doc.UnitOwnerData[y][x] = (byte)newPlayer;
        SaveFileWriter.WriteByteAt(doc, doc.Offsets.TileOffset(y, x), (byte)newPlayer);
    }

    public static void ChangePlayerTeam(SaveDocument doc, int playerId, uint teamId)
    {
        if (playerId < 0 || playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist");
        }

        var offset = doc.Offsets.PlayerOffsets[playerId];
        SaveFileWriter.WriteUInt32At(doc, offset + CountryData.TeamIdFieldOffset, teamId);
        doc.Players[playerId].TeamId = teamId;
    }

    public static int ConvertTeam(SaveDocument doc)
    {
        var mainTeamId = doc.Players[MainPlayer].TeamId;
        var convertedCount = 0;

        for (var i = 1; i < doc.Players.Count; i++)
        {
            if (doc.Players[i].TeamId != mainTeamId)
            {
                ChangePlayerTeam(doc, i, mainTeamId);
                convertedCount++;
            }
        }

        return convertedCount;
    }

    public static ConversionStats ConvertAllPlayers(SaveDocument doc)
    {
        var stats = new ConversionStats();

        for (var row = 0; row < doc.UnitOwnerData.Length; row++)
        {
            for (var col = 0; col < doc.UnitOwnerData[row].Length; col++)
            {
                var owner = doc.UnitOwnerData[row][col];
                if (owner != TileUnowned && owner != MainPlayer)
                {
                    doc.UnitOwnerData[row][col] = MainPlayer;
                    stats.AddChange(owner, MainPlayer, col, row);
                }
            }
        }

        WriteAllUnitOwners(doc);
        return stats;
    }

    public static ConversionStats ConvertAllAllies(SaveDocument doc)
    {
        var stats = new ConversionStats();
        var mainTeamId = doc.Players[MainPlayer].TeamId;

        for (var row = 0; row < doc.UnitOwnerData.Length; row++)
        {
            for (var col = 0; col < doc.UnitOwnerData[row].Length; col++)
            {
                var owner = doc.UnitOwnerData[row][col];
                if (owner == TileUnowned || owner == MainPlayer)
                {
                    continue;
                }
                if (doc.Players[owner].TeamId == mainTeamId)
                {
                    doc.UnitOwnerData[row][col] = MainPlayer;
                    stats.AddChange(owner, MainPlayer, col, row);
                }
            }
        }

        WriteAllUnitOwners(doc);
        return stats;
    }

    private static void WriteAllUnitOwners(SaveDocument doc)
    {
        var flattened = new byte[doc.UnitOwnerData.Length * doc.UnitOwnerData[0].Length];
        var i = 0;
        foreach (var row in doc.UnitOwnerData)
        {
            row.CopyTo(flattened, i);
            i += row.Length;
        }
        SaveFileWriter.WriteBytesAt(doc, doc.Offsets.UnitOwnerGridStart, flattened);
    }

    private static HashSet<ushort> CollectPlayerPositions(SaveDocument doc, int playerId)
    {
        var positions = new HashSet<ushort>();
        for (var row = 0; row < doc.UnitOwnerData.Length; row++)
        {
            for (var col = 0; col < doc.UnitOwnerData[row].Length; col++)
            {
                if (doc.UnitOwnerData[row][col] != (byte)playerId)
                {
                    continue;
                }
                var code = CoordinateMath.ConvertToCoordinateCode(row, col, (int)doc.Header.MapWidth, doc.Header.GameMode);
                positions.Add((ushort)code);
            }
        }
        return positions;
    }

    private static int ConvertUnitOwnershipMetadata(SaveDocument doc, HashSet<ushort> convertedPositions)
    {
        var activationTurn = (byte)doc.Header.TurnNumber;
        var converted = 0;

        for (var i = 0; i < doc.Units.Count; i++)
        {
            var unit = doc.Units[i];
            if (!convertedPositions.Contains(unit.CoordinateCode))
            {
                continue;
            }

            var offset = doc.Offsets.UnitOffsets[i];
            SaveFileWriter.WriteByteAt(doc, offset + UnitData.ActivationTurnOffset, activationTurn);
            SaveFileWriter.WriteByteAt(doc, offset + UnitData.ActionStateOffset, UnitData.ReadyActionState);
            unit.UnknownArr5[10] = activationTurn;
            unit.UnknownArr6[0] = UnitData.ReadyActionState;
            converted++;
        }

        return converted;
    }

    private static int ConvertLandmineOwners(SaveDocument doc, int oldPlayer, int newPlayer)
    {
        var oldMineOwner = (ushort)(oldPlayer + 1);
        var newMineOwner = (ushort)(newPlayer + 1);
        var converted = 0;

        for (var i = 0; i < doc.Landmines.Count; i++)
        {
            var mine = doc.Landmines[i];
            if (mine.Owner != oldMineOwner)
            {
                continue;
            }

            var offset = doc.Offsets.LandmineOffsets[i];
            SaveFileWriter.WriteUInt16At(doc, offset + LandmineOwnerOffset, newMineOwner);
            mine.Owner = newMineOwner;
            converted++;
        }

        return converted;
    }
}
