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

## Configuration Options

### StationOptionsBuilder

Used for finding stations via `client.FindStations()`.

| Method | Description |
|--------|-------------|
| `Nearby(lat, lon, radius)` | Filters stations within `radius` miles of the given coordinates. |
| `Unit(RadiusUnit)` | Sets the unit for the search radius (defaults to Miles). |
| `Type(StationType)` | Filters by station type (e.g., `noaago.StationTypeWaterLevels`). |

### TideOptionsBuilder

Used for fetching data via `client.GetTides()`.

| Method | Description |
|--------|-------------|
| `StationID(string)` | Sets the NOAA station ID (required). |
| `DateRange(start, end)` | Sets the time range for data. |
| `Product(ProductType)` | Sets the data product (e.g., Water Level, Wind). |
| `Datum(Datum)` | Sets the vertical datum (e.g., MLLW, MSL). |
| `Units(Units)` | Sets the units (Metric or English). |
| `TimeZone(TimeZone)` | Sets the time zone (GMT, LST, LST_LDT). |
| `Interval(Interval)` | Sets the data interval (e.g., High/Low, Hourly). |

## Client Configuration

The `Client` struct has exported fields that can be modified after initialization:

```go
client := noaago.NewClient()
client.UserAgent = "My-App/1.0" // Custom User-Agent header
client.HTTPClient = &http.Client{Timeout: 30 * time.Second} // Custom HTTP Client
```

### Available Constants

#### Product Types
- `ProductWaterLevel`: 6-minute water level data
- `ProductPredictions`: Tide predictions
- `ProductAirTemp`: Air temperature
- `ProductWaterTemp`: Water temperature
- `ProductWind`: Wind speed, direction, and gusts
- `ProductAirPressure`: Barometric pressure
- `ProductConductivity`: Conductivity
- `ProductVisibility`: Visibility
- `ProductHumidity`: Humidity
- `ProductSalinity`: Salinity
- `ProductHourlyHeight`: Hourly water levels
- `ProductHighLow`: High and Low water levels
- `ProductDailyMean`: Daily mean water level
- `ProductMonthlyMean`: Monthly mean water level
- `ProductOneMinuteWaterLevel`: One minute water level

#### Datums
- `DatumMLLW`: Mean Lower Low Water
- `DatumMHW`: Mean High Water
- `DatumMSL`: Mean Sea Level
- `DatumSTND`: Station Datum
- `DatumMHHW`: Mean Higher High Water
- `DatumMLW`: Mean Low Water
- `DatumNAVD`: North American Vertical Datum

#### Station Types
- `StationTypeWaterLevels`
- `StationTypeCurrents`

#### Units
- `UnitsMetric`: Celsius, Meters, etc.
- `UnitsEnglish`: Fahrenheit, Feet, etc.

#### Time Zones
- `TimeZoneGMT`: Greenwich Mean Time
- `TimeZoneLST`: Local Standard Time
- `TimeZoneLSTLDT`: Local Standard/Daylight Time

#### Intervals
- `IntervalHighLow`: High/Low tides
- `IntervalHourly`: Hourly data
- (Default): 6-minute intervals

#### Radius Units
- `RadiusUnitMiles` (Default)
- `RadiusUnitKilometers`
- `RadiusUnitNauticalMiles`

## Features

- **Fluent Interface**: Builder pattern for configuring requests.
- **Dual API Support**: Handles both Data (`datagetter`) and Metadata (`mdapi`) APIs transparently.
- **Type Safety**: Enums for Products, Datums, Units, etc.

## License

Apache 2.0


This is not an officially supported Google product. This project is not
eligible for the [Google Open Source Software Vulnerability Rewards
Program](https://bughunters.google.com/open-source-security).