package model

import "testing"

func TestStationDescription(t *testing.T) {
	tests := []struct {
		name    string
		station Station
		want    string
	}{
		{"empty", Station{}, ""},
		{"whitespace", Station{Tags: " ", Countrycode: " ", State: "\t", Language: " "}, ""},
		{"tags", Station{Tags: " jazz "}, "jazz"},
		{"unknown bitrate", Station{Bitrate: -1}, ""},
		{"location without country", Station{State: "Berlin", Language: "german"}, "Berlin, german"},
		{"metadata", Station{Tags: "jazz", Bitrate: 128, Countrycode: "DE", Language: "german"}, "jazz ┃ 128 kbps ┃ DE, german"},
		{"votes only", Station{Votes: 2}, "Votes: 2"},
		{"clicks only", Station{Clickcount: 3}, "Clicks: 3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.station.Description(); got != tt.want {
				t.Errorf("Description() = %q, want %q", got, tt.want)
			}
		})
	}
}
