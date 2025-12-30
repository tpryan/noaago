package noaago

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetTides_Success(t *testing.T) {
	mockResponse := `{
		"metadata": {
			"id": "9414290",
			"name": "San Francisco",
			"lat": "37.8063",
			"lon": "-122.4659"
		},
		"data": [
			{"t": "2023-01-01 00:00", "v": "2.5", "q": "p", "type": "H"}
		]
	}`

	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prod/datagetter" {
			t.Errorf("Expected path /api/prod/datagetter, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("station") != "9414290" {
			t.Errorf("Expected station 9414290, got %s", r.URL.Query().Get("station"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResponse)
	})
	defer server.Close()

	opts := NewTideOptionsBuilder().
		StationID("9414290").
		DateRange(time.Now(), time.Now()).
		Build()

	resp, err := client.GetTides(opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Metadata.Name != "San Francisco" {
		t.Errorf("Expected station name San Francisco, got %s", resp.Metadata.Name)
	}

	dataset := resp.GetData()
	if len(dataset) != 1 {
		t.Errorf("Expected 1 data point, got %d", len(dataset))
	}

	val, err := dataset[0].ValueFloat()
	if err != nil {
		t.Errorf("ValueFloat failed: %v", err)
	}
	if val != 2.5 {
		t.Errorf("Expected value 2.5, got %f", val)
	}
}
func TestGetTides_Error(t *testing.T) {
	mockResponse := `{
		"error": {
			"message": "No data found"
		}
	}`

	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResponse)
	})
	defer server.Close()

	opts := NewTideOptionsBuilder().StationID("invalid").Build()
	_, err := client.GetTides(opts)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "API error: No data found" {
		t.Errorf("Expected error message 'API error: No data found', got '%s'", err.Error())
	}
}

func TestFindStations_Success(t *testing.T) {
	mockResponse := `{
		"count": 2,
		"stations": [
			{
				"id": "9414290",
				"name": "San Francisco",
				"lat": 37.8063,
				"lng": -122.4659,
				"type": "waterlevels"
			},
			{
				"id": "1234567",
				"name": "Far Away",
				"lat": 0.0,
				"lng": 0.0,
				"type": "waterlevels"
			}
		]
	}`

	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mdapi/prod/webapi/stations.json" {
			t.Errorf("Expected path /mdapi/prod/webapi/stations.json, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResponse)
	})
	defer server.Close()

	// Test filtering: San Francisco is near (37.8, -122.4), Far Away is not.
	opts := NewStationOptionsBuilder().
		Nearby(37.8, -122.4, 10). // 10 miles radius
		Type(StationTypeWaterLevels).
		Build()

	resp, err := client.FindStations(opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Stations) != 1 {
		t.Errorf("Expected 1 station after filtering, got %d", len(resp.Stations))
	}
	if resp.Stations[0].ID != "9414290" {
		t.Errorf("Expected station ID 9414290, got %s", resp.Stations[0].ID)
	}
}

func TestFindStations_NoFilter(t *testing.T) {
	mockResponse := `{
		"count": 1,
		"stations": [
			{"id": "9414290", "name": "San Francisco"}
		]
	}`

	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockResponse)
	})
	defer server.Close()

	// No location provided
	opts := NewStationOptionsBuilder().Build()

	resp, err := client.FindStations(opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Stations) != 1 {
		t.Errorf("Expected 1 station, got %d", len(resp.Stations))
	}
}

func TestHaversine(t *testing.T) {
	// Distance between SF (37.7749, -122.4194) and NYC (40.7128, -74.0060)
	// roughly 2572 miles
	dist := haversine(37.7749, -122.4194, 40.7128, -74.0060)
	if dist < 2500 || dist > 2650 {
		t.Errorf("Expected distance around 2572 miles, got %f", dist)
	}

	// Distance to self should be 0
	dist = haversine(10, 10, 10, 10)
	if dist != 0 {
		t.Errorf("Expected distance 0, got %f", dist)
	}
}

func TestGetTides_NetworkError(t *testing.T) {
	// Create a client that will fail to connect (using a closed server port or invalid host)
	client := NewClient()
	client.dataHost = "invalid-host"

	opts := NewTideOptionsBuilder().StationID("123").Build()
	_, err := client.GetTides(opts)
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}
}

func TestGetTides_InvalidJSON(t *testing.T) {
	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{invalid-json")
	})
	defer server.Close()

	opts := NewTideOptionsBuilder().StationID("123").Build()
	_, err := client.GetTides(opts)
	if err == nil {
		t.Fatal("Expected JSON error, got nil")
	}
}

func TestGetTides_Non200(t *testing.T) {
	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	opts := NewTideOptionsBuilder().StationID("123").Build()
	_, err := client.GetTides(opts)
	if err == nil {
		t.Fatal("Expected status error, got nil")
	}
	if err.Error() != "API returned non-200 status code: 500" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFindStations_NetworkError(t *testing.T) {
	client := NewClient()
	client.metadataHost = "invalid-host"

	opts := NewStationOptionsBuilder().Build()
	_, err := client.FindStations(opts)
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}
}

func TestFindStations_InvalidJSON(t *testing.T) {
	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{invalid-json")
	})
	defer server.Close()

	opts := NewStationOptionsBuilder().Build()
	_, err := client.FindStations(opts)
	if err == nil {
		t.Fatal("Expected JSON error, got nil")
	}
}

func TestFindStations_Non200(t *testing.T) {
	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := NewStationOptionsBuilder().Build()
	_, err := client.FindStations(opts)
	if err == nil {
		t.Fatal("Expected status error, got nil")
	}
	if err.Error() != "API returned non-200 status code: 404" {
		t.Errorf("Unexpected error: %v", err)
	}
}
