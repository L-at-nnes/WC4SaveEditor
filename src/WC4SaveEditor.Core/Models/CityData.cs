using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class CityData
{
    public const int StructSize = 2 + 2 + 1 + 1 + 1 + 1 + 6 + 8 + 1 + 1 + 6 + 2;

    public const int BuildingTypeFieldOffset = 4;
    public const int TechLevelsFieldOffset = 2 + 2 + 1 + 1 + 1 + 1 + 6 + 8 + 1 + 1;

    public ushort CoordinateCode { get; init; }
    public ushort CityId { get; init; }
    public byte BuildingType { get; set; }
    public byte Apperance { get; init; }
    public byte UnknownByte1 { get; init; }
    public byte Wonders { get; init; }
    public required byte[] UnknownArr2 { get; init; }
    public required byte[] UnknownArr3 { get; init; }
    public byte AntiAirWeaponType { get; init; }
    public byte AntiAirRange { get; init; }
    public required byte[] TechLevels { get; set; }
    public required byte[] UnknownArr4 { get; init; }

    public CityData Clone() => new()
    {
        CoordinateCode = CoordinateCode,
        CityId = CityId,
        BuildingType = BuildingType,
        Apperance = Apperance,
        UnknownByte1 = UnknownByte1,
        Wonders = Wonders,
        UnknownArr2 = UnknownArr2,
        UnknownArr3 = UnknownArr3,
        AntiAirWeaponType = AntiAirWeaponType,
        AntiAirRange = AntiAirRange,
        TechLevels = (byte[])TechLevels.Clone(),
        UnknownArr4 = UnknownArr4,
    };

    public static CityData ReadFrom(ByteCursor c) => new()
    {
        CoordinateCode = c.ReadUInt16(),
        CityId = c.ReadUInt16(),
        BuildingType = c.ReadByte(),
        Apperance = c.ReadByte(),
        UnknownByte1 = c.ReadByte(),
        Wonders = c.ReadByte(),
        UnknownArr2 = c.ReadBytes(6),
        UnknownArr3 = c.ReadBytes(8),
        AntiAirWeaponType = c.ReadByte(),
        AntiAirRange = c.ReadByte(),
        TechLevels = c.ReadBytes(6),
        UnknownArr4 = c.ReadBytes(2),
    };
}
