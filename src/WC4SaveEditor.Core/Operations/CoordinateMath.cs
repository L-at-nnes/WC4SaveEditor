using WC4SaveEditor.Core.Data;

namespace WC4SaveEditor.Core.Operations;

public static class CoordinateMath
{
    public static (int Row, int Col, bool Valid) ConvertCoordinates(int coordinateCode, int mapWidth, int mapHeight, uint gameMode)
    {
        var row = coordinateCode / mapWidth;
        if (gameMode == GameModes.Conquest)
        {
            row -= 2;
        }
        var col = coordinateCode % mapWidth;

        var valid = row >= 0 && row < mapHeight && col >= 0 && col < mapWidth;
        return (row, col, valid);
    }

    public static int ConvertToCoordinateCode(int row, int col, int mapWidth, uint gameMode)
    {
        var adjustedRow = row;
        if (gameMode == GameModes.Conquest)
        {
            adjustedRow += 2;
        }
        return adjustedRow * mapWidth + col;
    }
}
