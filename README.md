# WC4SaveEditor

## Save File Location

Before using this tool, you need to locate your World Conqueror 4 save file:

**Windows (Microsoft Store version):**
```
C:\Users\<YourUserName>\AppData\Local\Packages\EasyTech.WorldConqueror4_<PackageIdentifier>\LocalState
```

Replace `<YourUserName>` with your actual Windows username. The `<PackageIdentifier>` is a unique string that may vary.

## Documentation

- [File Format Documentation](FILE_FORMAT.md) - Detailed description of the World Conqueror 4 save file format

## How to Use

There are various commands to modify the save file. 

Make sure you quit your current game and go to the main menu before overwriting the save file. If you overwrite the file while the game is still in progress, the game will overwrite the file when you leave and none of your new changes will apply.

Read Commands:
* list-players: Display all players with country information, team membership, and unit counts
* list-cities: Show all cities with their positions and ownership details
* list-units: Display all military units grouped by owner with detailed unit information
* list-units-by-map: Analyze unit distribution across the map
* list-tiles: Show city tile ownership analysis with territory control statistics
* list-generals: Display all units with assigned generals
* list-landmines: Show all landmines grouped by owner with position and health data
* list-teams: Display team analysis with player counts, territories, and alliance structure
* visualize-map: Generate an SVG map visualization showing territories, cities, and units with hexagonal tiles

Write Commands:
* max-money: Sets max currency to 9999.
* max-city-tech: Sets all city tech levels to level 4.
* max-city-level: Upgrades all city levels (defense points) up to level 4. Capital "star" cities are left untouched, and cities already at level 4 or higher are not downgraded.
* max-morale: Sets maximum morale for your combat units, excluding cities and fortifications.
* heal-allies: Heal all of your units and your allies units.
* weaken-enemies: Minimize all enemy money to 0, reduce all enemy city tech to 0, and reduce all enemy units to have 1 health and all enemy cities to have 0 health.
* conquer-player: Conquer all territories from a specific player (interactive). May crash game.
* conquer-province: Conquer one city province, including its units, ports, structures, and landmines (interactive).
* change-team: Change one player to any Team ID (interactive).
* transfer-tile: Convert one tile and assign ownership to another player (interactive). May crash game.
* conquer-allies: Conquer all allied territories and make them yours. May crash game.
* unite-team: Convert all players to be on the same team.
* conquer-all: Convert all tiles to be your tiles. May crash game.
* spawn-units: Spawn new special units (King Tiger, Stuka zu Fuss, Hawkeye Force, Richelieu, etc. - the units with an individual 1-9 level, as opposed to basic units which scale with city tech instead) into a province (interactive). Choose the owning player, then a spawn mode:
  * Manual: pick a province, a unit type, its level, and how many to create. Land units spawn on the nearest free tile belonging to the province; naval units spawn on the nearest free sea tile next to one of the province's ports. Stats (health, movement, etc.) are cloned from an existing unit of the same type in the save (scaled if only a different level is available), since the save format has no formula for them. If there aren't enough free tiles for the requested count, or no unit of that type exists anywhere in the save to use as a reference, nothing is spawned.
  * Fun: choose a fill percentage (10% to 100%, in steps of 10), then fills that share of the free tiles of ALL of the player's provinces with a random unit type at max level (9) in one go - land tiles get random land units, and free sea tiles near each province's ports get random ships (never placed on land). Which tiles get filled below 100% is picked at random, not nearest-first. Unit types with no existing unit in the save to clone stats from are skipped.

## Usage Examples

```bash
# Template: WC4SaveEditor.exe -input <savefile> -command <command>
```
