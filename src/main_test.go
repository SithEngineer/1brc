package main

import (
	"testing"
)

func Test_extractMeasurement(t *testing.T) {
	tests := map[string]struct {
		given string
		expectedKey string
		expectedValue float32
	}{
		"happy path": {
			given: "Buenos Aires;22.5",
			expectedKey: "Buenos Aires",
			expectedValue: 22.5,
		},
		"key with special chars": {
			given: "St. John's;15.2",
			expectedKey: "St. John's",
			expectedValue: 15.2,
		},
		"negative value": {
			given: "Cracow;-1",
			expectedKey: "Cracow",
			expectedValue: -1,
		},
	}
	
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			city, measurement := extractMeasurement(test.given)
			if city != test.expectedKey {
				t.Errorf("expected city to be %s, but got %s", test.expectedKey, city)
			}
			if measurement != test.expectedValue {
				t.Errorf("expected measurement to be %v, but got %v", test.expectedValue, measurement)
			}
		})
	}
}

func Test_printMeasurement(t *testing.T) {
	tests := map[string]struct {
		givenCity string
		givenMeasurement measurements
		expected string
	}{
		"simple test": {
			givenCity: "Buenos Aires",
			givenMeasurement: measurements{
				min: 22.5,
				max: 33.5,
				mean: 28.5,
			},
			expected: "Buenos Aires;22.5;28.5;33.5",
		},
	}
	
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := printMeasurement(test.givenCity, test.givenMeasurement)
			if result != test.expected {
				t.Errorf("expected result to be %s, but got %s", test.expected, result)
			}
		})
	}
}

func Test_process(t *testing.T){
	tests := map[string]struct {
		given []input
		expected map[string]measurements
	}{
		"simple processing": {
			given: []input{
				{"Cracow", 1},
				{"Cracow", 2},
				{"Cracow", 3},
			},
			expected: map[string]measurements{
				"Cracow": {min: 1, mean: 2.25, max: 3},
			},
		},
		"slightly more complex": {
			given: []input{
				{"Buenos Aires", 22.5,},
				{"St. John's", 15.2},
				{"Cracow", -1},
				{"Buenos Aires", 14.5},
				{"St. John's", 11.2},
				{"Cracow", 8.3},
			},
			expected: map[string]measurements{
				"Buenos Aires": {min: 14.5, max: 22.5, mean: 18.5},
				"St. John's": {min: 11.2, max: 15.2, mean: 13.2},
				"Cracow": {min: -1, max: 8.3, mean: 3.65},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := map[string]measurements{}
			for _, in := range test.given {
				process(in, result)
			}
			if len(result) != len(test.expected) {
				t.Errorf("expected result to have %d keys, but got %d", len(test.expected), len(result))
			}
			for k, v := range result {
				if v != test.expected[k] {
					t.Errorf("expected result to have value %v for key %s, but got %v", test.expected[k], k, v)
				}
			}
		})
	}
}
