package fileio

import (
	"fmt"
	"sort"
	"strings"
)

// UnitDisplayInfo represents a unit with its display information
type UnitDisplayInfo struct {
	Index int
	Unit  UnitData
	Row   int
	Col   int
	Owner byte
}

// TileInfo represents a tile position
type TileInfo struct {
	Row, Col int
}

// TileAnalysisResult contains the results of tile analysis
type TileAnalysisResult struct {
	CoordinateCodeCounts map[uint16]int
	TotalTiles           int
	OceanCount           int
	LandTiles            int
}

// GetSortedUnitTypes returns unit types sorted in ascending order
func GetSortedUnitTypes(unitsByType map[byte]int) []int {
	var unitTypes []int
	for unitType := range unitsByType {
		unitTypes = append(unitTypes, int(unitType))
	}
	sort.Ints(unitTypes)
	return unitTypes
}

// GetSortedPlayerIDs returns player IDs sorted by their index
func GetSortedPlayerIDs(playerData []CountryData) []int {
	playerIDs := make([]int, len(playerData))
	for i := range playerData {
		playerIDs[i] = i
	}
	return playerIDs
}

// DisplayHeader prints a formatted header with title and separator
func DisplayHeader(title string) {
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len(title)))
}

// DisplayPlayerInfo returns formatted player information string
func DisplayPlayerInfo(playerID int, player CountryData, count int) string {
	countryName, _ := GetCountryInfoFromData(player)
	return fmt.Sprintf("Player %d (%s - CountryId %d) owns %d units", playerID, countryName, player.CountryId, count)
}

// SortByCount sorts a slice of items by their count in descending order
func SortByCount[T any](items []T, getCount func(T) int) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if getCount(items[i]) < getCount(items[j]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// DisplayTypeAnalysis shows unit type analysis with percentages
func DisplayTypeAnalysis[T any](items []T, getType func(T) uint8, getCount func(T) int, itemName string) {
	typeCounts := make(map[uint8]int)
	for _, item := range items {
		typeCounts[getType(item)]++
	}

	type TypeCount struct {
		UnitType uint8
		Count    int
		Name     string
	}
	var typeStats []TypeCount
	for unitType, count := range typeCounts {
		typeStats = append(typeStats, TypeCount{
			UnitType: unitType,
			Count:    count,
			Name:     GetUnitTypeName(unitType),
		})
	}

	SortByCount(typeStats, func(tc TypeCount) int { return tc.Count })

	// Use table formatting for better display
	tableFormatter := NewTableFormatter()
	columns := []ColumnDef{
		{"Type", "int", "right"},
		{"Unit Type", "string", "left"},
		{"Count", "int", "right"},
		{"Percentage", "percent", "right"},
	}

	var rows [][]interface{}
	// Calculate total count by summing all type counts
	totalCount := 0
	for _, stat := range typeStats {
		totalCount += stat.Count
	}

	for _, stat := range typeStats {
		percentage := float64(stat.Count) / float64(totalCount)
		row := []interface{}{
			stat.UnitType,
			stat.Name,
			stat.Count,
			percentage,
		}
		rows = append(rows, row)
	}

	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
}

// DisplayCoordinatesGrid prints coordinates in a grid format (10 per row)
func DisplayCoordinatesGrid(coordinates []string) {
	countPerRow := 10

	for i, coord := range coordinates {
		if i%countPerRow == 0 {
			fmt.Printf("  ")
		}

		fmt.Printf("%s", coord)

		if (i+1)%countPerRow == 0 {
			fmt.Println()
		} else {
			fmt.Printf(" ")
		}
	}

	if len(coordinates)%countPerRow != 0 {
		fmt.Println()
	}
}

// ProcessAllUnits processes all units and returns valid units with coordinates
func ProcessAllUnits(saveOutput *WC4SaveOutput) []UnitDisplayInfo {
	var validUnits []UnitDisplayInfo

	for i, unit := range saveOutput.Units {
		row, col, valid := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			continue
		}
		owner := saveOutput.UnitOwnerData[row][col]

		validUnits = append(validUnits, UnitDisplayInfo{
			Index: i,
			Unit:  unit,
			Row:   row,
			Col:   col,
			Owner: owner,
		})
	}

	return validUnits
}

// GroupUnitsByOwner groups units by their owner
func GroupUnitsByOwner(units []UnitDisplayInfo) map[byte][]UnitDisplayInfo {
	unitsByOwner := make(map[byte][]UnitDisplayInfo)
	for _, unitInfo := range units {
		unitsByOwner[unitInfo.Owner] = append(unitsByOwner[unitInfo.Owner], unitInfo)
	}
	return unitsByOwner
}

// GroupPlayersByTeam groups players by their team ID
func GroupPlayersByTeam(playerData []CountryData) map[uint32][]int {
	playersByTeam := make(map[uint32][]int)
	for i := 0; i < len(playerData); i++ {
		teamId := playerData[i].TeamId
		playersByTeam[teamId] = append(playersByTeam[teamId], i)
	}
	return playersByTeam
}

// GroupPlayersByTeamWithStats groups players by team and calculates team statistics
func GroupPlayersByTeamWithStats(playerData []CountryData, unitCounts map[byte]int, cityCounts map[byte]int) map[uint32]*TeamAnalysisData {
	teamData := make(map[uint32]*TeamAnalysisData)

	for i, player := range playerData {
		teamId := player.TeamId

		// Initialize team data if not exists
		if teamData[teamId] == nil {
			teamData[teamId] = &TeamAnalysisData{
				TeamID:      int(teamId),
				PlayerCount: 0,
				CityCount:   0,
				UnitCount:   0,
				Players:     []string{},
			}
		}

		// Add player data to team
		countryName, _ := GetCountryInfoFromData(player)
		teamData[teamId].PlayerCount++
		teamData[teamId].CityCount += cityCounts[byte(i)]
		teamData[teamId].UnitCount += unitCounts[byte(i)]
		teamData[teamId].Players = append(teamData[teamId].Players, fmt.Sprintf("%d (%s)", i, countryName))
	}

	return teamData
}

func GroupPlayersByTeamWithTerritory(playerData []CountryData, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) map[uint32]*TeamTerritoryData {
	teamData := make(map[uint32]*TeamTerritoryData)

	for _, ownerID := range GetSortedPlayerIDs(playerData) {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			player := playerData[ownerID]
			teamId := player.TeamId

			// Calculate player tiles
			playerTiles := 0
			for _, cityInfo := range cities {
				playerTiles += result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
			}

			// Initialize team data if not exists
			if teamData[teamId] == nil {
				teamData[teamId] = &TeamTerritoryData{
					TeamID:      int(teamId),
					PlayerCount: 0,
					TileCount:   0,
					Players:     []string{},
				}
			}

			// Add player data to team
			countryName, _ := GetCountryInfoFromData(player)
			teamData[teamId].PlayerCount++
			teamData[teamId].TileCount += playerTiles
			teamData[teamId].Players = append(teamData[teamId].Players, fmt.Sprintf("%d (%s)", ownerID, countryName))
		}
	}

	return teamData
}

// CountPlayerTiles counts tiles owned by each player
func CountPlayerTiles(unitOwnerData [][]byte) map[byte]int {
	countMap := make(map[byte]int)
	for i := 0; i < len(unitOwnerData); i++ {
		for j := 0; j < len(unitOwnerData[i]); j++ {
			if unitOwnerData[i][j] == TileUnowned {
				continue
			}
			owner := unitOwnerData[i][j]
			existingCount, ok := countMap[owner]
			if !ok {
				countMap[owner] = 1
			} else {
				countMap[owner] = existingCount + 1
			}
		}
	}
	return countMap
}

// ListPlayers displays all players with their information
func ListPlayers(saveOutput *WC4SaveOutput) {
	countMap := CountPlayerTiles(saveOutput.UnitOwnerData)
	tableFormatter := NewTableFormatter()

	fmt.Println("Players:")
	fmt.Println("--------")

	// First, show linear list
	fmt.Println("Linear List:")

	// Inline player table logic
	columns := []ColumnDef{
		{"Player", "int", "right"},
		{"Country", "string", "left"},
		{"Country ID", "int", "right"},
		{"Team", "int", "right"},
		{"Units Owned", "int", "right"},
	}

	var rows [][]interface{}
	for i, player := range saveOutput.PlayerData {
		countryName, _ := GetCountryInfoFromData(player)
		row := []interface{}{
			i,
			countryName,
			int(player.CountryId),
			int(player.TeamId),
			countMap[byte(i)],
		}
		rows = append(rows, row)
	}

	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})

	// Then, group players by TeamId
	fmt.Println("\nGrouped by Team:")
	playersByTeam := GroupPlayersByTeam(saveOutput.PlayerData)

	// Display players grouped by team
	for teamId, playerIndices := range playersByTeam {
		fmt.Printf("\nTeam %d (%d players):\n", teamId, len(playerIndices))

		// Use custom table for team display (without Team column)
		columns := []ColumnDef{
			{"Player", "int", "right"},
			{"Country", "string", "left"},
			{"Country ID", "int", "right"},
			{"Units Owned", "int", "right"},
		}

		var rows [][]interface{}
		for _, playerIndex := range playerIndices {
			player := saveOutput.PlayerData[playerIndex]
			countryName, _ := GetCountryInfoFromData(player)
			row := []interface{}{
				playerIndex,
				countryName,
				int(player.CountryId),
				countMap[byte(playerIndex)],
			}
			rows = append(rows, row)
		}

		tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
	}
}

// ListCities displays cities in both linear and grouped format
func ListCities(saveOutput *WC4SaveOutput) {
	fmt.Printf("Cities (Map: %dx%d):\n", len(saveOutput.UnitOwnerData), len(saveOutput.UnitOwnerData[0]))
	fmt.Println("-------------------")

	// Process cities using the fileio package
	validCities := ProcessCitiesForDisplay(saveOutput)

	// First, show linear list
	fmt.Println("Linear List:")
	for _, cityInfo := range validCities {
		fmt.Printf("CityInfo[%d]: %+v\n", cityInfo.Index, cityInfo)
	}

	// Then, group cities by owner
	fmt.Println("\nGrouped by Owner:")
	citiesByOwner := GroupCitiesByOwner(validCities)

	// Display cities grouped by owner
	tableFormatter := NewTableFormatter()
	for owner := byte(0); owner < byte(len(saveOutput.PlayerData)); owner++ {
		if cities, exists := citiesByOwner[owner]; exists {
			countryName, _ := GetCountryInfoFromData(saveOutput.PlayerData[owner])
			fmt.Printf("\nPlayer %d (%s - CountryId %d) owns %d cities:\n", owner, countryName, saveOutput.PlayerData[owner].CountryId, len(cities))

			// Inline city table
			columns := []ColumnDef{
				{"City", "int", "right"},
				{"Name", "string", "left"},
				{"Position", "string", "right"},
			}
			var rows [][]interface{}
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					if cityInfo.City.CityId > 0 {
						cityName = fmt.Sprintf("City ID %d", cityInfo.City.CityId)
					} else {
						cityName = "Unknown"
					}
				}
				positionStr := fmt.Sprintf("(%d,%d)", cityInfo.Row, cityInfo.Col)
				rows = append(rows, []interface{}{cityInfo.Index, cityName, positionStr})
			}
			tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
		}
	}
}

// ListUnits displays units in both linear and grouped format
func ListUnits(saveOutput *WC4SaveOutput) {
	DisplayHeader(fmt.Sprintf("Units (Total: %d)", len(saveOutput.Units)))

	// Process all units using common function
	validUnits := ProcessAllUnits(saveOutput)

	// First, show linear list
	fmt.Println("Linear List:")
	for _, unitInfo := range validUnits {
		fmt.Printf("Unit %d (owner: %d): %+v\n", unitInfo.Index, unitInfo.Owner, unitInfo.Unit)
	}

	// Then, group units by owner
	fmt.Println("\nGrouped by Owner:")
	unitsByOwner := GroupUnitsByOwner(validUnits)

	// Display units grouped by owner using sorted player IDs
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if units, exists := unitsByOwner[byte(ownerID)]; exists {
			playerInfo := DisplayPlayerInfo(ownerID, saveOutput.PlayerData[ownerID], len(units))
			fmt.Printf("\n%s:\n", playerInfo)

			tableFormatter := NewTableFormatter()
			// Inline unit table
			columns := []ColumnDef{
				{"Unit", "int", "right"},
				{"Type", "string", "left"},
				{"Level", "int", "right"},
				{"Health", "string", "right"},
				{"Position", "string", "right"},
			}
			var rows [][]interface{}
			for _, unitInfo := range units {
				unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
				healthStr := fmt.Sprintf("%d/%d", unitInfo.Unit.CurrentHealth, unitInfo.Unit.MaxHealth)
				positionStr := fmt.Sprintf("(%d,%d)", unitInfo.Row, unitInfo.Col)
				rows = append(rows, []interface{}{unitInfo.Index, unitTypeName, int(unitInfo.Unit.Level), healthStr, positionStr})
			}
			tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
		}
	}

	// Add unit type analysis
	fmt.Println("\nUnit Type Analysis:")
	fmt.Println("------------------")
	DisplayTypeAnalysis(validUnits,
		func(u UnitDisplayInfo) uint8 { return u.Unit.UnitType },
		func(u UnitDisplayInfo) int { return 1 },
		"units")
}

// ListGenerals displays all units that have generals assigned, grouped by player
func ListGenerals(saveOutput *WC4SaveOutput) {
	DisplayHeader("Generals")

	// First, collect all generals with their information
	var generals []UnitDisplayInfo

	for i, unit := range saveOutput.Units {
		if unit.GeneralId > 0 {
			row, col, _ := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			owner := saveOutput.UnitOwnerData[row][col]

			generals = append(generals, UnitDisplayInfo{
				Index: i,
				Unit:  unit,
				Row:   row,
				Col:   col,
				Owner: owner,
			})
		}
	}

	if len(generals) == 0 {
		fmt.Println("No generals found.")
		return
	}

	// First, show linear list
	fmt.Println("Linear List:")
	for _, generalInfo := range generals {
		unitTypeName := GetUnitTypeName(generalInfo.Unit.UnitType)
		countryName := "Unowned"
		if int(generalInfo.Owner) < len(saveOutput.PlayerData) {
			countryName, _ = GetCountryInfoFromData(saveOutput.PlayerData[generalInfo.Owner])
		}
		generalName, hasGeneralName := GetGeneralName(generalInfo.Unit.GeneralId)

		generalInfoStr := fmt.Sprintf("GeneralId=%d", generalInfo.Unit.GeneralId)
		if hasGeneralName {
			generalInfoStr = fmt.Sprintf("GeneralId=%d (%s)", generalInfo.Unit.GeneralId, generalName)
		}

		fmt.Printf("General (unit %d, owner %d - %s): Type=%d (%s), %s, Level=%d, Health=%d/%d, Position=(%d,%d)\n",
			generalInfo.Index, generalInfo.Owner, countryName, generalInfo.Unit.UnitType, unitTypeName,
			generalInfoStr, generalInfo.Unit.Level, generalInfo.Unit.CurrentHealth,
			generalInfo.Unit.MaxHealth, generalInfo.Row, generalInfo.Col)
	}

	// Then, group generals by owner
	fmt.Println("\nGrouped by Owner:")
	generalsByOwner := GroupUnitsByOwner(generals)

	// Display generals grouped by owner using sorted player IDs
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if generalUnits, exists := generalsByOwner[byte(ownerID)]; exists {
			countryName, _ := GetCountryInfoFromData(saveOutput.PlayerData[ownerID])
			fmt.Printf("\nPlayer %d (%s - CountryId %d) has %d generals:\n",
				ownerID, countryName, saveOutput.PlayerData[ownerID].CountryId, len(generalUnits))

			tableFormatter := NewTableFormatter()
			// Inline general table
			columns := []ColumnDef{
				{"Unit", "int", "right"},
				{"Type", "string", "left"},
				{"General", "string", "left"},
				{"Level", "int", "right"},
				{"Health", "string", "right"},
				{"Position", "string", "right"},
			}
			var rows [][]interface{}
			for _, generalInfo := range generalUnits {
				unitTypeName := GetUnitTypeName(generalInfo.Unit.UnitType)
				generalName, hasGeneralName := GetGeneralName(generalInfo.Unit.GeneralId)
				healthStr := fmt.Sprintf("%d/%d", generalInfo.Unit.CurrentHealth, generalInfo.Unit.MaxHealth)
				positionStr := fmt.Sprintf("(%d,%d)", generalInfo.Row, generalInfo.Col)
				generalInfoStr := fmt.Sprintf("ID %d", generalInfo.Unit.GeneralId)
				if hasGeneralName {
					generalInfoStr = generalName
				}
				rows = append(rows, []interface{}{generalInfo.Index, unitTypeName, generalInfoStr, int(generalInfo.Unit.Level), healthStr, positionStr})
			}
			tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
		}
	}

	// Add general type analysis
	fmt.Println("\nGeneral Unit Type Analysis:")
	fmt.Println("---------------------------")
	DisplayTypeAnalysis(generals,
		func(g UnitDisplayInfo) uint8 { return g.Unit.UnitType },
		func(g UnitDisplayInfo) int { return 1 },
		"generals")

	fmt.Printf("\nTotal generals found: %d\n", len(generals))
}

// ListLandmines displays all landmines grouped by owner
func ListLandmines(saveOutput *WC4SaveOutput) {
	DisplayHeader(fmt.Sprintf("Landmines (Total: %d)", len(saveOutput.Landmines)))

	if len(saveOutput.Landmines) == 0 {
		fmt.Println("No landmines found.")
		return
	}

	// Process landmines and group by owner
	landminesByOwner := make(map[int][]LandmineData)

	for _, landmine := range saveOutput.Landmines {
		if landmine.Owner == 0 {
			continue // zero is not a valid player-owned mine
		}
		playerID := int(landmine.Owner) - 1
		landminesByOwner[playerID] = append(landminesByOwner[playerID], landmine)
	}

	// Display landmines grouped by owner
	tableFormatter := NewTableFormatter()
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if landmines, exists := landminesByOwner[ownerID]; exists {
			player := saveOutput.PlayerData[ownerID]
			countryName, _ := GetCountryInfoFromData(player)
			fmt.Printf("\nPlayer %d (%s - CountryId %d) owns %d landmines:\n",
				ownerID, countryName, player.CountryId, len(landmines))

			// Inline landmine table
			columns := []ColumnDef{
				{"ID", "int", "right"},
				{"Position", "string", "center"},
				{"Health", "string", "right"},
			}
			var rows [][]interface{}
			for i, landmine := range landmines {
				row, col, _ := ConvertCoordinates(int(landmine.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
				positionStr := fmt.Sprintf("(%d,%d)", row, col)
				healthStr := fmt.Sprintf("%d", landmine.Health)
				rows = append(rows, []interface{}{i, positionStr, healthStr})
			}
			tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
		}
	}

	// Add landmine summary
	fmt.Println("\nLandmine Summary:")
	fmt.Println("----------------")
	fmt.Printf("Total landmines: %d\n", len(saveOutput.Landmines))
}

// ListTeams displays team analysis with statistics
func ListTeams(saveOutput *WC4SaveOutput) {
	DisplayHeader("Team Analysis")

	// Count units by player for team statistics
	unitCounts := CountPlayerTiles(saveOutput.UnitOwnerData)

	// Count cities by player
	cityCounts := make(map[byte]int)
	for _, city := range saveOutput.Cities {
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if valid && row < len(saveOutput.UnitOwnerData) && col < len(saveOutput.UnitOwnerData[row]) {
			owner := saveOutput.UnitOwnerData[row][col]
			if int(owner) < len(saveOutput.PlayerData) {
				cityCounts[owner]++
			}
		}
	}

	// Group players by team and calculate team statistics
	teamData := GroupPlayersByTeamWithStats(saveOutput.PlayerData, unitCounts, cityCounts)

	// Convert to slice and calculate percentages
	var teamAnalysisData []TeamAnalysisData
	totalCities := 0
	totalUnits := 0
	totalTiles := 0

	// Calculate totals
	for _, team := range teamData {
		totalCities += team.CityCount
		totalUnits += team.UnitCount
	}

	// Count total tiles
	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] != TileUnowned {
				totalTiles++
			}
		}
	}

	for _, team := range teamData {
		landPercentage := float64(team.UnitCount) / float64(totalTiles)
		worldPercentage := float64(team.UnitCount) / float64(totalTiles)
		team.LandPercent = landPercentage
		team.WorldPercent = worldPercentage
		teamAnalysisData = append(teamAnalysisData, *team)
	}

	// Sort teams by unit count (power) descending
	sort.Slice(teamAnalysisData, func(i, j int) bool {
		return teamAnalysisData[i].UnitCount > teamAnalysisData[j].UnitCount
	})

	// Print team analysis table
	tableFormatter := NewTableFormatter()
	// Inline team analysis table
	columns := []ColumnDef{
		{"Team", "int", "right"},
		{"Players", "int", "right"},
		{"Cities", "int", "right"},
		{"Units", "int", "right"},
		{"Land %", "percent", "right"},
		{"World %", "percent", "right"},
	}
	var rows [][]interface{}
	for _, team := range teamAnalysisData {
		rows = append(rows, []interface{}{
			team.TeamID,
			team.PlayerCount,
			team.CityCount,
			team.UnitCount,
			team.LandPercent,
			team.WorldPercent,
		})
	}
	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})

	// Print team member details
	fmt.Println("\nTeam Member Details:")
	fmt.Println("-------------------")
	for _, team := range teamAnalysisData {
		fmt.Printf("Team %d (%d players, %d cities, %d units, %.1f%% land, %.1f%% world):\n",
			team.TeamID, team.PlayerCount, team.CityCount, team.UnitCount, team.LandPercent*100, team.WorldPercent*100)

		// Print players in groups of 5 for better readability
		for i := 0; i < len(team.Players); i += 5 {
			end := i + 5
			if end > len(team.Players) {
				end = len(team.Players)
			}
			playerGroup := team.Players[i:end]
			fmt.Printf("  %s\n", strings.Join(playerGroup, ", "))
		}
		fmt.Println()
	}
}

// ListUnitsByMap displays units by map analysis
func ListUnitsByMap(saveOutput *WC4SaveOutput) {
	fmt.Println("Units by Owner Analysis:")
	fmt.Println("----------------------------")

	// Analyze UnitOwnerData map directly
	ownerCounts := make(map[byte]int)
	totalUnits := 0

	// Count tiles by owner
	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			owner := saveOutput.UnitOwnerData[i][j]

			// Check if owner is valid
			if int(owner) >= len(saveOutput.PlayerData) {
				continue
			}

			ownerCounts[owner]++
			totalUnits++
		}
	}

	// Show statistics
	fmt.Printf("Total units: %d\n", totalUnits)
	fmt.Println()

	// Display units by owner using sorted player IDs
	fmt.Println("Units by Owner:")

	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if count, exists := ownerCounts[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			percentage := float64(count) / float64(totalUnits) * 100

			countryName, _ := GetCountryInfoFromData(player)
			fmt.Printf("\nPlayer %d (%s - CountryId %d, TeamId %d) owns %d units (%.1f%% of total)\n",
				ownerID, countryName, player.CountryId, player.TeamId, count, percentage)

			// Aggregate coordinates for this player
			coordinates := make([]string, 0)
			missingFromUnitList := 0

			for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
				for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
					if saveOutput.UnitOwnerData[i][j] == byte(ownerID) {
						// Derive coordinate code from (i, j)
						coordinateCode := ConvertToCoordinateCode(i, j, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))

						// Check if any unit has this coordinate code
						hasUnit := false
						for _, unit := range saveOutput.Units {
							if int(unit.CoordinateCode) == coordinateCode {
								hasUnit = true
								break
							}
						}

						if hasUnit {
							coordinates = append(coordinates, fmt.Sprintf("(%d,%d)", i, j))
						} else {
							coordinates = append(coordinates, fmt.Sprintf("[%d,%d]", i, j))
							missingFromUnitList++
						}
					}
				}
			}

			// Print coordinates in grid format
			DisplayCoordinatesGrid(coordinates)

			// Show summary of missing units
			if missingFromUnitList > 0 {
				fmt.Printf("  [Brackets] indicate %d tiles missing from unit list\n", missingFromUnitList)
			}
		}
	}
}

// AnalyzeCityTiles analyzes the city tiles array and returns statistics
func AnalyzeCityTiles(saveOutput *WC4SaveOutput) TileAnalysisResult {
	coordinateCodeCounts := make(map[uint16]int)
	totalTiles := 0

	for i := 0; i < len(saveOutput.CityTiles); i++ {
		for j := 0; j < len(saveOutput.CityTiles[i]); j++ {
			totalTiles++
			coordinateCode := saveOutput.CityTiles[i][j]
			coordinateCodeCounts[coordinateCode]++
		}
	}

	oceanCount := coordinateCodeCounts[65535]
	landTiles := totalTiles - oceanCount

	return TileAnalysisResult{
		CoordinateCodeCounts: coordinateCodeCounts,
		TotalTiles:           totalTiles,
		OceanCount:           oceanCount,
		LandTiles:            landTiles,
	}
}

func DisplayTileAnalysis(result TileAnalysisResult, saveOutput *WC4SaveOutput) {
	fmt.Println("\nCity Tiles Array Analysis:")
	fmt.Println("--------------------------")
	fmt.Printf("Total tiles analyzed: %d\n", result.TotalTiles)
	fmt.Printf("Unique coordinate codes found: %d\n", len(result.CoordinateCodeCounts))

	fmt.Printf("\nSummary: %d land tiles (%.1f%%), %d ocean tiles (%.1f%%)\n",
		result.LandTiles, float64(result.LandTiles)/float64(result.TotalTiles)*100,
		result.OceanCount, float64(result.OceanCount)/float64(result.TotalTiles)*100)
}

// ProcessCitiesForTiles processes cities and groups them by owner for tile analysis
func ProcessCitiesForTiles(saveOutput *WC4SaveOutput) map[byte][]CityDisplayInfo {
	citiesByOwner := make(map[byte][]CityDisplayInfo)

	for i, city := range saveOutput.Cities {
		row, col, _ := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		owner := saveOutput.UnitOwnerData[row][col]

		cityDisplayInfo := CityDisplayInfo{
			Index: i,
			City:  city,
			Row:   row,
			Col:   col,
		}
		citiesByOwner[owner] = append(citiesByOwner[owner], cityDisplayInfo)
	}

	return citiesByOwner
}

// DisplayCitiesByOwner shows cities grouped by owner with tile counts
func DisplayCitiesByOwner(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	fmt.Println("\nCity Tiles by Owner:")

	tableFormatter := NewTableFormatter()
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)

	// Print player summary table (inline)
	{
		columns := []ColumnDef{
			{"Player", "int", "right"},
			{"Country", "string", "left"},
			{"Cities", "int", "right"},
			{"Total Tiles", "int", "right"},
			{"Land %", "percent", "right"},
			{"World %", "percent", "right"},
		}
		var rows [][]interface{}
		for _, ownerID := range sortedPlayerIDs {
			if cities, exists := citiesByOwner[byte(ownerID)]; exists {
				player := saveOutput.PlayerData[ownerID]
				totalTilesForPlayer := 0
				for _, cityInfo := range cities {
					totalTilesForPlayer += result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				}
				landPercentage := float64(totalTilesForPlayer) / float64(result.LandTiles)
				worldPercentage := float64(totalTilesForPlayer) / float64(result.TotalTiles)
				countryName, _ := GetCountryInfoFromData(player)
				rows = append(rows, []interface{}{
					ownerID,
					countryName,
					len(cities),
					totalTilesForPlayer,
					landPercentage,
					worldPercentage,
				})
			}
		}
		tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
	}

	// Print detailed city tables for each player
	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			countryName, _ := GetCountryInfoFromData(player)

			// Calculate total tiles for this player
			totalTilesForPlayer := 0
			for _, cityInfo := range cities {
				tileCount := result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				totalTilesForPlayer += tileCount
			}

			landPercentage := float64(totalTilesForPlayer) / float64(result.LandTiles)
			worldPercentage := float64(totalTilesForPlayer) / float64(result.TotalTiles)

			// Show detailed player summary with better formatting
			fmt.Printf("\nPlayer %d: %s\n", ownerID, countryName)
			fmt.Printf("  Country ID: %d, Team ID: %d\n", player.CountryId, player.TeamId)
			fmt.Printf("  City Tiles: %d cities, %d total tiles\n", len(cities), totalTilesForPlayer)
			fmt.Printf("  Territory: %.1f%% of land, %.1f%% of world\n", landPercentage*100, worldPercentage*100)

			// Inline city tile table
			columns := []ColumnDef{
				{"City", "int", "right"},
				{"Name", "string", "left"},
				{"Position", "string", "center"},
				{"Coord Code", "int", "right"},
				{"Tile Count", "int", "right"},
			}
			var rows [][]interface{}
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				rows = append(rows, []interface{}{
					cityInfo.Index,
					cityName,
					fmt.Sprintf("(%d,%d)", cityInfo.Row, cityInfo.Col),
					int(cityInfo.City.CoordinateCode),
					result.CoordinateCodeCounts[cityInfo.City.CoordinateCode],
				})
			}
			tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
		}
	}
}

func DisplayTerritoryControl(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	fmt.Println("\nPlayer Territory Control:")
	fmt.Println("------------------------")

	tableFormatter := NewTableFormatter()
	totalOwnedTiles := 0
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)

	// Build and print territory control table inline (no struct)
	columns := []ColumnDef{
		{"Player", "int", "right"},
		{"Country", "string", "left"},
		{"Country ID", "int", "right"},
		{"Tiles", "int", "right"},
		{"Land %", "percent", "right"},
		{"World %", "percent", "right"},
	}
	var rows [][]interface{}
	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			playerTiles := 0
			for _, cityInfo := range cities {
				playerTiles += result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
			}
			totalOwnedTiles += playerTiles
			player := saveOutput.PlayerData[ownerID]
			countryName, _ := GetCountryInfoFromData(player)

			rows = append(rows, []interface{}{
				ownerID,
				countryName,
				int(player.CountryId),
				playerTiles,
				float64(playerTiles) / float64(result.LandTiles),
				float64(playerTiles) / float64(result.TotalTiles),
			})
		}
	}
	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})

	// Print team territory control summary
	fmt.Println("\nTeam Territory Control:")
	fmt.Println("----------------------")
	DisplayTeamTerritoryControl(saveOutput, citiesByOwner, result)

	// Print total summary
	landPercentage := float64(totalOwnedTiles) / float64(result.LandTiles)
	worldPercentage := float64(totalOwnedTiles) / float64(result.TotalTiles)
	fmt.Printf("\nTotal controlled: %d tiles (%.1f%% of land, %.1f%% of world)\n",
		totalOwnedTiles, landPercentage*100, worldPercentage*100)
}

// DisplayTeamTerritoryControl shows team territory control summary
func DisplayTeamTerritoryControl(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	tableFormatter := NewTableFormatter()

	// Group players by team and calculate team totals
	teamData := GroupPlayersByTeamWithTerritory(saveOutput.PlayerData, citiesByOwner, result)

	// Build rows inline and print (no table-specific struct)
	columns := []ColumnDef{
		{"Team", "int", "right"},
		{"Players", "int", "right"},
		{"Tiles", "int", "right"},
		{"Land %", "percent", "right"},
		{"World %", "percent", "right"},
	}

	// Flatten map to slice of keys for deterministic order by tiles desc
	type teamRow struct {
		id   int
		data *TeamTerritoryData
	}
	var rowsData []teamRow
	for _, t := range teamData {
		rowsData = append(rowsData, teamRow{id: t.TeamID, data: t})
	}
	sort.Slice(rowsData, func(i, j int) bool {
		return rowsData[i].data.TileCount > rowsData[j].data.TileCount
	})

	var rows [][]interface{}
	for _, tr := range rowsData {
		landPct := float64(tr.data.TileCount) / float64(result.LandTiles)
		worldPct := float64(tr.data.TileCount) / float64(result.TotalTiles)
		rows = append(rows, []interface{}{
			tr.id,
			tr.data.PlayerCount,
			tr.data.TileCount,
			landPct,
			worldPct,
		})
	}
	tableFormatter.PrintTable(TableData{Columns: columns, Rows: rows})
}

// ListTiles displays city tile ownership analysis and territory control
func ListTiles(saveOutput *WC4SaveOutput) {
	DisplayHeader("City Tile Ownership Analysis")

	// Analyze city tiles
	result := AnalyzeCityTiles(saveOutput)
	DisplayTileAnalysis(result, saveOutput)

	// Process cities and group by owner
	citiesByOwner := ProcessCitiesForTiles(saveOutput)

	fmt.Printf("Total cities: %d\n", len(saveOutput.Cities))

	// Display cities grouped by owner
	DisplayCitiesByOwner(saveOutput, citiesByOwner, result)

	// Show territory control summary
	DisplayTerritoryControl(saveOutput, citiesByOwner, result)
}
