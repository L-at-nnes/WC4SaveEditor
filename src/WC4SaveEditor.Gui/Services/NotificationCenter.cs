using System.Collections.ObjectModel;
using System.Windows;
using System.Windows.Threading;

namespace WC4SaveEditor.Gui.Services;

public static class NotificationCenter
{
    public static ObservableCollection<ToastViewModel> Toasts { get; } = [];

    public static void Info(string message) => Show(message, ToastKind.Info, TimeSpan.FromSeconds(3));

    public static void Success(string message) => Show(message, ToastKind.Success, TimeSpan.FromSeconds(3));

    public static void Error(string message) => Show(message, ToastKind.Error, TimeSpan.FromSeconds(5));

    private static void Show(string message, ToastKind kind, TimeSpan? autoDismiss)
    {
        void Add()
        {
            var toast = new ToastViewModel(message, kind);
            Toasts.Add(toast);

            if (autoDismiss is { } delay)
            {
                var timer = new DispatcherTimer { Interval = delay };
                timer.Tick += (_, _) =>
                {
                    timer.Stop();
                    Dismiss(toast);
                };
                timer.Start();
            }
        }

        if (Application.Current?.Dispatcher.CheckAccess() == false)
        {
            Application.Current.Dispatcher.Invoke(Add);
        }
        else
        {
            Add();
        }
    }

    public static void Dismiss(ToastViewModel toast)
    {
        if (Application.Current?.Dispatcher.CheckAccess() == false)
        {
            Application.Current.Dispatcher.Invoke(() => Toasts.Remove(toast));
        }
        else
        {
            Toasts.Remove(toast);
        }
    }
}
