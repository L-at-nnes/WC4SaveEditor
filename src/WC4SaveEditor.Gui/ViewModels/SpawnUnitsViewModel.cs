using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using WC4SaveEditor.Core.Data;
using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.Operations;
using WC4SaveEditor.Gui.Services;

namespace WC4SaveEditor.Gui.ViewModels;

public partial class SpawnUnitsViewModel : ObservableObject
{
    private readonly ChangeQueue _queue;
    private readonly Action _notifyChanged;
    private SaveDocument? _document;

    [ObservableProperty]
    private ObservableCollection<ProvinceInfo> _ownProvinces = [];

    [ObservableProperty]
    private ProvinceInfo? _selectedProvince;

    [ObservableProperty]
    private ObservableCollection<UnitTypeOption> _unitTypes = [];

    [ObservableProperty]
    private UnitTypeOption? _selectedUnitType;

    [ObservableProperty]
    private int _level = 9;

    [ObservableProperty]
    private int _count = 1;

    public SpawnUnitsViewModel(ChangeQueue queue, Action notifyChanged)
    {
        _queue = queue;
        _notifyChanged = notifyChanged;
        UnitTypes = new ObservableCollection<UnitTypeOption>(
            SpawnUnitsOperations.SpawnableUnitTypesList.Select(id => new UnitTypeOption(
                id,
                UnitTypes_Name(id),
                SpawnUnitsOperations.IsNavalUnitType(id))));
    }

    private static string UnitTypes_Name(byte id) => Core.Data.UnitTypes.TryGetName(id, out var name) ? name : $"Type {id}";

    public void Reset(SaveDocument document)
    {
        _document = document;
        OwnProvinces = new ObservableCollection<ProvinceInfo>(ProvinceLookup.GetPlayerProvinces(document, ConquestOperations.MainPlayer));
        SelectedProvince = null;
        SelectedUnitType = null;
        Level = 9;
        Count = 1;
    }

    [RelayCommand]
    private void ConfirmSpawn()
    {
        if (_document is null || SelectedProvince is null || SelectedUnitType is null)
        {
            NotificationCenter.Error("Pick a province and a unit type first.");
            return;
        }

        if (Level is < 1 or > 9)
        {
            NotificationCenter.Error("Level must be between 1 and 9.");
            return;
        }

        if (Count < 1)
        {
            NotificationCenter.Error("Count must be at least 1.");
            return;
        }

        if (!SpawnUnitsOperations.HasStatTemplate(_document.Units, SelectedUnitType.Id))
        {
            NotificationCenter.Error($"No existing {SelectedUnitType.Name} unit found in this save, so its stats can't be cloned - pick a unit type that already exists somewhere on the map.");
            return;
        }

        _queue.Add(new SpawnUnitChange(SelectedProvince.CoordinateCode, SelectedUnitType.Id, Level, Count));
        NotificationCenter.Info($"Will spawn {Count}x {SelectedUnitType.Name} (lvl {Level}) in {SelectedProvince.Name} on Save.");
        _notifyChanged();
    }
}
