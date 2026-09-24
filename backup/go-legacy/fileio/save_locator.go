package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)

// nonSaveFiles lists ".sav" files in the LocalState folder that are not
// map saves this tool can parse (player profile / tracking data instead).
var nonSaveFiles = map[string]bool{
	"headquarter.sav": true,
	"uuid.sav":        true,
	"prd.sav":         true,
}

// FindDefaultSaveFile locates the World Conqueror 4 LocalState folder for
// the Microsoft Store install on this PC and returns the path to the most
// recently modified save file in it, so the tool can be used without
// having to type the full AppData path by hand.
func FindDefaultSaveFile() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", fmt.Errorf("LOCALAPPDATA is not set (auto-detection only works on Windows)")
	}

	matches, err := filepath.Glob(filepath.Join(localAppData, "Packages", "EasyTech.WorldConqueror4_*", "LocalState"))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("could not find a World Conqueror 4 install under %s\\Packages - pass -input manually", localAppData)
	}

	localStateDir := matches[0]
	if len(matches) > 1 {
		localStateDir = mostRecentlyModifiedDir(matches)
	}

	savePath, err := mostRecentSaveFile(localStateDir)
	if err != nil {
		return "", err
	}

	fmt.Printf("Auto-detected save file: %s\n", savePath)
	return savePath, nil
}

func mostRecentlyModifiedDir(dirs []string) string {
	best := dirs[0]
	var bestModTime int64
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			continue
		}
		if modTime := info.ModTime().Unix(); modTime > bestModTime {
			bestModTime = modTime
			best = dir
		}
	}
	return best
}

func mostRecentSaveFile(localStateDir string) (string, error) {
	entries, err := os.ReadDir(localStateDir)
	if err != nil {
		return "", err
	}

	var bestPath string
	var bestModTime int64
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sav" || nonSaveFiles[entry.Name()] {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}
		if modTime := info.ModTime().Unix(); bestPath == "" || modTime > bestModTime {
			bestModTime = modTime
			bestPath = filepath.Join(localStateDir, entry.Name())
		}
	}

	if bestPath == "" {
		return "", fmt.Errorf("no save files found in %s", localStateDir)
	}
	return bestPath, nil
}
