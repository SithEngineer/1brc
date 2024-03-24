package main

import (
	"bytes"
	"testing"
)

// using ints -> 206,1 ns/op
// using float32 -> 309.5 ns/op
func Benchmark_printMeasurement(b *testing.B) {
	m := measurement{
		name:        []byte{'B', 'u', 'e', 'n', 'o', 's', ' ', 'A', 'i', 'r', 'e', 's'},
		min:         225,
		max:         335,
		sum:         560,
		nrSightings: 2,
	}
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		printMeasurement(m, &buf)
	}
}

func Test_printMeasurement(t *testing.T) {
	tests := []struct {
		name     string
		given    measurement
		expected string
	}{
		{
			name: "simple test",
			given: measurement{
				name:        []byte{'B', 'u', 'e', 'n', 'o', 's', ' ', 'A', 'i', 'r', 'e', 's'},
				min:         225,
				max:         335,
				sum:         560,
				nrSightings: 2,
			},
			expected: "Buenos Aires;22.5;28.0;33.5\n",
		},
		{
			name: "negative values",
			given: measurement{
				name:        []byte{'B', 'u', 'e', 'n', 'o', 's', ' ', 'A', 'i', 'r', 'e', 's'},
				min:         -335,
				max:         -225,
				sum:         -560,
				nrSightings: 2,
			},
			expected: "Buenos Aires;-33.5;-28.0;-22.5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printMeasurement(tt.given, &buf)
			if err != nil {
				t.Fatalf("got err: %s", err)
			}

			res := buf.String()
			if res != tt.expected {
				t.Errorf("expected result to be %s, but got %s", tt.expected, res)
			}
		})
	}
}
