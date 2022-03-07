package gocli

import (
	"fmt"
	"os"
	"strings"
)

// Argument of a CLI command (not a CLI option)
type Argument struct {
	Name        string
	Required    bool
	Description string
}

// CLI Command
type Command struct {
	// Name of the command (as referenced in the CLI)
	Name string

	// Description that is shown when "--help" is present
	LongDesc string

	// Description that is shown when "--help" is present
	// on the parent command
	ShortDesc string

	// Options
	Options *[]Option

	// Argument
	Argument Argument

	// Behavior of the command
	Behavior func(ctx Context)
}

type Context struct {
	// Parent's name.
	// For example, if "install" is the child of "run" which is the child of "root",
	// then "install" is run with the bash command `root run install`. In that case,
	// `root run install` is the Referrer
	Referrer string

	// Configurations for the child commands
	Children []*Command

	// Configuration for the command
	Command *Command

	// Configurations for the options
	Options []Option

	// Arguments in their raw forms
	StrArgs []string

	// A map of arguments that have been parsed and type-casted to their
	// respective types. An argument default to the value of <nil> if they
	// are not included in the cli command
	Args map[string]interface{}
}

func (c *Command) Run(args []string, parents []string, children []*Command) {

	if c.Options == nil {
		c.Options = &[]Option{}
	}
	temp := *c.Options
	(*c.Options) = append((*c.Options), DefaultOptions...)

	// build the context
	context := Context{
		Referrer: strings.Join(append(parents, c.Name), " "),
		Command:  c,
		Options:  *c.Options,
		StrArgs:  args,
		Children: children,
	}

	for _, arg := range args {
		if arg == "--help" {
			fmt.Println(context.HelpStr())
			return
		}
	}

	if c.Behavior == nil {
		fmt.Printf("Behavior method not configured for command '%s'", context.Referrer)
		return
	}

	populateArgs(&context)

	// run the behavior
	c.Behavior(context)

	c.Options = &temp
}

func (c *Command) RunUtil(args []string, childrenMap map[*Command][]*Command, parents []string) {
	if len(args) == 0 || string(args[0][0]) == "-" {
		c.Run(args, parents, childrenMap[c])
	} else {
		subCmd := &Command{}

		// check if the args match a child command
		for _, child := range childrenMap[c] {
			if child.Name == args[0] {
				subCmd = child
			}
		}
		if subCmd.Name == "" {
			c.Run(args, parents, childrenMap[c])
		} else {
			subCmd.RunUtil(args[1:], childrenMap, append(parents, c.Name))
		}
	}
}

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func paddedName(name string, width int) (p string) {
	p += name
	for i := 0; i < width-len(name); i++ {
		p += " "
	}
	return p
}

// Returns the help string for a command
func (c *Context) HelpStr() string {
	padding := 5
	txt := ""
	name := c.Referrer

	txt += fmt.Sprintf("Usage: %s", name)
	if len(c.Children) > 0 {
		txt += " [COMMAND]"
	}

	if len(c.Options) > 0 {
		txt += fmt.Sprintf(" [OPTIONS]")
	}

	if c.Command.Argument.Name != "" {
		txt += fmt.Sprintf(" %s", c.Command.Argument.Name)
	}

	txt += Sep()

	if c.Command.LongDesc != "" {
		txt += c.Command.LongDesc
		txt += Sep()
	}
	txt += Sep()

	if len(c.Children) > 0 {
		txt += "Commands:" + Sep()

		maxWidth := 0
		for _, child := range c.Children {
			maxWidth = max(len(child.Name), maxWidth)
		}
		width := maxWidth + padding
		for _, child := range c.Children {
			txt += "  " + paddedName(child.Name, width) + child.ShortDesc + Sep()
		}
		txt += Sep()
	}
	if options := c.Options; len(options) > 0 {
		txt += "Options:" + Sep()

		maxWidth := 0
		for _, option := range options {
			maxWidth = max(len(option.Name()), maxWidth)
		}
		width := maxWidth + padding
		for _, option := range options {
			required := "Optional"
			if option.Required {
				required = "Required"
			}
			txt += "  " + paddedName(option.Name(), width) + fmt.Sprintf("[%s, Type: %s] ", required, option.Type) + option.Description + Sep()
		}
		txt += Sep()
	}

	if arg := c.Command.Argument; arg.Name != "" {
		required := "Optional"
		if c.Command.Argument.Required {
			required = "Required"
		}
		txt += fmt.Sprintf("Argument: '%s' (%s)", arg.Name, required) + Sep()
		txt += arg.Description
	}

	return txt
}

// Populate an interface with argument values
func populateArgs(c *Context) {
	args, err := ParseArgs(*c.Command.Options, c.Command.Argument, c.StrArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.Args = args
}
