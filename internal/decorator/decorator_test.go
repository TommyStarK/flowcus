package decorator

import "testing"

func TestColorize(t *testing.T) {
	if "test" != Colorize("", "test") {
		t.Errorf("Should be white")
	}

	if "\x1b[31mtest\x1b[0m" != Colorize("red", "test") {
		t.Errorf("Wrong color. Expected red")
	}

	if "\x1b[32mtest\x1b[0m" != Colorize("green", "test") {
		t.Errorf("Wrong color. Expected green")
	}

	if "\x1b[33mtest\x1b[0m" != Colorize("yellow", "test") {
		t.Errorf("Wrong color. Expected yellow")
	}

	if "\x1b[34mtest\x1b[0m" != Colorize("blue", "test") {
		t.Errorf("Wrong color. Expected blue")
	}

	if "\x1b[35mtest\x1b[0m" != Colorize("purple", "test") {
		t.Errorf("Wrong color. Expected purple")
	}

	t.Log("Colorize test succeed")
}

func TestBoolToColorizedString(t *testing.T) {
	if "\x1b[32mtrue\x1b[0m" != BoolToColorizedString(true) {
		t.Errorf("Wrong color. Expected green")
	}

	if "\x1b[31mfalse\x1b[0m" != BoolToColorizedString(false) {
		t.Errorf("Wrong color. Expected red")
	}
}
