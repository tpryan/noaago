package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/noaago"
)

func main() {
	// Initialize the NOAA client
	client := noaago.NewClient()

	fmt.Println("Searching for stations near Alameda, CA...")

	// 1. Find nearby stations within 10 miles of Alameda
	// Alameda coordinates: 37.7652, -122.2416
	stationOpts := noaago.NewStationOptionsBuilder().
		Nearby(37.7652, -122.2416, 10).
		Type(noaago.StationTypeWaterLevels).
		Build()

	stations, err := client.FindStations(stationOpts)
	if err != nil {
		log.Fatalf("Error finding stations: %v", err)
	}

	if len(stations.Stations) == 0 {
		log.Fatal("No stations found in that area.")
	}

	// Use the first station found (the closest one)
	station := stations.Stations[0]
	fmt.Printf("Found station: %s (ID: %s) at %v, %v\n\n",
		station.Name, station.ID, station.Lat, station.Lng)

	fmt.Printf("Fetching water level data for %s (last 6 hours)...\n", station.ID)

	// 2. Get Tidal Data (Water Levels) for that station for the last 6 hours
	endTime := time.Now()
	startTime := endTime.Add(-6 * time.Hour)

	tideOpts := noaago.NewTideOptionsBuilder().
		StationID(station.ID).
		DateRange(startTime, endTime).
		Product(noaago.ProductWaterLevel).
		Datum(noaago.DatumMLLW).
		Units(noaago.UnitsEnglish).      // Using English units (feet)
		TimeZone(noaago.TimeZoneLSTLDT). // Local Standard/Daylight Time
		Build()

	data, err := client.GetTides(tideOpts)
	if err != nil {
		log.Fatalf("Error fetching tide data: %v", err)
	}

	// Print the results
	fmt.Printf("%-20s | %-10s | %-5s\n", "Time", "Level (ft)", "Quality")
	fmt.Println("--------------------------------------------------")
	for _, dp := range data.Data {
		val, _ := dp.ValueFloat()
		fmt.Printf("% -20s | %-10.3f | %-5s\n", dp.Time, val, dp.Quality)
	}

	if len(data.Data) == 0 {
		fmt.Println("No data points returned for this time range.")
	}
}
