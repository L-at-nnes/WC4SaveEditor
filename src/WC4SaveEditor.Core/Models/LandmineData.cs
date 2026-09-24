using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class LandmineData
{
    public const int StructSize = 2 + 2 + 2 + 2 + 4;
    public const int OwnerFieldOffset = 2;

    public ushort CoordinateCode { get; init; }
    public ushort Owner { get; set; }
    public required byte[] UnknownArr1 { get; init; }
    public ushort Health { get; init; }
    public required byte[] UnknownArr2 { get; init; }

    public LandmineData Clone() => new()
    {
        CoordinateCode = CoordinateCode,
        Owner = Owner,
        UnknownArr1 = UnknownArr1,
        Health = Health,
        UnknownArr2 = UnknownArr2,
    };

    public static LandmineData ReadFrom(ByteCursor c) => new()
    {
        CoordinateCode = c.ReadUInt16(),
        Owner = c.ReadUInt16(),
        UnknownArr1 = c.ReadBytes(2),
        Health = c.ReadUInt16(),
        UnknownArr2 = c.ReadBytes(4),
    };
}
