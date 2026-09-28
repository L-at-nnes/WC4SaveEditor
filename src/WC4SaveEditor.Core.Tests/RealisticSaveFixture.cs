using System.Buffers.Binary;
using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Tests;

/// <summary>
/// Builds a real on-disk save file byte-for-byte the way World Conqueror 4 would, then loads it
/// back through <see cref="SaveFileReader"/> so tests exercise the exact same offset math the
/// GUI uses - unlike the hand-built <see cref="SaveDocument"/> instances used elsewhere, which
/// skip the reader entirely and can hide bugs in how offsets are computed.
/// </summary>
public static class RealisticSaveFixture
{
    public const int MapWidth = 4;
    public const int MapHeight = 2;

    // Province A (tiles 0,1,4,5) capital city at coordinate 0, owned by player 1.
    // Province B (tiles 2,3,6,7) capital city at coordinate 2, owned by player 2.
    public const ushort ProvinceACode = 0;
    public const ushort ProvinceBCode = 2;

    public static SaveDocument Load(string path)
    {
        WriteTo(path);
        return SaveFileReader.Read(path);
    }

    private static void WriteTo(string path)
    {
        var players = new List<byte[]>
        {
            BuildPlayer(countryId: 1, teamId: 0),
            BuildPlayer(countryId: 2, teamId: 1),
            BuildPlayer(countryId: 3, teamId: 2),
        };

        var cityTiles = new ushort[MapHeight][];
        var unitOwners = new byte[MapHeight][];
        for (var row = 0; row < MapHeight; row++)
        {
            cityTiles[row] = new ushort[MapWidth];
            unitOwners[row] = new byte[MapWidth];
            for (var col = 0; col < MapWidth; col++)
            {
                var code = row * MapWidth + col;
                var inProvinceA = code is 0 or 1 or 4 or 5;
                cityTiles[row][col] = inProvinceA ? ProvinceACode : ProvinceBCode;
                unitOwners[row][col] = (byte)(inProvinceA ? 1 : 2);
            }
        }

        var cities = new List<byte[]>
        {
            BuildCity(coordinateCode: ProvinceACode, cityId: 1, buildingType: 11),
            BuildCity(coordinateCode: ProvinceBCode, cityId: 2, buildingType: 11),
        };

        var units = new List<byte[]>
        {
            BuildUnit(coordinateCode: ProvinceACode, unitType: UnitTypes.Id.City),
            BuildUnit(coordinateCode: ProvinceBCode, unitType: UnitTypes.Id.City),
            BuildUnit(coordinateCode: 1, unitType: 1), // regular unit for player 1
            BuildUnit(coordinateCode: 6, unitType: 1), // regular unit for player 2
        };

        var landmines = new List<byte[]>
        {
            BuildLandmine(coordinateCode: 1, owner: 2), // sits in province A, owned (1-based) by player 1
        };

        using var ms = new MemoryStream();
        WriteHeader(ms, players.Count, cities.Count, units.Count, landmines.Count);
        foreach (var p in players) ms.Write(p);
        Span<byte> tileCodeBuf = stackalloc byte[2];
        foreach (var row in cityTiles)
        {
            foreach (var code in row)
            {
                BinaryPrimitives.WriteUInt16LittleEndian(tileCodeBuf, code);
                ms.Write(tileCodeBuf);
            }
        }
        foreach (var row in unitOwners) ms.Write(row);
        foreach (var c in cities) ms.Write(c);
        foreach (var u in units) ms.Write(u);
        foreach (var l in landmines) ms.Write(l);

        File.WriteAllBytes(path, ms.ToArray());
    }

    private static void WriteHeader(MemoryStream ms, int countryCount, int cityCount, int unitCount, int landmineCount)
    {
        using var w = new BinaryWriter(ms, System.Text.Encoding.ASCII, leaveOpen: true);
        w.Write(new byte[4]); // Magic
        w.Write((uint)0); // UnknownInt1
        w.Write((uint)0); // MapId
        w.Write(GameModes.Campaign); // GameMode
        w.Write((uint)0); // UnknownInt2
        w.Write((uint)0); // UnknownInt3
        w.Write(0f); w.Write(0f); w.Write(0f); // Camera
        w.Write((uint)0); // UnknownInt4
        w.Write((uint)1); // TurnNumber
        w.Write(new byte[12]); // UnknownArr2
        w.Write((uint)0); w.Write((uint)0); w.Write((uint)0); w.Write((uint)0); w.Write((uint)0); // SaveTimestamp
        w.Write(new byte[16]); // UnknownArr3
        w.Write((uint)1); // UnknownInt7 - nonzero so the reader skips the campaign-only visibility grid
        w.Write((uint)0); // UnknownInt8
        w.Write((uint)MapWidth);
        w.Write((uint)MapHeight);
        w.Write((uint)countryCount);
        w.Write((uint)cityCount);
        w.Write((uint)unitCount);
        w.Write((uint)0); // UnknownCount1
        w.Write((uint)0); // UnknownCount2
        w.Write(new byte[8]); // UnknownArr4
        w.Write((uint)0); w.Write((uint)0); // TurnCount1/2
        w.Write((uint)0); // UnknownCount3
        w.Write((uint)0); // UnknownCount4
        w.Write((uint)0); // UnknownCount5
        w.Write((uint)0); // UnknownCount6
        w.Write((uint)0); // ImportantCityCount
        w.Write(new byte[4]); // UnknownArr5
        w.Write((uint)0); // UnknownInt9
        w.Write((uint)(MapWidth * MapHeight)); // UnknownInt10 - matches tile count so no mismatch padding is read
        w.Write(new byte[12]); // UnknownArr6
        w.Write((uint)landmineCount);
        w.Write(new byte[16]); // UnknownArr7
        w.Write((uint)0); // UnknownCount9
    }

    private static byte[] BuildPlayer(uint countryId, uint teamId)
    {
        using var ms = new MemoryStream();
        using var w = new BinaryWriter(ms);
        w.Write((uint)0); // TurnOrder
        w.Write(countryId);
        w.Write((uint)0); w.Write((uint)0); w.Write((uint)0); // Currency
        w.Write((uint)0); // BotFlag
        w.Write(teamId);
        w.Write(new byte[4]); // UnknownArr2
        w.Write(new byte[4]); w.Write(new byte[4]); // UnknownColor
        w.Write(new byte[4]); // PrimaryColor
        w.Write(new byte[16]); // UnknownArr4
        w.Write(new byte[460]); // UnknownArr5
        return ms.ToArray();
    }

    private static byte[] BuildCity(ushort coordinateCode, ushort cityId, byte buildingType)
    {
        using var ms = new MemoryStream();
        using var w = new BinaryWriter(ms);
        w.Write(coordinateCode);
        w.Write(cityId);
        w.Write(buildingType);
        w.Write((byte)0); w.Write((byte)0); w.Write((byte)0); // Apperance, UnknownByte1, Wonders
        w.Write(new byte[6]); // UnknownArr2
        w.Write(new byte[8]); // UnknownArr3
        w.Write((byte)0); w.Write((byte)0); // AntiAir
        w.Write(new byte[6]); // TechLevels
        w.Write(new byte[2]); // UnknownArr4
        return ms.ToArray();
    }

    private static byte[] BuildUnit(ushort coordinateCode, byte unitType)
    {
        using var ms = new MemoryStream();
        using var w = new BinaryWriter(ms);
        w.Write(coordinateCode);
        w.Write(unitType);
        w.Write((byte)0); w.Write((byte)0); w.Write((byte)0); // Level, Personnel, Direction
        w.Write((ushort)0); // Movement
        w.Write((ushort)0); // Experience
        w.Write((ushort)0); // UnknownHealth
        w.Write((ushort)100); // CurrentHealth
        w.Write((ushort)100); // MaxHealth
        w.Write((ushort)0); // GeneralId
        w.Write((byte)0); w.Write((byte)0); // GeneralMilitaryRank, GeneralTitle
        w.Write(new byte[3]); // GeneralBadges
        w.Write(new byte[5]); // GeneralSkillLevels
        w.Write(new byte[12]); // UnknownArr5
        w.Write((sbyte)0); // MoraleValue
        w.Write((ushort)0); // MoraleTurnsLeft
        w.Write(new byte[21]); // UnknownArr6
        return ms.ToArray();
    }

    private static byte[] BuildLandmine(ushort coordinateCode, ushort owner)
    {
        using var ms = new MemoryStream();
        using var w = new BinaryWriter(ms);
        w.Write(coordinateCode);
        w.Write(owner);
        w.Write(new byte[2]); // UnknownArr1
        w.Write((ushort)100); // Health
        w.Write(new byte[4]); // UnknownArr2
        return ms.ToArray();
    }
}
