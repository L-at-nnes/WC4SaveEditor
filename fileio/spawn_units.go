package fileio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"time"
)

// funRand drives "fun mode" unit type selection. Seeded once at startup so
// repeated fun spawns don't always pick the same sequence of types.
var funRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Maximum hex-tile distance from a province's port(s) to search for a free sea
// tile when spawning naval units. Keeps naval spawns "next to the coast"
// rather than picking an arbitrary far-away patch of open ocean.
const maxCoastSearchRadius = 6

// TileCodeOcean is the CityTiles value marking a tile as sea with no owning
// province (see MapVisualizer.getCityTileColor).
const TileCodeOcean = 65535

// hexPos identifies a map tile by row/col (matches UnitOwnerData/CityTiles indexing).
type hexPos struct {
	row, col int
}

// SpawnableUnitTypes returns every unit type this feature can create, in
// display order. Only "special" units are offered: basic units (Light
// Infantry, Medium Tank, Destroyer, Bunker, ...) are mass-produced and scale
// with a player's city tech instead of a chosen per-unit level, so they have
// no Level for this feature to set (confirmed across the whole save: every
// basic unit type tops out well under 9, while every special type reaches
// 9). UnitTypeCity is excluded too: it represents a city's own
// fortification, not a troop that can be spawned.
func SpawnableUnitTypes() []uint8 {
	return []uint8{
		UnitTypeBrandenburgers, UnitTypeHawkeyeForce, UnitTypeCombatMedic, UnitTypeM7Priest, UnitTypeHeavyGustav,
		UnitTypeBM21, UnitTypeT44, UnitTypeKingTiger, UnitTypeM26Pershing, UnitTypeTypeVIISubmarine, UnitTypeHMSPrinceOfWales,
		UnitTypeB4Howitzer, UnitTypeIS3HeavyTank, UnitTypeStukaZuFuss, UnitTypeRichelieu, UnitTypeEnterprise, UnitTypeRPGRocketeer,
		UnitTypeA41Centurion, UnitTypeAuF1, UnitTypeT72, UnitTypePhantomForce, UnitTypeM1A1Abrams, UnitTypeAH64Apache,
		UnitTypeM142Himars, UnitTypeDivineWrathMBT,
	}
}

// isNavalUnitType reports whether a unit type is a ship (spawns at sea, next
// to a port) rather than a land unit (spawns on the province's own tiles).
func isNavalUnitType(unitType uint8) bool {
	switch unitType {
	case UnitTypeTypeVIISubmarine, UnitTypeHMSPrinceOfWales, UnitTypeRichelieu, UnitTypeEnterprise:
		return true
	default:
		return false
	}
}

// landUnitTypes returns the spawnable unit types that place on land.
func landUnitTypes() []uint8 {
	var types []uint8
	for _, t := range SpawnableUnitTypes() {
		if !isNavalUnitType(t) {
			types = append(types, t)
		}
	}
	return types
}

// navalUnitTypes returns the spawnable unit types that place at sea.
func navalUnitTypes() []uint8 {
	var types []uint8
	for _, t := range SpawnableUnitTypes() {
		if isNavalUnitType(t) {
			types = append(types, t)
		}
	}
	return types
}

// usableTemplates resolves each of unitTypes to a stat template via
// findStatTemplate, dropping any type with no existing unit in the save to
// clone stats from. The returned templates are what fun mode picks from at
// random, so each pick already carries a valid UnitType and baseline stats.
func usableTemplates(units []UnitData, unitTypes []uint8, level int) []UnitData {
	var templates []UnitData
	for _, unitType := range unitTypes {
		if template, err := findStatTemplate(units, unitType, level); err == nil {
			templates = append(templates, template)
		}
	}
	return templates
}

// takeRandomShare shuffles positions and keeps a random percent% share of
// them (rounded to the nearest tile), so fun mode's partial fills are spread
// across the province instead of always favoring the tiles nearest the
// anchor.
func takeRandomShare(positions []hexPos, percent int) []hexPos {
	if percent >= 100 || len(positions) == 0 {
		return positions
	}
	count := int(math.Round(float64(len(positions)) * float64(percent) / 100.0))
	funRand.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })
	return positions[:count]
}

// FunSpawnResult reports how many units SpawnFunUnits placed.
type FunSpawnResult struct {
	LandSpawned  int
	NavalSpawned int
}

// SpawnFunUnits is "fun mode" for SpawnUnits: instead of one chosen unit
// type/level/count, it fills fillPercent% of the free tiles belonging to
// provinceCode with a random land unit at max level (9), and fillPercent% of
// the free sea tiles near the province's port(s) (see findFreeNavalTiles)
// with a random naval unit at max level - so ships only ever land on actual
// sea tiles next to a port, never on the province's own land tiles. Which
// tiles get filled is chosen at random rather than nearest-first, so a
// partial fill spreads across the province instead of clustering. fillPercent
// is clamped to [1, 100]. A unit type is only used if the save already has
// an existing unit of that type to clone stats from (same requirement as
// SpawnUnits). Missing pieces (no free land tiles, no port, no usable naval
// type, ...) just skip that half instead of failing the whole call; an error
// is only returned if nothing at all could be spawned.
func SpawnFunUnits(inputFilename string, saveOutput *WC4SaveOutput, provinceCode uint16, fillPercent int) (FunSpawnResult, error) {
	const funLevel = 9
	if fillPercent < 1 {
		fillPercent = 1
	}
	if fillPercent > 100 {
		fillPercent = 100
	}

	anchorRow, anchorCol, valid := ConvertCoordinates(int(provinceCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
	if !valid {
		return FunSpawnResult{}, fmt.Errorf("province coordinate code %d is not a valid map position", provinceCode)
	}
	provinceOwner := saveOutput.UnitOwnerData[anchorRow][anchorCol]

	landTemplates := usableTemplates(saveOutput.Units, landUnitTypes(), funLevel)
	navalTemplates := usableTemplates(saveOutput.Units, navalUnitTypes(), funLevel)

	currentTurn := byte(saveOutput.SaveHeader.TurnNumber)
	var newUnits []UnitData
	var claimTiles []hexPos
	var result FunSpawnResult

	if len(landTemplates) > 0 {
		if positions, err := findFreeLandTiles(saveOutput, hexPos{anchorRow, anchorCol}, provinceCode, math.MaxInt32); err == nil {
			positions = takeRandomShare(positions, fillPercent)
			for _, pos := range positions {
				template := landTemplates[funRand.Intn(len(landTemplates))]
				coordinateCode := ConvertToCoordinateCode(pos.row, pos.col, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
				newUnits = append(newUnits, buildNewUnit(template, funLevel, uint16(coordinateCode), currentTurn))
				claimTiles = append(claimTiles, pos)
			}
			result.LandSpawned = len(positions)
		}
	}

	if len(navalTemplates) > 0 {
		if positions, err := findFreeNavalTiles(saveOutput, provinceCode, math.MaxInt32); err == nil {
			positions = takeRandomShare(positions, fillPercent)
			for _, pos := range positions {
				template := navalTemplates[funRand.Intn(len(navalTemplates))]
				coordinateCode := ConvertToCoordinateCode(pos.row, pos.col, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
				newUnits = append(newUnits, buildNewUnit(template, funLevel, uint16(coordinateCode), currentTurn))
				claimTiles = append(claimTiles, pos)
			}
			result.NavalSpawned = len(positions)
		}
	}

	if len(newUnits) == 0 {
		return result, fmt.Errorf("no free tile or usable unit type found near this province; nothing was spawned")
	}

	// Same territory-claiming rule as SpawnUnits: occupying a free-but-still-
	// neutral tile claims it for the province's owner.
	for _, pos := range claimTiles {
		if saveOutput.UnitOwnerData[pos.row][pos.col] != provinceOwner {
			WriteUnitOwnerToFile(inputFilename, int(provinceOwner), pos.col, pos.row)
			saveOutput.UnitOwnerData[pos.row][pos.col] = provinceOwner
		}
	}

	if err := AppendNewUnits(inputFilename, saveOutput, newUnits); err != nil {
		return result, err
	}
	return result, nil
}

// SpawnUnits creates count new units of unitType at the given level, owned by
// whoever already owns provinceCode (the tile grid, not the unit itself,
// determines ownership - see ConvertCoordinates/UnitOwnerData). Land units are
// placed on the nearest free tile belonging to the province; naval units are
// placed on the nearest free sea tile next to one of the province's ports. If
// fewer free tiles than requested are found, or no stat reference for
// unitType exists anywhere in the save, nothing is written and an error is
// returned.
func SpawnUnits(inputFilename string, saveOutput *WC4SaveOutput, provinceCode uint16, unitType uint8, level int, count int) (int, error) {
	if count < 1 {
		return 0, fmt.Errorf("count must be at least 1")
	}
	if level < 1 || level > 9 {
		return 0, fmt.Errorf("level must be between 1 and 9")
	}

	template, err := findStatTemplate(saveOutput.Units, unitType, level)
	if err != nil {
		return 0, err
	}

	anchorRow, anchorCol, valid := ConvertCoordinates(int(provinceCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
	if !valid {
		return 0, fmt.Errorf("province coordinate code %d is not a valid map position", provinceCode)
	}
	provinceOwner := saveOutput.UnitOwnerData[anchorRow][anchorCol]

	var positions []hexPos
	if isNavalUnitType(unitType) {
		positions, err = findFreeNavalTiles(saveOutput, provinceCode, count)
	} else {
		positions, err = findFreeLandTiles(saveOutput, hexPos{anchorRow, anchorCol}, provinceCode, count)
	}
	if err != nil {
		return 0, err
	}
	if len(positions) < count {
		return 0, fmt.Errorf("only found %d free tile(s) near this province for %d requested unit(s); nothing was spawned", len(positions), count)
	}

	// Territory in this save is asserted by occupation: every owned tile we
	// saw was one with a unit standing on it. Spawning onto a free-but-still-
	// neutral tile inside the province's own footprint (or its coastline)
	// claims it for the province's owner, the same way ConvertTile does.
	for _, pos := range positions {
		if saveOutput.UnitOwnerData[pos.row][pos.col] != provinceOwner {
			WriteUnitOwnerToFile(inputFilename, int(provinceOwner), pos.col, pos.row)
			saveOutput.UnitOwnerData[pos.row][pos.col] = provinceOwner
		}
	}

	currentTurn := byte(saveOutput.SaveHeader.TurnNumber)
	newUnits := make([]UnitData, 0, count)
	for i := 0; i < count; i++ {
		coordinateCode := ConvertToCoordinateCode(positions[i].row, positions[i].col, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		newUnits = append(newUnits, buildNewUnit(template, level, uint16(coordinateCode), currentTurn))
	}

	if err := AppendNewUnits(inputFilename, saveOutput, newUnits); err != nil {
		return 0, err
	}
	return len(newUnits), nil
}

// findStatTemplate finds an existing unit of unitType to use as a baseline
// for stats (health, movement, personnel, ...) that this codebase has no
// formula for. An exact level match is used as-is; otherwise the closest
// level of the same type is scaled proportionally. Returns an error if no
// unit of that type exists anywhere in the save.
func findStatTemplate(units []UnitData, unitType uint8, level int) (UnitData, error) {
	var best UnitData
	found := false
	bestDiff := math.MaxInt32

	for _, u := range units {
		if u.UnitType != unitType {
			continue
		}
		diff := int(u.Level) - level
		if diff < 0 {
			diff = -diff
		}
		if !found || diff < bestDiff {
			best, bestDiff, found = u, diff, true
		}
		if diff == 0 {
			break
		}
	}

	if !found {
		return UnitData{}, fmt.Errorf("no existing %s unit found in this save; cannot determine stats for a new one", GetUnitTypeName(unitType))
	}
	return best, nil
}

// buildNewUnit derives a new unit from a same-type template: stats are kept
// (or scaled to the requested level), and general/experience/morale fields
// are cleared since this feature does not assign a general. The unit is
// stamped ready-to-act the same way ConvertProvince/ConvertPlayer stamp newly
// annexed units, so it isn't left frozen in an inactive AI state.
func buildNewUnit(template UnitData, level int, coordinateCode uint16, currentTurn byte) UnitData {
	newUnit := template
	newUnit.CoordinateCode = coordinateCode
	newUnit.Level = uint8(level)

	if int(template.Level) != level && template.Level > 0 {
		ratio := float64(level) / float64(template.Level)
		newUnit.MaxHealth = scaleUint16(template.MaxHealth, ratio)
		newUnit.Movement = scaleUint16(template.Movement, ratio)
	}
	newUnit.CurrentHealth = newUnit.MaxHealth

	newUnit.Experience = 0
	newUnit.GeneralId = 0
	newUnit.GeneralMilitaryRank = 0
	newUnit.GeneralTitle = 0
	newUnit.GeneralBadges = [3]byte{}
	newUnit.GeneralSkillLevels = [5]byte{}
	newUnit.MoraleValue = 0
	newUnit.MoraleTurnsLeft = 0

	newUnit.UnknownArr5[10] = currentTurn
	newUnit.UnknownArr6[0] = unitReadyActionState

	return newUnit
}

func scaleUint16(value uint16, ratio float64) uint16 {
	scaled := int(math.Round(float64(value) * ratio))
	if scaled < 1 {
		scaled = 1
	}
	if scaled > math.MaxUint16 {
		scaled = math.MaxUint16
	}
	return uint16(scaled)
}

// tileOccupancy indexes which tiles already have a unit on them.
type tileOccupancy map[hexPos]bool

// occupiedTiles returns every tile with a unit on it. If excludeCity is true,
// tiles occupied only by a UnitTypeCity entry are not counted as occupied -
// cities and a garrison unit are known to legitimately share a tile (e.g. a
// city's own fortification unit coexists with a defending general).
func occupiedTiles(saveOutput *WC4SaveOutput, excludeCity bool) tileOccupancy {
	occupied := make(tileOccupancy)
	for _, u := range saveOutput.Units {
		if excludeCity && u.UnitType == UnitTypeCity {
			continue
		}
		row, col, valid := ConvertCoordinates(int(u.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if valid {
			occupied[hexPos{row, col}] = true
		}
	}
	return occupied
}

// findFreeLandTiles returns up to count free tiles belonging to provinceCode,
// nearest first, using hex-tile distance from the province's own position.
func findFreeLandTiles(saveOutput *WC4SaveOutput, anchor hexPos, provinceCode uint16, count int) ([]hexPos, error) {
	occupied := occupiedTiles(saveOutput, true)

	var candidates []hexPos
	for row := range saveOutput.CityTiles {
		for col, code := range saveOutput.CityTiles[row] {
			if code != provinceCode || occupied[hexPos{row, col}] {
				continue
			}
			candidates = append(candidates, hexPos{row, col})
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no free tile available in this province; nothing was spawned")
	}

	sort.Slice(candidates, func(i, j int) bool {
		return hexDistance(anchor.col, anchor.row, candidates[i].col, candidates[i].row) <
			hexDistance(anchor.col, anchor.row, candidates[j].col, candidates[j].row)
	})
	if len(candidates) > count {
		candidates = candidates[:count]
	}
	return candidates, nil
}

// findFreeNavalTiles returns up to count free ocean tiles near one of
// provinceCode's ports, nearest first. Ocean tiles are identified the same
// way the map visualizer does: CityTiles == 65535 means no province owns it.
func findFreeNavalTiles(saveOutput *WC4SaveOutput, provinceCode uint16, count int) ([]hexPos, error) {
	var anchors []hexPos
	for _, city := range saveOutput.Cities {
		if !isPort(city.BuildingType) {
			continue
		}
		portProvince, attached := provinceCodeForPort(saveOutput, city)
		if !attached || portProvince != provinceCode {
			continue
		}
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if valid {
			anchors = append(anchors, hexPos{row, col})
		}
	}
	if len(anchors) == 0 {
		return nil, fmt.Errorf("this province has no port; naval units cannot spawn here")
	}

	occupied := occupiedTiles(saveOutput, false)

	type candidate struct {
		pos  hexPos
		dist int
	}
	var candidates []candidate
	for row := range saveOutput.CityTiles {
		for col, code := range saveOutput.CityTiles[row] {
			if code != TileCodeOcean || occupied[hexPos{row, col}] {
				continue
			}
			minDist := math.MaxInt32
			for _, anchor := range anchors {
				if d := hexDistance(anchor.col, anchor.row, col, row); d < minDist {
					minDist = d
				}
			}
			if minDist <= maxCoastSearchRadius {
				candidates = append(candidates, candidate{hexPos{row, col}, minDist})
			}
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no free water tile found next to this province's port(s); nothing was spawned")
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].dist < candidates[j].dist })
	if len(candidates) > count {
		candidates = candidates[:count]
	}
	result := make([]hexPos, len(candidates))
	for i, c := range candidates {
		result[i] = c.pos
	}
	return result, nil
}

// oddQToCube converts an odd-q offset hex coordinate (this map's layout - see
// MapVisualizer.getHexPosition, which shifts odd columns down) to cube
// coordinates, for exact hex-tile distance.
func oddQToCube(col, row int) (x, y, z int) {
	x = col
	z = row - (col-(col&1))/2
	y = -x - z
	return
}

func hexDistance(col1, row1, col2, row2 int) int {
	x1, y1, z1 := oddQToCube(col1, row1)
	x2, y2, z2 := oddQToCube(col2, row2)
	return (absInt(x1-x2) + absInt(y1-y2) + absInt(z1-z2)) / 2
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// AppendNewUnits inserts newUnits right after the save's existing unit
// records, shifting every following section (landmines, etc.) later in the
// file, and increments the header's UnitCount to match. saveOutput is
// updated in place to stay consistent with the file on disk.
func AppendNewUnits(inputFilename string, saveOutput *WC4SaveOutput, newUnits []UnitData) error {
	if len(newUnits) == 0 {
		return nil
	}

	data, err := os.ReadFile(inputFilename)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	unitSize := binary.Size(UnitData{})
	originalUnitCount := len(saveOutput.Units)
	insertionOffset := GetOffset(BuildUnitStartKey(originalUnitCount-1)) + unitSize

	unitBytes := new(bytes.Buffer)
	for _, u := range newUnits {
		if err := binary.Write(unitBytes, binary.LittleEndian, u); err != nil {
			return fmt.Errorf("failed to encode new unit: %w", err)
		}
	}

	newData := make([]byte, 0, len(data)+unitBytes.Len())
	newData = append(newData, data[:insertionOffset]...)
	newData = append(newData, unitBytes.Bytes()...)
	newData = append(newData, data[insertionOffset:]...)

	unitCountOffset, err := binaryFieldOffset(SaveHeader{}, "UnitCount")
	if err != nil {
		return err
	}
	newUnitCount := uint32(len(saveOutput.Units) + len(newUnits))
	binary.LittleEndian.PutUint32(newData[unitCountOffset:unitCountOffset+4], newUnitCount)

	if err := os.WriteFile(inputFilename, newData, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	// Register the newly inserted units in the file offset map so a later
	// AppendNewUnits call in the same run (e.g. fun mode spawning across many
	// provinces) can still resolve BuildUnitStartKey for them - they didn't
	// exist yet when ReadSaveFile built the map.
	for i := range newUnits {
		fileOffsetMap[BuildUnitStartKey(originalUnitCount+i)] = insertionOffset + i*unitSize
	}

	saveOutput.Units = append(saveOutput.Units, newUnits...)
	saveOutput.SaveHeader.UnitCount = newUnitCount
	return nil
}

// binaryFieldOffset returns the byte offset of fieldName within v as encoded
// by encoding/binary (i.e. fields packed in declaration order, ignoring Go's
// in-memory struct padding/alignment).
func binaryFieldOffset(v interface{}, fieldName string) (int, error) {
	t := reflect.TypeOf(v)
	offset := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			return offset, nil
		}
		size := binary.Size(reflect.New(field.Type).Elem().Interface())
		if size < 0 {
			return 0, fmt.Errorf("cannot determine binary size of field %s", field.Name)
		}
		offset += size
	}
	return 0, fmt.Errorf("field %s not found on %s", fieldName, t.Name())
}
