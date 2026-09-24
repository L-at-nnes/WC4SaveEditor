# World Conqueror 4 Save File Format

This document describes the binary format of World Conqueror 4 save files.

## File Structure Overview

The save file consists of several sections:
1. Save Header
2. Player/Country Data
3. City Tiles Data
4. Unit Owner Data
5. Cities Data
6. Units Data

## Save Header (SaveHeader)

The save header contains metadata about the save file and game state.

| Field | Type | Size | Description |
|-------|------|------|-------------|
| Magic | [4]byte | 4 | File magic identifier |
| UnknownInt1 | uint32 | 4 | Unknown integer value |
| MapId | uint32 | 4 | Map identifier |
| GameMode | uint32 | 4 | Game mode: 1=campaign, 2=conquest, 6=frontier |
| UnknownInt2 | uint32 | 4 | Unknown integer value |
| UnknownInt3 | uint32 | 4 | Unknown integer value |
| Camera | [3]float32 | 12 | Camera position (x, y, z) |
| UnknownInt4 | uint32 | 4 | Unknown integer value |
| TurnNumber | uint32 | 4 | Current turn number |
| UnknownArr2 | [12]byte | 12 | Unknown byte array |
| SaveTimestamp | [5]uint32 | 20 | Save file timestamp data |
| UnknownArr3 | [16]byte | 16 | Unknown byte array |
| UnknownInt7 | uint32 | 4 | Non-zero only in frontier mode last mission |
| UnknownInt8 | uint32 | 4 | Non-zero only in frontier mode last mission |
| MapWidth | uint32 | 4 | Map width in tiles |
| MapHeight | uint32 | 4 | Map height in tiles |
| CountryCount | uint32 | 4 | Number of countries/players |
| CityCount | uint32 | 4 | Number of cities |
| UnitCount | uint32 | 4 | Number of units |
| UnknownCount1 | uint32 | 4 | Unknown count |
| UnknownCount2 | uint32 | 4 | Unknown count |
| UnknownArr4 | [8]byte | 8 | Unknown byte array |
| TurnCount1 | uint32 | 4 | Turn count (variant 1) |
| TurnCount2 | uint32 | 4 | Turn count (variant 2) |
| UnknownCount3 | uint32 | 4 | Unknown count |
| UnknownCount4 | uint32 | 4 | Unknown count |
| UnknownCount5 | uint32 | 4 | Unknown count |
| UnknownCount6 | uint32 | 4 | Unknown count |
| ImportantCityCount | uint32 | 4 | Number of important cities |
| UnknownArr5 | [4]byte | 4 | Unknown byte array |
| UnknownInt9 | uint32 | 4 | Unknown integer value |
| UnknownInt10 | uint32 | 4 | Unknown integer value |
| UnknownArr6 | [12]byte | 12 | Unknown byte array |
| LandmineCount | uint32 | 4 | Number of landmines |
| UnknownArr7 | [16]byte | 16 | Unknown byte array |
| UnknownCount9 | uint32 | 4 | Unknown count |

**Total Size: 208 bytes**

## Country/Player Data (CountryData)

Each country/player has associated data.

| Field | Type | Size | Description |
|-------|------|------|-------------|
| TurnOrder | uint32 | 4 | Turn order in the game |
| CountryId | uint32 | 4 | Unique country identifier |
| Currency | [3]uint32 | 12 | Currency amounts (3 different types) |
| BotFlag | uint32 | 4 | Flag indicating if this is an AI player |
| TeamId | uint32 | 4 | Team identifier for alliances |
| UnknownArr2 | [4]byte | 4 | Unknown byte array |
| UnknownColor | [2][4]byte | 8 | Unknown color data (2 colors) |
| PrimaryColor | [4]byte | 4 | Primary country color (RGBA) |
| UnknownArr4 | [16]byte | 16 | Unknown byte array |
| UnknownArr5 | [460]byte | 460 | Large unknown byte array |

**Total Size: 520 bytes per country**

## City Data (CityData)

Information about each city on the map.

| Field | Type | Size | Description |
|-------|------|------|-------------|
| CoordinateCode | uint16 | 2 | Encoded position on the map |
| CityId | uint16 | 2 | Unique city identifier |
| BuildingType | uint8 | 1 | Type of building in the city. For ports, encodes level as 30+level (31-34). For regular cities, encodes level as 10+level (11-14, or 15+ for capital "star" cities). |
| Apperance | uint8 | 1 | Visual appearance of the city |
| UnknownByte1 | uint8 | 1 | Unknown byte value |
| Wonders | uint8 | 1 | Wonder buildings in the city |
| UnknownArr2 | [6]byte | 6 | Unknown byte array |
| UnknownArr3 | [8]byte | 8 | Unknown byte array |
| AntiAirWeaponType | uint8 | 1 | Type of anti-air weapon |
| AntiAirRange | uint8 | 1 | Range of anti-air weapon |
| TechLevels | [6]byte | 6 | Technology levels (6 different techs) |
| UnknownArr4 | [2]byte | 2 | Unknown byte array |

**Total Size: 32 bytes per city**

## Unit Data (UnitData)

Information about each unit on the map.

| Field | Type | Size | Description |
|-------|------|------|-------------|
| CoordinateCode | uint16 | 2 | Encoded position on the map |
| UnitType | uint8 | 1 | Type of unit |
| Level | uint8 | 1 | Unit level |
| Personnel | uint8 | 1 | Number of personnel |
| Direction | uint8 | 1 | Facing direction (0=left, 1=right) |
| Movement | uint16 | 2 | Movement points |
| Experience | uint16 | 2 | Unit experience |
| UnknownHealth | uint16 | 2 | Unknown health value |
| CurrentHealth | uint16 | 2 | Current health points |
| MaxHealth | uint16 | 2 | Maximum health points |
| GeneralId | uint16 | 2 | ID of assigned general (0 if none) |
| GeneralMilitaryRank | uint8 | 1 | General's military rank |
| GeneralTitle | uint8 | 1 | General's title |
| GeneralBadges | [3]byte | 3 | General's badges |
| GeneralSkillLevels | [5]byte | 5 | General's skill levels (5 skills) |
| UnknownArr5 | [12]byte | 12 | Unknown byte array |
| MoraleValue | int8 | 1 | Unit morale value |
| MoraleTurnsLeft | uint16 | 2 | Turns remaining for morale effect |
| UnknownArr6 | [21]byte | 21 | Unknown byte array |

**Total Size: 64 bytes per unit**

## Landmine Data (LandmineData)

Information about landmines on the map.

| Field | Type | Size | Description |
|-------|------|------|-------------|
| CoordinateCode | uint16 | 2 | Encoded position on the map |
| Owner | uint16 | 2 | One-based player ID of the owner (player 0 is stored as 1) |
| UnknownArr1 | [2]byte | 2 | Unknown byte array |
| Health | uint16 | 2 | Landmine health |
| UnknownArr2 | [4]byte | 4 | Unknown byte array |

**Total Size: 12 bytes per landmine**

## Notes

- All multi-byte integers are stored in little-endian format
