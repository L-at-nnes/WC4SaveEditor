package fileio

import (
	"fmt"
	"io"
	"log"
)

func buildUnitOwnerStartKey() string {
	return "UnitOwnerStart"
}

func buildUnitOwnerEndKey() string {
	return "UnitOwnerEnd"
}

func buildSingleUnitOwnerStartKey(x int, y int) string {
	return fmt.Sprintf("UnitOwner%v,%v", x, y)
}

func BuildPlayerStartKey(index int) string {
	return fmt.Sprintf("PlayerStart%v", index)
}

func BuildCityStartKey(index int) string {
	return fmt.Sprintf("CityStart%v", index)
}

func BuildUnitStartKey(index int) string {
	return fmt.Sprintf("UnitStart%v", index)
}

func BuildLandmineStartKey(index int) string {
	return fmt.Sprintf("LandmineStart%v", index)
}

func GetFileOffsetMap() map[string]int {
	return fileOffsetMap
}

// GetOffset safely retrieves an offset by key with validation
func GetOffset(offsetKey string) int {
	offset, ok := fileOffsetMap[offsetKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find offset for key '%s'. Command not run.", offsetKey))
	}
	return offset
}

func updateFileOffsetMap(fileOffsetMap map[string]int, streamReader *io.SectionReader, unitLocationKey string) {
	fileOffset, err := streamReader.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal(err)
	}
	fileOffsetMap[unitLocationKey] = int(fileOffset)
}

func updateFileOffsetMapForField(fileOffsetMap map[string]int, streamReader *io.SectionReader, unitLocationKey string, offset int) {
	fileOffset, err := streamReader.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal(err)
	}
	fileOffsetMap[unitLocationKey] = int(fileOffset) + offset
}
