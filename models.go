package noaago

import "strconv"

// TideResponse represents the Data API response
type TideResponse struct {
	Metadata    TideMetadata `json:"metadata"`
	Data        []DataPoint  `json:"data"`
	Predictions []DataPoint  `json:"predictions"`
	Error       *APIError    `json:"error,omitempty"`
}

// GetData returns the data points regardless of whether they came from "data" or "predictions"
func (r *TideResponse) GetData() []DataPoint {
	if len(r.Predictions) > 0 {
		return r.Predictions
	}
	return r.Data
}

// TideMetadata contains metadata about the station in the data response
type TideMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

// DataPoint represents a single data point from the API
type DataPoint struct {
	Time    string `json:"t"`              // Time
	Value   string `json:"v"`              // Value
	Quality string `json:"q,omitempty"`    // Quality
	Sigma   string `json:"s,omitempty"`    // Sigma (standard deviation)
	Flags   string `json:"f,omitempty"`    // Flags
	Type    string `json:"type,omitempty"` // Type (e.g. H/L for High/Low)
}

// ValueFloat returns the value as a float64
func (d *DataPoint) ValueFloat() (float64, error) {
	return strconv.ParseFloat(d.Value, 64)
}

// StationResponse represents the Metadata API response
type StationResponse struct {
	Count    int       `json:"count"`
	Stations []Station `json:"stations"`
}

// Station represents a single station's metadata
type Station struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	State       string  `json:"state"`
	Type        string  `json:"type"`
	Affiliation string  `json:"affiliation"`
}

// APIError represents an error returned by the NOAA API
type APIError struct {
	Message string `json:"message"`
}
