package fileio

import (
	"os"
	"testing"
)

func TestProcessPlayerCityTechUpgradesPortAsPort(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_port_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	data := make([]byte, 64)
	data[4] = PortLevel1
	for i := 24; i < 30; i++ {
		data[i] = 4 // damaged by the old max-city-tech implementation
	}
	if _, err := tempFile.Write(data); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{BuildCityStartKey(0): 0}
	save := &WC4SaveOutput{
		SaveHeader:    SaveHeader{GameMode: GameModeCampaign},
		PlayerData:    []CountryData{{}},
		UnitOwnerData: [][]byte{{0}},
		Cities:        []CityData{{CoordinateCode: 0, BuildingType: PortLevel1}},
	}

	if got := processPlayerCityTech(tempFile.Name(), save, 0, 4); got != 1 {
		t.Fatalf("modified %d ports, want 1", got)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), 4); got != PortLevel4 {
		t.Errorf("port building type = %d, want %d", got, PortLevel4)
	}
	for offset := 24; offset < 30; offset++ {
		if got := ReadUint8AtFileOffset(tempFile.Name(), offset); got != 0 {
			t.Errorf("port tech byte at %d = %d, want unchanged 0", offset, got)
		}
	}
}

func TestProcessPlayerCityLevelUpgradesLowLevelCities(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_city_level_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	// 3 cities of 32 bytes each: level II city, capital "star" city, port
	data := make([]byte, 32*3)
	data[0*32+4] = CityLevelBase + 2 // Bordeaux-like: level II
	data[1*32+4] = CityLevelBase + 5 // Paris-like: capital, level 5 ("star")
	data[2*32+4] = PortLevel1        // port, must be skipped
	if _, err := tempFile.Write(data); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{
		BuildCityStartKey(0): 0,
		BuildCityStartKey(1): 32,
		BuildCityStartKey(2): 64,
	}
	save := &WC4SaveOutput{
		SaveHeader:    SaveHeader{GameMode: GameModeCampaign},
		PlayerData:    []CountryData{{}},
		UnitOwnerData: [][]byte{{0, 0, 0}},
		Cities: []CityData{
			{CoordinateCode: 0, BuildingType: CityLevelBase + 2},
			{CoordinateCode: 1, BuildingType: CityLevelBase + 5},
			{CoordinateCode: 2, BuildingType: PortLevel1},
		},
	}

	if got := processPlayerCityLevel(tempFile.Name(), save, 0, 4); got != 1 {
		t.Fatalf("modified %d cities, want 1", got)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), 4); got != CityLevelBase+4 {
		t.Errorf("level II city building type = %d, want %d", got, CityLevelBase+4)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), 32+4); got != CityLevelBase+5 {
		t.Errorf("capital city building type changed to %d, want unchanged %d", got, CityLevelBase+5)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), 64+4); got != PortLevel1 {
		t.Errorf("port building type changed to %d, want unchanged %d", got, PortLevel1)
	}
}

func TestConvertOwnershipMetadata(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_ownership_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(make([]byte, 128)); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{
		BuildUnitStartKey(0):     0,
		BuildLandmineStartKey(0): 64,
		BuildLandmineStartKey(1): 76,
	}
	save := &WC4SaveOutput{
		SaveHeader: SaveHeader{TurnNumber: 7},
		PlayerData: []CountryData{{TeamId: 1}, {TeamId: 2}},
		Units:      []UnitData{{CoordinateCode: 42}},
		Landmines: []LandmineData{
			{Owner: 1}, // already belongs to player 0
			{Owner: 2}, // belongs to player 1 (one-based encoding)
		},
	}

	positions := map[uint16]struct{}{42: {}}
	if got := convertUnitOwnershipMetadata(tempFile.Name(), save, positions); got != 1 {
		t.Fatalf("converted %d units, want 1", got)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), unitActivationTurnOffset); got != 7 {
		t.Errorf("activation turn = %d, want 7", got)
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), unitActionStateOffset); got != unitReadyActionState {
		t.Errorf("unit action state = %d, want %d", got, unitReadyActionState)
	}

	if got := convertLandmineOwners(tempFile.Name(), save, 1, 0); got != 1 {
		t.Fatalf("converted %d landmines, want 1", got)
	}
	if got := ReadUint16AtFileOffset(tempFile.Name(), 64+landmineOwnerOffset); got != 0 {
		t.Errorf("unrelated file bytes changed to %d, want 0", got)
	}
	if got := ReadUint16AtFileOffset(tempFile.Name(), 76+landmineOwnerOffset); got != 1 {
		t.Errorf("landmine owner = %d, want encoded player 0 value 1", got)
	}
}

func TestCollectPlayerPositionsIncludesPreviouslyConvertedUnits(t *testing.T) {
	save := &WC4SaveOutput{
		SaveHeader: SaveHeader{GameMode: GameModeCampaign},
		UnitOwnerData: [][]byte{
			{0, 1},
			{0, TileUnowned},
		},
	}

	positions := collectPlayerPositions(save, 0)
	for _, coordinate := range []uint16{0, 2} {
		if _, ok := positions[coordinate]; !ok {
			t.Errorf("player 0 position %d was not collected", coordinate)
		}
	}
	if len(positions) != 2 {
		t.Errorf("collected %d positions, want 2", len(positions))
	}
}

func TestConvertProvinceTransfersOnlySelectedTerritory(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_province_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(make([]byte, 64)); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{
		buildUnitOwnerStartKey(): 0,
		BuildLandmineStartKey(0): 16,
		BuildLandmineStartKey(1): 28,
	}
	save := &WC4SaveOutput{
		SaveHeader: SaveHeader{GameMode: GameModeCampaign},
		PlayerData: []CountryData{{}, {}, {}},
		CityTiles: [][]uint16{
			{10, 10},
			{20, 65535},
		},
		UnitOwnerData: [][]byte{
			{1, 2},
			{1, TileUnowned},
		},
		Landmines: []LandmineData{
			{CoordinateCode: 0, Owner: 2},
			{CoordinateCode: 2, Owner: 2},
		},
	}

	stats, mines, err := ConvertProvince(tempFile.Name(), save, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalChanged != 2 {
		t.Errorf("changed %d occupied tiles, want 2", stats.TotalChanged)
	}
	if mines != 1 {
		t.Errorf("changed %d mines, want 1", mines)
	}
	if save.UnitOwnerData[0][0] != 0 || save.UnitOwnerData[0][1] != 0 {
		t.Errorf("selected province owners = %v, want [0 0]", save.UnitOwnerData[0])
	}
	if save.UnitOwnerData[1][0] != 1 {
		t.Errorf("other province owner = %d, want 1", save.UnitOwnerData[1][0])
	}
	if save.Landmines[0].Owner != 1 || save.Landmines[1].Owner != 2 {
		t.Errorf("mine owners = [%d %d], want [1 2]", save.Landmines[0].Owner, save.Landmines[1].Owner)
	}
}

func TestConvertProvinceIncludesAttachedPort(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_province_port_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(make([]byte, 64)); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{buildUnitOwnerStartKey(): 0}
	save := &WC4SaveOutput{
		SaveHeader: SaveHeader{GameMode: GameModeCampaign},
		PlayerData: []CountryData{{}, {}},
		CityTiles: [][]uint16{
			{65535, 65535, 65535},
			{4, 4, 4},
			{4, 4, 4},
		},
		UnitOwnerData: [][]byte{
			{TileUnowned, 1, TileUnowned},
			{TileUnowned, 1, TileUnowned},
			{TileUnowned, TileUnowned, TileUnowned},
		},
		Cities: []CityData{
			{CoordinateCode: 4, BuildingType: 11},
			{CoordinateCode: 1, BuildingType: PortLevel1},
		},
	}

	stats, _, err := ConvertProvince(tempFile.Name(), save, 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalChanged != 2 {
		t.Errorf("changed %d occupied tiles, want city and port", stats.TotalChanged)
	}
	if save.UnitOwnerData[0][1] != 0 {
		t.Errorf("attached port owner = %d, want 0", save.UnitOwnerData[0][1])
	}
}

func TestChangePlayerTeamChangesOnlySelectedPlayer(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_team_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(make([]byte, 64)); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{BuildPlayerStartKey(1): 0}
	save := &WC4SaveOutput{PlayerData: []CountryData{{TeamId: 1}, {TeamId: 2}}}
	if err := ChangePlayerTeam(tempFile.Name(), save, 1, 7); err != nil {
		t.Fatal(err)
	}
	if save.PlayerData[0].TeamId != 1 || save.PlayerData[1].TeamId != 7 {
		t.Errorf("team IDs = [%d %d], want [1 7]", save.PlayerData[0].TeamId, save.PlayerData[1].TeamId)
	}
	if got := ReadUint32AtFileOffset(tempFile.Name(), 24); got != 7 {
		t.Errorf("written Team ID = %d, want 7", got)
	}
}

func TestSetPlayerMaxMoraleOnlyChangesCombatUnits(t *testing.T) {
	tempFile, err := os.CreateTemp("", "wc4_morale_*.sav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(make([]byte, 256)); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileOffsetMap = map[string]int{
		BuildUnitStartKey(0): 0,
		BuildUnitStartKey(1): 64,
		BuildUnitStartKey(2): 128,
		BuildUnitStartKey(3): 192,
	}
	save := &WC4SaveOutput{
		SaveHeader:    SaveHeader{GameMode: GameModeCampaign},
		PlayerData:    []CountryData{{}, {}},
		UnitOwnerData: [][]byte{{0, 0, 0, 1}},
		Units: []UnitData{
			{CoordinateCode: 0, UnitType: UnitTypeLightInfantry, MoraleValue: -1, MoraleTurnsLeft: 1},
			{CoordinateCode: 1, UnitType: UnitTypeLandFort},
			{CoordinateCode: 2, UnitType: UnitTypeCity},
			{CoordinateCode: 3, UnitType: UnitTypeMediumTank, MoraleValue: -1, MoraleTurnsLeft: 1},
		},
	}

	if got := SetPlayerMaxMorale(tempFile.Name(), save, 0); got != 1 {
		t.Fatalf("updated %d units, want 1", got)
	}
	if save.Units[0].MoraleValue != maxMoraleValue || save.Units[0].MoraleTurnsLeft != maxMoraleTurns {
		t.Errorf("combat unit morale = %d/%d, want %d/%d", save.Units[0].MoraleValue, save.Units[0].MoraleTurnsLeft, maxMoraleValue, maxMoraleTurns)
	}
	for _, index := range []int{1, 2} {
		if save.Units[index].MoraleValue != 0 || save.Units[index].MoraleTurnsLeft != 0 {
			t.Errorf("fortification unit %d morale was unexpectedly changed", index)
		}
	}
	if save.Units[3].MoraleValue != -1 || save.Units[3].MoraleTurnsLeft != 1 {
		t.Error("enemy unit morale was unexpectedly changed")
	}
	if got := ReadUint8AtFileOffset(tempFile.Name(), unitMoraleValueOffset); got != maxMoraleValue {
		t.Errorf("written morale value = %d, want %d", got, maxMoraleValue)
	}
	if got := ReadUint16AtFileOffset(tempFile.Name(), unitMoraleTurnsOffset); got != maxMoraleTurns {
		t.Errorf("written morale turns = %d, want %d", got, maxMoraleTurns)
	}
}

func TestConvertCoordinates(t *testing.T) {
	// Create test unit owner data (3x3 grid)
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	tests := []struct {
		name           string
		coordinateCode int
		gameMode       int
		expectedRow    int
		expectedCol    int
		expectedValid  bool
	}{
		{
			name:           "Valid coordinates - Campaign mode, top-left",
			coordinateCode: 0,
			gameMode:       GameModeCampaign,
			expectedRow:    0,
			expectedCol:    0,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Campaign mode, center",
			coordinateCode: 4,
			gameMode:       GameModeCampaign,
			expectedRow:    1,
			expectedCol:    1,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Campaign mode, bottom-right",
			coordinateCode: 8,
			gameMode:       GameModeCampaign,
			expectedRow:    2,
			expectedCol:    2,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Conquest mode, top-left",
			coordinateCode: 6, // 0 + 2*3
			gameMode:       GameModeConquest,
			expectedRow:    0,
			expectedCol:    0,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Conquest mode, center",
			coordinateCode: 10, // 4 + 2*3
			gameMode:       GameModeConquest,
			expectedRow:    1,
			expectedCol:    1,
			expectedValid:  true,
		},
		{
			name:           "Invalid coordinates - negative",
			coordinateCode: -1,
			gameMode:       GameModeCampaign,
			expectedRow:    -1,
			expectedCol:    2,
			expectedValid:  false,
		},
		{
			name:           "Invalid coordinates - too high for campaign",
			coordinateCode: 9,
			gameMode:       GameModeCampaign,
			expectedRow:    3,
			expectedCol:    0,
			expectedValid:  false,
		},
		{
			name:           "Invalid coordinates - too high for conquest",
			coordinateCode: 15,
			gameMode:       GameModeConquest,
			expectedRow:    3,
			expectedCol:    0,
			expectedValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, valid := ConvertCoordinates(tt.coordinateCode, unitOwnerData, tt.gameMode)
			if valid != tt.expectedValid {
				t.Errorf("ConvertCoordinates(%d, %d) valid = %v, want %v", tt.coordinateCode, tt.gameMode, valid, tt.expectedValid)
			}
			if valid {
				if row != tt.expectedRow {
					t.Errorf("ConvertCoordinates(%d, %d) row = %d, want %d", tt.coordinateCode, tt.gameMode, row, tt.expectedRow)
				}
				if col != tt.expectedCol {
					t.Errorf("ConvertCoordinates(%d, %d) col = %d, want %d", tt.coordinateCode, tt.gameMode, col, tt.expectedCol)
				}
			}
		})
	}
}

func TestConvertToCoordinateCode(t *testing.T) {
	// Create test unit owner data (3x3 grid)
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	tests := []struct {
		name         string
		row          int
		col          int
		gameMode     int
		expectedCode int
	}{
		{
			name:         "Campaign mode - top-left",
			row:          0,
			col:          0,
			gameMode:     GameModeCampaign,
			expectedCode: 0,
		},
		{
			name:         "Campaign mode - center",
			row:          1,
			col:          1,
			gameMode:     GameModeCampaign,
			expectedCode: 4,
		},
		{
			name:         "Campaign mode - bottom-right",
			row:          2,
			col:          2,
			gameMode:     GameModeCampaign,
			expectedCode: 8,
		},
		{
			name:         "Conquest mode - top-left",
			row:          0,
			col:          0,
			gameMode:     GameModeConquest,
			expectedCode: 6, // 0 + 2*3
		},
		{
			name:         "Conquest mode - center",
			row:          1,
			col:          1,
			gameMode:     GameModeConquest,
			expectedCode: 10, // 4 + 2*3
		},
		{
			name:         "Conquest mode - bottom-right",
			row:          2,
			col:          2,
			gameMode:     GameModeConquest,
			expectedCode: 14, // 8 + 2*3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCoordinateCode(tt.row, tt.col, unitOwnerData, tt.gameMode)
			if result != tt.expectedCode {
				t.Errorf("ConvertToCoordinateCode(%d, %d, %d) = %d, want %d", tt.row, tt.col, tt.gameMode, result, tt.expectedCode)
			}
		})
	}
}

// Test round-trip conversion (coordinate -> code -> coordinate)
func TestCoordinateRoundTrip(t *testing.T) {
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	gameModes := []int{GameModeCampaign, GameModeConquest}

	for _, gameMode := range gameModes {
		for row := 0; row < len(unitOwnerData); row++ {
			for col := 0; col < len(unitOwnerData[0]); col++ {
				// Convert row,col to coordinate code
				code := ConvertToCoordinateCode(row, col, unitOwnerData, gameMode)

				// Convert coordinate code back to row,col
				resultRow, resultCol, valid := ConvertCoordinates(code, unitOwnerData, gameMode)

				if !valid {
					t.Errorf("Round-trip conversion failed: (%d,%d) -> %d -> invalid", row, col, code)
					continue
				}

				if resultRow != row || resultCol != col {
					t.Errorf("Round-trip conversion failed: (%d,%d) -> %d -> (%d,%d)",
						row, col, code, resultRow, resultCol)
				}
			}
		}
	}
}

// Test with different map sizes
func TestConvertCoordinatesDifferentSizes(t *testing.T) {
	testCases := []struct {
		name          string
		width         int
		height        int
		code          int
		gameMode      int
		expectedRow   int
		expectedCol   int
		expectedValid bool
	}{
		{
			name:          "1x1 map - Campaign",
			width:         1,
			height:        1,
			code:          0,
			gameMode:      GameModeCampaign,
			expectedRow:   0,
			expectedCol:   0,
			expectedValid: true,
		},
		{
			name:          "1x1 map - Conquest",
			width:         1,
			height:        1,
			code:          2, // 0 + 2*1
			gameMode:      GameModeConquest,
			expectedRow:   0,
			expectedCol:   0,
			expectedValid: true,
		},
		{
			name:          "2x2 map - Campaign",
			width:         2,
			height:        2,
			code:          3,
			gameMode:      GameModeCampaign,
			expectedRow:   1,
			expectedCol:   1,
			expectedValid: true,
		},
		{
			name:          "2x2 map - Conquest",
			width:         2,
			height:        2,
			code:          7, // 3 + 2*2
			gameMode:      GameModeConquest,
			expectedRow:   1,
			expectedCol:   1,
			expectedValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create unit owner data with specified dimensions
			unitOwnerData := make([][]byte, tc.height)
			for i := range unitOwnerData {
				unitOwnerData[i] = make([]byte, tc.width)
			}

			row, col, valid := ConvertCoordinates(tc.code, unitOwnerData, tc.gameMode)
			if valid != tc.expectedValid {
				t.Errorf("ConvertCoordinates(%d, %d) valid = %v, want %v", tc.code, tc.gameMode, valid, tc.expectedValid)
			}
			if valid {
				if row != tc.expectedRow {
					t.Errorf("ConvertCoordinates(%d, %d) row = %d, want %d", tc.code, tc.gameMode, row, tc.expectedRow)
				}
				if col != tc.expectedCol {
					t.Errorf("ConvertCoordinates(%d, %d) col = %d, want %d", tc.code, tc.gameMode, col, tc.expectedCol)
				}
			}
		})
	}
}
