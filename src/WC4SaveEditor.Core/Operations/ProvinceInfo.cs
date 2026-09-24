namespace WC4SaveEditor.Core.Operations;

public sealed class ProvinceInfo
{
    public required int CityIndex { get; init; }
    public required ushort CoordinateCode { get; init; }
    public required string Name { get; init; }
    public required byte Owner { get; init; }
    public required int TileCount { get; init; }
    public int UnitCount { get; set; }
    public int PortCount { get; set; }
}
