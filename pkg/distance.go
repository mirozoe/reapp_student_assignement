package reapp_students_assignement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// country is appended to all city names what are used for constructing distanceMatrix
const country = "Czechia"

// NewDistanceMatrix is ctor
func NewDistanceMatrix(cities []string, dmi DistanceInterface) DistanceMatrix {
	if len(cities) <= 1 {
		log.Error().Msgf("Cannot have only one city defined: %+v", cities)
		return DistanceMatrix{}
	}
	dm := DistanceMatrix{
		cities:            cities,
		DistanceInterface: dmi,
	}
	dm.distanceMatrix = dm.initiateDistanceMatrix()
	log.Debug().Msgf("Complete distanceMatrix %+v", dm.distanceMatrix)
	return dm
}

// GetMatrixFor get distance matrix for specific cities
func (dm DistanceMatrix) GetMatrixFor(from, to []string) (map[string]Distance, error) {
	var distanceMatrix map[string]Distance
	gmResponse, err := parseGMapJson(dm.GetDistance(from, to))
	if err != nil {
		return distanceMatrix, err
	}
	distanceMatrix, err = prepareDistanceMatrix(gmResponse)
	if err != nil {
		return distanceMatrix, err
	}

	return distanceMatrix, nil
}

// GetDistanceFromDM returns Distance between two cities
func (dm DistanceMatrix) GetDistanceFromDM(from, to string) Distance {
	if from == to {
		return Distance{
			DistanceM: 1,
			DurationS: 1,
		}
	}
	return dm.findCitiesKey(from, to, 0)
}

// CalculateLogDistance calculates fitness for distance between two cities
// Hint: It provides log2 of distance only, fitness is calculated in ABC further
func (dm DistanceMatrix) CalculateLogDistance(attr ABCAttributes, cityW, cityS string) float64 {
	distance := dm.GetDistanceFromDM(cityW, cityS)
	if distance.DistanceM > attr.ReferenceMaxDistance || distance.DistanceM == 0 {
		log.Debug().Msgf("MaxRefDistance between %s and %s; %d", cityW, cityS, distance.DistanceM)
		return float64(attr.ReferenceMaxDistance)
	}
	return math.Log2(float64(distance.DistanceM))
}

// initiateDistanceMatrix iterate through all cities and create mash and store it as distanceMatrix
func (dm DistanceMatrix) initiateDistanceMatrix() map[string]Distance {
	distanceMatrix := make(map[string]Distance)

	citiesWithCountry := []string{}
	for _, c := range dm.cities {
		citiesWithCountry = append(citiesWithCountry, appendCountry(c, country))
	}

	for i, city := range citiesWithCountry {
		if i+1 == len(citiesWithCountry) {
			return distanceMatrix
		}

		newDM, err := dm.GetMatrixFor([]string{city}, citiesWithCountry[i+1:])
		if err != nil {
			log.Error().Err(err).Msgf("Issue with getting distances from %s to %+v", city, citiesWithCountry[i+1:])
		}
		distanceMatrix = appendToDistanceMatrix(distanceMatrix, newDM)
	}
	return distanceMatrix
}

func (dm DistanceMatrix) findCitiesKey(from, to string, run int) Distance {
	if run > 1 {
		log.Warn().Msgf("Didn't find distance between %s and %s", from, to)
		return Distance{}
	}
	if val, ok := dm.distanceMatrix[stringToHex(prepareCitiesString([]string{from, to}, false))]; ok {
		return val
	}
	run += 1
	return dm.findCitiesKey(to, from, run)
}

// GetDistance asks GMaps Distance API and returns response
func (g GMapGetter) GetDistance(from, to []string) []byte {
	originCities := prepareCitiesString(from, true)
	destinationCities := prepareCitiesString(to, true)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=KEY_PLACEHOLDER", originCities, destinationCities)
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get response from GMaps with call: %s", url)
		return []byte{}
	}
	defer resp.Body.Close()

	var completeResp []byte
	for {
		buff := make([]byte, 1024)
		if _, err := io.ReadFull(resp.Body, buff); err == nil {
			completeResp = append(completeResp, buff...)
		} else if err == io.EOF {
			break
			// error is returned when read is smaller than len(buff)
		} else {
			buff = bytes.Trim(buff, "\x00")
			completeResp = append(completeResp, buff...)
			break
		}

	}
	return completeResp
}

func prepareCitiesString(cities []string, replaceWhitespace bool) string {
	var updatedCities []string
	for _, city := range cities {
		if replaceWhitespace {
			updatedCities = append(updatedCities, replaceWhitespaces(city))
		} else {
			updatedCities = append(updatedCities, city)
		}
	}
	return strings.Join(updatedCities, "|")
}

func parseGMapJson(jsonData []byte) (GMapResponse, error) {
	var parsedGMapResponse GMapResponse
	if err := json.Unmarshal(jsonData, &parsedGMapResponse); err != nil {
		log.Error().Err(err).Msgf("Cannot unmarshall JSON data %s", string(jsonData))
	}
	if wrong := validateResponse(parsedGMapResponse); len(wrong) > 0 {
		mapStringSlice(wrong, func(s string) bool {
			log.Warn().Msgf("This city doesn't properly fits that it is in Czechia %s", s)
			return true
		})
		return parsedGMapResponse, ErrorFactory(CityValidationError, strings.Join(wrong, " "))
	}
	return parsedGMapResponse, nil
}

func validateResponse(resp GMapResponse) []string {
	dCities := mapStringSlice(resp.DestinationAddresses, func(s string) bool {
		m, err := regexp.MatchString(`.*Czechia`, s)
		if err != nil {
			log.Error().Err(err).Msgf("Error validating name, doesn't contain Czechia: %s", s)
			return false
		}
		return m
	})
	oCities := mapStringSlice(resp.OriginAddreses, func(s string) bool {
		m, err := regexp.MatchString(`.*Czechia`, s)
		if err != nil {
			log.Error().Err(err).Msgf("Error validating name, doesn't contain Czechia: %s", s)
			return false
		}
		return m
	})
	if len(dCities) > 0 {
		return dCities
	} else if len(oCities) > 0 {
		return oCities
	}
	return []string{}
}

func prepareDistanceMatrix(dm GMapResponse) (map[string]Distance, error) {
	matrix := make(map[string]Distance)
	for i, ocity := range dm.OriginAddreses {
		for j, dcity := range dm.DestinationAddresses {
			if dm.Rows[i].Elements[j].Status != "OK" {
				return matrix, ErrorFactory(CityValidationError, dcity)
			}
			key := stringToHex(prepareMatrixKey(stripString(ocity, GMapCitySeparator), stripString(dcity, GMapCitySeparator)))
			matrix[key] = Distance{dm.Rows[i].Elements[j].Distance.Value, dm.Rows[i].Elements[j].Duration.Value}
		}
	}
	return matrix, nil
}

// appendCountry adds country to city name (Most, Czechia) as it is better for GoogleMapsAPI
func appendCountry(city, country string) string {
	return fmt.Sprintf("%s, %s", city, country)
}
