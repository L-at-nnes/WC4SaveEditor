namespace WC4SaveEditor.Core.Data;

public static class UnitTypes
{
    public static class Id
    {
        public const byte LightInfantry = 1;
        public const byte AssaultInfantry = 2;
        public const byte MotorizedInfantry = 3;
        public const byte MechanizedInfantry = 4;
        public const byte Commando = 5;
        public const byte ArmoredCar = 6;
        public const byte LightTank = 7;
        public const byte MediumTank = 8;
        public const byte HeavyTank = 9;
        public const byte SuperTank = 10;
        public const byte FieldArtillery = 11;
        public const byte Howitzer = 12;
        public const byte RocketArtillery = 13;
        public const byte SuperArtillery = 14;
        public const byte Submarine = 15;
        public const byte Destroyer = 16;
        public const byte Cruiser = 17;
        public const byte Carrier = 18;
        public const byte SuperCarrier = 19;
        public const byte Bunker = 35;
        public const byte LandFort = 36;
        public const byte CoastalArtillery = 37;
        public const byte RocketLauncher = 38;
        public const byte City = 39;
        public const byte Brandenburgers = 48;
        public const byte HawkeyeForce = 49;
        public const byte CombatMedic = 50;
        public const byte M7Priest = 51;
        public const byte HeavyGustav = 52;
        public const byte BM21 = 53;
        public const byte T44 = 54;
        public const byte KingTiger = 55;
        public const byte M26Pershing = 56;
        public const byte TypeVIISubmarine = 60;
        public const byte HMSPrinceOfWales = 61;
        public const byte B4Howitzer = 64;
        public const byte IS3HeavyTank = 65;
        public const byte StukaZuFuss = 66;
        public const byte Richelieu = 67;
        public const byte Enterprise = 70;
        public const byte RPGRocketeer = 73;
        public const byte A41Centurion = 74;
        public const byte AuF1 = 79;
        public const byte T72 = 80;
        public const byte PhantomForce = 82;
        public const byte M1A1Abrams = 83;
        public const byte AH64Apache = 85;
        public const byte M142Himars = 91;
        public const byte DivineWrathMBT = 92;
    }

    private static readonly Dictionary<byte, string> IdToName = new()
    {
        [Id.LightInfantry] = "Light Infantry",
        [Id.AssaultInfantry] = "Assault Infantry",
        [Id.MotorizedInfantry] = "Motorized Infantry",
        [Id.MechanizedInfantry] = "Mechanized Infantry",
        [Id.Commando] = "Commando",
        [Id.ArmoredCar] = "Armored Car",
        [Id.LightTank] = "Light Tank",
        [Id.MediumTank] = "Medium Tank",
        [Id.HeavyTank] = "Heavy Tank",
        [Id.SuperTank] = "Super Tank",
        [Id.FieldArtillery] = "Field Artillery",
        [Id.Howitzer] = "Howitzer",
        [Id.RocketArtillery] = "Rocket Artillery",
        [Id.SuperArtillery] = "Super Artillery",
        [Id.Submarine] = "Submarine",
        [Id.Destroyer] = "Destroyer",
        [Id.Cruiser] = "Cruiser",
        [Id.Carrier] = "Carrier",
        [Id.SuperCarrier] = "Super Carrier",
        [Id.Bunker] = "Bunker",
        [Id.LandFort] = "Land Fort",
        [Id.CoastalArtillery] = "Coastal Artillery",
        [Id.RocketLauncher] = "Rocket Launcher",
        [Id.City] = "City",
        [Id.Brandenburgers] = "Brandenburgers",
        [Id.HawkeyeForce] = "Hawkeye Force",
        [Id.CombatMedic] = "Combat Medic",
        [Id.M7Priest] = "M7 Priest",
        [Id.HeavyGustav] = "Heavy Gustav",
        [Id.BM21] = "BM-21",
        [Id.T44] = "T-44",
        [Id.KingTiger] = "King Tiger",
        [Id.M26Pershing] = "M26 Pershing",
        [Id.TypeVIISubmarine] = "Type VII Submarine",
        [Id.HMSPrinceOfWales] = "HMS Prince of Wales",
        [Id.B4Howitzer] = "B-4 Howitzer",
        [Id.IS3HeavyTank] = "IS-3 Heavy Tank",
        [Id.StukaZuFuss] = "Stuka zu Fuss",
        [Id.Richelieu] = "Richelieu",
        [Id.Enterprise] = "Enterprise",
        [Id.RPGRocketeer] = "RPG Rocketeer",
        [Id.A41Centurion] = "A41 Centurion",
        [Id.AuF1] = "AuF1",
        [Id.T72] = "T-72",
        [Id.PhantomForce] = "Phantom Force",
        [Id.M1A1Abrams] = "M1A1 Abrams",
        [Id.AH64Apache] = "AH-64 Apache",
        [Id.M142Himars] = "M142 Himars",
        [Id.DivineWrathMBT] = "Divine Wrath MBT",
    };

    public static bool TryGetName(byte unitType, out string name) => IdToName.TryGetValue(unitType, out name!);
}
