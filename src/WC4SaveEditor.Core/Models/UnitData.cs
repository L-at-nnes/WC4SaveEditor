using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class UnitData
{
    public const int StructSize = 2 + 1 + 1 + 1 + 1 + 2 + 2 + 2 + 2 + 2 + 2 + 1 + 1 + 3 + 5 + 12 + 1 + 2 + 21;

    public const int CurrentHealthFieldOffset = 2 + 1 + 1 + 1 + 1 + 2 + 2 + 2;
    public const int MaxHealthFieldOffset = CurrentHealthFieldOffset + 2;
    public const int MoraleValueFieldOffset = 2 + 1 + 1 + 1 + 1 + 2 + 2 + 2 + 2 + 2 + 2 + 1 + 1 + 3 + 5 + 12;
    public const int MoraleTurnsLeftFieldOffset = MoraleValueFieldOffset + 1;
    public const int UnknownArr5FieldOffset = 2 + 1 + 1 + 1 + 1 + 2 + 2 + 2 + 2 + 2 + 2 + 1 + 1 + 3 + 5;
    public const int UnknownArr6FieldOffset = MoraleTurnsLeftFieldOffset + 2;

    public const int ActivationTurnOffset = UnknownArr5FieldOffset + 10;
    public const int ActionStateOffset = UnknownArr6FieldOffset + 0;
    public const byte ReadyActionState = 2;

    public ushort CoordinateCode { get; set; }
    public byte UnitType { get; set; }
    public byte Level { get; set; }
    public byte Personnel { get; set; }
    public byte Direction { get; set; }
    public ushort Movement { get; set; }
    public ushort Experience { get; set; }
    public ushort UnknownHealth { get; set; }
    public ushort CurrentHealth { get; set; }
    public ushort MaxHealth { get; set; }
    public ushort GeneralId { get; set; }
    public byte GeneralMilitaryRank { get; set; }
    public byte GeneralTitle { get; set; }
    public required byte[] GeneralBadges { get; set; }
    public required byte[] GeneralSkillLevels { get; set; }
    public required byte[] UnknownArr5 { get; set; }
    public sbyte MoraleValue { get; set; }
    public ushort MoraleTurnsLeft { get; set; }
    public required byte[] UnknownArr6 { get; set; }

    public static UnitData ReadFrom(ByteCursor c) => new()
    {
        CoordinateCode = c.ReadUInt16(),
        UnitType = c.ReadByte(),
        Level = c.ReadByte(),
        Personnel = c.ReadByte(),
        Direction = c.ReadByte(),
        Movement = c.ReadUInt16(),
        Experience = c.ReadUInt16(),
        UnknownHealth = c.ReadUInt16(),
        CurrentHealth = c.ReadUInt16(),
        MaxHealth = c.ReadUInt16(),
        GeneralId = c.ReadUInt16(),
        GeneralMilitaryRank = c.ReadByte(),
        GeneralTitle = c.ReadByte(),
        GeneralBadges = c.ReadBytes(3),
        GeneralSkillLevels = c.ReadBytes(5),
        UnknownArr5 = c.ReadBytes(12),
        MoraleValue = c.ReadSByte(),
        MoraleTurnsLeft = c.ReadUInt16(),
        UnknownArr6 = c.ReadBytes(21),
    };

    public byte[] ToBytes()
    {
        var buffer = new byte[StructSize];
        var offset = 0;
        void PutU16(ushort v) { System.Buffers.Binary.BinaryPrimitives.WriteUInt16LittleEndian(buffer.AsSpan(offset, 2), v); offset += 2; }
        void PutU8(byte v) { buffer[offset++] = v; }
        void PutBytes(byte[] v) { v.CopyTo(buffer.AsSpan(offset, v.Length)); offset += v.Length; }

        PutU16(CoordinateCode);
        PutU8(UnitType);
        PutU8(Level);
        PutU8(Personnel);
        PutU8(Direction);
        PutU16(Movement);
        PutU16(Experience);
        PutU16(UnknownHealth);
        PutU16(CurrentHealth);
        PutU16(MaxHealth);
        PutU16(GeneralId);
        PutU8(GeneralMilitaryRank);
        PutU8(GeneralTitle);
        PutBytes(GeneralBadges);
        PutBytes(GeneralSkillLevels);
        PutBytes(UnknownArr5);
        buffer[offset++] = unchecked((byte)MoraleValue);
        PutU16(MoraleTurnsLeft);
        PutBytes(UnknownArr6);
        return buffer;
    }

    public UnitData Clone() => new()
    {
        CoordinateCode = CoordinateCode,
        UnitType = UnitType,
        Level = Level,
        Personnel = Personnel,
        Direction = Direction,
        Movement = Movement,
        Experience = Experience,
        UnknownHealth = UnknownHealth,
        CurrentHealth = CurrentHealth,
        MaxHealth = MaxHealth,
        GeneralId = GeneralId,
        GeneralMilitaryRank = GeneralMilitaryRank,
        GeneralTitle = GeneralTitle,
        GeneralBadges = (byte[])GeneralBadges.Clone(),
        GeneralSkillLevels = (byte[])GeneralSkillLevels.Clone(),
        UnknownArr5 = (byte[])UnknownArr5.Clone(),
        MoraleValue = MoraleValue,
        MoraleTurnsLeft = MoraleTurnsLeft,
        UnknownArr6 = (byte[])UnknownArr6.Clone(),
    };
}
