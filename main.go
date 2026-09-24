package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/samuelyuan/WC4SaveEditor/fileio"
)

func warnf(format string, a ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", a...)
}

func parseInt(s string, fieldName string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error: Invalid %s '%s'. Must be a number.\n", fieldName, s)
		return -1
	}
	return val
}

// writeCommands are commands that modify the save file and therefore need
// a backup taken first; every other command is read-only.
var writeCommands = map[string]bool{
	"max-money":        true,
	"max-city-tech":    true,
	"max-city-level":   true,
	"max-morale":       true,
	"heal-allies":      true,
	"weaken-enemies":   true,
	"conquer-player":   true,
	"conquer-province": true,
	"change-team":      true,
	"transfer-tile":    true,
	"spawn-units":      true,
	"conquer-allies":   true,
	"unite-team":       true,
	"conquer-all":      true,
}

func showHelp() {
	fmt.Println("WC4SaveEditor - World Conqueror 4 Save File Editor")
	fmt.Println("")
	fmt.Println("Usage: WC4SaveEditor.exe -input <savefile> -command <command>")
	fmt.Println("If -input is omitted, the save file is auto-detected from the")
	fmt.Println("Microsoft Store install's LocalState folder.")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list-players, list-cities, list-units, list-units-by-map, list-tiles, list-generals, list-landmines, list-teams")
	fmt.Println("  max-money, max-city-tech, max-city-level, max-morale, heal-allies, weaken-enemies")
	fmt.Println("  conquer-player, conquer-province, transfer-tile (interactive)")
	fmt.Println("  spawn-units (interactive)")
	fmt.Println("  change-team (interactive)")
	fmt.Println("  conquer-allies, unite-team, conquer-all")
	fmt.Println("  visualize-map")
	fmt.Println("")
	fmt.Println("Use -help to show this message")
}

func main() {
	inputFilenamePtr := flag.String("input", "", "Path to the World Conqueror 4 save file")
	commandPtr := flag.String("command", "", "Command to execute (use -help for list of commands)")
	helpPtr := flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Show help if requested or no command provided
	if *helpPtr || *commandPtr == "" {
		showHelp()
		return
	}

	inputFilename := *inputFilenamePtr
	command := *commandPtr

	if inputFilename == "" {
		detectedFilename, err := fileio.FindDefaultSaveFile()
		if err != nil {
			fmt.Printf("Error: -input flag was not set and auto-detection failed: %v\n", err)
			fmt.Println("Use -help for usage information")
			return
		}
		inputFilename = detectedFilename
	}

	saveOutput, err := fileio.ReadSaveFile(inputFilename)
	if err != nil {
		fmt.Printf("Error reading save file '%s': %v\n", inputFilename, err)
		fmt.Println("Make sure the file exists and is a valid World Conqueror 4 save file")
		return
	}

	if writeCommands[command] {
		if err := fileio.EnsureBackup(inputFilename); err != nil {
			fmt.Printf("Error creating backup before writing to '%s': %v\n", inputFilename, err)
			return
		}
	}

	if command == "list-players" {
		fileio.ListPlayers(saveOutput)
	} else if command == "list-cities" {
		fileio.ListCities(saveOutput)
	} else if command == "list-units" {
		fileio.ListUnits(saveOutput)
	} else if command == "list-units-by-map" {
		fileio.ListUnitsByMap(saveOutput)
	} else if command == "list-tiles" {
		fileio.ListTiles(saveOutput)
	} else if command == "list-generals" {
		fileio.ListGenerals(saveOutput)
	} else if command == "list-landmines" {
		fileio.ListLandmines(saveOutput)
	} else if command == "list-teams" {
		fileio.ListTeams(saveOutput)
	} else if command == "max-money" {
		fileio.SetPlayerMaxCurrency(inputFilename, 0, 9999)
	} else if command == "max-city-tech" {
		fileio.SetPlayerMaxCityTech(inputFilename, saveOutput, 0, 4)
	} else if command == "max-city-level" {
		fileio.SetPlayerMaxCityLevel(inputFilename, saveOutput, 0, 4)
	} else if command == "max-morale" {
		fileio.SetPlayerMaxMorale(inputFilename, saveOutput, 0)
	} else if command == "heal-allies" {
		fileio.RestoreAlliesHealth(inputFilename, saveOutput, 0)
	} else if command == "weaken-enemies" {
		fileio.WeakenEnemies(inputFilename, saveOutput, 0)
	} else if command == "conquer-player" {
		fileio.InteractiveConvertPlayer(inputFilename, saveOutput)
	} else if command == "conquer-province" {
		fileio.InteractiveConvertProvince(inputFilename, saveOutput)
	} else if command == "change-team" {
		fileio.InteractiveChangeTeam(inputFilename, saveOutput)
	} else if command == "transfer-tile" {
		fileio.InteractiveConvertTile(inputFilename, saveOutput)
	} else if command == "spawn-units" {
		fileio.InteractiveSpawnUnits(inputFilename, saveOutput)
	} else if command == "conquer-allies" {
		stats := fileio.ConvertAllAllies(inputFilename, saveOutput)
		stats.PrintSummary("Conquer Allies", saveOutput)
	} else if command == "unite-team" {
		fileio.ConvertTeam(inputFilename, saveOutput)
	} else if command == "conquer-all" {
		stats := fileio.ConvertAllPlayers(inputFilename, saveOutput)
		stats.PrintSummary("Conquer All Players", saveOutput)
	} else if command == "visualize-map" {
		outputPath := "map_visualization.svg"
		err := fileio.VisualizeMap(saveOutput, outputPath)
		if err != nil {
			fmt.Printf("Error creating map visualization: %v\n", err)
			return
		}
		fmt.Printf("Map visualization saved to: %s\n", outputPath)
	} else {
		fmt.Printf("Error: Unrecognized command '%s'\n", command)
		fmt.Println("Use -help to see available commands")
		return
	}
}
