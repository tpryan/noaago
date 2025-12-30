package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/noaago"
)

func main() {
	client := noaago.NewClient()

	fmt.Println("Finding a station near Alameda, CA...")

	// Alameda coordinates
	lat, lon := 37.7652, -122.2416

	// 1. Find a station
	stationOpts := noaago.NewStationOptionsBuilder().
		Nearby(lat, lon, 10).
		Type(noaago.StationTypeWaterLevels).
		Build()

	stations, err := client.FindStations(stationOpts)
	if err != nil {
		log.Fatalf("Error finding stations: %v", err)
	}

	if len(stations.Stations) == 0 {
		log.Fatal("No stations found.")
	}

	station := stations.Stations[0]
	fmt.Printf("Using station: %s (%s)\n", station.Name, station.ID)

	// 2. Calculate dates for 1 month from now
	targetDate := time.Now().AddDate(0, 1, 0) // 1 month from now
	startDate := targetDate
	endDate := targetDate.AddDate(0, 0, 5) // 5 days of predictions

	fmt.Printf("Fetching high/low tide predictions for the period of %s to %s...\n", 
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// 3. Get Predictions
	tideOpts := noaago.NewTideOptionsBuilder().
		StationID(station.ID).
		DateRange(startDate, endDate).
		Product(noaago.ProductPredictions).
		Datum(noaago.DatumMLLW).
		Units(noaago.UnitsEnglish).
		TimeZone(noaago.TimeZoneLSTLDT).
		Interval(noaago.IntervalHighLow).
		Build()

	// Note: We need to specify "High/Low" interval if the client supports it,
	// or filter the predictions ourselves if the API returns all intervals.
	// The standard predictions product returns interval data (e.g. every 6 mins) unless specified.
	// However, NOAA has a specific product "high_low" for observed data,
	// but for PREDICTIONS, the interval parameter controls this.
	// Our current TideOptions doesn't have an "Interval" field.
	// By default, 'predictions' returns 6-minute intervals.
	// To get High/Low predictions specifically, we usually add `&interval=hilo`.

	// Let's quickly check if we need to update the client to support 'interval'.
	// Yes, to get just high/lows from the predictions product, we need that param.
	// But let's see what happens with just ProductPredictions first, or if we can use ProductHighLow.
	// ProductHighLow is usually for verified *observed* highs and lows.

	// For now, let's try to fetch the predictions and see.
	// If we get too much data, we might need to add Interval support to the client.

	data, err := client.GetTides(tideOpts)
	if err != nil {
		log.Fatalf("Error fetching tides: %v", err)
	}

	fmt.Printf("\nPredictions for %s:\n", station.Name)
	fmt.Printf("%-20s | %-10s | %-5s\n", "Time", "Level (ft)", "Type")
	fmt.Println("--------------------------------------------------")

	count := 0
	dataset := data.GetData()
	for _, dp := range dataset {
		// If we receive all 6-minute data, this list will be huge.
		// If we receive High/Low, the 'Type' field (H or L) should be populated.

		// If the client doesn't support sending 'interval=hilo', we'll get 6-minute data.
		// Let's filter client-side if necessary or just show a sample.
		// But ideally, we should update the client.

		val, _ := dp.ValueFloat()

		// Use a simple check to see if we got H/L data or raw series
		// The 'Type' field in DataPoint is populated for H/L data.
		tType := dp.Type
		if tType == "" {
			tType = "-"
		}

		fmt.Printf("% -20s | % -10.3f | % -5s\n", dp.Time, val, tType)

		count++
		if count >= 20 {
			fmt.Println("... (truncating output) ...")
			break
		}
	}
}
