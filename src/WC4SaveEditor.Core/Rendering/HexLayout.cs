namespace WC4SaveEditor.Core.Rendering;

public static class HexLayout
{
    public static (double X, double Y) GetHexCenter(int col, int row, double hexSize)
    {
        var hexWidth = hexSize * 2;
        var hexHeight = hexSize * 1.732;

        var x = col * hexWidth * 0.75;
        var y = row * hexHeight;
        if (col % 2 == 1)
        {
            y += hexHeight / 2;
        }

        return (x + hexSize, y + hexSize);
    }

    public static (double X, double Y)[] GetHexPoints(double centerX, double centerY, double hexSize)
    {
        var points = new (double X, double Y)[6];
        for (var i = 0; i < 6; i++)
        {
            var angle = i * 60.0 * Math.PI / 180.0;
            points[i] = (centerX + hexSize * Math.Cos(angle), centerY + hexSize * Math.Sin(angle));
        }
        return points;
    }
}
