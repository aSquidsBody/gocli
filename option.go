package gocli

import "strings"

type Option struct {
	// A short description of the option
	Description string

	// Short name of the option (e.g. "v")
	Short string

	// Long name of the option (e.g. "verbose")
	Long string

	// Indicates whether the option is required or not
	Required bool

	// Type of the option: "string", "bool", or "float", "int"
	Type string
}

func (o *Option) Name() string {
	if o.Short != "" && o.Long != "" {
		return strings.Join([]string{"-" + o.Short, "--" + o.Long}, ",")
	}
	if o.Short != "" {
		return "-" + o.Short
	}
	if o.Long != "" {
		return "--" + o.Long
	}
	return ""
}

// default options
var HelpOption = Option{
	Description: "Print a help string",
	Short:       "",
	Long:        "help",
	Required:    false,
	Type:        "bool",
}

var DefaultOptions = []Option{
	HelpOption,
}
