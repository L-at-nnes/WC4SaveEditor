using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;

namespace WC4SaveEditor.Gui.Services;

public enum ToastKind
{
    Info,
    Success,
    Error,
}

public partial class ToastViewModel(string message, ToastKind kind) : ObservableObject
{
    public Guid Id { get; } = Guid.NewGuid();
    public string Message { get; } = message;
    public ToastKind Kind { get; } = kind;

    [RelayCommand]
    private void Close() => NotificationCenter.Dismiss(this);
}
