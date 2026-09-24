namespace WC4SaveEditor.Core.SaveFile;

public sealed class OffsetTable
{
    public List<int> PlayerOffsets { get; } = [];
    public List<int> CityOffsets { get; } = [];
    public List<int> UnitOffsets { get; } = [];
    public List<int> LandmineOffsets { get; } = [];
    public int UnitOwnerGridStart { get; set; }
    public int MapWidth { get; set; }

    public int TileOffset(int row, int col) => UnitOwnerGridStart + row * MapWidth + col;

    public OffsetTable Clone()
    {
        var clone = new OffsetTable { UnitOwnerGridStart = UnitOwnerGridStart, MapWidth = MapWidth };
        clone.PlayerOffsets.AddRange(PlayerOffsets);
        clone.CityOffsets.AddRange(CityOffsets);
        clone.UnitOffsets.AddRange(UnitOffsets);
        clone.LandmineOffsets.AddRange(LandmineOffsets);
        return clone;
    }
}
