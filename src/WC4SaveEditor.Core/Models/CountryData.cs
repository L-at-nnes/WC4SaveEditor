using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class CountryData
{
    public const int StructSize = 4 + 4 + 12 + 4 + 4 + 4 + 8 + 4 + 16 + 460;

    public const int CurrencyFieldOffset = 8;
    public const int TeamIdFieldOffset = 4 + 4 + 12 + 4;

    public uint TurnOrder { get; init; }
    public uint CountryId { get; init; }
    public required uint[] Currency { get; init; }
    public uint BotFlag { get; init; }
    public uint TeamId { get; set; }
    public required byte[] UnknownArr2 { get; init; }
    public required byte[][] UnknownColor { get; init; }
    public required byte[] PrimaryColor { get; init; }
    public required byte[] UnknownArr4 { get; init; }
    public required byte[] UnknownArr5 { get; init; }

    public CountryData Clone() => new()
    {
        TurnOrder = TurnOrder,
        CountryId = CountryId,
        Currency = (uint[])Currency.Clone(),
        BotFlag = BotFlag,
        TeamId = TeamId,
        UnknownArr2 = UnknownArr2,
        UnknownColor = UnknownColor,
        PrimaryColor = PrimaryColor,
        UnknownArr4 = UnknownArr4,
        UnknownArr5 = UnknownArr5,
    };

    public static CountryData ReadFrom(ByteCursor c) => new()
    {
        TurnOrder = c.ReadUInt32(),
        CountryId = c.ReadUInt32(),
        Currency = [c.ReadUInt32(), c.ReadUInt32(), c.ReadUInt32()],
        BotFlag = c.ReadUInt32(),
        TeamId = c.ReadUInt32(),
        UnknownArr2 = c.ReadBytes(4),
        UnknownColor = [c.ReadBytes(4), c.ReadBytes(4)],
        PrimaryColor = c.ReadBytes(4),
        UnknownArr4 = c.ReadBytes(16),
        UnknownArr5 = c.ReadBytes(460),
    };
}
