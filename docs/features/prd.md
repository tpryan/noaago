The following design document outlines the architecture for the `noaago` client, mirroring the structure and patterns found in `openmeteogo` while adapting to the specific requirements of the NOAA Tides and Currents API.

### **Design Document: NOAA Tides & Currents Go Client (`noaago`)**

#### **1. Package Overview**

* **Name**: `noaago`
* **Purpose**: To provide a fluent, easy-to-use Go client for fetching tidal data and finding station metadata from NOAA.
* **Base URLs**:
* Data API: `https://api.tidesandcurrents.noaa.gov/api/prod/datagetter`
* Metadata API: `https://api.tidesandcurrents.noaa.gov/mdapi/prod/webapi`



#### **2. Core Components**

The architecture consists of three main parts:

1. **Client**: Handles HTTP communication, URL construction, and response decoding.
2. **Options**: Structs and Builders to configure requests (fluent interface).
3. **Data Models**: Structs mapping the JSON responses from NOAA.

---

### **3. Detailed Implementation Plan**

#### **File: `client.go**`

The entry point for the library. It manages two distinct base URLs (one for data, one for metadata).

```go
package noaago

import (
    "net/http"
    "time"
)

const (
    defaultDataHost     = "api.tidesandcurrents.noaa.gov"
    defaultMetadataHost = "api.tidesandcurrents.noaa.gov" // Same host, different path
    defaultUserAgent    = "NOAAGo-Client"
)

type Client struct {
    HTTPClient   *http.Client
    UserAgent    string
    dataHost     string
    metadataHost string
}

func NewClient() *Client {
    return &Client{
        HTTPClient:   http.DefaultClient,
        UserAgent:    defaultUserAgent,
        dataHost:     defaultDataHost,
        metadataHost: defaultMetadataHost,
    }
}

// GetTides fetches tidal/water level data.
func (c *Client) GetTides(o *TideOptions) (*TideResponse, error) {
    // 1. Validate Options (StationID is required)
    // 2. Build URL using c.dataUrl(o)
    // 3. Execute Request
    // 4. Decode JSON into TideResponse
}

// FindStations searches for stations based on criteria.
func (c *Client) FindStations(o *StationOptions) (*StationResponse, error) {
    // 1. Build URL using c.metadataUrl(o)
    // 2. Execute Request
    // 3. Decode JSON into StationResponse
}

// Internal helper to build Data API URL
func (c *Client) dataUrl(o *TideOptions) string {
    // Path: /api/prod/datagetter
    // Query Params: product, station, begin_date, end_date, units, time_zone, datum, format=json
}

// Internal helper to build Metadata API URL
func (c *Client) metadataUrl(o *StationOptions) string {
    // Path: /mdapi/prod/webapi/stations.json
    // Query Params: lat, lon, radius, type
}

```

#### **File: `options.go**`

Defines the parameters for requests. We need two separate option sets because the APIs are distinct.

**TideOptions (Data API)**

```go
type TideOptions struct {
    StationID string
    BeginDate time.Time
    EndDate   time.Time
    Product   ProductType // e.g., WaterLevel, Predictions
    Datum     Datum       // e.g., MLLW, MSL
    Units     Units       // e.g., English, Metric
    TimeZone  TimeZone    // e.g., GMT, LST_LDT
}

type TideOptionsBuilder struct {
    options *TideOptions
}

func NewTideOptionsBuilder() *TideOptionsBuilder { ... }
func (b *TideOptionsBuilder) StationID(id string) *TideOptionsBuilder { ... }
func (b *TideOptionsBuilder) DateRange(start, end time.Time) *TideOptionsBuilder { ... }
func (b *TideOptionsBuilder) Product(p ProductType) *TideOptionsBuilder { ... }
func (b *TideOptionsBuilder) Build() *TideOptions { ... }

```

**StationOptions (Metadata API)**

```go
type StationOptions struct {
    Latitude  float64
    Longitude float64
    Radius    float64 // In kilometers or miles (check API docs, usually miles for NOAA)
    Type      StationType // e.g., WaterLevels, Currents
}

type StationOptionsBuilder struct {
    options *StationOptions
}

func NewStationOptionsBuilder() *StationOptionsBuilder { ... }
func (b *StationOptionsBuilder) Nearby(lat, lon, radius float64) *StationOptionsBuilder { ... }
func (b *StationOptionsBuilder) Type(t StationType) *StationOptionsBuilder { ... }
func (b *StationOptionsBuilder) Build() *StationOptions { ... }

```

#### **File: `types.go**`

Enums and constants to ensure type safety, similar to `metrics.go` in the reference code.

```go
type ProductType string
const (
    ProductWaterLevel   ProductType = "water_level"
    ProductPredictions  ProductType = "predictions"
    ProductAirTemp      ProductType = "air_temperature"
    // ... other products
)

type Datum string
const (
    DatumMLLW Datum = "MLLW" // Mean Lower Low Water
    DatumMSL  Datum = "MSL"  // Mean Sea Level
    // ... other datums
)

type StationType string
const (
    StationTypeWaterLevels StationType = "waterlevels"
    StationTypeCurrents    StationType = "currents"
)

```

#### **File: `models.go**`

Structs to hold the deserialized JSON data.

```go
// TideResponse represents the Data API response
type TideResponse struct {
    Metadata TideMetadata `json:"metadata"` 
    Data     []DataPoint  `json:"data"`
}

type DataPoint struct {
    Time  string `json:"t"` // Time
    Value string `json:"v"` // Value (string in JSON, might need parsing to float)
    Quality string `json:"q,omitempty"` 
    // ... other fields (s, f, etc depending on product)
}

// StationResponse represents the Metadata API response
type StationResponse struct {
    Count    int       `json:"count"`
    Stations []Station `json:"stations"`
}

type Station struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Lat       float64 `json:"lat"`
    Lng       float64 `json:"lng"`
    State     string  `json:"state"`
    Type      string  `json:"type"`
    // ... other metadata fields
}

```

### **4. Example Usage**

```go
func main() {
    client := noaago.NewClient()

    // 1. Find nearby stations
    stationOpts := noaago.NewStationOptionsBuilder().
        Nearby(37.8086, -122.4764, 10). // San Francisco
        Type(noaago.StationTypeWaterLevels).
        Build()

    stations, _ := client.FindStations(stationOpts)
    fmt.Printf("Found station: %s (%s)\n", stations.Stations[0].Name, stations.Stations[0].ID)

    // 2. Get Tidal Data for that station
    tideOpts := noaago.NewTideOptionsBuilder().
        StationID(stations.Stations[0].ID).
        DateRange(time.Now().AddDate(0, 0, -1), time.Now()).
        Product(noaago.ProductWaterLevel).
        Datum(noaago.DatumMLLW).
        Build()

    data, _ := client.GetTides(tideOpts)
    for _, dp := range data.Data {
        fmt.Printf("Time: %s, Level: %s\n", dp.Time, dp.Value)
    }
}

```

### **5. Key Differences from `openmeteogo**`

* **Dual API Endpoints**: Unlike Open-Meteo which largely uses one base URL structure, NOAA uses a `datagetter` endpoint for values and a `mdapi` endpoint for station metadata. The `Client` struct handles this routing transparently.
* **String Values**: NOAA often returns numerical data (like water levels) as strings in JSON (e.g., `"v": "2.45"`). The `DataPoint` struct reflects this; you may want helper methods to convert these to `float64`.
* **DateFormat**: NOAA requires specific date formatting (`yyyyMMdd` or `yyyyMMdd HH:mm`). The `client.url` builder method must handle this conversion from `time.Time`.

This video covers the basics of using the NOAA Tides and Currents API which will help visualize the endpoints you are building:
[Is There An API For NOAA Tides & Currents And How Do I Use It?](https://www.youtube.com/watch?v=istmgPNcCrI)

The video explains the available endpoints and request parameters, confirming the structure needed for your `TideOptions` and `StationOptions`.