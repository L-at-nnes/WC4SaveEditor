package fileio

import "fmt"

type CountryInfo struct {
	ID   uint8
	Name string
}

var countryCodes = []CountryInfo{
	{0x01, "UK"},
	{0x02, "France"},
	{0x03, "Germany"},
	{0x04, "Germany"},
	{0x05, "Soviet Union"},
	{0x06, "USA"},
	{0x07, "Italy"},
	{0x08, "ROC"},
	{0x09, "PRC"},
	{0x0a, "Japan"},
	{0x0b, "Finland"},
	{0x0c, "Poland"},
	{0x0d, "Yugoslavia"},
	{0x0e, "Canada"},
	{0x0f, "Australia"},
	{0x10, "Norway"},
	{0x11, "Sweden"},
	{0x12, "Denmark"},
	{0x13, "Netherlands"},
	{0x14, "Belgium"},
	{0x15, "Spain"},
	{0x16, "Portugal"},
	{0x17, "Hungary"},
	{0x18, "Romania"},
	{0x19, "Bulgaria"},
	{0x1a, "Switzerland"},
	{0x1b, "Greece"},
	{0x1c, "Turkey"},
	{0x1d, "Saudi Arabia"},
	{0x1e, "Iraq"},
	{0x1f, "Iran"},
	{0x20, "India"},
	{0x21, "Thailand"},
	{0x22, "Mongolia"},
	{0x23, "North Korea"},
	{0x24, "South Korea"},
	{0x25, "Mexico"},
	{0x26, "Cuba"},
	{0x27, "Colombia"},
	{0x28, "Brazil"},
	{0x29, "Bolivia"},
	{0x2a, "Venezuela"},
	{0x2b, "Peru"},
	{0x2c, "Chile"},
	{0x2d, "Argentina"},
	{0x2e, "Egypt"},
	{0x2f, "Liberia"},
	{0x30, "Mysterious Forces"},
	{0x31, "East Germany"},
	{0x33, "Americas Scorpion"},
	{0x34, "Asia-Pacific Scorpion"},
	{0x35, "European Scorpion"},
	{0x36, "Middle East Scorpion"},
	{0x37, "African Scorpion"},
}

var countryIDToName map[uint8]string

func init() {
	countryIDToName = make(map[uint8]string)
	for _, country := range countryCodes {
		countryIDToName[country.ID] = country.Name
	}
}

func GetCountryName(countryID uint8) (string, bool) {
	name, exists := countryIDToName[countryID]
	return name, exists
}

func GetAllCountries() []CountryInfo {
	return countryCodes
}

func GetCountryCount() int {
	return len(countryCodes)
}

func ColorBytesToHex(colorBytes [4]byte) string {
	return fmt.Sprintf("%02x%02x%02x", colorBytes[0], colorBytes[1], colorBytes[2])
}

func GetCountryInfoFromData(countryData CountryData) (string, string) {
	name, _ := GetCountryName(uint8(countryData.CountryId))
	color := ColorBytesToHex(countryData.PrimaryColor)
	return name, color
}

func GetCountryColorFromData(countryData CountryData) string {
	return ColorBytesToHex(countryData.PrimaryColor)
}
