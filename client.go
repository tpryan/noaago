package noaago

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

const (
	defaultDataHost     = "api.tidesandcurrents.noaa.gov"
	defaultMetadataHost = "api.tidesandcurrents.noaa.gov"
	defaultUserAgent    = "NOAAGo-Client"

	dataPath     = "/api/prod/datagetter"
	metadataPath = "/mdapi/prod/webapi/stations.json"
)

// Client handles communication with the NOAA API
type Client struct {
	HTTPClient   *http.Client
	UserAgent    string
	dataHost     string
	metadataHost string

	// Cache for station lists, keyed by station type (e.g. "waterlevels")
	stationCache map[string]*StationResponse
	cacheMutex   sync.RWMutex
}

// NewClient returns a new Client with default settings
func NewClient() *Client {
	return &Client{
		HTTPClient:   http.DefaultClient,
		UserAgent:    defaultUserAgent,
		dataHost:     defaultDataHost,
		metadataHost: defaultMetadataHost,
		stationCache: make(map[string]*StationResponse),
	}
}

// GetTides fetches tidal/water level data
func (c *Client) GetTides(o *TideOptions) (*TideResponse, error) {
	if o.StationID == "" {
		return nil, fmt.Errorf("station ID is required")
	}

	u := c.dataUrl(o)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var tideResp TideResponse
	if err := json.NewDecoder(resp.Body).Decode(&tideResp); err != nil {
		return nil, err
	}

	if tideResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", tideResp.Error.Message)
	}

	return &tideResp, nil
}

// FindStations searches for stations based on criteria
func (c *Client) FindStations(o *StationOptions) (*StationResponse, error) {
	// Check cache first
	cacheKey := string(o.Type)

	c.cacheMutex.RLock()
	cachedResp, found := c.stationCache[cacheKey]
	c.cacheMutex.RUnlock()

	var stationResp *StationResponse

	if found {
		// Create a copy to avoid modifying the cached data when filtering
		// Shallow copy of the struct is fine, but we'll create a new StationResponse
		// holding the same stations slice for now, but since we filter below by creating *new* slice,
		// it should be safe.
		stationResp = &StationResponse{
			Count:    cachedResp.Count,
			Stations: cachedResp.Stations,
		}
	} else {
		// Fetch from API
		u := c.metadataUrl(o)

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", c.UserAgent)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		}

		stationResp = &StationResponse{}
		if err := json.NewDecoder(resp.Body).Decode(stationResp); err != nil {
			return nil, err
		}

		// Update cache
		c.cacheMutex.Lock()
		c.stationCache[cacheKey] = stationResp
		c.cacheMutex.Unlock()
	}

	// Filter by distance if lat/lon/radius are provided
	if o.Latitude != 0 && o.Longitude != 0 && o.Radius > 0 {
		var filtered []Station
		for _, s := range stationResp.Stations {
			dist := haversine(o.Latitude, o.Longitude, s.Lat, s.Lng)
			if dist <= o.Radius {
				filtered = append(filtered, s)
			}
		}

		// Sort by distance
		sort.Slice(filtered, func(i, j int) bool {
			d1 := haversine(o.Latitude, o.Longitude, filtered[i].Lat, filtered[i].Lng)
			d2 := haversine(o.Latitude, o.Longitude, filtered[j].Lat, filtered[j].Lng)
			return d1 < d2
		})

		// Return a new response object with filtered results
		return &StationResponse{
			Count:    len(filtered),
			Stations: filtered,
		}, nil
	}

	return stationResp, nil
}

// dataUrl builds the URL for the Data API
func (c *Client) dataUrl(o *TideOptions) string {
	u := url.URL{
		Scheme: "https",
		Host:   c.dataHost,
		Path:   dataPath,
	}

	q := u.Query()
	q.Set("station", o.StationID)

	// Format dates. API accepts yyyyMMdd or yyyyMMdd HH:mm
	// We'll use yyyyMMdd HH:mm for better precision if hours/mins are non-zero,
	// but mostly yyyyMMdd is enough for daily ranges.
	// Let's use 20060102 15:04
	const layout = "20060102 15:04"
	q.Set("begin_date", o.BeginDate.Format(layout))
	q.Set("end_date", o.EndDate.Format(layout))

	q.Set("product", string(o.Product))

	if o.Datum != "" {
		q.Set("datum", string(o.Datum))
	}

	if o.Units != "" {
		q.Set("units", string(o.Units))
	}

	if o.TimeZone != "" {
		q.Set("time_zone", string(o.TimeZone))
	}

	q.Set("format", "json")
	q.Set("application", c.UserAgent) // Good practice to identify app

	u.RawQuery = q.Encode()
	return u.String()
}

// metadataUrl builds the URL for the Metadata API
func (c *Client) metadataUrl(o *StationOptions) string {
	u := url.URL{
		Scheme: "https",
		Host:   c.metadataHost,
		Path:   metadataPath,
	}

	q := u.Query()
	// Note: We do NOT send lat/lon/radius to the API because the API
	// does not support searching by location for the list endpoint.
	// We filter client-side in FindStations.

	if o.Type != "" {
		q.Set("type", string(o.Type))
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// haversine calculates the distance between two points in miles
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 3958.8 // Earth radius in miles

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
