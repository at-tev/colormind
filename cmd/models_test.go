package cmd

import (
	"testing"
	"time"
)

func TestGetData(t *testing.T) {
	tests := [3]input{
		{},
		{"ui", nil},
		{"default", [][3]int{{122, 34, 56}}},
	}

	for i, test := range tests {
		data, err := getData(test)
		if data == nil || err != nil {
			t.Errorf("%d. got data: %v, err: %t", i, data, err)
		}
	}
}

func TestUpdated(t *testing.T) {
	cur := time.Now().UTC()

	tests := [10]struct {
		m    modelsList
		want bool
	}{
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 14, 0, 0, 0, time.UTC)}, true},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 6, 0, 0, 0, time.UTC)}, false},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 15, 0, 0, 0, time.UTC)}, true},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 4, 0, 0, 0, time.UTC)}, false},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 12, 0, 0, 0, time.UTC)}, true},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 5, 0, 0, 0, time.UTC)}, false},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 11, 0, 0, 0, time.UTC)}, true},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 1, 0, 0, 0, time.UTC)}, false},
		{modelsList{lastUpdate: time.Date(cur.Year(), cur.Month(), cur.Day(), 7, 0, 30, 0, time.UTC)}, true},
		{modelsList{lastUpdate: cur.AddDate(0, 0, -1)}, false},
	}

	for _, test := range tests {
		got := test.m.updated()
		if got != test.want {
			t.Errorf("%+v: want %t, got %t", test, test.want, got)
		}
	}
}
