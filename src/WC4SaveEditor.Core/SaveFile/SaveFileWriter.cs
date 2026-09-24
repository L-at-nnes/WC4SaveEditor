using System.Buffers.Binary;
using WC4SaveEditor.Core.Models;

namespace WC4SaveEditor.Core.SaveFile;

public static class SaveFileWriter
{
    public static void WriteByteAt(SaveDocument doc, int offset, byte value) => doc.RawBuffer[offset] = value;

    public static void WriteUInt16At(SaveDocument doc, int offset, ushort value) =>
        BinaryPrimitives.WriteUInt16LittleEndian(doc.RawBuffer.AsSpan(offset, 2), value);

    public static void WriteUInt32At(SaveDocument doc, int offset, uint value) =>
        BinaryPrimitives.WriteUInt32LittleEndian(doc.RawBuffer.AsSpan(offset, 4), value);

    public static void WriteBytesAt(SaveDocument doc, int offset, byte[] value) =>
        value.CopyTo(doc.RawBuffer.AsSpan(offset, value.Length));

    public static void InsertBytesAt(SaveDocument doc, int offset, byte[] insertedData)
    {
        var newBuffer = new byte[doc.RawBuffer.Length + insertedData.Length];
        doc.RawBuffer.AsSpan(0, offset).CopyTo(newBuffer);
        insertedData.CopyTo(newBuffer.AsSpan(offset, insertedData.Length));
        doc.RawBuffer.AsSpan(offset).CopyTo(newBuffer.AsSpan(offset + insertedData.Length));
        doc.RawBuffer = newBuffer;
    }

    public static void EnsureBackup(string filePath)
    {
        var backupPath = filePath + ".bak";
        if (File.Exists(backupPath))
        {
            return;
        }

        File.Copy(filePath, backupPath);
    }

    public static void Commit(SaveDocument doc)
    {
        EnsureBackup(doc.FilePath);
        File.WriteAllBytes(doc.FilePath, doc.RawBuffer);
    }
}
