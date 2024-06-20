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

func TestSlice(t *testing.T) {
	tests := []struct {
		in          string
		start, stop int
		want        string
	}{
		{"", 0, 0, ""},
		{"abc", 0, 0, ""},
		{"abc", 1, 1, ""},

		{"abc", 0, 1, "a"},
		{"abc", 1, 2, "b"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 0, 10,
			"┌─ \x1b[1mAshidee"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 0, 37,
			"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 0, 38,
			"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 0, 308,
			"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 1, 38,
			"─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 1, 308,
			"─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐"},

		{"┌─ \x1b[1mAshideena +2\x1b[0m ────┬─ \x1b[1mAttributes\x1b[0m ───┐", 3, 15,
			"Ashideena +2"},

		{"1\x1b[1m23\x1b[0m", 0, 2, "1\x1b[1m2"},
		{"\x1b[1m123\x1b[0m", 0, 2, "\x1b[1m12"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Slice(tt.in, tt.start, tt.stop)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		in                     string
		w                      int
		prefix                 string
		wantWrap, wantWordWrap string
	}{
		{"", 10, "", "", ""},
		{"Hello, world!", 15, "", "Hello, world!", "Hello, world!"},
		{"Hello, world!", 15, "XX", "Hello, world!", "Hello, world!"},

		{"Hello, world!", 10, "", "Hello, wor\nld!", "Hello,\nworld!"},
		{"Hello, world!", 10, "XX", "Hello, wor\nXXld!", "Hello,\nXXworld!"},
		{
			"https://www.fastmail.help/hc/en-us/articles/360058753614-Why-messages-bounce-back",
			80, "",
			"https://www.fastmail.help/hc/en-us/articles/360058753614-Why-messages-bounce-bac\nk",
			"https://www.fastmail.help/hc/en-us/articles/360058753614-Why-messages-bounce-back",
		},

		{`
This crude wooden club burns with the raging spirit of the demon forever trapped within by the powerful enchantments placed on the weapon. Occasionally, however, the demon's wrath escapes in a fiery blast.

STATISTICS:

Combat abilities:
– 20% chance per hit target will take an additional 10 points of fire damage
– 7% chance per hit a 15-ft. radius fireball will automatically detonate (5d6 fire damage; Save vs. Spell for half)
`, 75, "",
			// Wrap()
			`
This crude wooden club burns with the raging spirit of the demon forever tr
apped within by the powerful enchantments placed on the weapon. Occasionall
y, however, the demon's wrath escapes in a fiery blast.

STATISTICS:

Combat abilities:
– 20% chance per hit target will take an additional 10 points of fire damag
e
– 7% chance per hit a 15-ft. radius fireball will automatically detonate (5
d6 fire damage; Save vs. Spell for half)
`,
			// WordWrap()
			`
This crude wooden club burns with the raging spirit of the demon forever
trapped within by the powerful enchantments placed on the weapon.
Occasionally, however, the demon's wrath escapes in a fiery blast.

STATISTICS:

Combat abilities:
– 20% chance per hit target will take an additional 10 points of fire 
damage
– 7% chance per hit a 15-ft. radius fireball will automatically detonate
(5d6 fire damage; Save vs. Spell for half)
`,
		},

		{
			`Priests of Horus-Re in far off Mulhorand were the first to create this life-saving potion. A foul-smelling brew, it is made by boiling used mummy wrappings along with naturally desiccated animal remains. The resulting tea can then be consumed to neutralize the effects of disease and certain afflictions.

STATISTICS:

Special: Cures blindness, deafness, and disease

Weight: 1`,
			59,
			"",
			`Priests of Horus-Re in far off Mulhorand were the first to 
create this life-saving potion. A foul-smelling brew, it is
made by boiling used mummy wrappings along with naturally d
esiccated animal remains. The resulting tea can then be con
sumed to neutralize the effects of disease and certain affl
ictions.

STATISTICS:

Special: Cures blindness, deafness, and disease

Weight: 1`,
			`Priests of Horus-Re in far off Mulhorand were the first to
create this life-saving potion. A foul-smelling brew, it is
made by boiling used mummy wrappings along with naturally
desiccated animal remains. The resulting tea can then be
consumed to neutralize the effects of disease and certain 
afflictions.

STATISTICS:

Special: Cures blindness, deafness, and disease

Weight: 1`,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Wrap(tt.in, tt.w, tt.prefix)
			if have != tt.wantWrap {
				t.Errorf("Wrap() wrong\nhave:\n%s\nwant:\n%s", have, tt.wantWrap)
			}

			have = WordWrap(tt.in, tt.w, tt.prefix)
			if have != tt.wantWordWrap {
				t.Errorf("WordWrap() wrong\nhave:\n%s\nwant:\n%s", have, tt.wantWordWrap)
			}
		})
	}
}

func BenchmarkWidth(b *testing.B) {
	for _, s := range []string{
		"",
		"Hello",
		strings.Repeat("Just some plain ASCII text", 10),
		"More \x1b[1mcomplex\x1b[0m text 🙃 \t 🙃 \t 🙃",
	} {
		b.Run("", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Width(s)
			}
		})
	}
}
