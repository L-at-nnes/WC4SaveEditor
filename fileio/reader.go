package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Debug flag to control verbose output during deserialization
const isDebug = false

// debugPrint prints the message only if debug mode is enabled
func debugPrint(format string, args ...interface{}) {
	if isDebug {
		fmt.Printf(format, args...)
	}
}

// debugPrintSection prints a section header only if debug mode is enabled
func debugPrintSection(sectionNum int, title string) {
	if isDebug {
		fmt.Printf("\n[%d] %s\n", sectionNum, title)
		fmt.Println(strings.Repeat("-", len(title)+4))
	}
}

var (
	fileOffsetMap = make(map[string]int)
)

type SaveHeader struct {
	Magic              [4]byte
	UnknownInt1        uint32
	MapId              uint32
	GameMode           uint32 // 1 - campaign, 2 - conquest, 6 - frontier
	UnknownInt2        uint32
	UnknownInt3        uint32
	Camera             [3]float32
	UnknownInt4        uint32
	TurnNumber         uint32
	UnknownArr2        [12]byte
	SaveTimestamp      [5]uint32
	UnknownArr3        [16]byte
	UnknownInt7        uint32 // only seems to be set to non-zero value in frontier mode last mission
	UnknownInt8        uint32 // only seems to be set to non-zero value in frontier mode last mission
	MapWidth           uint32
	MapHeight          uint32
	CountryCount       uint32
	CityCount          uint32
	UnitCount          uint32
	UnknownCount1      uint32
	UnknownCount2      uint32
	UnknownArr4        [8]byte
	TurnCount1         uint32
	TurnCount2         uint32
	UnknownCount3      uint32
	UnknownCount4      uint32
	UnknownCount5      uint32
	UnknownCount6      uint32
	ImportantCityCount uint32
	UnknownArr5        [4]byte
	UnknownInt9        uint32
	UnknownInt10       uint32
	UnknownArr6        [12]byte
	LandmineCount      uint32
	UnknownArr7        [16]byte
	UnknownCount9      uint32
}

type CountryData struct {
	TurnOrder    uint32
	CountryId    uint32
	Currency     [3]uint32
	BotFlag      uint32
	TeamId       uint32
	UnknownArr2  [4]byte
	UnknownColor [2][4]byte
	PrimaryColor [4]byte
	UnknownArr4  [16]byte
	UnknownArr5  [460]byte
}

type CityData struct {
	CoordinateCode    uint16
	CityId            uint16
	BuildingType      uint8
	Apperance         uint8
	UnknownByte1      uint8
	Wonders           uint8
	UnknownArr2       [6]byte
	UnknownArr3       [8]byte
	AntiAirWeaponType uint8
	AntiAirRange      uint8
	TechLevels        [6]byte
	UnknownArr4       [2]byte
}

type UnitData struct {
	CoordinateCode      uint16
	UnitType            uint8
	Level               uint8
	Personnel           uint8
	Direction           uint8 // 0 - left, 1 - right
	Movement            uint16
	Experience          uint16
	UnknownHealth       uint16
	CurrentHealth       uint16
	MaxHealth           uint16
	GeneralId           uint16
	GeneralMilitaryRank uint8
	GeneralTitle        uint8
	GeneralBadges       [3]byte
	GeneralSkillLevels  [5]byte
	UnknownArr5         [12]byte
	MoraleValue         int8
	MoraleTurnsLeft     uint16
	UnknownArr6         [21]byte
}

type LandmineData struct {
	CoordinateCode uint16
	Owner          uint16 // One-based player ID: player 0 is stored as 1
	UnknownArr1    [2]byte
	Health         uint16
	UnknownArr2    [4]byte
}

type WC4SaveOutput struct {
	SaveHeader    SaveHeader
	PlayerData    []CountryData
	CityTiles     [][]uint16
	UnitOwnerData [][]byte
	Cities        []CityData
	Units         []UnitData
	Landmines     []LandmineData
}

func DeserializeMapHeaderFromBytes(streamReader *io.SectionReader) SaveHeader {
	mapHeaderInput := SaveHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &mapHeaderInput); err != nil {
		log.Fatal("Failed to load MapHeaderInput: ", err)
	}
	debugPrint("Map Header Input: %+v\n", mapHeaderInput)
	return mapHeaderInput
}

func DeserializeCountryDataFromBytes(streamReader *io.SectionReader, count int) []CountryData {
	allPlayerData := make([]CountryData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMap(fileOffsetMap, streamReader, BuildPlayerStartKey(i))

		countryData := CountryData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &countryData); err != nil {
			log.Fatal("Failed to load country data: ", err)
		}
		allPlayerData[i] = countryData
		debugPrint("Player %v data: %+v\n", i, countryData)
	}
	return allPlayerData
}

func DeserializeCityTileOwnershipFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) [][]uint16 {
	allCityTiles := make([][]uint16, 0)
	for i := 0; i < mapHeight; i++ {
		cityRow := make([]uint16, 0)
		for j := 0; j < mapWidth; j++ {
			cityCoordinates := uint16(0)
			if err := binary.Read(streamReader, binary.LittleEndian, &cityCoordinates); err != nil {
				log.Fatal("Failed to load city tile ownership: ", err)
			}
			cityRow = append(cityRow, cityCoordinates)
		}
		allCityTiles = append(allCityTiles, cityRow)
		debugPrint("City Tile Owner Row %d: %v\n", i, cityRow)
	}
	return allCityTiles
}

func DeserializeUnknownCampaignBlockFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) {
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			unknownBlock := make([]byte, 16)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load city tile ownership: ", err)
			}
			debugPrint("Unknown block: %v\n", unknownBlock)
		}
	}
}

func DeserializeUnitOwnerDataFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) [][]byte {
	unitOwnerData := make([][]byte, 0)

	for i := 0; i < mapHeight; i++ {
		unitOwnerRow := make([]byte, mapWidth*1)

		for j := 0; j < mapWidth; j++ {
			updateFileOffsetMap(fileOffsetMap, streamReader, buildSingleUnitOwnerStartKey(j, i))

			if err := binary.Read(streamReader, binary.LittleEndian, &unitOwnerRow[j]); err != nil {
				log.Fatal("Failed to load unit owner: ", err)
			}
		}

		unitOwnerData = append(unitOwnerData, unitOwnerRow)
		debugPrint("Unit Owner Row %d: %v\n", i, unitOwnerRow)
	}

	return unitOwnerData
}

func DeserializeCityDataFromBytes(streamReader *io.SectionReader, count int) []CityData {
	allCities := make([]CityData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMap(fileOffsetMap, streamReader, BuildCityStartKey(i))

		cityData := CityData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &cityData); err != nil {
			log.Fatal("Failed to load cityData: ", err)
		}

		allCities[i] = cityData
		debugPrint("City data: %+v\n", cityData)

		if i > 0 && cityData.CoordinateCode == 0 {
			log.Fatal("Invalid city data")
		}
	}
	return allCities
}

func DeserializeUnitDataFromBytes(streamReader *io.SectionReader, count int) []UnitData {
	allUnits := make([]UnitData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMapForField(fileOffsetMap, streamReader, BuildUnitStartKey(i), 0)

		unitData := UnitData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &unitData); err != nil {
			log.Fatal("Failed to load unit data: ", err)
		}
		allUnits[i] = unitData
		debugPrint("Unit data: %+v\n", unitData)
	}
	return allUnits
}

func DeserializeLandmineDataFromBytes(streamReader *io.SectionReader, count int) []LandmineData {
	var landmines []LandmineData
	for i := 0; i < count; i++ {
		updateFileOffsetMap(fileOffsetMap, streamReader, BuildLandmineStartKey(i))

		landmineData := LandmineData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &landmineData); err != nil {
			log.Fatal("Failed to load landmine data: ", err)
		}
		debugPrint("Landmine: %+v\n", landmineData)
		landmines = append(landmines, landmineData)
	}
	return landmines
}

func DeserializeUnknownData2FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 16)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Unknown block 2: %v\n", unknownBlock)
	}
}

func DeserializeUnknownData3FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 44)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Unknown block 3: %v\n", unknownBlock)
	}
}

func DeserializeUnknownData4FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 80)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Unknown block 4: %v\n", unknownBlock)
	}
}

func DeserializeUnknownData5FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 8)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Unknown block 5: %v\n", unknownBlock)
	}
}

func DeserializeImportantCityDataFromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 4)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Important city: %v\n", unknownBlock)
	}
}

func DeserializeUnknownData7FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 16)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		debugPrint("Unknown block 7: %v\n", unknownBlock)
	}
}

func ReadSaveFile(inputFilename string) (*WC4SaveOutput, error) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	fi, err := inputFile.Stat()
	if err != nil {
		return nil, err
	}
	fileLength := fi.Size()
	streamReader := io.NewSectionReader(inputFile, int64(0), fileLength)

	debugPrint("\n%s\n", strings.Repeat("=", 60))
	debugPrint("READING SAVE FILE SECTIONS\n")
	debugPrint("%s\n", strings.Repeat("=", 60))

	debugPrintSection(1, "SAVE HEADER")
	saveHeader := DeserializeMapHeaderFromBytes(streamReader)

	debugPrintSection(2, "PLAYER DATA")
	allPlayerData := DeserializeCountryDataFromBytes(streamReader, int(saveHeader.CountryCount))

	isConquest := (int(saveHeader.GameMode) == GameModeConquest)
	if !isConquest {
		if saveHeader.UnknownInt7 == 0 {
			debugPrintSection(3, "CAMPAIGN BLOCK")
			DeserializeUnknownCampaignBlockFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))
		}
	}

	debugPrintSection(4, "CITY TILES")
	allCityTiles := DeserializeCityTileOwnershipFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))

	if !isConquest {
		// required for some maps because the data is shifted
		if int(saveHeader.MapWidth)*int(saveHeader.MapHeight) != int(saveHeader.UnknownInt10) {
			debugPrintSection(5, "UNKNOWN BLOCK 1")
			unknownBlock := make([]byte, 8)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load unknownBlock: ", err)
			}
		}
	}

	debugPrintSection(6, "UNIT OWNER DATA")
	updateFileOffsetMap(fileOffsetMap, streamReader, buildUnitOwnerStartKey())
	unitOwnerData := DeserializeUnitOwnerDataFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))
	updateFileOffsetMap(fileOffsetMap, streamReader, buildUnitOwnerEndKey())

	if !isConquest {
		if int(saveHeader.MapWidth)*int(saveHeader.MapHeight) != int(saveHeader.UnknownInt10) {
			debugPrintSection(7, "UNKNOWN BLOCK 2")
			unknownBlock := make([]byte, 4)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load unknownBlock: ", err)
			}
		}
	}

	debugPrintSection(8, "CITIES")
	allCities := DeserializeCityDataFromBytes(streamReader, int(saveHeader.CityCount))

	debugPrintSection(9, "UNITS")
	allUnits := DeserializeUnitDataFromBytes(streamReader, int(saveHeader.UnitCount))

	debugPrintSection(10, "LANDMINES")
	allLandmines := DeserializeLandmineDataFromBytes(streamReader, int(saveHeader.LandmineCount))

	debugPrintSection(11, "UNKNOWN DATA SECTIONS")
	DeserializeUnknownData2FromBytes(streamReader, int(saveHeader.UnknownCount1))
	DeserializeUnknownData3FromBytes(streamReader, int(saveHeader.UnknownCount2))
	DeserializeUnknownData4FromBytes(streamReader, int(saveHeader.UnknownCount3))
	DeserializeUnknownData5FromBytes(streamReader, int(saveHeader.UnknownCount5))
	DeserializeUnknownData5FromBytes(streamReader, int(saveHeader.UnknownCount6))
	DeserializeImportantCityDataFromBytes(streamReader, int(saveHeader.ImportantCityCount))
	DeserializeUnknownData7FromBytes(streamReader, int(saveHeader.UnknownCount9))

	saveOutput := &WC4SaveOutput{
		SaveHeader:    saveHeader,
		PlayerData:    allPlayerData,
		CityTiles:     allCityTiles,
		UnitOwnerData: unitOwnerData,
		Cities:        allCities,
		Units:         allUnits,
		Landmines:     allLandmines,
	}

	// Print file summary with byte ranges
	printFileSummary(saveHeader, fileLength)

	return saveOutput, nil
}

func printFileSummary(header SaveHeader, fileLength int64) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FILE SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total file size: %d bytes (%.2f KB)\n", fileLength, float64(fileLength)/1024)
	fmt.Printf("Game mode: %d (%s)\n", header.GameMode, GetGameModeName(header.GameMode))
	fmt.Printf("Map dimensions: %dx%d\n", header.MapWidth, header.MapHeight)
	fmt.Printf("Turn number: %d\n", header.TurnNumber)
	fmt.Println()

	fmt.Println(strings.Repeat("=", 60))
}
