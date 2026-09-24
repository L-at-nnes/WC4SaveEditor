using CommunityToolkit.Mvvm.ComponentModel;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class QuickTogglesViewModel : ObservableObject
{
    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private bool _suppressToggleHandling;

    [ObservableProperty] private bool _conquerAll;
    [ObservableProperty] private bool _maxMoney;
    [ObservableProperty] private bool _maxCityTech;
    [ObservableProperty] private bool _maxCityLevel;
    [ObservableProperty] private bool _maxMorale;
    [ObservableProperty] private bool _healAllies;
    [ObservableProperty] private bool _weakenEnemies;

    public QuickTogglesViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
    }

    public void Reset(SaveDocument document)
    {
        _suppressToggleHandling = true;
        ConquerAll = false;
        MaxMoney = false;
        MaxCityTech = false;
        MaxCityLevel = false;
        MaxMorale = false;
        HealAllies = false;
        WeakenEnemies = false;
        _suppressToggleHandling = false;
    }

    private void SetToggle<T>(bool isOn, string label) where T : PendingChange, new()
    {
        if (_suppressToggleHandling)
        {
            return;
        }
        _queue.RemoveAll<T>();
        if (isOn)
        {
            _queue.Add(new T());
            NotificationCenter.Info($"Will apply \"{label}\" on Save.");
        }
        _notifyChanged();
    }

    partial void OnConquerAllChanged(bool value) => SetToggle<ConquerAllChange>(value, "Conquer all");
    partial void OnMaxMoneyChanged(bool value) => SetToggle<MaxMoneyChange>(value, "Max money");
    partial void OnMaxCityTechChanged(bool value) => SetToggle<MaxCityTechChange>(value, "Max city tech");
    partial void OnMaxCityLevelChanged(bool value) => SetToggle<MaxCityLevelChange>(value, "Max city level");
    partial void OnMaxMoraleChanged(bool value) => SetToggle<MaxMoraleChange>(value, "Max morale");
    partial void OnHealAlliesChanged(bool value) => SetToggle<HealAlliesChange>(value, "Heal allies");
    partial void OnWeakenEnemiesChanged(bool value) => SetToggle<WeakenEnemiesChange>(value, "Weaken enemies");
}
