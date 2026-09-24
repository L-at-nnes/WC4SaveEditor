# WC4 Save Editor

A save file editor for World Conqueror 4 (Microsoft Store version), with a dark, native Windows GUI.

## Usage

1. Run `build.ps1` (PowerShell) to produce a single self-contained `WC4SaveEditor.exe` at the repository root. No .NET runtime needs to be installed on the machine that runs it.
2. Launch `WC4SaveEditor.exe`. It auto-detects your World Conqueror 4 save folder (`%LOCALAPPDATA%\Packages\EasyTech.WorldConqueror4_*\LocalState`) and lists every save file found there along the top bar.
3. Pick a save file, configure whatever changes you want across the panels below, then click **Save** in the top-right corner. Nothing is written to disk until you click Save — a `.bak` copy of the original file is created automatically the first time you save.

## Panels

- **Quick actions** — one-click toggles for the main player (player 0): Conquer All, Max Money, Max City Tech, Max City Level, Max Morale, Heal Allies, Weaken Enemies.
- **Teams** — drag a country card between team columns to reassign it to that team. "Unite team" moves every player onto your team in one click.
- **Conquer player** — pick one player to annex from the left list, check as many target players as you want on the right; every tile, unit and landmine they own transfers to the annexer.
- **Conquer province** — same layout, but clicking a target opens a picker listing their individual provinces, so you can annex just part of a country.
- **Spawn units** — pick one of your own provinces, then a unit type/level/count to queue a new unit spawn (only unit types that already exist somewhere in the save can be spawned, since their base stats are cloned from an existing unit).
- **Transfer tile** — an interactive, zoomable hex map at the bottom of the window. Scroll to zoom, drag to pan, click a tile to reassign its owner.

## Project layout

```
src/WC4SaveEditor.Core        Save file format: reading, writing, all game logic (no UI dependency)
src/WC4SaveEditor.Core.Tests  xUnit tests pinning the binary format's invariants
src/WC4SaveEditor.Gui         WPF application
build.ps1                     Publishes the self-contained single-file exe to the repo root
backup/go-legacy/             The original Go CLI version of this tool, kept for reference
```

## Legacy CLI

This project started as a Go command-line tool. That implementation still exists, unmodified, under `backup/go-legacy/` (with its own `README.md`), in case the CLI workflow or the Go source is useful as a reference. It is not maintained going forward — the WPF app in `src/` is the actively developed version and re-implements the same save file format and operations.

## File format

The save file format itself (structs, offsets, coordinate math) is documented in `backup/go-legacy/FILE_FORMAT.md`; the C# port in `src/WC4SaveEditor.Core` follows it field-for-field.
