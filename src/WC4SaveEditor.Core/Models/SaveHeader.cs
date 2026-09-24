using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class SaveHeader
{
    public required byte[] Magic { get; init; }
    public uint UnknownInt1 { get; init; }
    public uint MapId { get; init; }
    public uint GameMode { get; init; }
    public uint UnknownInt2 { get; init; }
    public uint UnknownInt3 { get; init; }
    public required float[] Camera { get; init; }
    public uint UnknownInt4 { get; init; }
    public uint TurnNumber { get; init; }
    public required byte[] UnknownArr2 { get; init; }
    public required uint[] SaveTimestamp { get; init; }
    public required byte[] UnknownArr3 { get; init; }
    public uint UnknownInt7 { get; init; }
    public uint UnknownInt8 { get; init; }
    public uint MapWidth { get; init; }
    public uint MapHeight { get; init; }
    public uint CountryCount { get; init; }
    public uint CityCount { get; init; }
    public uint UnitCount { get; set; }
    public uint UnknownCount1 { get; init; }
    public uint UnknownCount2 { get; init; }
    public required byte[] UnknownArr4 { get; init; }
    public uint TurnCount1 { get; init; }
    public uint TurnCount2 { get; init; }
    public uint UnknownCount3 { get; init; }
    public uint UnknownCount4 { get; init; }
    public uint UnknownCount5 { get; init; }
    public uint UnknownCount6 { get; init; }
    public uint ImportantCityCount { get; init; }
    public required byte[] UnknownArr5 { get; init; }
    public uint UnknownInt9 { get; init; }
    public uint UnknownInt10 { get; init; }
    public required byte[] UnknownArr6 { get; init; }
    public uint LandmineCount { get; init; }
    public required byte[] UnknownArr7 { get; init; }
    public uint UnknownCount9 { get; init; }

    public SaveHeader Clone() => new()
    {
        Magic = Magic,
        UnknownInt1 = UnknownInt1,
        MapId = MapId,
        GameMode = GameMode,
        UnknownInt2 = UnknownInt2,
        UnknownInt3 = UnknownInt3,
        Camera = Camera,
        UnknownInt4 = UnknownInt4,
        TurnNumber = TurnNumber,
        UnknownArr2 = UnknownArr2,
        SaveTimestamp = SaveTimestamp,
        UnknownArr3 = UnknownArr3,
        UnknownInt7 = UnknownInt7,
        UnknownInt8 = UnknownInt8,
        MapWidth = MapWidth,
        MapHeight = MapHeight,
        CountryCount = CountryCount,
        CityCount = CityCount,
        UnitCount = UnitCount,
        UnknownCount1 = UnknownCount1,
        UnknownCount2 = UnknownCount2,
        UnknownArr4 = UnknownArr4,
        TurnCount1 = TurnCount1,
        TurnCount2 = TurnCount2,
        UnknownCount3 = UnknownCount3,
        UnknownCount4 = UnknownCount4,
        UnknownCount5 = UnknownCount5,
        UnknownCount6 = UnknownCount6,
        ImportantCityCount = ImportantCityCount,
        UnknownArr5 = UnknownArr5,
        UnknownInt9 = UnknownInt9,
        UnknownInt10 = UnknownInt10,
        UnknownArr6 = UnknownArr6,
        LandmineCount = LandmineCount,
        UnknownArr7 = UnknownArr7,
        UnknownCount9 = UnknownCount9,
    };

    public const int UnitCountFieldOffset = 4 + 4 + 4 + 4 + 4 + 4 + 12 + 4 + 4 + 12 + 20 + 16 + 4 + 4 + 4 + 4 + 4 + 4;

    public static SaveHeader ReadFrom(ByteCursor c) => new()
    {
        Magic = c.ReadBytes(4),
        UnknownInt1 = c.ReadUInt32(),
        MapId = c.ReadUInt32(),
        GameMode = c.ReadUInt32(),
        UnknownInt2 = c.ReadUInt32(),
        UnknownInt3 = c.ReadUInt32(),
        Camera = [c.ReadSingle(), c.ReadSingle(), c.ReadSingle()],
        UnknownInt4 = c.ReadUInt32(),
        TurnNumber = c.ReadUInt32(),
        UnknownArr2 = c.ReadBytes(12),
        SaveTimestamp = [c.ReadUInt32(), c.ReadUInt32(), c.ReadUInt32(), c.ReadUInt32(), c.ReadUInt32()],
        UnknownArr3 = c.ReadBytes(16),
        UnknownInt7 = c.ReadUInt32(),
        UnknownInt8 = c.ReadUInt32(),
        MapWidth = c.ReadUInt32(),
        MapHeight = c.ReadUInt32(),
        CountryCount = c.ReadUInt32(),
        CityCount = c.ReadUInt32(),
        UnitCount = c.ReadUInt32(),
        UnknownCount1 = c.ReadUInt32(),
        UnknownCount2 = c.ReadUInt32(),
        UnknownArr4 = c.ReadBytes(8),
        TurnCount1 = c.ReadUInt32(),
        TurnCount2 = c.ReadUInt32(),
        UnknownCount3 = c.ReadUInt32(),
        UnknownCount4 = c.ReadUInt32(),
        UnknownCount5 = c.ReadUInt32(),
        UnknownCount6 = c.ReadUInt32(),
        ImportantCityCount = c.ReadUInt32(),
        UnknownArr5 = c.ReadBytes(4),
        UnknownInt9 = c.ReadUInt32(),
        UnknownInt10 = c.ReadUInt32(),
        UnknownArr6 = c.ReadBytes(12),
        LandmineCount = c.ReadUInt32(),
        UnknownArr7 = c.ReadBytes(16),
        UnknownCount9 = c.ReadUInt32(),
    };
}
