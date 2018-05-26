package decorator

import (
	"fmt"
)

func Colorize(color, text string) string {
	switch color {
	case "red":
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", text)
	case "green":
		return fmt.Sprintf("\x1b[32m%s\x1b[0m", text)
	case "yellow":
		return fmt.Sprintf("\x1b[33m%s\x1b[0m", text)
	case "blue":
		return fmt.Sprintf("\x1b[34m%s\x1b[0m", text)
	case "purple":
		return fmt.Sprintf("\x1b[35m%s\x1b[0m", text)
	}

	return text
}

func BoolToColorizedString(b bool) string {
	if b {
		return Colorize("green", fmt.Sprintf("%t", b))
	}

	return Colorize("red", fmt.Sprintf("%t", b))
}
