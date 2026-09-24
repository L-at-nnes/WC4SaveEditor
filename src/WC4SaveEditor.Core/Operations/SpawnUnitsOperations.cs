using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Operations;

public sealed class FunSpawnResult
{
    public int LandSpawned { get; set; }
    public int NavalSpawned { get; set; }
}

public static class SpawnUnitsOperations
{
    private const int MaxCoastSearchRadius = 6;
    private static readonly Random FunRandom = new();

    public static readonly byte[] SpawnableUnitTypesList =
    [
        UnitTypes.Id.Brandenburgers, UnitTypes.Id.HawkeyeForce, UnitTypes.Id.CombatMedic, UnitTypes.Id.M7Priest, UnitTypes.Id.HeavyGustav,
        UnitTypes.Id.BM21, UnitTypes.Id.T44, UnitTypes.Id.KingTiger, UnitTypes.Id.M26Pershing, UnitTypes.Id.TypeVIISubmarine, UnitTypes.Id.HMSPrinceOfWales,
        UnitTypes.Id.B4Howitzer, UnitTypes.Id.IS3HeavyTank, UnitTypes.Id.StukaZuFuss, UnitTypes.Id.Richelieu, UnitTypes.Id.Enterprise, UnitTypes.Id.RPGRocketeer,
        UnitTypes.Id.A41Centurion, UnitTypes.Id.AuF1, UnitTypes.Id.T72, UnitTypes.Id.PhantomForce, UnitTypes.Id.M1A1Abrams, UnitTypes.Id.AH64Apache,
        UnitTypes.Id.M142Himars, UnitTypes.Id.DivineWrathMBT,
    ];

    public static bool IsNavalUnitType(byte unitType) => unitType is UnitTypes.Id.TypeVIISubmarine or UnitTypes.Id.HMSPrinceOfWales
        or UnitTypes.Id.Richelieu or UnitTypes.Id.Enterprise;

    public static IEnumerable<byte> LandUnitTypes() => SpawnableUnitTypesList.Where(t => !IsNavalUnitType(t));
    public static IEnumerable<byte> NavalUnitTypes() => SpawnableUnitTypesList.Where(IsNavalUnitType);

    /// <summary>
    /// Whether at least one existing unit of this type exists anywhere in the save to clone
    /// stats from - spawning is impossible without one. Check this before queueing a spawn.
    /// </summary>
    public static bool HasStatTemplate(List<UnitData> units, byte unitType) => units.Any(u => u.UnitType == unitType);

    public static int SpawnUnits(SaveDocument doc, ushort provinceCode, byte unitType, int level, int count)
    {
        if (count < 1)
        {
            throw new ArgumentException("Count must be at least 1");
        }
        if (level is < 1 or > 9)
        {
            throw new ArgumentException("Level must be between 1 and 9");
        }

        var template = FindStatTemplate(doc.Units, unitType, level);

        var (anchorRow, anchorCol, valid) = CoordinateMath.ConvertCoordinates(provinceCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
        if (!valid)
        {
            throw new ArgumentException($"Province coordinate code {provinceCode} is not a valid map position");
        }
        var provinceOwner = doc.UnitOwnerData[anchorRow][anchorCol];

        var positions = IsNavalUnitType(unitType)
            ? FindFreeNavalTiles(doc, provinceCode, count)
            : FindFreeLandTiles(doc, new HexPos(anchorRow, anchorCol), provinceCode, count);

        if (positions.Count < count)
        {
            throw new InvalidOperationException($"Only found {positions.Count} free tile(s) near this province for {count} requested unit(s); nothing was spawned");
        }

        ClaimTiles(doc, positions, provinceOwner);

        var currentTurn = (byte)doc.Header.TurnNumber;
        var newUnits = new List<UnitData>(count);
        for (var i = 0; i < count; i++)
        {
            var coordinateCode = (ushort)CoordinateMath.ConvertToCoordinateCode(positions[i].Row, positions[i].Col, (int)doc.Header.MapWidth, doc.Header.GameMode);
            newUnits.Add(BuildNewUnit(template, level, coordinateCode, currentTurn));
        }

        AppendNewUnits(doc, newUnits);
        return newUnits.Count;
    }

    public static FunSpawnResult SpawnFunUnits(SaveDocument doc, ushort provinceCode, int fillPercent)
    {
        const int funLevel = 9;
        fillPercent = Math.Clamp(fillPercent, 1, 100);

        var (anchorRow, anchorCol, valid) = CoordinateMath.ConvertCoordinates(provinceCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
        if (!valid)
        {
            throw new ArgumentException($"Province coordinate code {provinceCode} is not a valid map position");
        }
        var provinceOwner = doc.UnitOwnerData[anchorRow][anchorCol];

        var landTemplates = UsableTemplates(doc.Units, LandUnitTypes(), funLevel);
        var navalTemplates = UsableTemplates(doc.Units, NavalUnitTypes(), funLevel);

        var currentTurn = (byte)doc.Header.TurnNumber;
        var newUnits = new List<UnitData>();
        var claimTiles = new List<HexPos>();
        var result = new FunSpawnResult();

        if (landTemplates.Count > 0)
        {
            var positions = TryFindFreeLandTiles(doc, new HexPos(anchorRow, anchorCol), provinceCode, int.MaxValue);
            positions = TakeRandomShare(positions, fillPercent);
            foreach (var pos in positions)
            {
                var template = landTemplates[FunRandom.Next(landTemplates.Count)];
                var coordinateCode = (ushort)CoordinateMath.ConvertToCoordinateCode(pos.Row, pos.Col, (int)doc.Header.MapWidth, doc.Header.GameMode);
                newUnits.Add(BuildNewUnit(template, funLevel, coordinateCode, currentTurn));
                claimTiles.Add(pos);
            }
            result.LandSpawned = positions.Count;
        }

        if (navalTemplates.Count > 0)
        {
            var positions = TryFindFreeNavalTiles(doc, provinceCode, int.MaxValue);
            positions = TakeRandomShare(positions, fillPercent);
            foreach (var pos in positions)
            {
                var template = navalTemplates[FunRandom.Next(navalTemplates.Count)];
                var coordinateCode = (ushort)CoordinateMath.ConvertToCoordinateCode(pos.Row, pos.Col, (int)doc.Header.MapWidth, doc.Header.GameMode);
                newUnits.Add(BuildNewUnit(template, funLevel, coordinateCode, currentTurn));
                claimTiles.Add(pos);
            }
            result.NavalSpawned = positions.Count;
        }

        if (newUnits.Count == 0)
        {
            throw new InvalidOperationException("No free tile or usable unit type found near this province; nothing was spawned");
        }

        ClaimTiles(doc, claimTiles, provinceOwner);
        AppendNewUnits(doc, newUnits);
        return result;
    }

    private static List<UnitData> UsableTemplates(List<UnitData> units, IEnumerable<byte> unitTypes, int level)
    {
        var templates = new List<UnitData>();
        foreach (var unitType in unitTypes)
        {
            try
            {
                templates.Add(FindStatTemplate(units, unitType, level));
            }
            catch (InvalidOperationException)
            {
            }
        }
        return templates;
    }

    private static List<HexPos> TakeRandomShare(List<HexPos> positions, int percent)
    {
        if (percent >= 100 || positions.Count == 0)
        {
            return positions;
        }
        var count = (int)Math.Round(positions.Count * percent / 100.0);
        for (var i = positions.Count - 1; i > 0; i--)
        {
            var j = FunRandom.Next(i + 1);
            (positions[i], positions[j]) = (positions[j], positions[i]);
        }
        return positions.GetRange(0, count);
    }

    private static void ClaimTiles(SaveDocument doc, List<HexPos> positions, byte provinceOwner)
    {
        foreach (var pos in positions)
        {
            if (doc.UnitOwnerData[pos.Row][pos.Col] != provinceOwner)
            {
                doc.UnitOwnerData[pos.Row][pos.Col] = provinceOwner;
                SaveFileWriter.WriteByteAt(doc, doc.Offsets.TileOffset(pos.Row, pos.Col), provinceOwner);
            }
        }
    }

    private static UnitData FindStatTemplate(List<UnitData> units, byte unitType, int level)
    {
        UnitData? best = null;
        var bestDiff = int.MaxValue;

        foreach (var u in units)
        {
            if (u.UnitType != unitType)
            {
                continue;
            }
            var diff = Math.Abs(u.Level - level);
            if (best is null || diff < bestDiff)
            {
                best = u;
                bestDiff = diff;
            }
            if (diff == 0)
            {
                break;
            }
        }

        if (best is null)
        {
            var name = UnitTypes.TryGetName(unitType, out var n) ? n : $"Type {unitType}";
            throw new InvalidOperationException($"No existing {name} unit found in this save; cannot determine stats for a new one");
        }
        return best;
    }

    private static UnitData BuildNewUnit(UnitData template, int level, ushort coordinateCode, byte currentTurn)
    {
        var newUnit = template.Clone();
        newUnit.CoordinateCode = coordinateCode;
        newUnit.Level = (byte)level;

        if (template.Level != level && template.Level > 0)
        {
            var ratio = (double)level / template.Level;
            newUnit.MaxHealth = ScaleUInt16(template.MaxHealth, ratio);
            newUnit.Movement = ScaleUInt16(template.Movement, ratio);
        }
        newUnit.CurrentHealth = newUnit.MaxHealth;

        newUnit.Experience = 0;
        newUnit.GeneralId = 0;
        newUnit.GeneralMilitaryRank = 0;
        newUnit.GeneralTitle = 0;
        newUnit.GeneralBadges = new byte[3];
        newUnit.GeneralSkillLevels = new byte[5];
        newUnit.MoraleValue = 0;
        newUnit.MoraleTurnsLeft = 0;

        newUnit.UnknownArr5[10] = currentTurn;
        newUnit.UnknownArr6[0] = UnitData.ReadyActionState;

        return newUnit;
    }

    private static ushort ScaleUInt16(ushort value, double ratio)
    {
        var scaled = (int)Math.Round(value * ratio);
        if (scaled < 1)
        {
            scaled = 1;
        }
        if (scaled > ushort.MaxValue)
        {
            scaled = ushort.MaxValue;
        }
        return (ushort)scaled;
    }

    private static HashSet<HexPos> OccupiedTiles(SaveDocument doc, bool excludeCity)
    {
        var occupied = new HashSet<HexPos>();
        foreach (var u in doc.Units)
        {
            if (excludeCity && u.UnitType == UnitTypes.Id.City)
            {
                continue;
            }
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(u.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (valid)
            {
                occupied.Add(new HexPos(row, col));
            }
        }
        return occupied;
    }

    private static List<HexPos> FindFreeLandTiles(SaveDocument doc, HexPos anchor, ushort provinceCode, int count)
    {
        var result = TryFindFreeLandTiles(doc, anchor, provinceCode, count);
        if (result.Count == 0)
        {
            throw new InvalidOperationException("No free tile available in this province; nothing was spawned");
        }
        return result;
    }

    private static List<HexPos> TryFindFreeLandTiles(SaveDocument doc, HexPos anchor, ushort provinceCode, int count)
    {
        var occupied = OccupiedTiles(doc, excludeCity: true);

        var candidates = new List<HexPos>();
        for (var row = 0; row < doc.CityTiles.Length; row++)
        {
            for (var col = 0; col < doc.CityTiles[row].Length; col++)
            {
                if (doc.CityTiles[row][col] != provinceCode || occupied.Contains(new HexPos(row, col)))
                {
                    continue;
                }
                candidates.Add(new HexPos(row, col));
            }
        }

        candidates.Sort((a, b) => HexMath.HexDistance(anchor.Col, anchor.Row, a.Col, a.Row)
            .CompareTo(HexMath.HexDistance(anchor.Col, anchor.Row, b.Col, b.Row)));

        return candidates.Count > count ? candidates.GetRange(0, count) : candidates;
    }

    private static List<HexPos> FindFreeNavalTiles(SaveDocument doc, ushort provinceCode, int count)
    {
        var result = TryFindFreeNavalTiles(doc, provinceCode, count);
        if (result.Count == 0)
        {
            throw new InvalidOperationException(FindFreeNavalTilesErrorMessage(doc, provinceCode));
        }
        return result;
    }

    private static string FindFreeNavalTilesErrorMessage(SaveDocument doc, ushort provinceCode) =>
        FindPortAnchors(doc, provinceCode).Count == 0
            ? "This province has no port; naval units cannot spawn here"
            : "No free water tile found next to this province's port(s); nothing was spawned";

    private static List<HexPos> FindPortAnchors(SaveDocument doc, ushort provinceCode)
    {
        var anchors = new List<HexPos>();
        foreach (var city in doc.Cities)
        {
            if (!ProvinceLookup.IsPort(city.BuildingType))
            {
                continue;
            }
            var (portProvince, attached) = ProvinceLookup.ProvinceCodeForPort(doc, city);
            if (!attached || portProvince != provinceCode)
            {
                continue;
            }
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (valid)
            {
                anchors.Add(new HexPos(row, col));
            }
        }
        return anchors;
    }

    private static List<HexPos> TryFindFreeNavalTiles(SaveDocument doc, ushort provinceCode, int count)
    {
        var anchors = FindPortAnchors(doc, provinceCode);
        if (anchors.Count == 0)
        {
            return [];
        }

        var occupied = OccupiedTiles(doc, excludeCity: false);
        var candidates = new List<(HexPos Pos, int Dist)>();

        for (var row = 0; row < doc.CityTiles.Length; row++)
        {
            for (var col = 0; col < doc.CityTiles[row].Length; col++)
            {
                if (doc.CityTiles[row][col] != SaveDocument.OceanTileCode || occupied.Contains(new HexPos(row, col)))
                {
                    continue;
                }
                var minDist = int.MaxValue;
                foreach (var anchor in anchors)
                {
                    var d = HexMath.HexDistance(anchor.Col, anchor.Row, col, row);
                    if (d < minDist)
                    {
                        minDist = d;
                    }
                }
                if (minDist <= MaxCoastSearchRadius)
                {
                    candidates.Add((new HexPos(row, col), minDist));
                }
            }
        }

        candidates.Sort((a, b) => a.Dist.CompareTo(b.Dist));
        var take = candidates.Count > count ? candidates.GetRange(0, count) : candidates;
        return take.Select(c => c.Pos).ToList();
    }

    private static void AppendNewUnits(SaveDocument doc, List<UnitData> newUnits)
    {
        if (newUnits.Count == 0)
        {
            return;
        }

        var unitSize = UnitData.StructSize;
        var originalUnitCount = doc.Units.Count;
        var insertionOffset = doc.Offsets.UnitOffsets[originalUnitCount - 1] + unitSize;

        var unitBytes = new byte[newUnits.Count * unitSize];
        for (var i = 0; i < newUnits.Count; i++)
        {
            newUnits[i].ToBytes().CopyTo(unitBytes, i * unitSize);
        }

        SaveFileWriter.InsertBytesAt(doc, insertionOffset, unitBytes);

        var newUnitCount = (uint)(doc.Units.Count + newUnits.Count);
        SaveFileWriter.WriteUInt32At(doc, SaveHeader.UnitCountFieldOffset, newUnitCount);

        for (var i = 0; i < newUnits.Count; i++)
        {
            doc.Offsets.UnitOffsets.Add(insertionOffset + i * unitSize);
        }

        // Landmines (and anything else) physically after the insertion point shift with it.
        for (var i = 0; i < doc.Offsets.LandmineOffsets.Count; i++)
        {
            if (doc.Offsets.LandmineOffsets[i] >= insertionOffset)
            {
                doc.Offsets.LandmineOffsets[i] += unitBytes.Length;
            }
        }

        doc.Units.AddRange(newUnits);
        doc.Header.UnitCount = newUnitCount;
    }
}
