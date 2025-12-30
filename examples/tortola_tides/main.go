package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/noaago"
)

func main() {
	client := noaago.NewClient()

	fmt.Println("Finding a station near Tortola, BVI...")

	// Tortola coordinates: 18.4167, -64.5833
	lat, lon := 18.4167, -64.5833

	// 1. Find a station
	// Using a larger radius as BVI might be served by nearby stations
	stationOpts := noaago.NewStationOptionsBuilder().
		Nearby(lat, lon, 100). 
		Type(noaago.StationTypeWaterLevels).
		Build()

	stations, err := client.FindStations(stationOpts)
	if err != nil {
		log.Fatalf("Error finding stations: %v", err)
	}

	if len(stations.Stations) == 0 {
		log.Fatal("No stations found near Tortola.")
	}

	station := stations.Stations[0]
	fmt.Printf("Using station: %s (%s)\n", station.Name, station.ID)

	// 2. Set dates for May 1 to May 4, 2026
	startDate := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.May, 4, 23, 59, 0, 0, time.UTC)

	fmt.Printf("Fetching high/low tide predictions for %s to %s...\n",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// 3. Get Predictions
	tideOpts := noaago.NewTideOptionsBuilder().
		StationID(station.ID).
		DateRange(startDate, endDate).
		Product(noaago.ProductPredictions).
		Datum(noaago.DatumMLLW).
		Units(noaago.UnitsEnglish).
		TimeZone(noaago.TimeZoneLST). // Use Local Standard Time
		Interval(noaago.IntervalHighLow).
		Build()

	data, err := client.GetTides(tideOpts)
	if err != nil {
		log.Fatalf("Error fetching tides: %v", err)
	}

	fmt.Printf("\nPredictions for %s:\n", station.Name)
	fmt.Printf("% -20s | % -10s | % -5s\n", "Time (Local)", "Level (ft)", "Type")
	fmt.Println("--------------------------------------------------")

	dataset := data.GetData()
	for _, dp := range dataset {
		val, _ := dp.ValueFloat()
		fmt.Printf("% -20s | % -10.3f | % -5s\n", dp.Time, val, dp.Type)
	}

	if len(dataset) == 0 {
		fmt.Println("No data points returned for this time range.")
	}
}
