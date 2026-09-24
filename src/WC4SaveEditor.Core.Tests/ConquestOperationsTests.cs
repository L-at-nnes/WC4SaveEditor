using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Core.SaveFile;
using Xunit;

namespace WC4SaveEditor.Core.Tests;

public class ConquestOperationsTests
{
    private static SaveDocument BuildMinimalDocument(int width, int height, byte[][] owners)
    {
        var offsets = new OffsetTable { UnitOwnerGridStart = 0, MapWidth = width };
        var rawBuffer = new byte[width * height];
        for (var row = 0; row < height; row++)
        {
            for (var col = 0; col < width; col++)
            {
                rawBuffer[row * width + col] = owners[row][col];
            }
        }

        var players = new List<CountryData>
        {
            new() { Currency = [0, 0, 0], UnknownArr2 = new byte[4], UnknownColor = [new byte[4], new byte[4]], PrimaryColor = new byte[4], UnknownArr4 = new byte[16], UnknownArr5 = new byte[460] },
            new() { Currency = [0, 0, 0], UnknownArr2 = new byte[4], UnknownColor = [new byte[4], new byte[4]], PrimaryColor = new byte[4], UnknownArr4 = new byte[16], UnknownArr5 = new byte[460] },
        };

        return new SaveDocument
        {
            FilePath = "test.sav",
            RawBuffer = rawBuffer,
            Header = new SaveHeader
            {
                Magic = new byte[4],
                Camera = new float[3],
                UnknownArr2 = new byte[12],
                SaveTimestamp = new uint[5],
                UnknownArr3 = new byte[16],
                UnknownArr4 = new byte[8],
                UnknownArr5 = new byte[4],
                UnknownArr6 = new byte[12],
                UnknownArr7 = new byte[16],
                GameMode = GameModes.Campaign,
                MapWidth = (uint)width,
                MapHeight = (uint)height,
                CountryCount = (uint)players.Count,
            },
            Players = players,
            CityTiles = Enumerable.Range(0, height).Select(_ => new ushort[width]).ToArray(),
            UnitOwnerData = owners,
            Cities = [],
            Units = [],
            Landmines = [],
            Offsets = offsets,
        };
    }

    [Fact]
    public void ConvertTile_WritesNewOwnerToRawBufferAtComputedOffset()
    {
        var owners = new[] { new byte[] { 0, 1 }, new byte[] { 1, ConquestOperations.TileUnowned } };
        var doc = BuildMinimalDocument(2, 2, owners);

        ConquestOperations.ConvertTile(doc, x: 1, y: 0, newPlayer: 0);

        Assert.Equal(0, doc.UnitOwnerData[0][1]);
        Assert.Equal(0, doc.RawBuffer[doc.Offsets.TileOffset(0, 1)]);
    }

    [Fact]
    public void ConvertTile_ThrowsWhenTileIsUnowned()
    {
        var owners = new[] { new byte[] { 0, 1 }, new byte[] { 1, ConquestOperations.TileUnowned } };
        var doc = BuildMinimalDocument(2, 2, owners);

        Assert.Throws<InvalidOperationException>(() => ConquestOperations.ConvertTile(doc, x: 1, y: 1, newPlayer: 0));
    }

    [Fact]
    public void ConvertAllPlayers_ClaimsEveryOwnedTileForMainPlayer()
    {
        var owners = new[] { new byte[] { 0, 1 }, new byte[] { 1, ConquestOperations.TileUnowned } };
        var doc = BuildMinimalDocument(2, 2, owners);

        var stats = ConquestOperations.ConvertAllPlayers(doc);

        Assert.Equal(2, stats.TotalChanged);
        Assert.Equal(0, doc.UnitOwnerData[0][0]);
        Assert.Equal(0, doc.UnitOwnerData[0][1]);
        Assert.Equal(0, doc.UnitOwnerData[1][0]);
        Assert.Equal(ConquestOperations.TileUnowned, doc.UnitOwnerData[1][1]);
    }

    [Fact]
    public void ChangePlayerTeam_WritesTeamIdAtCorrectOffset()
    {
        var owners = new[] { new byte[] { 0, 1 }, new byte[] { 1, 1 } };
        var doc = BuildMinimalDocument(2, 2, owners);
        doc.Offsets.PlayerOffsets.Add(1000);
        doc.Offsets.PlayerOffsets.Add(1000 + CountryData.StructSize);
        doc.RawBuffer = new byte[2000];

        ConquestOperations.ChangePlayerTeam(doc, 1, teamId: 7);

        Assert.Equal(7u, doc.Players[1].TeamId);
        var writtenOffset = 1000 + CountryData.StructSize + CountryData.TeamIdFieldOffset;
        Assert.Equal(7u, System.Buffers.Binary.BinaryPrimitives.ReadUInt32LittleEndian(doc.RawBuffer.AsSpan(writtenOffset, 4)));
    }
}
