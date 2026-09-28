using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using Xunit;

namespace WC4SaveEditor.Core.Tests;

public class ChangeQueueOrderingTests
{
    private static string TempSavePath() => Path.Combine(Path.GetTempPath(), $"wc4-test-{Guid.NewGuid():N}.sav");

    [Fact]
    public void Commit_AppliesMaxCityLevel_AfterAnnexation_SoNewlyAnnexedCitiesAreIncluded()
    {
        var path = TempSavePath();
        try
        {
            var doc = RealisticSaveFixture.Load(path);
            var annexedCityIndex = doc.Cities.FindIndex(c => c.CoordinateCode == RealisticSaveFixture.ProvinceBCode);

            var queue = new ChangeQueue();
            queue.Add(new ConquerPlayerChange(oldPlayer: 2, newPlayer: ConquestOperations.MainPlayer));
            queue.Add(new MaxCityLevelChange());

            queue.Commit(doc);

            Assert.Equal(StatOperations.CityLevelBase + 4, doc.Cities[annexedCityIndex].BuildingType);
        }
        finally
        {
            File.Delete(path);
        }
    }
}
