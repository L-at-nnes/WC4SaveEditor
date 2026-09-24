package fileio

// CityInfo represents a city with its ID and name
type CityInfo struct {
	ID   uint16
	Name string
}

// City codes for World Conqueror 4
// Format: (city_id_hex, name)
var cityCodes = []CityInfo{
	{0x01, "London"},
	{0x02, "Birmingham"},
	{0x03, "Manchester"},
	{0x04, "Dublin"},
	{0x05, "Plymouth"},
	{0x06, "Liverpool"},
	{0x07, "Glasgow"},
	{0x08, "Galway"},
	{0x09, "Stockholm"},
	{0x0a, "Gothenburg"},
	{0x0b, "Oslo"},
	{0x0c, "Bergen"},
	{0x0d, "Paris"},
	{0x0e, "Bordeaux"},
	{0x0f, "Marseille"},
	{0x10, "Lyon"},
	{0x11, "Brest"},
	{0x12, "Metz"},
	{0x13, "Strasbourg"},
	{0x14, "Dijon"},
	{0x15, "Toulouse"},
	{0x16, "Palma"},
	{0x17, "Ajaccio"},
	{0x18, "Madrid"},
	{0x19, "Barcelona"},
	{0x1a, "Seville"},
	{0x1b, "Zaragoza"},
	{0x1c, "A Coruña"},
	{0x1d, "Badajoz"},
	{0x1e, "Valencia"},
	{0x1f, "Lisbon"},
	{0x20, "Porto"},
	{0x21, "Brussels"},
	{0x22, "Antwerp"},
	{0x23, "Amsterdam"},
	{0x24, "Berlin"},
	{0x25, "Hamburg"},
	{0x26, "Cologne"},
	{0x27, "Munich"},
	{0x28, "Dortmund"},
	{0x29, "Frankfurt"},
	{0x2a, "Nuremberg"},
	{0x2b, "Zurich"},
	{0x2c, "Rome"},
	{0x2d, "Milan"},
	{0x2e, "Naples"},
	{0x2f, "Turin"},
	{0x30, "Venice"},
	{0x31, "Palermo"},
	{0x32, "Vienna"},
	{0x33, "Warsaw"},
	{0x34, "Poznan"},
	{0x35, "Bialystok"},
	{0x36, "Lublin"},
	{0x37, "Krakow"},
	{0x38, "Copenhagen"},
	{0x39, "Budapest"},
	{0x3a, "Belgrade"},
	{0x3b, "Sarajevo"},
	{0x3c, "Sofia"},
	{0x3d, "Athens"},
	{0x3e, "Prague"},
	{0x3f, "Bucharest"},
	{0x40, "Cluj-Napoca"},
	{0x41, "Königsberg"},
	{0x42, "Riga"},
	{0x43, "Helsinki"},
	{0x44, "Turku"},
	{0x45, "Joensuu"},
	{0x46, "Ankara"},
	{0x47, "Istanbul"},
	{0x48, "Trabzon"},
	{0x49, "Izmir"},
	{0x4a, "Aleppo"},
	{0x4b, "Jerusalem"},
	{0x4c, "Damascus"},
	{0x4d, "Baghdad"},
	{0x4e, "Basra"},
	{0x4f, "Riyadh"},
	{0x50, "Medina"},
	{0x51, "Jeddah"},
	{0x52, "Dubai"},
	{0x53, "Sanaa"},
	{0x54, "Salalah"},
	{0x55, "Tehran"},
	{0x56, "Isfahan"},
	{0x57, "Mashhad"},
	{0x58, "Almaty"},
	{0x59, "Astana"},
	{0x5a, "Aktyubinsk"},
	{0x5b, "Ashgabat"},
	{0x5c, "Moscow"},
	{0x5d, "Yekaterinburg"},
	{0x5e, "Novosibirsk"},
	{0x5f, "Leningrad"},
	{0x60, "Smolensk"},
	{0x61, "Gorky"},
	{0x62, "Kuybyshev"},
	{0x63, "Saratov"},
	{0x64, "Kazan"},
	{0x65, "Orenbug"},
	{0x66, "Ufa"},
	{0x67, "Chelyabinsk"},
	{0x68, "Omsk"},
	{0x69, "Krasnoyarsk"},
	{0x6a, "Vladivostok"},
	{0x6b, "Chita"},
	{0x6c, "Khabarovsk"},
	{0x6d, "Petropavlovsk"},
	{0x6e, "Petrozavodsk"},
	{0x6f, "Voronezh"},
	{0x70, "Yakutsk"},
	{0x71, "Minsk"},
	{0x72, "Gomel"},
	{0x73, "Kiev"},
	{0x74, "Kharkov"},
	{0x75, "Donetsk"},
	{0x76, "Sevastopol"},
	{0x77, "Stalingrad"},
	{0x78, "Tbilisi"},
	{0x79, "Tashkent"},
	{0x7a, "Dushanbe"},
	{0x7b, "Uliastai"},
	{0x7c, "Choibalsan"},
	{0x7d, "Ulaanbaatar"},
	{0x7e, "Beijing"},
	{0x7f, "Nanjing"},
	{0x80, "Shanghai"},
	{0x81, "Wuhan"},
	{0x82, "Chongqing"},
	{0x83, "Hong Kong"},
	{0x84, "Guihui"},
	{0x85, "Lanzhou"},
	{0x86, "Kunming"},
	{0x87, "Changsha"},
	{0x88, "Urumqi"},
	{0x89, "Lhasa"},
	{0x8a, "Kashgar"},
	{0x8b, "Taipei"},
	{0x8d, "Changchun"},
	{0x8e, "Shenyang"},
	{0x8f, "Harbin"},
	{0x90, "New Delhi"},
	{0x91, "Kolkata"},
	{0x92, "Mumbai"},
	{0x93, "Lucknow"},
	{0x94, "Jaipur"},
	{0x95, "Bhopal"},
	{0x96, "Karachi"},
	{0x97, "Chennai"},
	{0x98, "Mandalay"},
	{0x99, "Yangon"},
	{0x9a, "Bangkok"},
	{0x9b, "Chiang Mai"},
	{0x9c, "Singapore"},
	{0x9d, "Kuala Lumpur"},
	{0x9e, "Hanoi"},
	{0x9f, "Phnom Penh"},
	{0xa0, "Tokyo"},
	{0xa1, "Yokohama"},
	{0xa2, "Osaka"},
	{0xa3, "Sendai"},
	{0xa4, "Sapporo"},
	{0xa5, "Nagasaki"},
	{0xa6, "Seoul"},
	{0xa7, "Pyongyang"},
	{0xa8, "Busan"},
	{0xa9, "Manila"},
	{0xaa, "Jakarta"},
	{0xab, "Port Moresby"},
	{0xac, "Balikpapan"},
	{0xad, "Medan"},
	{0xae, "Kupang"},
	{0xaf, "Darwin"},
	{0xb0, "Brisbane"},
	{0xb1, "Broome"},
	{0xb2, "Port Hedland"},
	{0xb4, "Alice Springs"},
	{0xb5, "Longreach"},
	{0xb6, "Casablanca"},
	{0xb7, "Algiers"},
	{0xb8, "Adrar"},
	{0xb9, "Tunis"},
	{0xba, "Tripoli"},
	{0xbb, "Benghazi"},
	{0xbc, "Al Jawf"},
	{0xbd, "Cairo"},
	{0xbe, "Alexandria"},
	{0xbf, "Aswan"},
	{0xc0, "Addis Ababa"},
	{0xc1, "Mogadishu"},
	{0xc2, "Khartoum"},
	{0xc3, "Port Sudan"},
	{0xc4, "Juba"},
	{0xc5, "Brazzaville"},
	{0xc6, "Yaoundé"},
	{0xc7, "N'Djamena"},
	{0xc8, "Bangui"},
	{0xc9, "Dakar"},
	{0xca, "Bamako"},
	{0xcb, "Freetown"},
	{0xcd, "Niamey"},
	{0xce, "Ottawa"},
	{0xcf, "Montreal"},
	{0xd0, "Toronto"},
	{0xd1, "Edmonton"},
	{0xd2, "Calgary"},
	{0xd3, "Whitehorse"},
	{0xd5, "Vancouver"},
	{0xd6, "Yellowknife"},
	{0xd7, "Winnipeg"},
	{0xda, "Prince George"},
	{0xdb, "Baker Lake"},
	{0xdc, "Gillam"},
	{0xdd, "Kuujjuaq"},
	{0xde, "Rankin Inlet"},
	{0xdf, "Washington"},
	{0xe0, "New York"},
	{0xe1, "Miami"},
	{0xe2, "Houston"},
	{0xe3, "Denver"},
	{0xe4, "San Francisco"},
	{0xe5, "Los Angeles"},
	{0xe6, "Billings"},
	{0xe7, "Minneapolis"},
	{0xe8, "Phoenix"},
	{0xe9, "Dallas"},
	{0xea, "Chicago"},
	{0xeb, "Detroit"},
	{0xec, "Atlanta"},
	{0xed, "Hilo"},
	{0xee, "Honolulu"},
	{0xef, "Henderson Field"},
	{0xf0, "Mexico City"},
	{0xf1, "Monterrey"},
	{0xf2, "Guadalajara"},
	{0xf3, "Mérida"},
	{0xf4, "Havana"},
	{0xf5, "Santiago"},
	{0xf6, "Port-au-Prince"},
	{0xf7, "Bogotá"},
	{0xf8, "Cali"},
	{0xf9, "Caracas"},
	{0xfa, "Paramaribo"},
	{0xfb, "Lima"},
	{0xfc, "Quito"},
	{0xfd, "Santa Cruz"},
	{0xfe, "La Paz"},
	{0xff, "Rio de Janeiro"},
	{0x100, "Fortaleza"},
	{0x104, "Anchorage"},
	{0x105, "Gdansk"},
	{0x106, "Rotterdam"},
	{0x107, "Amiens"},
	{0x108, "Tobruk"},
	{0x109, "Vyazma"},
	{0x10a, "Bryansk"},
	{0x10b, "Rostov"},
	{0x10c, "Luhansk"},
	{0x10d, "Kaluga"},
	{0x10e, "Mariupol"},
	{0x10f, "Swansea"},
	{0x110, "Nantes"},
	{0x111, "Rouen"},
	{0x112, "Malaga"},
	{0x113, "Irkutsk"},
	{0x114, "Jinan"},
	{0x115, "Yinchuan"},
	{0x116, "Xiamen"},
	{0x117, "Zhengzhou"},
	{0x118, "Xi'an"},
	{0x119, "Qingdao"},
	{0x11a, "Xining"},
	{0x11b, "Yanji"},
	{0x11c, "Hiroshima"},
	{0x11d, "Fukushima"},
	{0x11e, "Ho Chi Minh City"},
	{0x11f, "Sandakan"},
	{0x120, "Midway Island"},
	{0x121, "Buenos Aires"},
	{0x122, "Cordoba"},
	{0x123, "Brasília"},
	{0x124, "Manaus"},
	{0x125, "Panama"},
	{0x126, "Albuquerque"},
	{0x127, "Las Vegas"},
	{0x128, "Kansas City"},
	{0x129, "Boston"},
	{0x12a, "Jacksonville"},
	{0x12b, "Halifax"},
	{0x12c, "Baker Lake"},
	{0x12d, "Sudbury"},
	{0x12e, "Nairobi"},
	{0x12f, "Kampala"},
	{0x130, "Monrovia"},
	{0x131, "Seattle"},
	{0x132, "Portland"},
	{0x133, "Kuantan"},
	{0x134, "Baguio"},
	{0x13b, "Songkhla"},
	{0x13c, "Saigon"},
	{0x13d, "Hanseong"},
	{0x13e, "Changsha"},
	{0x13f, "Kursk"},
	{0x140, "Normandy"},
	{0x141, "Baku"},
	{0x142, "Islamabad"},
	{0x19a, "Kaliningrad"},
	{0x1c5, "East Berlin"},
	{0x1c6, "West Berlin"},
	{0x1c9, "Charouine"},
}

// Maps for faster lookups
var cityIDToName map[uint16]string
var cityNameToID map[string]uint16

// Initialize the lookup maps
func init() {
	cityIDToName = make(map[uint16]string)
	cityNameToID = make(map[string]uint16)

	for _, city := range cityCodes {
		cityIDToName[city.ID] = city.Name
		cityNameToID[city.Name] = city.ID
	}
}

// GetCityName returns the English city name by city ID
func GetCityName(cityID uint16) (string, bool) {
	name, exists := cityIDToName[cityID]
	return name, exists
}

// GetCityID returns the city ID by English city name (case insensitive)
func GetCityID(cityName string) (uint16, bool) {
	// Try exact match first
	if id, exists := cityNameToID[cityName]; exists {
		return id, true
	}

	// Try case insensitive match
	for name, id := range cityNameToID {
		if name == cityName {
			return id, true
		}
	}

	return 0, false
}

// GetAllCities returns all city codes
func GetAllCities() []CityInfo {
	return cityCodes
}

// GetCityCount returns the total number of cities
func GetCityCount() int {
	return len(cityCodes)
}

// CityDisplayInfo represents a city with its display information
type CityDisplayInfo struct {
	Index int
	City  CityData
	Row   int
	Col   int
	Owner byte
}

// ProcessCitiesForDisplay processes cities and returns display information
func ProcessCitiesForDisplay(saveOutput *WC4SaveOutput) []CityDisplayInfo {
	var validCities []CityDisplayInfo
	for i := 0; i < len(saveOutput.Cities); i++ {
		city := saveOutput.Cities[i]
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			continue // Skip invalid coordinates
		}
		owner := saveOutput.UnitOwnerData[row][col]
		validCities = append(validCities, CityDisplayInfo{
			Index: i,
			City:  city,
			Row:   row,
			Col:   col,
			Owner: owner,
		})
	}
	return validCities
}

// GroupCitiesByOwner groups cities by their owner
func GroupCitiesByOwner(cities []CityDisplayInfo) map[byte][]CityDisplayInfo {
	citiesByOwner := make(map[byte][]CityDisplayInfo)
	for _, cityInfo := range cities {
		citiesByOwner[cityInfo.Owner] = append(citiesByOwner[cityInfo.Owner], cityInfo)
	}
	return citiesByOwner
}
