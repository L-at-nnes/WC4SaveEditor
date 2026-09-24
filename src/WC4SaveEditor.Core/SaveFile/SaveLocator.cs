using System.Text.RegularExpressions;

namespace WC4SaveEditor.Core.SaveFile;

public static partial class SaveLocator
{
    private static readonly HashSet<string> NonSaveFiles = new(StringComparer.OrdinalIgnoreCase)
    {
        "headquarter.sav",
        "uuid.sav",
        "prd.sav",
    };

    public static IReadOnlyList<string> FindLocalStateFolders()
    {
        var localAppData = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
        if (string.IsNullOrEmpty(localAppData))
        {
            return [];
        }

        var packagesDir = Path.Combine(localAppData, "Packages");
        if (!Directory.Exists(packagesDir))
        {
            return [];
        }

        return Directory.GetDirectories(packagesDir, "EasyTech.WorldConqueror4_*")
            .Select(dir => Path.Combine(dir, "LocalState"))
            .Where(Directory.Exists)
            .ToList();
    }

    [GeneratedRegex(@"^conquest(\d+)$", RegexOptions.IgnoreCase)]
    private static partial Regex ConquestSlotPattern();

    private static (int Rank, int SubRank, string Name) SortKey(string filePath)
    {
        var name = Path.GetFileNameWithoutExtension(filePath);

        if (string.Equals(name, "campaign", StringComparison.OrdinalIgnoreCase))
        {
            return (0, 0, name);
        }

        var match = ConquestSlotPattern().Match(name);
        if (match.Success)
        {
            return (1, int.Parse(match.Groups[1].Value), name);
        }

        return (2, 0, name);
    }

    public static IReadOnlyList<string> FindSaveFiles()
    {
        var folders = FindLocalStateFolders();
        if (folders.Count == 0)
        {
            return [];
        }

        var localStateDir = folders.Count == 1
            ? folders[0]
            : folders.OrderByDescending(f => Directory.GetLastWriteTimeUtc(f)).First();

        return Directory.GetFiles(localStateDir, "*.sav")
            .Where(f => !NonSaveFiles.Contains(Path.GetFileName(f)))
            .Select(f => (Path: f, Key: SortKey(f)))
            .OrderBy(x => x.Key.Rank).ThenBy(x => x.Key.SubRank).ThenBy(x => x.Key.Name, StringComparer.OrdinalIgnoreCase)
            .Select(x => x.Path)
            .ToList();
    }

    public static string? FindDefaultSaveFile()
    {
        var files = FindSaveFiles();
        return files.Count == 0 ? null : files.OrderByDescending(File.GetLastWriteTimeUtc).First();
    }
}
