using WC4SaveEditor.Core.Data;
using Xunit;

namespace WC4SaveEditor.Core.Tests;

public class CountriesTests
{
    [Theory]
    [InlineData(new byte[] { 255, 0, 0, 0 }, "ff0000")]
    [InlineData(new byte[] { 0, 255, 0, 0 }, "00ff00")]
    [InlineData(new byte[] { 0, 0, 255, 0 }, "0000ff")]
    [InlineData(new byte[] { 250, 250, 150, 0 }, "fafa96")]
    public void ColorBytesToHex_MatchesGoReference(byte[] colorBytes, string expected)
    {
        Assert.Equal(expected, Countries.ColorBytesToHex(colorBytes));
    }

    [Fact]
    public void TryGetName_ResolvesKnownCountries()
    {
        Assert.True(Countries.TryGetName(0x01, out var uk));
        Assert.Equal("UK", uk);

        Assert.True(Countries.TryGetName(0x0a, out var japan));
        Assert.Equal("Japan", japan);

        Assert.True(Countries.TryGetName(0x30, out var mysterious));
        Assert.Equal("Mysterious Forces", mysterious);
    }

    [Fact]
    public void TryGetName_RejectsUnknownIds()
    {
        Assert.False(Countries.TryGetName(0x00, out _));
        Assert.False(Countries.TryGetName(0x32, out _));
        Assert.False(Countries.TryGetName(0x38, out _));
    }

    [Fact]
    public void All_StartsAndEndsWithExpectedEntries()
    {
        Assert.NotEmpty(Countries.All);
        Assert.Equal((byte)0x01, Countries.All[0].Id);
        Assert.Equal("UK", Countries.All[0].Name);

        var last = Countries.All[^1];
        Assert.Equal((byte)0x37, last.Id);
        Assert.Equal("African Scorpion", last.Name);
    }
}
