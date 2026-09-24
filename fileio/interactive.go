package fileio

import (
	"fmt"
)

// UnitTileInfo holds information about a tile with a unit
type UnitTileInfo struct {
	X        int
	Y        int
	Owner    byte
	UnitType string
	UnitName string
}

type ProvinceInfo struct {
	CityIndex      int
	CoordinateCode uint16
	Name           string
	Owner          byte
	TileCount      int
	UnitCount      int
	PortCount      int
}

// InteractiveConvertPlayer provides an interactive interface for player conversion
func InteractiveConvertPlayer(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Player Conversion ===")

	// Show available players and get selection
	oldPlayer := getPlayerSelection("old player ID to convert FROM", saveOutput)
	if oldPlayer == -1 {
		return
	}

	newPlayer := getPlayerSelection("new player ID to convert TO", saveOutput)
	if newPlayer == -1 {
		return
	}

	// Show preview and get confirmation
	if showPlayerConversionPreview(saveOutput, oldPlayer, newPlayer) {
		performPlayerConversion(inputFilename, saveOutput, oldPlayer, newPlayer)
	}
}

// getPlayerSelection shows available players and gets user selection
func getPlayerSelection(prompt string, saveOutput *WC4SaveOutput) int {
	showAvailablePlayers(saveOutput)

	fmt.Printf("\nEnter %s: ", prompt)
	var playerStr string
	fmt.Scanln(&playerStr)
	player := parseInt(playerStr, prompt)
	if player == -1 {
		return -1
	}

	// Validate player exists
	if player < 0 || player >= len(saveOutput.PlayerData) {
		fmt.Printf("Error: Player ID %d does not exist\n", player)
		return -1
	}

	return player
}

// showAvailablePlayers displays all available players
func showAvailablePlayers(saveOutput *WC4SaveOutput) {
	fmt.Println("\nAvailable Players:")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, playerID := range sortedPlayerIDs {
		player := saveOutput.PlayerData[playerID]
		countryName, _ := GetCountryInfoFromData(player)
		fmt.Printf("  %d: %s (CountryId %d, TeamId %d)\n", playerID, countryName, player.CountryId, player.TeamId)
	}
}

// showPlayerConversionPreview shows conversion preview and returns true if user confirms
func showPlayerConversionPreview(saveOutput *WC4SaveOutput, oldPlayer, newPlayer int) bool {
	oldPlayerData := saveOutput.PlayerData[oldPlayer]
	newPlayerData := saveOutput.PlayerData[newPlayer]
	oldCountryName, _ := GetCountryInfoFromData(oldPlayerData)
	newCountryName, _ := GetCountryInfoFromData(newPlayerData)

	fmt.Printf("\n=== Conversion Preview ===\n")
	fmt.Printf("FROM: Player %d (%s - CountryId %d, TeamId %d)\n", oldPlayer, oldCountryName, oldPlayerData.CountryId, oldPlayerData.TeamId)
	fmt.Printf("TO:   Player %d (%s - CountryId %d, TeamId %d)\n", newPlayer, newCountryName, newPlayerData.CountryId, newPlayerData.TeamId)

	// Count affected tiles
	affectedCount := countPlayerTiles(saveOutput, oldPlayer)
	fmt.Printf("This will convert %d units/tiles from Player %d to Player %d\n", affectedCount, oldPlayer, newPlayer)

	// Ask for confirmation
	fmt.Print("\nProceed with conversion? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	return confirm == "y" || confirm == "Y" || confirm == "yes" || confirm == "Yes"
}

// countPlayerTiles counts how many tiles a player owns
func countPlayerTiles(saveOutput *WC4SaveOutput, playerID int) int {
	count := 0
	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] == byte(playerID) {
				count++
			}
		}
	}
	return count
}

// performPlayerConversion executes the player conversion
func performPlayerConversion(inputFilename string, saveOutput *WC4SaveOutput, oldPlayer, newPlayer int) {
	fmt.Println("\nPerforming conversion...")
	stats := ConvertPlayer(inputFilename, saveOutput, oldPlayer, newPlayer)
	stats.PrintSummary("Convert Player", saveOutput)
	fmt.Println("Conversion completed!")
}

func InteractiveConvertProvince(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Province Conversion ===")

	oldPlayer := getPlayerSelection("player ID whose province will be annexed", saveOutput)
	if oldPlayer == -1 {
		return
	}

	provinces := getPlayerProvinces(saveOutput, oldPlayer)
	if len(provinces) == 0 {
		fmt.Printf("Player %d has no province.\n", oldPlayer)
		return
	}

	fmt.Println("\nAvailable Provinces:")
	for i, province := range provinces {
		fmt.Printf("  %d: %s (position code %d, %d tiles, %d units/structures, %d ports)\n",
			i+1, province.Name, province.CoordinateCode, province.TileCount, province.UnitCount, province.PortCount)
	}

	fmt.Printf("\nSelect province (1-%d): ", len(provinces))
	var provinceStr string
	fmt.Scanln(&provinceStr)
	provinceIndex := parseInt(provinceStr, "province index")
	if provinceIndex < 1 || provinceIndex > len(provinces) {
		fmt.Printf("Error: Invalid province index %d\n", provinceIndex)
		return
	}
	province := provinces[provinceIndex-1]

	newPlayer := getPlayerSelection("new player ID to convert TO", saveOutput)
	if newPlayer == -1 || newPlayer == oldPlayer {
		if newPlayer == oldPlayer {
			fmt.Println("Error: Old player and new player cannot be the same")
		}
		return
	}

	fmt.Printf("\nAnnex %s from Player %d to Player %d?\n", province.Name, oldPlayer, newPlayer)
	fmt.Printf("This will transfer %d occupied tiles, %d ports and mines located in the province.\n", province.UnitCount, province.PortCount)
	fmt.Print("Proceed with conversion? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "Yes" {
		return
	}

	stats, minesConverted, err := ConvertProvince(inputFilename, saveOutput, province.CoordinateCode, newPlayer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Province converted: %d occupied tiles and %d landmines transferred.\n", stats.TotalChanged, minesConverted)
}

func getPlayerProvinces(saveOutput *WC4SaveOutput, playerID int) []ProvinceInfo {
	tileCounts := make(map[uint16]int)
	for row := range saveOutput.CityTiles {
		for _, coordinateCode := range saveOutput.CityTiles[row] {
			if coordinateCode != 65535 {
				tileCounts[coordinateCode]++
			}
		}
	}

	var provinces []ProvinceInfo
	for cityIndex, city := range saveOutput.Cities {
		tileCount := tileCounts[city.CoordinateCode]
		if tileCount == 0 {
			continue
		}
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid || saveOutput.UnitOwnerData[row][col] != byte(playerID) {
			continue
		}

		name, exists := GetCityName(city.CityId)
		if !exists {
			name = fmt.Sprintf("Province at (%d,%d)", row, col)
		}
		province := ProvinceInfo{
			CityIndex:      cityIndex,
			CoordinateCode: city.CoordinateCode,
			Name:           name,
			Owner:          byte(playerID),
			TileCount:      tileCount,
		}
		for tileRow := range saveOutput.CityTiles {
			for tileCol, coordinateCode := range saveOutput.CityTiles[tileRow] {
				if coordinateCode == city.CoordinateCode && saveOutput.UnitOwnerData[tileRow][tileCol] != TileUnowned {
					province.UnitCount++
				}
			}
		}
		for _, candidate := range saveOutput.Cities {
			if !isPort(candidate.BuildingType) {
				continue
			}
			portProvince, attached := provinceCodeForPort(saveOutput, candidate)
			if attached && portProvince == city.CoordinateCode {
				province.PortCount++
			}
		}
		provinces = append(provinces, province)
	}
	return provinces
}

func InteractiveChangeTeam(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Team Change ===")
	playerID := getPlayerSelection("player ID whose team will change", saveOutput)
	if playerID == -1 {
		return
	}

	fmt.Print("Enter new Team ID: ")
	var teamStr string
	fmt.Scanln(&teamStr)
	teamID := parseInt(teamStr, "Team ID")
	if teamID < 0 {
		return
	}

	oldTeamID := saveOutput.PlayerData[playerID].TeamId
	fmt.Printf("Change Player %d from Team %d to Team %d? (y/N): ", playerID, oldTeamID, teamID)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "Yes" {
		return
	}

	if err := ChangePlayerTeam(inputFilename, saveOutput, playerID, uint32(teamID)); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Player %d changed from Team %d to Team %d.\n", playerID, oldTeamID, teamID)
}

// InteractiveConvertTile provides an interactive interface for tile conversion
func InteractiveConvertTile(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Tile Conversion ===")

	// Find and display all tiles with units
	unitTiles := findAllUnitTiles(saveOutput)
	if len(unitTiles) == 0 {
		fmt.Println("No units found on the map.")
		return
	}

	// Show tiles and get user selection
	selectedTile := selectTileFromList(unitTiles, saveOutput)
	if selectedTile == nil {
		return
	}

	// Get new player and show preview
	newPlayer := getPlayerSelection("new player ID to convert TO", saveOutput)
	if newPlayer == -1 {
		return
	}

	// Show preview and get confirmation
	if showTileConversionPreview(selectedTile, newPlayer, saveOutput) {
		performTileConversion(inputFilename, saveOutput, selectedTile, newPlayer)
	}
}

// findAllUnitTiles scans the map and finds all tiles with units
func findAllUnitTiles(saveOutput *WC4SaveOutput) []UnitTileInfo {
	var unitTiles []UnitTileInfo

	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			owner := saveOutput.UnitOwnerData[i][j]
			if owner != TileUnowned { // Skip unowned tiles
				unitInfo := identifyUnitAtLocation(saveOutput, i, j, owner)
				unitTiles = append(unitTiles, unitInfo)
			}
		}
	}

	return unitTiles
}

// identifyUnitAtLocation identifies what unit/city is at a specific location
func identifyUnitAtLocation(saveOutput *WC4SaveOutput, x, y int, owner byte) UnitTileInfo {
	coordinateCode := ConvertToCoordinateCode(x, y, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))

	// Check if there's a unit at this coordinate
	for _, unit := range saveOutput.Units {
		if int(unit.CoordinateCode) == coordinateCode {
			unitTypeName := GetUnitTypeName(unit.UnitType)
			unitType := unitTypeName
			if unitType == "" {
				unitType = fmt.Sprintf("Type %d", unit.UnitType)
			}
			return UnitTileInfo{
				X:        x,
				Y:        y,
				Owner:    owner,
				UnitType: unitType,
				UnitName: unitType,
			}
		}
	}

	// Check if there's a city at this coordinate
	for _, city := range saveOutput.Cities {
		if int(city.CoordinateCode) == coordinateCode {
			cityName, exists := GetCityName(city.CityId)
			unitName := cityName
			if !exists {
				unitName = fmt.Sprintf("City (ID %d)", city.CityId)
			}
			return UnitTileInfo{
				X:        x,
				Y:        y,
				Owner:    owner,
				UnitType: "City",
				UnitName: unitName,
			}
		}
	}

	// Default to structure if nothing else found
	return UnitTileInfo{
		X:        x,
		Y:        y,
		Owner:    owner,
		UnitType: "Structure",
		UnitName: "Structure",
	}
}

// selectTileFromList displays tiles and gets user selection
func selectTileFromList(unitTiles []UnitTileInfo, saveOutput *WC4SaveOutput) *UnitTileInfo {
	// Show available tiles
	fmt.Printf("\nFound %d tiles with units:\n", len(unitTiles))

	// Create table data
	tableFormatter := NewTableFormatter()
	columns := []ColumnDef{
		{"Index", "int", "right"},
		{"Coordinates", "string", "center"},
		{"Owner", "string", "left"},
		{"Unit Type/Name", "string", "left"},
	}

	var rows [][]interface{}
	for i, tile := range unitTiles {
		ownerName := getOwnerDisplayName(tile.Owner, saveOutput)
		displayName := getUnitDisplayName(tile)
		row := []interface{}{
			i + 1,
			fmt.Sprintf("(%d,%d)", tile.X, tile.Y),
			ownerName,
			displayName,
		}
		rows = append(rows, row)
	}

	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})

	// Get tile selection
	fmt.Printf("\nSelect tile to convert (1-%d): ", len(unitTiles))
	var tileStr string
	fmt.Scanln(&tileStr)
	tileIndex := parseInt(tileStr, "tile index")
	if tileIndex == -1 {
		return nil
	}

	if tileIndex < 1 || tileIndex > len(unitTiles) {
		fmt.Printf("Error: Invalid tile index %d\n", tileIndex)
		return nil
	}

	return &unitTiles[tileIndex-1]
}

// getOwnerDisplayName returns a formatted string for the owner
func getOwnerDisplayName(owner byte, saveOutput *WC4SaveOutput) string {
	if int(owner) < len(saveOutput.PlayerData) {
		player := saveOutput.PlayerData[owner]
		countryName, _ := GetCountryInfoFromData(player)
		return fmt.Sprintf("P%d (%s)", owner, countryName)
	}
	return "Unknown"
}

// getUnitDisplayName returns the display name for a unit
func getUnitDisplayName(tile UnitTileInfo) string {
	if tile.UnitType == "City" {
		return tile.UnitName
	}
	return tile.UnitType
}

// showTileConversionPreview shows conversion preview and returns true if user confirms
func showTileConversionPreview(selectedTile *UnitTileInfo, newPlayer int, saveOutput *WC4SaveOutput) bool {
	newPlayerData := saveOutput.PlayerData[newPlayer]
	newCountryName, _ := GetCountryInfoFromData(newPlayerData)

	fmt.Printf("\n=== Conversion Preview ===\n")
	displayName := getUnitDisplayName(*selectedTile)
	fmt.Printf("Tile: (%d, %d) - %s\n", selectedTile.X, selectedTile.Y, displayName)
	fmt.Printf("Current owner: Player %d\n", selectedTile.Owner)
	fmt.Printf("New owner: Player %d (%s - CountryId %d, TeamId %d)\n", newPlayer, newCountryName, newPlayerData.CountryId, newPlayerData.TeamId)

	// Ask for confirmation
	fmt.Print("\nProceed with conversion? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	return confirm == "y" || confirm == "Y" || confirm == "yes" || confirm == "Yes"
}

// performTileConversion executes the tile conversion
func performTileConversion(inputFilename string, saveOutput *WC4SaveOutput, selectedTile *UnitTileInfo, newPlayer int) {
	fmt.Println("\nPerforming conversion...")
	err := ConvertTile(inputFilename, saveOutput, selectedTile.Y, selectedTile.X, newPlayer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Conversion completed!")
	}
}

// InteractiveSpawnUnits provides an interactive interface for spawning new
// units into a province.
func InteractiveSpawnUnits(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Unit Spawn ===")

	playerID := getPlayerSelection("player ID who will own the new units", saveOutput)
	if playerID == -1 {
		return
	}

	provinces := getPlayerProvinces(saveOutput, playerID)
	if len(provinces) == 0 {
		fmt.Printf("Player %d has no province.\n", playerID)
		return
	}

	fmt.Println("\nSpawn mode:")
	fmt.Println("  1: Manual (choose province, unit type, level, count)")
	fmt.Println("  2: Fun (fill every free tile of ALL your provinces with random units at max level; coastal water gets random ships)")
	fmt.Print("Select mode (1-2): ")
	var modeStr string
	fmt.Scanln(&modeStr)
	mode := parseInt(modeStr, "mode")
	if mode == 2 {
		fmt.Println("\nFill percentage:")
		for p := 10; p <= 100; p += 10 {
			fmt.Printf("  %d: %d%%\n", p/10, p)
		}
		fmt.Print("Select fill percentage (1-10): ")
		var pctStr string
		fmt.Scanln(&pctStr)
		pctChoice := parseInt(pctStr, "fill percentage")
		if pctChoice < 1 || pctChoice > 10 {
			fmt.Printf("Error: Invalid selection %d\n", pctChoice)
			return
		}
		performFunSpawn(inputFilename, saveOutput, playerID, provinces, pctChoice*10)
		return
	}
	if mode != 1 {
		fmt.Printf("Error: Invalid mode %d\n", mode)
		return
	}

	fmt.Println("\nAvailable Provinces:")
	for i, province := range provinces {
		fmt.Printf("  %d: %s (%d tiles, %d ports)\n", i+1, province.Name, province.TileCount, province.PortCount)
	}
	fmt.Printf("\nSelect province (1-%d): ", len(provinces))
	var provinceStr string
	fmt.Scanln(&provinceStr)
	provinceIndex := parseInt(provinceStr, "province index")
	if provinceIndex < 1 || provinceIndex > len(provinces) {
		fmt.Printf("Error: Invalid province index %d\n", provinceIndex)
		return
	}
	province := provinces[provinceIndex-1]

	unitTypes := SpawnableUnitTypes()
	fmt.Println("\nAvailable Unit Types:")
	for i, unitType := range unitTypes {
		kind := "land"
		if isNavalUnitType(unitType) {
			kind = "naval"
		}
		fmt.Printf("  %d: %s (%s)\n", i+1, GetUnitTypeName(unitType), kind)
	}
	fmt.Printf("\nSelect unit type (1-%d): ", len(unitTypes))
	var unitTypeStr string
	fmt.Scanln(&unitTypeStr)
	unitTypeIndex := parseInt(unitTypeStr, "unit type index")
	if unitTypeIndex < 1 || unitTypeIndex > len(unitTypes) {
		fmt.Printf("Error: Invalid unit type index %d\n", unitTypeIndex)
		return
	}
	unitType := unitTypes[unitTypeIndex-1]

	fmt.Print("\nEnter level (1-9): ")
	var levelStr string
	fmt.Scanln(&levelStr)
	level := parseInt(levelStr, "level")
	if level < 1 || level > 9 {
		fmt.Printf("Error: Invalid level %d. Must be between 1 and 9\n", level)
		return
	}

	fmt.Print("Enter number of units to spawn: ")
	var countStr string
	fmt.Scanln(&countStr)
	count := parseInt(countStr, "count")
	if count < 1 {
		fmt.Printf("Error: Invalid count %d. Must be at least 1\n", count)
		return
	}

	fmt.Printf("\nSpawn %d x %s (level %d) in %s for Player %d?\n", count, GetUnitTypeName(unitType), level, province.Name, playerID)
	fmt.Print("Proceed? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "Yes" {
		return
	}

	spawned, err := SpawnUnits(inputFilename, saveOutput, province.CoordinateCode, unitType, level, count)
	if err != nil {
		fmt.Printf("Error: %v. Nothing was spawned.\n", err)
		return
	}
	fmt.Printf("Spawned %d x %s (level %d) in %s for Player %d.\n", spawned, GetUnitTypeName(unitType), level, province.Name, playerID)
}

// performFunSpawn runs "fun mode" (see SpawnFunUnits) across every province
// playerID owns: fillPercent% of each province's free tiles get a random
// unit at max level, and the same share of coastal water near its port(s)
// gets random ships.
func performFunSpawn(inputFilename string, saveOutput *WC4SaveOutput, playerID int, provinces []ProvinceInfo, fillPercent int) {
	fmt.Printf("\nFun mode will fill %d%% of the free tiles of Player %d's %d province(s) with random units at max level (9), including nearby sea tiles with random ships.\n", fillPercent, playerID, len(provinces))
	fmt.Print("Proceed? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "Yes" {
		return
	}

	fmt.Println("\nSpawning...")
	totalLand, totalNaval := 0, 0
	for _, province := range provinces {
		result, err := SpawnFunUnits(inputFilename, saveOutput, province.CoordinateCode, fillPercent)
		if err != nil {
			fmt.Printf("  %s: %v\n", province.Name, err)
			continue
		}
		fmt.Printf("  %s: %d land unit(s), %d naval unit(s)\n", province.Name, result.LandSpawned, result.NavalSpawned)
		totalLand += result.LandSpawned
		totalNaval += result.NavalSpawned
	}
	fmt.Printf("\nFun mode completed: %d land unit(s) and %d naval unit(s) spawned across %d province(s) for Player %d.\n", totalLand, totalNaval, len(provinces), playerID)
}

// parseInt is a helper function to parse integer input with error handling
func parseInt(input, fieldName string) int {
	var result int
	_, err := fmt.Sscanf(input, "%d", &result)
	if err != nil {
		fmt.Printf("Error: Invalid %s '%s'. Please enter a valid number.\n", fieldName, input)
		return -1
	}
	return result
}
