package gocli

import color "github.com/fatih/color"

func Blue(s string) string {
	return color.New(color.FgBlue).SprintFunc()(s)
}

func Red(s string) string {
	return color.New(color.FgRed).SprintFunc()(s)
}

func Yellow(s string) string {
	return color.New(color.FgYellow).SprintFunc()(s)
}

func Green(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}

func White(s string) string {
	return color.New(color.FgWhite).SprintFunc()(s)
}

func Cyan(s string) string {
	return color.New(color.FgCyan).SprintFunc()(s)
}

func Magenta(s string) string {
	return color.New(color.FgMagenta).SprintFunc()(s)
}
