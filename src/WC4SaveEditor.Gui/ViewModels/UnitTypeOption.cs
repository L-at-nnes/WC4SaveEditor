namespace WC4SaveEditor.Gui.ViewModels;

public sealed class UnitTypeOption(byte id, string name, bool isNaval)
{
    public byte Id { get; } = id;
    public string Name { get; } = name;
    public bool IsNaval { get; } = isNaval;
    public string DisplayName => IsNaval ? $"{Name} (naval)" : $"{Name} (land)";
}
