using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Models;

public sealed class SaveDocument
{
    public required string FilePath { get; init; }
    public required byte[] RawBuffer { get; set; }
    public required SaveHeader Header { get; init; }
    public required List<CountryData> Players { get; init; }
    public required ushort[][] CityTiles { get; init; }
    public required byte[][] UnitOwnerData { get; init; }
    public required List<CityData> Cities { get; init; }
    public required List<UnitData> Units { get; init; }
    public required List<LandmineData> Landmines { get; init; }
    public required OffsetTable Offsets { get; init; }

    public const ushort OceanTileCode = 65535;

    public byte GetTileOwner(int row, int col) => UnitOwnerData[row][col];

    public void SetTileOwner(int row, int col, byte owner) => UnitOwnerData[row][col] = owner;

    /// <summary>
    /// Deep-enough copy for previewing pending changes without disk I/O and without risking
    /// any mutation bleeding back into the document being edited: independent RawBuffer,
    /// UnitOwnerData, and every model field any operation writes to in place.
    /// </summary>
    public SaveDocument Clone() => new()
    {
        FilePath = FilePath,
        RawBuffer = (byte[])RawBuffer.Clone(),
        Header = Header.Clone(),
        Players = Players.Select(p => p.Clone()).ToList(),
        CityTiles = CityTiles,
        UnitOwnerData = UnitOwnerData.Select(row => (byte[])row.Clone()).ToArray(),
        Cities = Cities.Select(c => c.Clone()).ToList(),
        Units = Units.Select(u => u.Clone()).ToList(),
        Landmines = Landmines.Select(l => l.Clone()).ToList(),
        Offsets = Offsets.Clone(),
    };
}
