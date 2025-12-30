package noaago

import (
	"net/url"
	"testing"
	"time"
)

func TestDataUrl(t *testing.T) {
	client := NewClient()

	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	opts := NewTideOptionsBuilder().
		StationID("9414290").
		DateRange(start, end).
		Product(ProductWaterLevel).
		Datum(DatumMLLW).
		Units(UnitsMetric).
		TimeZone(TimeZoneGMT).
		Build()

	uStr := client.dataUrl(opts)
	u, err := url.Parse(uStr)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	q := u.Query()
	if q.Get("station") != "9414290" {
		t.Errorf("Expected station 9414290, got %s", q.Get("station"))
	}
	if q.Get("product") != string(ProductWaterLevel) {
		t.Errorf("Expected product %s, got %s", ProductWaterLevel, q.Get("product"))
	}
	if q.Get("datum") != string(DatumMLLW) {
		t.Errorf("Expected datum %s, got %s", DatumMLLW, q.Get("datum"))
	}
	if q.Get("units") != string(UnitsMetric) {
		t.Errorf("Expected units %s, got %s", UnitsMetric, q.Get("units"))
	}
	if q.Get("time_zone") != string(TimeZoneGMT) {
		t.Errorf("Expected time_zone %s, got %s", TimeZoneGMT, q.Get("time_zone"))
	}
	if q.Get("format") != "json" {
		t.Errorf("Expected format json, got %s", q.Get("format"))
	}
	// Check dates
	if q.Get("begin_date") != "20230101 00:00" {
		t.Errorf("Expected begin_date 20230101 00:00, got %s", q.Get("begin_date"))
	}
	if q.Get("end_date") != "20230102 00:00" {
		t.Errorf("Expected end_date 20230102 00:00, got %s", q.Get("end_date"))
	}
}

func TestMetadataUrl(t *testing.T) {
	client := NewClient()

	opts := NewStationOptionsBuilder().
		Nearby(37.8086, -122.4764, 10).
		Type(StationTypeWaterLevels).
		Build()

	uStr := client.metadataUrl(opts)
	u, err := url.Parse(uStr)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	q := u.Query()
	// lat, lon, radius are now handled client-side, so they shouldn't be in the URL
	if q.Get("lat") != "" {
		t.Errorf("Expected lat to be empty, got %s", q.Get("lat"))
	}
	if q.Get("lon") != "" {
		t.Errorf("Expected lon to be empty, got %s", q.Get("lon"))
	}
	if q.Get("radius") != "" {
		t.Errorf("Expected radius to be empty, got %s", q.Get("radius"))
	}
	if q.Get("type") != string(StationTypeWaterLevels) {
		t.Errorf("Expected type %s, got %s", StationTypeWaterLevels, q.Get("type"))
	}
}
