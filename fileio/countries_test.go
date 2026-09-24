package fileio

import (
	"testing"
)

func TestGetCountryName(t *testing.T) {
	tests := []struct {
		name        string
		countryID   uint8
		expected    string
		shouldExist bool
	}{
		{
			name:        "Valid country ID - UK",
			countryID:   0x01,
			expected:    "UK",
			shouldExist: true,
		},
		{
			name:        "Valid country ID - France",
			countryID:   0x02,
			expected:    "France",
			shouldExist: true,
		},
		{
			name:        "Valid country ID - USA",
			countryID:   0x06,
			expected:    "USA",
			shouldExist: true,
		},
		{
			name:        "Valid country ID - Japan",
			countryID:   0x0a,
			expected:    "Japan",
			shouldExist: true,
		},
		{
			name:        "Valid country ID - Last country",
			countryID:   0x30,
			expected:    "Mysterious Forces",
			shouldExist: true,
		},
		{
			name:        "Invalid country ID - too high",
			countryID:   0x31,
			expected:    "",
			shouldExist: false,
		},
		{
			name:        "Invalid country ID - zero",
			countryID:   0x00,
			expected:    "",
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := GetCountryName(tt.countryID)
			if exists != tt.shouldExist {
				t.Errorf("GetCountryName(%d) exists = %v, want %v", tt.countryID, exists, tt.shouldExist)
			}
			if exists && result != tt.expected {
				t.Errorf("GetCountryName(%d) = %v, want %v", tt.countryID, result, tt.expected)
			}
		})
	}
}

func TestGetAllCountries(t *testing.T) {
	countries := GetAllCountries()

	// Test that we get the expected number of countries
	expectedCount := 53 // Based on the countryCodes slice length
	if len(countries) != expectedCount {
		t.Errorf("GetAllCountries() returned %d countries, want %d", len(countries), expectedCount)
	}

	// Test that the first and last countries are correct
	if len(countries) > 0 {
		first := countries[0]
		if first.ID != 0x01 || first.Name != "UK" {
			t.Errorf("First country = {ID: 0x%02x, Name: %s}, want {ID: 0x01, Name: UK}",
				first.ID, first.Name)
		}

		last := countries[len(countries)-1]
		if last.ID != 0x37 || last.Name != "African Scorpion" {
			t.Errorf("Last country = {ID: 0x%02x, Name: %s}, want {ID: 0x37, Name: African Scorpion}",
				last.ID, last.Name)
		}
	}
}

func TestGetCountryCount(t *testing.T) {
	count := GetCountryCount()
	expectedCount := 53 // Based on the countryCodes slice length

	if count != expectedCount {
		t.Errorf("GetCountryCount() = %d, want %d", count, expectedCount)
	}
}

// Test ColorBytesToHex function
func TestColorBytesToHex(t *testing.T) {
	tests := []struct {
		name       string
		colorBytes [4]byte
		expected   string
	}{
		{
			name:       "Red color",
			colorBytes: [4]byte{255, 0, 0, 0},
			expected:   "ff0000",
		},
		{
			name:       "Green color",
			colorBytes: [4]byte{0, 255, 0, 0},
			expected:   "00ff00",
		},
		{
			name:       "Blue color",
			colorBytes: [4]byte{0, 0, 255, 0},
			expected:   "0000ff",
		},
		{
			name:       "White color",
			colorBytes: [4]byte{255, 255, 255, 0},
			expected:   "ffffff",
		},
		{
			name:       "Black color",
			colorBytes: [4]byte{0, 0, 0, 0},
			expected:   "000000",
		},
		{
			name:       "Mixed color",
			colorBytes: [4]byte{128, 64, 192, 0},
			expected:   "8040c0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorBytesToHex(tt.colorBytes)
			if result != tt.expected {
				t.Errorf("ColorBytesToHex(%v) = %v, want %v", tt.colorBytes, result, tt.expected)
			}
		})
	}
}

// Test GetCountryInfoFromData function
func TestGetCountryInfoFromData(t *testing.T) {
	tests := []struct {
		name          string
		countryData   CountryData
		expectedName  string
		expectedColor string
	}{
		{
			name: "UK country data",
			countryData: CountryData{
				CountryId:    0x01,
				PrimaryColor: [4]byte{250, 250, 150, 0}, // fafa96
			},
			expectedName:  "UK",
			expectedColor: "fafa96",
		},
		{
			name: "France country data",
			countryData: CountryData{
				CountryId:    0x02,
				PrimaryColor: [4]byte{150, 230, 255, 0}, // 96e6ff
			},
			expectedName:  "France",
			expectedColor: "96e6ff",
		},
		{
			name: "Unknown country data",
			countryData: CountryData{
				CountryId:    0xFF, // Invalid country ID
				PrimaryColor: [4]byte{128, 64, 192, 0},
			},
			expectedName:  "",
			expectedColor: "8040c0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, color := GetCountryInfoFromData(tt.countryData)
			if name != tt.expectedName {
				t.Errorf("GetCountryInfoFromData() name = %v, want %v", name, tt.expectedName)
			}
			if color != tt.expectedColor {
				t.Errorf("GetCountryInfoFromData() color = %v, want %v", color, tt.expectedColor)
			}
		})
	}
}

// Test GetCountryColorFromData function
func TestGetCountryColorFromData(t *testing.T) {
	tests := []struct {
		name        string
		countryData CountryData
		expected    string
	}{
		{
			name: "Red color",
			countryData: CountryData{
				CountryId:    0x01,
				PrimaryColor: [4]byte{255, 0, 0, 0},
			},
			expected: "ff0000",
		},
		{
			name: "Green color",
			countryData: CountryData{
				CountryId:    0x02,
				PrimaryColor: [4]byte{0, 255, 0, 0},
			},
			expected: "00ff00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCountryColorFromData(tt.countryData)
			if result != tt.expected {
				t.Errorf("GetCountryColorFromData() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test that the lookup maps are properly initialized
func TestLookupMapsInitialization(t *testing.T) {
	// Test that all countries in countryCodes are in the lookup maps
	for _, country := range countryCodes {
		name, exists := GetCountryName(country.ID)
		if !exists {
			t.Errorf("Country ID 0x%02x (%s) not found in name lookup map", country.ID, country.Name)
		}
		if name != country.Name {
			t.Errorf("Country ID 0x%02x name = %s, want %s", country.ID, name, country.Name)
		}
	}
}
