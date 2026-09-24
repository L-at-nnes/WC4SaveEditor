using System.Buffers.Binary;

namespace WC4SaveEditor.Core.SaveFile;

public sealed class ByteCursor(byte[] buffer, int position = 0)
{
    public byte[] Buffer { get; } = buffer;
    public int Position { get; private set; } = position;

    public byte ReadByte() => Buffer[Position++];

    public sbyte ReadSByte() => (sbyte)Buffer[Position++];

    public ushort ReadUInt16()
    {
        var value = BinaryPrimitives.ReadUInt16LittleEndian(Buffer.AsSpan(Position, 2));
        Position += 2;
        return value;
    }

    public uint ReadUInt32()
    {
        var value = BinaryPrimitives.ReadUInt32LittleEndian(Buffer.AsSpan(Position, 4));
        Position += 4;
        return value;
    }

    public float ReadSingle()
    {
        var value = BinaryPrimitives.ReadSingleLittleEndian(Buffer.AsSpan(Position, 4));
        Position += 4;
        return value;
    }

    public byte[] ReadBytes(int count)
    {
        var slice = Buffer.AsSpan(Position, count).ToArray();
        Position += count;
        return slice;
    }

    public void Skip(int count) => Position += count;
}
