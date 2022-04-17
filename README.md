termtext is a Go package to deal with monospace text as interpreted by
terminals.

Mostly intended to neatly align/wrap stuff terminal. There are a few tricky bits
with this:

1. Multiple codepoints can be combined to render one character (or "grapheme
   cluster" in Unicode speak).
2. Some characters are rendered as double-width, such as East-Asian characters,
   some emojis, and a number of others.
3. A single tab can render as multiple spaces, and the number of spaces depends
   on its position in the string.
4. Some characters aren't actually printed to the screen, such as the zero-width
   space and escape sequences to set the colour.

This package takes care of these issue.

Import as `zgo.at/termtext` – godoc: https://godocs.io/zgo.at/termtext

---

The main function is `termtext.Width()`:

    Width("\ta")                → 9    Tab expands to 8 spaces, followed by "a".
    Width("a\t")                → 8    Tab expands to 7 spaces.
	Width("🧑\t🧑")             → 10   🧑 is double-width
    Width("\x1b[1mbold\x1b[0m") → 4    Escape sequences are ignored.

You can configure the tab width by setting `termtext.TabWidth`.

There are a few other functions too:

    Expand()            Expand tabs.
    Slice()             Slice a string by display width, like str[n:m].

    AlignLeft()         Align a string, filling up the remainder with spaces.
    AlignRight()
    AlignCenter()

    Wrap()              Wrap a string. This is a simple wrap which just breaks
                        if a line's width is too long.
    WordWrap()          Word-wrap a string: lines are at most w wide, but don't
                        break in the middle of words.
