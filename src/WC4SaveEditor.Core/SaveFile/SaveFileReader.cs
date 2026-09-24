using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;

namespace WC4SaveEditor.Core.SaveFile;

public static class SaveFileReader
{
    public static SaveDocument Read(string filePath)
    {
        var buffer = File.ReadAllBytes(filePath);
        var c = new ByteCursor(buffer);
        var offsets = new OffsetTable();

        var header = SaveHeader.ReadFrom(c);
        var isConquest = header.GameMode == GameModes.Conquest;

        var players = new List<CountryData>((int)header.CountryCount);
        for (var i = 0; i < header.CountryCount; i++)
        {
            offsets.PlayerOffsets.Add(c.Position);
            players.Add(CountryData.ReadFrom(c));
        }

        var mapWidth = (int)header.MapWidth;
        var mapHeight = (int)header.MapHeight;

        if (!isConquest && header.UnknownInt7 == 0)
        {
            c.Skip(mapWidth * mapHeight * 16);
        }

        var cityTiles = new ushort[mapHeight][];
        for (var row = 0; row < mapHeight; row++)
        {
            cityTiles[row] = new ushort[mapWidth];
            for (var col = 0; col < mapWidth; col++)
            {
                cityTiles[row][col] = c.ReadUInt16();
            }
        }

        var tileCountMismatch = mapWidth * mapHeight != (int)header.UnknownInt10;
        if (!isConquest && tileCountMismatch)
        {
            c.Skip(8);
        }

        offsets.UnitOwnerGridStart = c.Position;
        offsets.MapWidth = mapWidth;
        var unitOwnerData = new byte[mapHeight][];
        for (var row = 0; row < mapHeight; row++)
        {
            unitOwnerData[row] = new byte[mapWidth];
            for (var col = 0; col < mapWidth; col++)
            {
                unitOwnerData[row][col] = c.ReadByte();
            }
        }

        if (!isConquest && tileCountMismatch)
        {
            c.Skip(4);
        }

        var cities = new List<CityData>((int)header.CityCount);
        for (var i = 0; i < header.CityCount; i++)
        {
            offsets.CityOffsets.Add(c.Position);
            var city = CityData.ReadFrom(c);
            if (i > 0 && city.CoordinateCode == 0)
            {
                throw new InvalidDataException($"Invalid city data at index {i}");
            }
            cities.Add(city);
        }

        var units = new List<UnitData>((int)header.UnitCount);
        for (var i = 0; i < header.UnitCount; i++)
        {
            offsets.UnitOffsets.Add(c.Position);
            units.Add(UnitData.ReadFrom(c));
        }

        var landmines = new List<LandmineData>((int)header.LandmineCount);
        for (var i = 0; i < header.LandmineCount; i++)
        {
            offsets.LandmineOffsets.Add(c.Position);
            landmines.Add(LandmineData.ReadFrom(c));
        }

        c.Skip((int)header.UnknownCount1 * 16);
        c.Skip((int)header.UnknownCount2 * 44);
        c.Skip((int)header.UnknownCount3 * 80);
        c.Skip((int)header.UnknownCount5 * 8);
        c.Skip((int)header.UnknownCount6 * 8);
        c.Skip((int)header.ImportantCityCount * 4);
        c.Skip((int)header.UnknownCount9 * 16);

        return new SaveDocument
        {
            FilePath = filePath,
            RawBuffer = buffer,
            Header = header,
            Players = players,
            CityTiles = cityTiles,
            UnitOwnerData = unitOwnerData,
            Cities = cities,
            Units = units,
            Landmines = landmines,
            Offsets = offsets,
        };
    }
}
