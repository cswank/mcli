package colors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aybabtme/rgbterm"
	ui "github.com/jroimartin/gocui"
)

var (
	ansiColors = map[string]string{
		"black":   "30",
		"red":     "31",
		"green":   "32",
		"yellow":  "33",
		"blue":    "34",
		"magenta": "35",
		"cyan":    "36",
		"white":   "37",
	}

	lookup = map[string]Colorer{
		"black": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["black"], ansiColors["black"]), s)
		},
		"red": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["red"], ansiColors["red"]), s)
		},
		"green": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["green"], ansiColors["green"]), s)
		},
		"yellow": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["yellow"], ansiColors["yellow"]), s)
		},
		"blue": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["blue"], ansiColors["blue"]), s)
		},
		"magenta": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["magenta"], ansiColors["magenta"]), s)
		},
		"cyan": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["cyan"], ansiColors["cyan"]), s)
		},
		"white": func(s string) string {
			return fmt.Sprintf(fmt.Sprintf("\033[%sm%%s\033[%sm", ansiColors["white"], ansiColors["white"]), s)
		},
	}

	background = map[string]ui.Attribute{
		"black":   ui.ColorBlack,
		"red":     ui.ColorRed,
		"green":   ui.ColorGreen,
		"yellow":  ui.ColorYellow,
		"blue":    ui.ColorBlue,
		"magenta": ui.ColorMagenta,
		"cyan":    ui.ColorCyan,
		"white":   ui.ColorWhite,
	}
)

//GetBackground sets the background color for the ui.
func GetBackground(s string) ui.Attribute {
	c, ok := background[s]
	if !ok {
		return background["black"]
	}
	return c
}

//Colorer wraps a string with ansi color escape codes.
type Colorer func(string) string

//Get fetches the colorer func for the given color
func Get(s string) Colorer {
	if strings.Contains(s, ",") {
		return getRGB(s)
	}
	return lookup[s]
}

func getRGB(c string) Colorer {
	parts := strings.Split(c, ",")
	r, g, b := rgbFromStrings(parts)
	return func(s string) string {
		return rgbterm.FgString(s, r, g, b)
	}
}

func rgbFromStrings(s []string) (uint8, uint8, uint8) {
	if len(s) != 3 {
		s = []string{"255", "255", "255"}
	}
	r, err := strconv.ParseUint(s[0], 10, 8)
	if err != nil {
		return rgbFromStrings([]string{})
	}

	g, err := strconv.ParseUint(s[1], 10, 8)
	if err != nil {
		return rgbFromStrings([]string{})
	}

	b, err := strconv.ParseUint(s[2], 10, 8)
	if err != nil {
		return rgbFromStrings([]string{})
	}

	return uint8(r), uint8(g), uint8(b)
}
