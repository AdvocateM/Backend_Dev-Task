package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

type Spot struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var spots = []Spot{
	{
		Name:      "Spot 1",
		Latitude:  37.7749,
		Longitude: -122.4194,
	},
	{
		Name:      "Spot 2",
		Latitude:  37.7833,
		Longitude: -122.4167,
	},
	{
		Name:      "Spot 3",
		Latitude:  37.7936,
		Longitude: -122.3987,
	},

}

func main() {
	http.HandleFunc("/spots", getSpotsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getSpotsHandler(w http.ResponseWriter, r *http.Request) {
	latitudeStr := r.URL.Query().Get("latitude")
	longitudeStr := r.URL.Query().Get("longitude")
	radiusStr := r.URL.Query().Get("radius")

	latitude, err := parseFloat64(latitudeStr)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := parseFloat64(longitudeStr)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius, err := parseFloat64(radiusStr)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	spotsInArea := filterSpotsInArea(latitude, longitude, radius)

	response, err := json.Marshal(spotsInArea)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func filterSpotsInArea(latitude, longitude, radius float64) []Spot {
	var spotsInArea []Spot

	for _, spot := range spots {
		if isSpotInArea(latitude, longitude, radius, spot) {
			spotsInArea = append(spotsInArea, spot)
		}
	}

	return spotsInArea
}

func isSpotInArea(latitude, longitude, radius float64, spot Spot) bool {
	lat1 := toRadians(latitude)
	lon1 := toRadians(longitude)
	lat2 := toRadians(spot.Latitude)
	lon2 := toRadians(spot.Longitude)

	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := c * 6371000 // Earth's radius

	return distance <= radius
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func parseFloat64(value string) (float64, error) {
	if value == "" {
		return 0, fmt.Errorf("value is empty")
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float64: %w", err)
	}

	return result, nil
}
