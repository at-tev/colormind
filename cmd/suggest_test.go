package cmd

import (
	"fmt"
	"testing"
)

func TestToHexScheme(t *testing.T) {
	tests := [10]struct {
		scheme []string
		want   [][3]int
	}{
		{[]string{"ff24", "61"}, [][3]int{{255, 36, 0}, {97, 0, 0}}},
		{[]string{"a9fw67"}, nil},
		{[]string{"hgjff2400", "610000"}, [][3]int{{255, 36, 0}, {97, 0, 0}}},
		{[]string{"FFF"}, nil},
		{[]string{"ffff", "858205"}, [][3]int{{255, 255, 0}, {133, 130, 5}}},
		{[]string{"123, 90, 89", "a3b7c8"}, nil},
		{[]string{"212101"}, [][3]int{{33, 33, 1}}},
		{[]string{"b9c8c8", ""}, nil},
		{[]string{"ffffxz", "kl858205h"}, [][3]int{{255, 255, 0}, {133, 130, 5}}},
		{[]string{"a3c7b66", "000000"}, nil},
	}

	for _, test := range tests {
		got, err := toHexScheme(test.scheme)
		if err != nil {
			fmt.Println("err: ", err)
			continue
		}

		switch test.want {
		case nil:
			if got != nil {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		default:
			for i := range got {
				if test.want[i] != got[i] {
					t.Errorf("got: %v, want: %v", got, test.want)
				}
			}
		}
	}
}

func TestToScheme(t *testing.T) {
	tests := [10]struct {
		scheme []string
		want   [][3]int
	}{
		{[]string{"1, 2, 3", "2, 3, 4", "3, 4, 5"}, [][3]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}},
		{[]string{"jhfwh"}, nil},
		{[]string{"255.255:255"}, [][3]int{{255, 255, 255}}},
		{[]string{"2345, 0, 0"}, nil},
		{[]string{"10, 23, 240, 75"}, [][3]int{{10, 23, 240}}},
		{[]string{"256", "2", "3"}, nil},
		{[]string{"56"}, [][3]int{{56, 0, 0}}},
		{[]string{"890, 34"}, nil},
		{[]string{"-128|0|255"}, [][3]int{{128, 0, 255}}},
		{[]string{"45, 34, 78", "bugkj.;kh"}, nil},
	}

	for _, test := range tests {
		got, err := toScheme(test.scheme)
		if err != nil {
			fmt.Println("err: ", err)
			continue
		}

		switch test.want {
		case nil:
			if got != nil {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		default:
			for i := range got {
				if test.want[i] != got[i] {
					t.Errorf("got: %v, want: %v", got, test.want)
				}
			}
		}
	}
}

func TestIsValidRGB(t *testing.T) {
	tests := [10]struct {
		rgb  [3]int
		want bool
	}{
		{[3]int{1, 2, 3}, true}, {[3]int{1000, 0, 0}, false},
		{[3]int{0, 0, 0}, true}, {[3]int{1, 256, 3}, false},
		{[3]int{255, 255, 255}, true}, {[3]int{1, 2, 344}, false},
		{[3]int{128, 10, 3}, true}, {[3]int{111, 222, 333}, false},
		{[3]int{56, 24, 63}, true}, {[3]int{100, 200, 300}, false},
	}

	for _, test := range tests {
		got := isValidRGB(test.rgb)
		if got != test.want {
			t.Errorf("want %t, got %t", test.want, got)
		}
	}
}

func TestIsValidHex(t *testing.T) {
	tests := [10]struct {
		color string
		want  bool
	}{
		{"ab", true}, {"00h000", false},
		{"FFFF", true}, {"123", false},
		{"ff6800", true}, {"2f9b4d6", false},
		{"000000", true}, {"FFF", false},
		{"858205", true}, {"a5nnc7", false},
	}

	for _, test := range tests {
		got := isValidHex(test.color)
		if got != test.want {
			t.Errorf("want %t, got %t", test.want, got)
		}
	}
}
