using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Operations;

public static class StatOperations
{
    public const int CityLevelBase = 10;
    public const int CityLevelCapital = 5;
    public const byte MaxMoraleValue = 1;
    public const ushort MaxMoraleTurns = 3;

    public static void SetPlayerMaxCurrency(SaveDocument doc, int playerId, int currencyValue)
    {
        if (playerId < 0)
        {
            throw new ArgumentException("Player ID cannot be negative");
        }
        if (currencyValue is < 0 or > 65535)
        {
            throw new ArgumentException($"Currency value {currencyValue} is invalid. Must be between 0 and 65535");
        }
        if (playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        var offset = doc.Offsets.PlayerOffsets[playerId];
        var currencyStartOffset = offset + CountryData.CurrencyFieldOffset;
        var value = (uint)currencyValue;

        for (var i = 0; i < 3; i++)
        {
            SaveFileWriter.WriteUInt32At(doc, currencyStartOffset + i * 4, value);
        }
        doc.Players[playerId].Currency[0] = value;
        doc.Players[playerId].Currency[1] = value;
        doc.Players[playerId].Currency[2] = value;
    }

    public static int SetPlayerMaxCityTech(SaveDocument doc, int playerId, int techLevel)
    {
        if (playerId < 0)
        {
            throw new ArgumentException("Player ID cannot be negative");
        }
        if (techLevel is < 0 or > 255)
        {
            throw new ArgumentException($"Tech level {techLevel} is invalid. Must be between 0 and 255");
        }
        if (playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        return ProcessPlayerCityTech(doc, playerId, techLevel);
    }

    public static int SetPlayerMaxCityLevel(SaveDocument doc, int playerId, int targetLevel)
    {
        if (playerId < 0)
        {
            throw new ArgumentException("Player ID cannot be negative");
        }
        if (targetLevel < 1 || targetLevel >= CityLevelCapital)
        {
            throw new ArgumentException($"City level {targetLevel} is invalid. Must be between 1 and {CityLevelCapital - 1}");
        }
        if (playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        return ProcessPlayerCityLevel(doc, playerId, targetLevel);
    }

    public static int SetPlayerMaxMorale(SaveDocument doc, int playerId)
    {
        if (playerId < 0 || playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        var updated = 0;
        for (var i = 0; i < doc.Units.Count; i++)
        {
            var unit = doc.Units[i];
            if (unit.UnitType is >= UnitTypes.Id.Bunker and <= UnitTypes.Id.City)
            {
                continue;
            }
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(unit.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid || doc.UnitOwnerData[row][col] != (byte)playerId)
            {
                continue;
            }

            var offset = doc.Offsets.UnitOffsets[i];
            SaveFileWriter.WriteByteAt(doc, offset + UnitData.MoraleValueFieldOffset, MaxMoraleValue);
            SaveFileWriter.WriteUInt16At(doc, offset + UnitData.MoraleTurnsLeftFieldOffset, MaxMoraleTurns);
            unit.MoraleValue = (sbyte)MaxMoraleValue;
            unit.MoraleTurnsLeft = MaxMoraleTurns;
            updated++;
        }

        return updated;
    }

    public static int RestoreAlliesHealth(SaveDocument doc, int playerId)
    {
        if (playerId < 0)
        {
            throw new ArgumentException("Player ID cannot be negative");
        }
        if (playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        var playerTeamId = doc.Players[playerId].TeamId;
        var restored = 0;

        foreach (var (index, unit, row, col, owner) in ProcessAllUnits(doc))
        {
            if (owner >= doc.Players.Count || doc.Players[owner].TeamId != playerTeamId)
            {
                continue;
            }

            var offset = doc.Offsets.UnitOffsets[index] + UnitData.CurrentHealthFieldOffset;
            SaveFileWriter.WriteUInt16At(doc, offset, unit.MaxHealth);
            unit.CurrentHealth = unit.MaxHealth;
            restored++;
        }

        return restored;
    }

    public static (int EnemiesWeakened, int CitiesWeakened, int UnitsWeakened) WeakenEnemies(SaveDocument doc, int playerId)
    {
        if (playerId < 0)
        {
            throw new ArgumentException("Player ID cannot be negative");
        }
        if (playerId >= doc.Players.Count)
        {
            throw new ArgumentException($"Player ID {playerId} does not exist. Valid range: 0-{doc.Players.Count - 1}");
        }

        var playerTeamId = doc.Players[playerId].TeamId;

        var enemiesWeakened = 0;
        for (var i = 0; i < doc.Players.Count; i++)
        {
            if (i == playerId || doc.Players[i].TeamId == playerTeamId)
            {
                continue;
            }
            SetPlayerMaxCurrency(doc, i, 0);
            enemiesWeakened++;
        }

        var citiesWeakened = ProcessEnemyCityTech(doc, playerId, playerTeamId, 0);

        var unitsWeakened = 0;
        foreach (var (index, unit, _, _, owner) in ProcessAllUnits(doc))
        {
            if (owner >= doc.Players.Count || doc.Players[owner].TeamId == playerTeamId)
            {
                continue;
            }

            var newHealth = (ushort)(unit.UnitType == UnitTypes.Id.City ? 0 : 1);
            var offset = doc.Offsets.UnitOffsets[index] + UnitData.CurrentHealthFieldOffset;
            SaveFileWriter.WriteUInt16At(doc, offset, newHealth);
            unit.CurrentHealth = newHealth;
            unitsWeakened++;
        }

        return (enemiesWeakened, citiesWeakened, unitsWeakened);
    }

    private static int ProcessPlayerCityTech(SaveDocument doc, int playerId, int techLevel)
    {
        var citiesModified = 0;
        for (var i = 0; i < doc.Cities.Count; i++)
        {
            var city = doc.Cities[i];
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid)
            {
                continue;
            }

            var owner = doc.UnitOwnerData[row][col];
            if (owner != (byte)playerId)
            {
                continue;
            }

            if (ProvinceLookup.IsPort(city.BuildingType))
            {
                UpdatePortLevel(doc, i);
            }
            else
            {
                UpdateCityTechLevels(doc, i, techLevel);
            }

            citiesModified++;
        }
        return citiesModified;
    }

    private static int ProcessPlayerCityLevel(SaveDocument doc, int playerId, int targetLevel)
    {
        var citiesModified = 0;
        for (var i = 0; i < doc.Cities.Count; i++)
        {
            var city = doc.Cities[i];
            if (ProvinceLookup.IsPort(city.BuildingType))
            {
                continue;
            }

            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid)
            {
                continue;
            }

            var owner = doc.UnitOwnerData[row][col];
            if (owner != (byte)playerId)
            {
                continue;
            }

            var currentLevel = city.BuildingType - CityLevelBase;
            if (currentLevel >= targetLevel)
            {
                continue;
            }

            var offset = doc.Offsets.CityOffsets[i];
            var newBuildingType = (byte)(CityLevelBase + targetLevel);
            SaveFileWriter.WriteByteAt(doc, offset + CityData.BuildingTypeFieldOffset, newBuildingType);
            city.BuildingType = newBuildingType;

            citiesModified++;
        }
        return citiesModified;
    }

    private static int ProcessEnemyCityTech(SaveDocument doc, int playerId, uint playerTeamId, int techLevel)
    {
        var citiesModified = 0;
        for (var i = 0; i < doc.Cities.Count; i++)
        {
            var city = doc.Cities[i];
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid)
            {
                continue;
            }

            var owner = doc.UnitOwnerData[row][col];
            if (owner >= doc.Players.Count)
            {
                continue;
            }
            if (owner == playerId || doc.Players[owner].TeamId == playerTeamId)
            {
                continue;
            }

            UpdateCityTechLevels(doc, i, techLevel);
            citiesModified++;
        }
        return citiesModified;
    }

    private static void UpdatePortLevel(SaveDocument doc, int cityIndex)
    {
        var offset = doc.Offsets.CityOffsets[cityIndex];
        SaveFileWriter.WriteByteAt(doc, offset + CityData.BuildingTypeFieldOffset, ProvinceLookup.PortLevel4);
        SaveFileWriter.WriteBytesAt(doc, offset + CityData.TechLevelsFieldOffset, new byte[6]);

        var city = doc.Cities[cityIndex];
        city.BuildingType = ProvinceLookup.PortLevel4;
        for (var j = 0; j < 6; j++)
        {
            city.TechLevels[j] = 0;
        }
    }

    private static void UpdateCityTechLevels(SaveDocument doc, int cityIndex, int techLevel)
    {
        var offset = doc.Offsets.CityOffsets[cityIndex];
        var techData = new byte[6];
        for (var j = 0; j < 6; j++)
        {
            techData[j] = (byte)techLevel;
        }
        SaveFileWriter.WriteBytesAt(doc, offset + CityData.TechLevelsFieldOffset, techData);

        var city = doc.Cities[cityIndex];
        for (var j = 0; j < 6; j++)
        {
            city.TechLevels[j] = (byte)techLevel;
        }
    }

    private static IEnumerable<(int Index, UnitData Unit, int Row, int Col, byte Owner)> ProcessAllUnits(SaveDocument doc)
    {
        for (var i = 0; i < doc.Units.Count; i++)
        {
            var unit = doc.Units[i];
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(unit.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid)
            {
                continue;
            }
            yield return (i, unit, row, col, doc.UnitOwnerData[row][col]);
        }
    }
}
