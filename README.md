# NOAA Tides & Currents Go Client (noaago)

`noaago` is a Go client for fetching tidal data and finding station metadata from NOAA.

## Installation

```bash
go get github.com/tpryan/noaago
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tpryan/noaago"
)

func main() {
	client := noaago.NewClient()

	// 1. Find nearby stations (e.g., San Francisco)
	stationOpts := noaago.NewStationOptionsBuilder().
		Nearby(37.8086, -122.4764, 10). // 10 mile radius
		Type(noaago.StationTypeWaterLevels).
		Build()

	stations, err := client.FindStations(stationOpts)
	if err != nil {
		log.Fatal(err)
	}

	if len(stations.Stations) == 0 {
		log.Fatal("No stations found")
	}

	station := stations.Stations[0]
	fmt.Printf("Found station: %s (%s)\n", station.Name, station.ID)

	// 2. Get Tidal Data for that station (last 24 hours)
	tideOpts := noaago.NewTideOptionsBuilder().
		StationID(station.ID).
		DateRange(time.Now().Add(-24*time.Hour), time.Now()).
		Product(noaago.ProductWaterLevel).
		Datum(noaago.DatumMLLW).
		Units(noaago.UnitsMetric).
		TimeZone(noaago.TimeZoneGMT).
		Build()

	data, err := client.GetTides(tideOpts)
	if err != nil {
		log.Fatal(err)
	}

	for _, dp := range data.Data {
		val, _ := dp.ValueFloat()
		fmt.Printf("Time: %s, Level: %.3f\n", dp.Time, val)
	}
}
```

## Features

- **Fluent Interface**: Builder pattern for configuring requests.
- **Dual API Support**: Handles both Data (`datagetter`) and Metadata (`mdapi`) APIs transparently.
- **Type Safety**: Enums for Products, Datums, Units, etc.

## License

MIT
