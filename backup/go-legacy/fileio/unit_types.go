package fileio

import "fmt"

const (
	UnitTypeLightInfantry      = 1
	UnitTypeAssaultInfantry    = 2
	UnitTypeMotorizedInfantry  = 3
	UnitTypeMechanizedInfantry = 4
	UnitTypeCommando           = 5
	UnitTypeArmoredCar         = 6
	UnitTypeLightTank          = 7
	UnitTypeMediumTank         = 8
	UnitTypeHeavyTank          = 9
	UnitTypeSuperTank          = 10
	UnitTypeFieldArtillery     = 11
	UnitTypeHowitzer           = 12
	UnitTypeRocketArtillery    = 13
	UnitTypeSuperArtillery     = 14
	UnitTypeSubmarine          = 15
	UnitTypeDestroyer          = 16
	UnitTypeCruiser            = 17
	UnitTypeCarrier            = 18
	UnitTypeSuperCarrier       = 19
	UnitTypeBunker             = 35
	UnitTypeLandFort           = 36
	UnitTypeCoastalArtillery   = 37
	UnitTypeRocketLauncher     = 38
	UnitTypeCity               = 39
	UnitTypeBrandenburgers     = 48
	UnitTypeHawkeyeForce       = 49
	UnitTypeCombatMedic        = 50
	UnitTypeM7Priest           = 51
	UnitTypeHeavyGustav        = 52
	UnitTypeBM21               = 53
	UnitTypeT44                = 54
	UnitTypeKingTiger          = 55
	UnitTypeM26Pershing        = 56
	UnitTypeTypeVIISubmarine   = 60
	UnitTypeHMSPrinceOfWales   = 61
	UnitTypeB4Howitzer         = 64
	UnitTypeIS3HeavyTank       = 65
	UnitTypeStukaZuFuss        = 66
	UnitTypeRichelieu          = 67
	UnitTypeEnterprise         = 70
	UnitTypeRPGRocketeer       = 73
	UnitTypeA41Centurion       = 74
	UnitTypeAuF1               = 79
	UnitTypeT72                = 80
	UnitTypePhantomForce       = 82
	UnitTypeM1A1Abrams         = 83
	UnitTypeAH64Apache         = 85
	UnitTypeM142Himars         = 91
	UnitTypeDivineWrathMBT     = 92
)

// GetUnitTypeName returns a human-readable name for the unit type
func GetUnitTypeName(unitType uint8) string {
	switch unitType {
	case UnitTypeLightInfantry:
		return "Light Infantry"
	case UnitTypeAssaultInfantry:
		return "Assault Infantry"
	case UnitTypeMotorizedInfantry:
		return "Motorized Infantry"
	case UnitTypeMechanizedInfantry:
		return "Mechanized Infantry"
	case UnitTypeCommando:
		return "Commando"
	case UnitTypeArmoredCar:
		return "Armored Car"
	case UnitTypeLightTank:
		return "Light Tank"
	case UnitTypeMediumTank:
		return "Medium Tank"
	case UnitTypeHeavyTank:
		return "Heavy Tank"
	case UnitTypeSuperTank:
		return "Super Tank"
	case UnitTypeFieldArtillery:
		return "Field Artillery"
	case UnitTypeHowitzer:
		return "Howitzer"
	case UnitTypeRocketArtillery:
		return "Rocket Artillery"
	case UnitTypeSuperArtillery:
		return "Super Artillery"
	case UnitTypeSubmarine:
		return "Submarine"
	case UnitTypeDestroyer:
		return "Destroyer"
	case UnitTypeCruiser:
		return "Cruiser"
	case UnitTypeCarrier:
		return "Carrier"
	case UnitTypeSuperCarrier:
		return "Super Carrier"
	case UnitTypeBunker:
		return "Bunker"
	case UnitTypeLandFort:
		return "Land Fort"
	case UnitTypeCoastalArtillery:
		return "Coastal Artillery"
	case UnitTypeRocketLauncher:
		return "Rocket Launcher"
	case UnitTypeCity:
		return "City"
	case UnitTypeBrandenburgers:
		return "Brandenburgers"
	case UnitTypeHawkeyeForce:
		return "Hawkeye Force"
	case UnitTypeCombatMedic:
		return "Combat Medic"
	case UnitTypeM7Priest:
		return "M7 Priest"
	case UnitTypeHeavyGustav:
		return "Heavy Gustav"
	case UnitTypeBM21:
		return "BM-21"
	case UnitTypeT44:
		return "T-44"
	case UnitTypeKingTiger:
		return "King Tiger"
	case UnitTypeM26Pershing:
		return "M26 Pershing"
	case UnitTypeTypeVIISubmarine:
		return "Type VII Submarine"
	case UnitTypeHMSPrinceOfWales:
		return "HMS Prince of Wales"
	case UnitTypeB4Howitzer:
		return "B-4 Howitzer"
	case UnitTypeIS3HeavyTank:
		return "IS-3 Heavy Tank"
	case UnitTypeStukaZuFuss:
		return "Stuka zu Fuss"
	case UnitTypeRichelieu:
		return "Richelieu"
	case UnitTypeEnterprise:
		return "Enterprise"
	case UnitTypeRPGRocketeer:
		return "RPG Rocketeer"
	case UnitTypeA41Centurion:
		return "A41 Centurion"
	case UnitTypeAuF1:
		return "AuF1"
	case UnitTypeT72:
		return "T-72"
	case UnitTypePhantomForce:
		return "Phantom Force"
	case UnitTypeM1A1Abrams:
		return "M1A1 Abrams"
	case UnitTypeAH64Apache:
		return "AH-64 Apache"
	case UnitTypeM142Himars:
		return "M142 Himars"
	case UnitTypeDivineWrathMBT:
		return "Divine Wrath MBT"
	default:
		return fmt.Sprintf("Unknown Type %d", unitType)
	}
}
