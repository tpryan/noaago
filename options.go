package noaago

import (
	"time"
)

// TideOptions defines the parameters for data requests
type TideOptions struct {
	StationID string
	BeginDate time.Time
	EndDate   time.Time
	Product   ProductType
	Datum     Datum
	Units     Units
	TimeZone  TimeZone
	Format    string
}

// TideOptionsBuilder is a builder for TideOptions
type TideOptionsBuilder struct {
	options *TideOptions
}

// NewTideOptionsBuilder creates a new instance of TideOptionsBuilder
func NewTideOptionsBuilder() *TideOptionsBuilder {
	return &TideOptionsBuilder{
		options: &TideOptions{
			Units:    UnitsMetric, // Default
			TimeZone: TimeZoneGMT, // Default
			Format:   "json",      // Default
		},
	}
}

// StationID sets the station ID
func (b *TideOptionsBuilder) StationID(id string) *TideOptionsBuilder {
	b.options.StationID = id
	return b
}

// DateRange sets the begin and end dates
func (b *TideOptionsBuilder) DateRange(start, end time.Time) *TideOptionsBuilder {
	b.options.BeginDate = start
	b.options.EndDate = end
	return b
}

// Product sets the product type
func (b *TideOptionsBuilder) Product(p ProductType) *TideOptionsBuilder {
	b.options.Product = p
	return b
}

// Datum sets the datum
func (b *TideOptionsBuilder) Datum(d Datum) *TideOptionsBuilder {
	b.options.Datum = d
	return b
}

// Units sets the units
func (b *TideOptionsBuilder) Units(u Units) *TideOptionsBuilder {
	b.options.Units = u
	return b
}

// TimeZone sets the time zone
func (b *TideOptionsBuilder) TimeZone(tz TimeZone) *TideOptionsBuilder {
	b.options.TimeZone = tz
	return b
}

// Build returns the constructed TideOptions
func (b *TideOptionsBuilder) Build() *TideOptions {
	return b.options
}

// StationOptions defines the parameters for metadata requests
type StationOptions struct {
	Latitude  float64
	Longitude float64
	Radius    float64
	Type      StationType
}

// StationOptionsBuilder is a builder for StationOptions
type StationOptionsBuilder struct {
	options *StationOptions
}

// NewStationOptionsBuilder creates a new instance of StationOptionsBuilder
func NewStationOptionsBuilder() *StationOptionsBuilder {
	return &StationOptionsBuilder{
		options: &StationOptions{
			Radius: 10, // Default radius might be good, or let user set it.
		},
	}
}

// Nearby sets the location and search radius
func (b *StationOptionsBuilder) Nearby(lat, lon, radius float64) *StationOptionsBuilder {
	b.options.Latitude = lat
	b.options.Longitude = lon
	b.options.Radius = radius
	return b
}

// Type sets the station type
func (b *StationOptionsBuilder) Type(t StationType) *StationOptionsBuilder {
	b.options.Type = t
	return b
}

// Build returns the constructed StationOptions
func (b *StationOptionsBuilder) Build() *StationOptions {
	return b.options
}
