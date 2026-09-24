namespace WC4SaveEditor.Core.Operations;

public readonly record struct HexPos(int Row, int Col);

public static class HexMath
{
    public static (int X, int Y, int Z) OddQToCube(int col, int row)
    {
        var x = col;
        var z = row - (col - (col & 1)) / 2;
        var y = -x - z;
        return (x, y, z);
    }

    public static int HexDistance(int col1, int row1, int col2, int row2)
    {
        var (x1, y1, z1) = OddQToCube(col1, row1);
        var (x2, y2, z2) = OddQToCube(col2, row2);
        return (Math.Abs(x1 - x2) + Math.Abs(y1 - y2) + Math.Abs(z1 - z2)) / 2;
    }
}
