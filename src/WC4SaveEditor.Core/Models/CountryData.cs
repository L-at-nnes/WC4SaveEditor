using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class CountryData
{
    public const int StructSize = 4 + 4 + 12 + 4 + 4 + 4 + 8 + 4 + 16 + 460;

    public const int CurrencyFieldOffset = 8;
    public const int TeamIdFieldOffset = 4 + 4 + 12 + 4;

    private const int UnknownArr5StartOffset = 4 + 4 + 12 + 4 + 4 + 4 + 8 + 4 + 16;
    private const int EliminatedIndexInUnknownArr5 = 456;
    public const int EliminatedFieldOffset = UnknownArr5StartOffset + EliminatedIndexInUnknownArr5;

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

    /// <summary>
    /// The game itself flips this to true once a country owns zero tiles, false while it still
    /// holds any - confirmed by diffing a real save before/after an in-game elimination. Any
    /// code that removes a country's last tile must set this too, or the game's own turn/AI
    /// logic (which reads this flag directly) chokes on a country it still thinks is alive but
    /// that owns nothing.
    /// </summary>
    public bool IsEliminated
    {
        get => UnknownArr5[EliminatedIndexInUnknownArr5] != 0;
        set => UnknownArr5[EliminatedIndexInUnknownArr5] = (byte)(value ? 1 : 0);
    }

    public CountryData Clone() => new()
    {
        TurnOrder = TurnOrder,
        CountryId = CountryId,
        Currency = (uint[])Currency.Clone(),
        BotFlag = BotFlag,
        TeamId = TeamId,
        UnknownArr2 = (byte[])UnknownArr2.Clone(),
        UnknownColor = UnknownColor.Select(c => (byte[])c.Clone()).ToArray(),
        PrimaryColor = (byte[])PrimaryColor.Clone(),
        UnknownArr4 = (byte[])UnknownArr4.Clone(),
        UnknownArr5 = (byte[])UnknownArr5.Clone(),
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
