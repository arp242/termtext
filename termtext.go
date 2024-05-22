// Package termtext deals with monospace text as interpreted by terminals.
package termtext

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/rivo/uniseg"
	"zgo.at/runewidth"
)

// Number of spaces to count a tab as.
var TabWidth = 8

// Widths sets the width for the given runes.
//
// This can be useful in some cases where the written text is different. For
// example when encoding JSON some characters get written as a \-escape. To
// account for this, you can override the width for those:
//
//	termtext.Escapes = map[rune]int{
//		'\b': 2,
//		'\f': 2,
//		'\n': 2,
//		'\r': 2,
//		'\t': 2,
//	}
var Widths map[rune]int

// Width gets the display width of a string.
//
// The "display width" is the number of columns a string will occupy in a
// monospace terminal.
func Width(s string) int {
	var (
		g   = uniseg.NewGraphemes(s)
		l   int
		esc bool
	)
	for g.Next() {
		runes := g.Runes()

		/// More than one codepoint: use the width of the first one non-zero one.
		if len(runes) > 1 {
			l += clusterWidth(runes)
			continue
		}

		if ll, ok := Widths[runes[0]]; ok {
			l += ll
			continue
		}

		/// One codepoint: check for tab and escapes.
		switch r := runes[0]; {
		case r == '\t':
			l += TabWidth - l%TabWidth
		case r == '\x1b':
			esc = true
		case esc:
			if r == 'm' {
				esc = false
			}
		default:
			l += runewidth.RuneWidth(r)
		}
	}
	return l
}

// Expand tabs to spaces.
func Expand(s string) string {
	var (
		g   = uniseg.NewGraphemes(s)
		b   strings.Builder
		l   int
		esc bool
	)
	b.Grow(len(s))
	for g.Next() {
		runes := g.Runes()

		if len(runes) > 1 {
			l += clusterWidth(runes)
			b.WriteString(string(runes))
			continue
		}

		/// One codepoint: check for tab and escapes.
		switch r := runes[0]; {
		case r == '\t':
			tw := TabWidth - l%TabWidth
			b.WriteString(fill(tw))
			l += tw
		case r == '\x1b':
			esc = true
			b.WriteRune(r)
		case esc:
			b.WriteRune(r)
			if r == 'm' {
				esc = false
			}
		default:
			b.WriteRune(r)
			l += runewidth.RuneWidth(r)
		}
	}
	return b.String()
}

// AlignLeft left-aligns a string, filling up any remaining width with spaces.
//
// Tabs will be expanded to spaces.
func AlignLeft(s string, w int) string {
	s = Expand(s)
	count := w - Width(s)
	if count <= 0 {
		return s
	}
	return s + fill(count)
}

// AlignRight right-aligns a string, filling up any remaining width with spaces.
//
// Tabs will be expanded to spaces.
func AlignRight(s string, w int) string {
	s = Expand(s)
	count := w - Width(s)
	if count <= 0 {
		return s
	}
	return fill(count) + s
}

// AlignCenter centre-aligns a string, filling up any remaining width with spaces.
//
// Tabs will be expanded to spaces.
func AlignCenter(s string, w int) string {
	if s == "" {
		return fill(w)
	}

	s = Expand(s)
	count := w - Width(s)
	if count <= 0 {
		return s
	}
	pad := fill(count / 2)
	if w%2 == 0 {
		return pad + s + pad + " "
	}
	return pad + s + pad
}

// Slice a string by character index. This works the same as str[n:m] slicing.
//
// Tabs will be expanded to spaces.
func Slice(s string, start, stop int) string {
	if start == stop {
		return ""
	}

	s = Expand(s)
	var (
		g                 = uniseg.NewGraphemes(s)
		startOff, stopOff int
		pos               int
		esc               bool
	)
	if stop == 0 {
		stopOff = len(s)
	}
	for g.Next() {
		if pos == start {
			startOff, _ = g.Positions()
			if stop == 0 {
				break
			}
		}
		if stop > 0 && pos == stop {
			stopOff, _ = g.Positions()
			break
		}

		runes := g.Runes()
		if len(runes) > 1 {
			pos += clusterWidth(runes)
			continue
		}

		/// One codepoint: check for tab and escapes.
		switch r := runes[0]; {
		case r == '\x1b':
			esc = true
		case esc:
			if r == 'm' {
				esc = false
			}
		default:
			pos += runewidth.RuneWidth(r)
		}
	}
	if stopOff == 0 {
		stopOff = len(s)
	}
	return s[startOff:stopOff]
}

// Wrap lines to be at most w characters wide.
//
// Lines will be prefixed with prefix. The prefix isn't counted in line length
// calculations.
//
// This does not use word wrap, use WordWrap() for this instead.
//
// Tabs will be expanded to spaces.
func Wrap(s string, w int, prefix string) string {
	s = Expand(s)
	var (
		g = uniseg.NewGraphemes(s)
		l = 0
		b strings.Builder
	)
	b.Grow(len(s))
	for g.Next() {
		runes := g.Runes()

		if len(runes) == 1 && runes[0] == '\n' {
			b.WriteRune(runes[0])
			b.WriteString(prefix)
			l = 0
			continue
		}

		cw := clusterWidth(runes)

		if l+cw > w {
			b.WriteByte('\n')
			b.WriteString(prefix)
			if len(runes) == 1 && unicode.IsSpace(runes[0]) { /// No spaces at start of line.
				l = 0
			} else {
				b.WriteString(string(runes))
				l = cw
			}
			continue
		}

		if l == 0 && len(runes) == 1 && unicode.IsSpace(runes[0]) { /// No spaces at start of line.
			continue
		}

		b.WriteString(string(runes))
		l += cw
	}
	return b.String()
}

// WordWrap wraps lines to be at most w characters wide, but doesn't break in
// the middle of words.
//
// Lines will be prefixed with prefix. The prefix isn't counted in line length
// calculations.
//
// Tabs will be expanded to spaces.
func WordWrap(s string, w int, prefix string) string {
	s = Expand(s)
	var (
		g       = uniseg.NewGraphemes(s)
		b       bytes.Buffer
		lineLen int
		wordLen int
		word    strings.Builder
	)
	b.Grow(len(s))
	for g.Next() {
		runes := g.Runes()

		/// Line break in input: reset word; start a new line.
		if len(runes) == 1 && runes[0] == '\n' {
			if lineLen+wordLen > w {
				b.WriteByte('\n')
				b.WriteString(prefix)
			}
			b.WriteString(word.String())
			b.WriteByte('\n')
			b.WriteString(prefix)

			word.Reset()
			wordLen, lineLen = 0, 0
			continue
		}

		runesLen := clusterWidth(runes)

		/// Note that unicode word breaks are actually quite a bit more complex,
		/// but yeah, this should be "good enough".
		if len(runes) == 1 && isBreak(runes[0]) {
			if lineLen > w { /// Break and write word on next line.
				b.Truncate(b.Len() - 1) /// Trailing space.
				b.WriteByte('\n')
				b.WriteString(prefix)
				b.WriteString(word.String())
				b.WriteString(string(runes))

				lineLen = wordLen + runesLen
				wordLen = 0
				word.Reset()
			} else { /// Write word on current line.
				lineLen += runesLen
				b.WriteString(word.String())
				b.WriteString(string(runes))

				word.Reset()
				wordLen = 0
			}

			continue
		}

		lineLen += runesLen
		wordLen += runesLen
		word.WriteString(string(runes))
	}

	/// Last word.
	if lineLen > 0 {
		if lineLen > w {
			if b.Len() > 0 {
				b.Truncate(b.Len() - 1) /// Trailing space.
				b.WriteByte('\n')
			}
			b.WriteString(prefix)
		}
		b.WriteString(word.String())
	}

	return b.String()
}

func isBreak(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}

func clusterWidth(runes []rune) int {
	// Our best guess is to use the width of the first non-zero-width rune.
	for _, r := range runes {
		if w := runewidth.RuneWidth(r); w > 0 {
			return w
		}
	}
	return 0
}

func fill(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
