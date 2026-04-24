package ui

import "strings"

type BigFont struct {
	glyphs map[rune][]string
	rows   int
}

var defaultBigGlyphs = map[rune][]string{
	'0': {
		"█████",
		"█   █",
		"█   █",
		"█   █",
		"█████",
	},
	'1': {
		"  ██ ",
		" ███ ",
		"  ██ ",
		"  ██ ",
		"█████",
	},
	'2': {
		"█████",
		"    █",
		"█████",
		"█    ",
		"█████",
	},
	'3': {
		"█████",
		"    █",
		"█████",
		"    █",
		"█████",
	},
	'4': {
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
	},
	'5': {
		"█████",
		"█    ",
		"█████",
		"    █",
		"█████",
	},
	'6': {
		"█████",
		"█    ",
		"█████",
		"█   █",
		"█████",
	},
	'7': {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
	},
	'8': {
		"█████",
		"█   █",
		"█████",
		"█   █",
		"█████",
	},
	'9': {
		"█████",
		"█   █",
		"█████",
		"    █",
		"█████",
	},
	':': {
		"     ",
		"  █  ",
		"     ",
		"  █  ",
		"     ",
	},
	' ': {
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	},
}

func NewBigFont() *BigFont {
	return &BigFont{glyphs: defaultBigGlyphs, rows: 5}
}

func (b *BigFont) Render(s string) string {
	lines := make([]strings.Builder, b.rows)
	for i, r := range s {
		glyph, ok := b.glyphs[r]
		if !ok {
			glyph = b.glyphs[' ']
		}
		for row := 0; row < b.rows; row++ {
			if i > 0 {
				lines[row].WriteString(" ")
			}
			lines[row].WriteString(glyph[row])
		}
	}
	out := make([]string, b.rows)
	for i := range lines {
		out[i] = lines[i].String()
	}
	return strings.Join(out, "\n")
}
