package colors

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

const (
	fgTpl = "\033[38;5;%dm%s\u001b[0m"
	bgTpl = "\033[48;5;%d;38;5;%dm%s\u001b[0m"
)

//Colorer wraps a string with ansi color escape codes.
type Colorer func(string) string

func (c *Colorer) Decode(value string) error {
	ss := strings.Split(value, ",")
	if len(ss) == 1 {
		i := defaultuint(value, 0)
		*c = func(s string) string {
			return fmt.Sprintf(fgTpl, i, s)
		}
	} else if len(ss) == 2 {
		fg := defaultuint(ss[0], 0)
		bg := defaultuint(ss[1], 1)
		*c = func(s string) string {
			return fmt.Sprintf(bgTpl, bg, fg, s)
		}
	} else {
		*c = func(s string) string {
			return s
		}
	}

	return nil
}

func defaultuint(s string, i uint8) uint8 {
	out, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return i
	}
	return uint8(out)
}

func Get() Colors {
	var c Colors
	err := envconfig.Process("MCLI", &c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Colors provide foreground and background colors for the ui.
// If the default value is comma separated the first value is
// the foreground color and the second is the background color.
// If only one value is supplied then it's the foreground color.
type Colors struct {
	C1 Colorer `default:"252"`
	C2 Colorer `default:"2"`
	C3 Colorer `default:"11"`
}
