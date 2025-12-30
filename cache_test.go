package noaago

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFindStations_Cache(t *testing.T) {
	requestCount := 0
	mockResponse := `{
		"count": 1,
		"stations": [
			{"id": "9414290", "name": "San Francisco", "type": "waterlevels"}
		]
	}`

	server, client := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResponse)
	})
	defer server.Close()

	opts := NewStationOptionsBuilder().Type(StationTypeWaterLevels).Build()

	// 1. First call: should hit the server
	_, err := client.FindStations(opts)
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 request, got %d", requestCount)
	}

	// 2. Second call: should use cache
	_, err = client.FindStations(opts)
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 request (cached), got %d", requestCount)
	}

	// 3. Different type: should hit the server
	optsCurrents := NewStationOptionsBuilder().Type(StationTypeCurrents).Build()
	_, err = client.FindStations(optsCurrents)
	if err != nil {
		t.Fatalf("Third call failed: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("Expected 2 requests, got %d", requestCount)
	}
}
