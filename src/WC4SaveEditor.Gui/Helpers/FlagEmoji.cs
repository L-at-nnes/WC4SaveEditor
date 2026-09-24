namespace WC4SaveEditor.Gui.Helpers;

public static class FlagEmoji
{
    private static readonly Dictionary<string, string> NameToIso = new()
    {
        ["UK"] = "GB",
        ["France"] = "FR",
        ["Germany"] = "DE",
        ["Soviet Union"] = "RU",
        ["USA"] = "US",
        ["Italy"] = "IT",
        ["ROC"] = "TW",
        ["PRC"] = "CN",
        ["Japan"] = "JP",
        ["Finland"] = "FI",
        ["Poland"] = "PL",
        ["Yugoslavia"] = "RS",
        ["Canada"] = "CA",
        ["Australia"] = "AU",
        ["Norway"] = "NO",
        ["Sweden"] = "SE",
        ["Denmark"] = "DK",
        ["Netherlands"] = "NL",
        ["Belgium"] = "BE",
        ["Spain"] = "ES",
        ["Portugal"] = "PT",
        ["Hungary"] = "HU",
        ["Romania"] = "RO",
        ["Bulgaria"] = "BG",
        ["Switzerland"] = "CH",
        ["Greece"] = "GR",
        ["Turkey"] = "TR",
        ["Saudi Arabia"] = "SA",
        ["Iraq"] = "IQ",
        ["Iran"] = "IR",
        ["India"] = "IN",
        ["Thailand"] = "TH",
        ["Mongolia"] = "MN",
        ["North Korea"] = "KP",
        ["South Korea"] = "KR",
        ["Mexico"] = "MX",
        ["Cuba"] = "CU",
        ["Colombia"] = "CO",
        ["Brazil"] = "BR",
        ["Bolivia"] = "BO",
        ["Venezuela"] = "VE",
        ["Peru"] = "PE",
        ["Chile"] = "CL",
        ["Argentina"] = "AR",
        ["Egypt"] = "EG",
        ["Liberia"] = "LR",
    };

    public static string? TryGetFlag(string countryName)
    {
        if (!NameToIso.TryGetValue(countryName, out var iso))
        {
            return null;
        }
        return FromIsoCode(iso);
    }

    private static string FromIsoCode(string iso)
    {
        var first = char.ConvertFromUtf32(0x1F1E6 + (iso[0] - 'A'));
        var second = char.ConvertFromUtf32(0x1F1E6 + (iso[1] - 'A'));
        return first + second;
    }
}
