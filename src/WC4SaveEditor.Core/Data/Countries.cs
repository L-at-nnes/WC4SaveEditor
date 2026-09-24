namespace WC4SaveEditor.Core.Data;

public static class Countries
{
    public readonly record struct CountryInfo(byte Id, string Name);

    public static readonly IReadOnlyList<CountryInfo> All = new CountryInfo[]
    {
        new(0x01, "UK"),
        new(0x02, "France"),
        new(0x03, "Germany"),
        new(0x04, "Germany"),
        new(0x05, "Soviet Union"),
        new(0x06, "USA"),
        new(0x07, "Italy"),
        new(0x08, "ROC"),
        new(0x09, "PRC"),
        new(0x0a, "Japan"),
        new(0x0b, "Finland"),
        new(0x0c, "Poland"),
        new(0x0d, "Yugoslavia"),
        new(0x0e, "Canada"),
        new(0x0f, "Australia"),
        new(0x10, "Norway"),
        new(0x11, "Sweden"),
        new(0x12, "Denmark"),
        new(0x13, "Netherlands"),
        new(0x14, "Belgium"),
        new(0x15, "Spain"),
        new(0x16, "Portugal"),
        new(0x17, "Hungary"),
        new(0x18, "Romania"),
        new(0x19, "Bulgaria"),
        new(0x1a, "Switzerland"),
        new(0x1b, "Greece"),
        new(0x1c, "Turkey"),
        new(0x1d, "Saudi Arabia"),
        new(0x1e, "Iraq"),
        new(0x1f, "Iran"),
        new(0x20, "India"),
        new(0x21, "Thailand"),
        new(0x22, "Mongolia"),
        new(0x23, "North Korea"),
        new(0x24, "South Korea"),
        new(0x25, "Mexico"),
        new(0x26, "Cuba"),
        new(0x27, "Colombia"),
        new(0x28, "Brazil"),
        new(0x29, "Bolivia"),
        new(0x2a, "Venezuela"),
        new(0x2b, "Peru"),
        new(0x2c, "Chile"),
        new(0x2d, "Argentina"),
        new(0x2e, "Egypt"),
        new(0x2f, "Liberia"),
        new(0x30, "Mysterious Forces"),
        new(0x31, "East Germany"),
        new(0x33, "Americas Scorpion"),
        new(0x34, "Asia-Pacific Scorpion"),
        new(0x35, "European Scorpion"),
        new(0x36, "Middle East Scorpion"),
        new(0x37, "African Scorpion"),
    };

    private static readonly Dictionary<byte, string> IdToName = All.ToDictionary(c => c.Id, c => c.Name);

    public static bool TryGetName(byte id, out string name) => IdToName.TryGetValue(id, out name!);
    public static int Count => All.Count;

    public static string ColorBytesToHex(ReadOnlySpan<byte> colorBytes) =>
        $"{colorBytes[0]:x2}{colorBytes[1]:x2}{colorBytes[2]:x2}";
}
