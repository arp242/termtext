package termtext

import (
	"fmt"
	"strings"
	"testing"
)

func TestWidth(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},

		// Tabs.
		{"\t", 8},
		{"\ta", 9},
		{"a\t", 8},
		{"aaaa\tx", 9},
		{"aaaaaaa\tx", 9},
		{"\t\t", 16},
		{"a\ta\t", 16},
		{"a\ta\ta", 17},
		{"\t", 8},
		{"\tሧ", 9},
		{"ሧ\t", 8},
		{"ሧሧሧሧ\tx", 9},
		{"ሧሧሧሧሧሧሧ\tx", 9},
		{"\t\t", 16},
		{"ሧ\tሧ\t", 16},
		{"ሧ\tሧ\tሧ", 17},

		// Escape characters
		{"\x1b123]m asd\x1b0m", 4},
		{"a\x05", 1},

		// Combining
		{"Mö\u0308hr", 4},

		{"\u200B\ufeff", 0},

		// Emojis.
		{"🧑\u200d🚒", 2},
		{"🧑🏽\u200d🚒", 2},
		{"🧑\t🧑", 10},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Width(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %d\nwant: %d", have, tt.want)
			}
		})
	}
}

func TestExpand(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"\t", strings.Repeat(" ", 8)},
		{"xx\t", "xx" + strings.Repeat(" ", 6)},
		{"🧑\t", "🧑" + strings.Repeat(" ", 6)},
		{"\x1b123]mMö\u0308h🧑\x1b0m\t", "\x1b123]mMö\u0308h🧑\x1b0m   "},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			have := Expand(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestAlign(t *testing.T) {
	tests := []struct {
		in                  string
		n                   int
		left, right, center string
	}{
		{"", 4, "    ", "    ", "    "},
		{"a", 4, "a   ", "   a", " a  "},

		{"Hello", 4, "Hello", "Hello", "Hello"},
		{"Hello", -2, "Hello", "Hello", "Hello"},

		{"Hello", 6, "Hello ", " Hello", "Hello "},
		{"Hello", 7, "Hello  ", "  Hello", " Hello "},
		{"Hello", 8, "Hello   ", "   Hello", " Hello  "},
		{"Hello", 9, "Hello    ", "    Hello", "  Hello  "},
		{"Hello", 10, "Hello     ", "     Hello", "  Hello   "},
		{"Hello", 11, "Hello      ", "      Hello", "   Hello   "},

		{"\x1b123]mMö\u0308h🧑\x1b0m\t", 8,
			"\x1b123]mMö\u0308h🧑\x1b0m   ",
			"\x1b123]mMö\u0308h🧑\x1b0m   ",
			"\x1b123]mMö\u0308h🧑\x1b0m   ",
		},

		{"\x1b123]mMö\u0308h🧑\x1b0m\t", 9,
			"\x1b123]mMö\u0308h🧑\x1b0m    ",
			" \x1b123]mMö\u0308h🧑\x1b0m   ",
			"\x1b123]mMö\u0308h🧑\x1b0m   ",
		},

		{"\x1b123]mMö\u0308h🧑\x1b0m\t", 10,
			"\x1b123]mMö\u0308h🧑\x1b0m     ",
			"  \x1b123]mMö\u0308h🧑\x1b0m   ",
			" \x1b123]mMö\u0308h🧑\x1b0m     ",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%d", tt.in, tt.n), func(t *testing.T) {
			left := AlignLeft(tt.in, tt.n)
			right := AlignRight(tt.in, tt.n)
			center := AlignCenter(tt.in, tt.n)

			if left != tt.left {
				t.Errorf("left wrong\ngot:  %q\nwant: %q", left, tt.left)
			}
			if right != tt.right {
				t.Errorf("right wrong\ngot:  %q\nwant: %q", right, tt.right)
			}
			if center != tt.center {
				t.Errorf("center wrong\ngot:  %q\nwant: %q", center, tt.center)
			}
		})
	}
}
