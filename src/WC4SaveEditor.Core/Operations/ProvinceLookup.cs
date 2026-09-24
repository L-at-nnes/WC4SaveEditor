using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;

namespace WC4SaveEditor.Core.Operations;

public static class ProvinceLookup
{
    public const byte PortLevel1 = 31;
    public const byte PortLevel4 = 34;

    public static bool IsPort(byte buildingType) => buildingType is >= PortLevel1 and <= PortLevel4;

    public static (ushort Code, bool Attached) ProvinceCodeForPort(SaveDocument doc, CityData port)
    {
        var gameMode = doc.Header.GameMode;
        var (row, col, valid) = CoordinateMath.ConvertCoordinates(port.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, gameMode);
        if (!valid)
        {
            return (0, false);
        }

        var directCode = doc.CityTiles[row][col];
        if (directCode != SaveDocument.OceanTileCode)
        {
            return (directCode, true);
        }

        var portOwner = doc.UnitOwnerData[row][col];
        ushort nearestCode = 0;
        var hasNearest = false;

        for (var radius = 1; radius <= 3; radius++)
        {
            for (var candidateRow = row - radius; candidateRow <= row + radius; candidateRow++)
            {
                for (var candidateCol = col - radius; candidateCol <= col + radius; candidateCol++)
                {
                    if (candidateRow < 0 || candidateRow >= doc.CityTiles.Length || candidateCol < 0 || candidateCol >= doc.CityTiles[candidateRow].Length)
                    {
                        continue;
                    }
                    if (candidateRow != row - radius && candidateRow != row + radius && candidateCol != col - radius && candidateCol != col + radius)
                    {
                        continue;
                    }

                    var code = doc.CityTiles[candidateRow][candidateCol];
                    if (code == SaveDocument.OceanTileCode)
                    {
                        continue;
                    }
                    if (!hasNearest)
                    {
                        nearestCode = code;
                        hasNearest = true;
                    }

                    var (provinceRow, provinceCol, provinceValid) = CoordinateMath.ConvertCoordinates(code, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, gameMode);
                    if (provinceValid && doc.UnitOwnerData[provinceRow][provinceCol] == portOwner)
                    {
                        return (code, true);
                    }
                }
            }
        }

        return (nearestCode, hasNearest);
    }

    public static string GetCityNameAtPosition(SaveDocument doc, int row, int col)
    {
        foreach (var city in doc.Cities)
        {
            var (cityRow, cityCol, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, doc.Header.GameMode);
            if (valid && cityRow == row && cityCol == col)
            {
                return Cities.TryGetName(city.CityId, out var name) ? name : "Unknown";
            }
        }
        return "Unknown";
    }

    public static List<ProvinceInfo> GetPlayerProvinces(SaveDocument doc, int playerId)
    {
        var tileCounts = new Dictionary<ushort, int>();
        foreach (var row in doc.CityTiles)
        {
            foreach (var code in row)
            {
                if (code != SaveDocument.OceanTileCode)
                {
                    tileCounts[code] = tileCounts.GetValueOrDefault(code) + 1;
                }
            }
        }

        var gameMode = doc.Header.GameMode;
        var provinces = new List<ProvinceInfo>();

        for (var cityIndex = 0; cityIndex < doc.Cities.Count; cityIndex++)
        {
            var city = doc.Cities[cityIndex];
            if (!tileCounts.TryGetValue(city.CoordinateCode, out var tileCount) || tileCount == 0)
            {
                continue;
            }

            var (row, col, valid) = CoordinateMath.ConvertCoordinates(city.CoordinateCode, (int)doc.Header.MapWidth, (int)doc.Header.MapHeight, gameMode);
            if (!valid || doc.UnitOwnerData[row][col] != (byte)playerId)
            {
                continue;
            }

            var name = Cities.TryGetName(city.CityId, out var cityName) ? cityName : $"Province at ({row},{col})";

            var province = new ProvinceInfo
            {
                CityIndex = cityIndex,
                CoordinateCode = city.CoordinateCode,
                Name = name,
                Owner = (byte)playerId,
                TileCount = tileCount,
            };

            for (var tileRow = 0; tileRow < doc.CityTiles.Length; tileRow++)
            {
                for (var tileCol = 0; tileCol < doc.CityTiles[tileRow].Length; tileCol++)
                {
                    if (doc.CityTiles[tileRow][tileCol] == city.CoordinateCode && doc.UnitOwnerData[tileRow][tileCol] != ConquestOperations.TileUnowned)
                    {
                        province.UnitCount++;
                    }
                }
            }

            foreach (var candidate in doc.Cities)
            {
                if (!IsPort(candidate.BuildingType))
                {
                    continue;
                }
                var (portProvince, attached) = ProvinceCodeForPort(doc, candidate);
                if (attached && portProvince == city.CoordinateCode)
                {
                    province.PortCount++;
                }
            }

            provinces.Add(province);
        }

        return provinces;
    }

    public static string IdentifyUnitAtTile(SaveDocument doc, int row, int col)
    {
        var coordinateCode = CoordinateMath.ConvertToCoordinateCode(row, col, (int)doc.Header.MapWidth, doc.Header.GameMode);

        foreach (var unit in doc.Units)
        {
            if (unit.CoordinateCode == coordinateCode)
            {
                return UnitTypes.TryGetName(unit.UnitType, out var unitName) ? unitName : $"Type {unit.UnitType}";
            }
        }

        foreach (var city in doc.Cities)
        {
            if (city.CoordinateCode == coordinateCode)
            {
                return Cities.TryGetName(city.CityId, out var cityName) ? cityName : $"City (ID {city.CityId})";
            }
        }

        return "Structure";
    }
}
