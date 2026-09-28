using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Core.SaveFile;
using Xunit;

namespace WC4SaveEditor.Core.Tests;

/// <summary>
/// Exercises annexation end to end through the real byte-level reader/writer instead of a
/// hand-built <see cref="Models.SaveDocument"/>, so any bug in offset math (the kind that
/// corrupts a real save file) shows up here instead of only in the field.
/// </summary>
public class AnnexationIntegrationTests
{
    private static string TempSavePath() => Path.Combine(Path.GetTempPath(), $"wc4-test-{Guid.NewGuid():N}.sav");

    [Fact]
    public void ConvertPlayer_MarksFullyStrippedPlayerAsEliminated()
    {
        // Root cause of a real "annexing breaks the save" report: the actual game keeps a
        // per-country "eliminated" byte (CountryData.IsEliminated) and misbehaves if a country
        // shows zero tiles but isn't flagged eliminated - confirmed by diffing a real save
        // before/after annexing a country through this tool. ConvertPlayer must keep this flag
        // in sync whenever a player's last tile is taken.
        var path = TempSavePath();
        try
        {
            var doc = RealisticSaveFixture.Load(path);
            Assert.False(doc.Players[2].IsEliminated);

            ConquestOperations.ConvertPlayer(doc, oldPlayer: 2, newPlayer: 1);

            Assert.True(doc.Players[2].IsEliminated);
            Assert.False(doc.Players[1].IsEliminated);

            SaveFileWriter.Commit(doc);
            var reloaded = SaveFileReader.Read(path);
            Assert.True(reloaded.Players[2].IsEliminated);
            Assert.False(reloaded.Players[1].IsEliminated);
        }
        finally
        {
            File.Delete(path);
            File.Delete(path + ".bak");
        }
    }

    [Fact]
    public void ConvertPlayer_ThenReload_ProducesValidSaveWithTransferredOwnership()
    {
        var path = TempSavePath();
        try
        {
            var doc = RealisticSaveFixture.Load(path);

            ConquestOperations.ConvertPlayer(doc, oldPlayer: 2, newPlayer: 1);
            SaveFileWriter.Commit(doc);

            var reloaded = SaveFileReader.Read(path);

            for (var row = 0; row < RealisticSaveFixture.MapHeight; row++)
            {
                for (var col = 0; col < RealisticSaveFixture.MapWidth; col++)
                {
                    Assert.NotEqual(2, reloaded.UnitOwnerData[row][col]);
                }
            }

            Assert.All(reloaded.Landmines, m => Assert.Equal(2, m.Owner)); // player 1, one-based
        }
        finally
        {
            File.Delete(path);
            File.Delete(path + ".bak");
        }
    }

    [Fact]
    public void ConvertProvince_ThenReload_ProducesValidSaveWithTransferredOwnership()
    {
        var path = TempSavePath();
        try
        {
            var doc = RealisticSaveFixture.Load(path);

            ConquestOperations.ConvertProvince(doc, RealisticSaveFixture.ProvinceBCode, newPlayer: 1);
            SaveFileWriter.Commit(doc);

            var reloaded = SaveFileReader.Read(path);

            for (var row = 0; row < RealisticSaveFixture.MapHeight; row++)
            {
                for (var col = 0; col < RealisticSaveFixture.MapWidth; col++)
                {
                    var code = row * RealisticSaveFixture.MapWidth + col;
                    var inProvinceB = code is 2 or 3 or 6 or 7;
                    if (inProvinceB)
                    {
                        Assert.Equal(1, reloaded.UnitOwnerData[row][col]);
                    }
                }
            }
        }
        finally
        {
            File.Delete(path);
            File.Delete(path + ".bak");
        }
    }

    [Fact]
    public void ConvertPlayer_NeverStealsMainPlayerTiles_AndLeavesTheirLandminesAlone()
    {
        // ConvertPlayer refuses to move any tile currently owned by the main player (id 0),
        // even when player 0 is passed as the annexed "oldPlayer" - a deliberate safeguard so
        // you can never accidentally annex yourself away. Landmines must honor the same
        // safeguard, or a save comes out with landmines owned by someone whose territory
        // never actually changed.
        var path = TempSavePath();
        try
        {
            var doc = RealisticSaveFixture.Load(path);
            for (var row = 0; row < RealisticSaveFixture.MapHeight; row++)
            {
                for (var col = 0; col < RealisticSaveFixture.MapWidth; col++)
                {
                    doc.UnitOwnerData[row][col] = 0;
                }
            }
            doc.Landmines[0].Owner = 1; // one-based owner id for player 0

            ConquestOperations.ConvertPlayer(doc, oldPlayer: 0, newPlayer: 1);

            for (var row = 0; row < RealisticSaveFixture.MapHeight; row++)
            {
                for (var col = 0; col < RealisticSaveFixture.MapWidth; col++)
                {
                    Assert.Equal(0, doc.UnitOwnerData[row][col]);
                }
            }
            Assert.Equal(1, doc.Landmines[0].Owner);
        }
        finally
        {
            File.Delete(path);
            File.Delete(path + ".bak");
        }
    }
}
