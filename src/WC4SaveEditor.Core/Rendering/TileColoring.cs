using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;

namespace WC4SaveEditor.Core.Rendering;

public static class TileColoring
{
    public const string OceanColor = "#3D6A6A";
    public const string DefaultGray = "#C8C8C8";
    public const string InvalidGray = "#646464";

    public static Dictionary<byte, string> BuildPlayerColors(SaveDocument doc)
    {
        var colors = new Dictionary<byte, string>();
        for (var i = 0; i < doc.Players.Count; i++)
        {
            colors[(byte)i] = "#" + Countries.ColorBytesToHex(doc.Players[i].PrimaryColor);
        }
        return colors;
    }

    public static Dictionary<ushort, byte> DetermineCityOwnership(SaveDocument doc)
    {
        var ownership = new Dictionary<ushort, byte>();
        foreach (var city in doc.Cities)
        {
            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (!valid || row >= doc.UnitOwnerData.Length || col >= doc.UnitOwnerData[row].Length)
            {
                continue;
            }
            var owner = doc.UnitOwnerData[row][col];
            if (owner < doc.Players.Count)
            {
                ownership[city.CoordinateCode] = owner;
            }
        }
        return ownership;
    }

    public static string GetTileColor(SaveDocument doc, int row, int col, Dictionary<ushort, byte> cityOwnership, Dictionary<byte, string> playerColors)
    {
        if (row >= doc.CityTiles.Length || col >= doc.CityTiles[row].Length)
        {
            return OceanColor;
        }

        var cityTileId = doc.CityTiles[row][col];
        if (cityTileId == SaveDocument.OceanTileCode)
        {
            return OceanColor;
        }

        if (cityOwnership.TryGetValue(cityTileId, out var ownerId) && playerColors.TryGetValue(ownerId, out var color))
        {
            return color;
        }

        return DefaultGray;
    }
}
