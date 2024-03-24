package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

// These tests use the `weather_stations.csv` file, used to generate the rows in the `measurements.txt` file.

func Test_gomaphash(t *testing.T) {
	f, err := os.Open("weather_stations.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	type nameCount struct {
		name  string
		count int
	}

	mapStations := make(map[uint64]nameCount)

	bufR := bufio.NewReaderSize(f, 16_000)
	for {
		line, _, err := bufR.ReadLine()
		if err != nil {
			break
		}
		// skip comments
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}

		station := bytes.SplitN(line, []byte(";"), 1)

		// compute the hash of each station name and ensure we have 44689 unique stations in the map
		stHash := gomaphash(station[0])
		foundSt, ok := mapStations[stHash]
		if !ok {
			foundSt = nameCount{name: string(station[0]), count: 1}
			mapStations[stHash] = foundSt
		} else {
			foundSt.count++
			mapStations[stHash] = foundSt
		}
	}

	const totalStations = 44689
	if len(mapStations) != totalStations {
		t.Errorf("expected %d stations, but got %d", totalStations, len(mapStations))
	}

	for _, nCnt := range mapStations {
		if len(nCnt.name) != 0 && nCnt.count == 0 {
			t.Errorf("expected a non zero value but got zero for station %s", nCnt.name)
		}
	}
}

// 5,083 ns/op
func Benchmark_gomaphash(b *testing.B) {
	data := []byte{'s', 't', 'a', 't', 'i', 'o', 'n', ' ', 'n', 'a', 'm', 'e'}
	for i := 0; i < b.N; i++ {
		gomaphash(data)
	}
}

// 2,520 ns/op
// while this is faster, using the station name as a map key is slower than using the hash...
func Benchmark_bytesToString(b *testing.B) {
	data := []byte{'s', 't', 'a', 't', 'i', 'o', 'n', ' ', 'n', 'a', 'm', 'e'}
	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

func Test_parseStationLine(t *testing.T) {
	tests := []struct {
		name            string
		line            []byte
		wantMeasurement int16
		wantStation     []byte
	}{
		{
			name:            "Test case 1",
			line:            []byte("Station1;12.3"),
			wantMeasurement: 123,
			wantStation:     []byte("Station1"),
		},
		{
			name:            "Test case 2",
			line:            []byte("Station2;-45.6"),
			wantMeasurement: -456,
			wantStation:     []byte("Station2"),
		},
		{
			name:            "Test case 3",
			line:            []byte("Station3;0"),
			wantMeasurement: 0,
			wantStation:     []byte("Station3"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotMeasurement, gotStation := parseStationLine(tt.line)
			if gotMeasurement != tt.wantMeasurement {
				t.Errorf("parseStationLine() gotMeasurement = %v, want %v", gotMeasurement, tt.wantMeasurement)
			}
			if !bytes.Equal(gotStation, tt.wantStation) {
				t.Errorf("parseStationLine() gotStation = %v, want %v", gotStation, tt.wantStation)
			}
		})
	}
}

// 33,24 ns/op using gomaphash(stationName)
func Benchmark_parseStationLine(b *testing.B) {
	var data [][]byte
	data = append(data, []byte("Station1;123"))
	data = append(data, []byte("Station2;-456"))
	data = append(data, []byte("Station3;0"))
	for i := 0; i < b.N; i++ {
		for _, line := range data {
			_, _, _ = parseStationLine(line)
		}
	}
}

func Test_parseMeasurement(t *testing.T) {
	tests := []struct {
		name       string
		linePart   []byte
		wantResult int16
	}{
		{
			name:       "Positive number",
			linePart:   []byte("12.3"),
			wantResult: 123,
		},
		{
			name:       "Negative number",
			linePart:   []byte("-45.6"),
			wantResult: -456,
		},
		{
			name:       "Zero",
			linePart:   []byte("0"),
			wantResult: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := parseMeasurement(tt.linePart)
			if gotResult != tt.wantResult {
				t.Errorf("parseMeasurement() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

// 7,86 ns/op
func Benchmark_parseMeasurement(b *testing.B) {
	var data [][]byte
	data = append(data, []byte("12.3"))
	data = append(data, []byte("-45.6"))
	data = append(data, []byte("0"))
	for i := 0; i < b.N; i++ {
		for _, line := range data {
			parseMeasurement(line)
		}
	}
}

func Test_minMeasurements(t *testing.T) {
	tests := []struct {
		name     string
		a        int16
		b        int16
		expected int16
	}{
		{
			name:     "a is smaller than b",
			a:        -10,
			b:        10,
			expected: -10,
		},
		{
			name:     "b is smaller than a",
			a:        15,
			b:        8,
			expected: 8,
		},
		{
			name:     "a and b are equal",
			a:        20,
			b:        20,
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minMeasurements(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("minMeasurements() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// a < b -> a | a > b -> b = 0,9536 ns/op
// (a + ((b - a) & ((b - a) >> 15))) = 0.8010 ns/op
func Benchmark_minMeasurements(b *testing.B) {
	var v1 int16 = 357
	var v2 int16 = -618
	for i := 0; i < b.N; i++ {
		minMeasurements(v1, v2)
	}
}

func Test_maxMeasurements(t *testing.T) {
	tests := []struct {
		name     string
		a        int16
		b        int16
		expected int16
	}{
		{
			name:     "a is greater than b",
			a:        10,
			b:        -10,
			expected: 10,
		},
		{
			name:     "b is greater than a",
			a:        8,
			b:        15,
			expected: 15,
		},
		{
			name:     "a and b are equal",
			a:        20,
			b:        20,
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxMeasurements(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("maxMeasurements() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// a > b -> a | a < b -> b = 0,9538 ns/op
// (a - ((a - b) & ((a - b) >> 15))) = 0,8022 ns/op
func Benchmark_maxMeasurements(b *testing.B) {
	var v1 int16 = 357
	var v2 int16 = -618
	for i := 0; i < b.N; i++ {
		maxMeasurements(v1, v2)
	}
}
