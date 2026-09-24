package fileio

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Color constants
const (
	OceanColor  = "#3D6A6A" // Teal for ocean tiles
	DefaultGray = "#C8C8C8" // Light gray for default/unowned
	InvalidGray = "#646464" // Dark gray for invalid owners
)

// MapVisualizer handles creating visual representations of the game map
type MapVisualizer struct {
	SaveOutput *WC4SaveOutput
	HexSize    float64
}

// CityOwnership represents which player owns a city
type CityOwnership struct {
	CoordinateCode uint16
	PlayerID       byte
	CountryID      uint32
}

// NewMapVisualizer creates a new map visualizer
func NewMapVisualizer(saveOutput *WC4SaveOutput) *MapVisualizer {
	return &MapVisualizer{
		SaveOutput: saveOutput,
		HexSize:    20.0, // Size of each hexagon
	}
}

// GenerateSVGMap creates an SVG visualization of the map with hexagonal tiles
func (mv *MapVisualizer) GenerateSVGMap(filename string) error {
	// Get map dimensions
	height := len(mv.SaveOutput.UnitOwnerData)
	width := len(mv.SaveOutput.UnitOwnerData[0])

	// Calculate SVG dimensions
	hexWidth := mv.HexSize * 2
	hexHeight := mv.HexSize * 1.732              // sqrt(3) for proper hex proportions
	svgWidth := float64(width) * hexWidth * 0.75 // 0.75 for hex spacing
	svgHeight := float64(height)*hexHeight + mv.HexSize

	// Start building SVG
	var svg strings.Builder
	svg.WriteString(fmt.Sprintf(`<svg width="%.0f" height="%.0f" xmlns="http://www.w3.org/2000/svg">`, svgWidth, svgHeight))
	svg.WriteString("\n")

	// Add style definitions
	svg.WriteString(`<defs>
		<style>
			.hex { stroke: #000; stroke-width: 1; }
			.city-center { fill: white; stroke: black; stroke-width: 2; }
			.unit-marker { stroke: black; stroke-width: 0.5; }
		</style>
	</defs>`)
	svg.WriteString("\n")

	// Color map for players/teams
	playerColors := mv.generatePlayerColors()

	// Process city tiles to determine city boundaries and ownership
	cityOwnership := mv.determineCityOwnership()

	// Process each tile as a hexagon
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			// Calculate hex position
			hexX, hexY := mv.getHexPosition(col, row)

			// Get tile color based on city ownership
			tileColor := mv.getCityTileColor(row, col, cityOwnership, playerColors)

			// Draw hexagon
			hexPoints := mv.generateHexPoints(hexX, hexY)
			svg.WriteString(fmt.Sprintf(`<polygon points="%s" fill="%s" class="hex"/>`, hexPoints, tileColor))
			svg.WriteString("\n")
		}
	}

	// Add city markers
	mv.drawCityMarkers(&svg)

	// Add unit markers
	mv.drawUnitMarkers(&svg, playerColors)

	// Close SVG
	svg.WriteString("</svg>")

	// Save SVG file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(svg.String())
	return err
}

// generatePlayerColors creates a color palette for players/teams
func (mv *MapVisualizer) generatePlayerColors() map[byte]string {
	colors := make(map[byte]string)

	// Use country colors for each player
	for i := 0; i < len(mv.SaveOutput.PlayerData); i++ {
		player := mv.SaveOutput.PlayerData[i]
		countryColor := GetCountryColorFromData(player)
		// Add # prefix to make it a proper hex color
		colors[byte(i)] = "#" + countryColor
	}

	return colors
}

// determineCityOwnership analyzes city tiles and unit owner data to determine city ownership
func (mv *MapVisualizer) determineCityOwnership() map[uint16]CityOwnership {
	cityOwnership := make(map[uint16]CityOwnership)

	// Process each city
	for _, city := range mv.SaveOutput.Cities {
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), mv.SaveOutput.UnitOwnerData, int(mv.SaveOutput.SaveHeader.GameMode))
		if !valid {
			continue
		}

		// Get the owner of the city center
		if row >= len(mv.SaveOutput.UnitOwnerData) || col >= len(mv.SaveOutput.UnitOwnerData[row]) {
			continue
		}
		owner := mv.SaveOutput.UnitOwnerData[row][col]

		// Get player info
		if int(owner) < len(mv.SaveOutput.PlayerData) {
			player := mv.SaveOutput.PlayerData[owner]
			cityOwnership[city.CoordinateCode] = CityOwnership{
				CoordinateCode: city.CoordinateCode,
				PlayerID:       owner,
				CountryID:      player.CountryId,
			}
		}
	}

	return cityOwnership
}

// getCityTileColor determines the color for a tile based on city ownership
func (mv *MapVisualizer) getCityTileColor(row, col int, cityOwnership map[uint16]CityOwnership, playerColors map[byte]string) string {
	// Check if we have city tiles data
	if row >= len(mv.SaveOutput.CityTiles) || col >= len(mv.SaveOutput.CityTiles[row]) {
		return OceanColor // Ocean color for out of bounds
	}

	cityTileID := mv.SaveOutput.CityTiles[row][col]

	// Check if tile is ocean (65535)
	if cityTileID == 65535 {
		return OceanColor // Ocean color
	}

	// The city tile ID is the coordinate code of the city that owns this tile
	// Look up the ownership for this coordinate code
	if ownership, exists := cityOwnership[cityTileID]; exists {
		// Get player color
		if playerColor, exists := playerColors[ownership.PlayerID]; exists {
			return playerColor
		}
	}

	// Default color for unowned cities
	return DefaultGray
}

// getHexPosition calculates the SVG position for a hexagon
func (mv *MapVisualizer) getHexPosition(col, row int) (float64, float64) {
	hexWidth := mv.HexSize * 2
	hexHeight := mv.HexSize * 1.732 // sqrt(3)

	x := float64(col) * hexWidth * 0.75
	y := float64(row) * hexHeight
	if col%2 == 1 {
		y += hexHeight / 2 // Offset odd columns
	}

	return x + mv.HexSize, y + mv.HexSize
}

// generateHexPoints creates the SVG points for a hexagon
func (mv *MapVisualizer) generateHexPoints(centerX, centerY float64) string {
	points := make([]string, 6)
	for i := 0; i < 6; i++ {
		angle := float64(i) * 60.0 * 3.14159 / 180.0 // Convert to radians
		x := centerX + mv.HexSize*math.Cos(angle)
		y := centerY + mv.HexSize*math.Sin(angle)
		points[i] = fmt.Sprintf("%.1f,%.1f", x, y)
	}
	return strings.Join(points, " ")
}

// drawCityMarkers draws city markers on the map
func (mv *MapVisualizer) drawCityMarkers(svg *strings.Builder) {
	for _, city := range mv.SaveOutput.Cities {
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), mv.SaveOutput.UnitOwnerData, int(mv.SaveOutput.SaveHeader.GameMode))
		if !valid {
			continue
		}

		// Calculate hex position
		hexX, hexY := mv.getHexPosition(col, row)

		// Draw city center marker (white circle)
		svg.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" class="city-center"/>`, hexX, hexY, mv.HexSize*0.3))
		svg.WriteString("\n")

		// Draw port icon for some building types
		if city.BuildingType == 32 || city.BuildingType == 33 || city.BuildingType == 34 {
			mv.drawPortIcon(svg, hexX, hexY)
		}

		// Add city name label for cities with CityId > 0
		if city.CityId > 0 {
			cityName, exists := GetCityName(city.CityId)
			if !exists {
				cityName = fmt.Sprintf("City ID %d", city.CityId)
			}
			// Position the text slightly below the city marker
			textY := hexY + mv.HexSize*0.8
			fontSize := mv.HexSize * 0.5
			svg.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-family="Arial, sans-serif" font-size="%.1f" font-weight="bold" fill="black">%s</text>`,
				hexX, textY, fontSize, cityName))
			svg.WriteString("\n")
		}
	}
}

// drawPortIcon draws a port icon (square) for port cities
func (mv *MapVisualizer) drawPortIcon(svg *strings.Builder, centerX, centerY float64) {
	// Draw a simple square icon
	iconSize := mv.HexSize * 0.2

	// Square positioned at city center
	left := centerX - iconSize
	top := centerY - iconSize

	svg.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="white" stroke="black" stroke-width="1"/>`,
		left, top, iconSize*2, iconSize*2))
	svg.WriteString("\n")
}

// drawUnitMarkers draws unit markers on the map
func (mv *MapVisualizer) drawUnitMarkers(svg *strings.Builder, playerColors map[byte]string) {

	for i := 0; i < len(mv.SaveOutput.Units); i++ {
		unit := mv.SaveOutput.Units[i]
		row, col, valid := ConvertCoordinates(int(unit.CoordinateCode), mv.SaveOutput.UnitOwnerData, int(mv.SaveOutput.SaveHeader.GameMode))
		if !valid {
			continue
		}

		// Calculate hex position
		hexX, hexY := mv.getHexPosition(col, row)

		// Get unit owner from the unit owner data
		var unitColor string
		if row >= len(mv.SaveOutput.UnitOwnerData) || col >= len(mv.SaveOutput.UnitOwnerData[row]) {
			unitColor = DefaultGray // Default gray for out of bounds
		} else {
			owner := mv.SaveOutput.UnitOwnerData[row][col]
			if int(owner) < len(mv.SaveOutput.PlayerData) {
				if playerColor, exists := playerColors[byte(owner)]; exists {
					unitColor = playerColor
				} else {
					unitColor = DefaultGray // Default gray
				}
			} else {
				unitColor = InvalidGray // Invalid owner - dark gray
			}
		}

		// Draw unit marker based on type
		if unit.UnitType != UnitTypeCity {
			// Other units get a small circle
			svg.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" class="unit-marker" fill="%s" stroke="black" stroke-width="0.5"/>`,
				hexX, hexY, mv.HexSize*0.15, unitColor))
		}
		svg.WriteString("\n")
	}
}

// VisualizeMap creates an SVG map visualization
func VisualizeMap(saveOutput *WC4SaveOutput, outputPath string) error {
	visualizer := NewMapVisualizer(saveOutput)

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Generate the map
	return visualizer.GenerateSVGMap(outputPath)
}
