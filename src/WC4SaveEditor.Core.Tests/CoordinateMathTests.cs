using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Operations;
using Xunit;

namespace WC4SaveEditor.Core.Tests;

public class CoordinateMathTests
{
    [Theory]
    [InlineData(0, GameModes.Campaign, 0, 0, true)]
    [InlineData(4, GameModes.Campaign, 1, 1, true)]
    [InlineData(8, GameModes.Campaign, 2, 2, true)]
    [InlineData(6, GameModes.Conquest, 0, 0, true)]
    [InlineData(10, GameModes.Conquest, 1, 1, true)]
    [InlineData(9, GameModes.Campaign, 3, 0, false)]
    [InlineData(15, GameModes.Conquest, 3, 0, false)]
    [InlineData(-1, GameModes.Campaign, 0, 0, false)]
    public void ConvertCoordinates_MatchesGoReference(int code, uint gameMode, int expectedRow, int expectedCol, bool expectedValid)
    {
        var (row, col, valid) = CoordinateMath.ConvertCoordinates(code, 3, 3, gameMode);

        Assert.Equal(expectedValid, valid);
        if (valid)
        {
            Assert.Equal(expectedRow, row);
            Assert.Equal(expectedCol, col);
        }
    }

    [Theory]
    [InlineData(0, 0, GameModes.Campaign, 0)]
    [InlineData(1, 1, GameModes.Campaign, 4)]
    [InlineData(2, 2, GameModes.Campaign, 8)]
    [InlineData(0, 0, GameModes.Conquest, 6)]
    [InlineData(1, 1, GameModes.Conquest, 10)]
    [InlineData(2, 2, GameModes.Conquest, 14)]
    public void ConvertToCoordinateCode_MatchesGoReference(int row, int col, uint gameMode, int expectedCode)
    {
        var code = CoordinateMath.ConvertToCoordinateCode(row, col, 3, gameMode);
        Assert.Equal(expectedCode, code);
    }

    [Theory]
    [InlineData(GameModes.Campaign)]
    [InlineData(GameModes.Conquest)]
    public void CoordinateRoundTrip_HoldsForEveryTileOn3x3Map(uint gameMode)
    {
        for (var row = 0; row < 3; row++)
        {
            for (var col = 0; col < 3; col++)
            {
                var code = CoordinateMath.ConvertToCoordinateCode(row, col, 3, gameMode);
                var (resultRow, resultCol, valid) = CoordinateMath.ConvertCoordinates(code, 3, 3, gameMode);

                Assert.True(valid, $"round-trip for ({row},{col}) produced an invalid coordinate");
                Assert.Equal(row, resultRow);
                Assert.Equal(col, resultCol);
            }
        }
    }
}
