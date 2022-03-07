// An example usage of the framework
package gocli

import (
	"fmt"
	"os"
)

//////////////////////////////////////////
// OPTIONS
//////////////////////////////////////////

// number of something
var num = Option{
	Short:       "n",
	Description: "Indicates how many steps up the file system the example will take.",
	Required:    true,
	Type:        "int",
}

var label = Option{
	Short:       "l",
	Long:        "label",
	Description: "Label for the output.",
	Required:    false,
	Type:        "string",
}

// run in verbose mode
var verbose = Option{
	Short:       "v",
	Long:        "verbose",
	Description: "Run in verbose mode",
	Required:    false,
	Type:        "bool",
}

//////////////////////////////////////////
// Arguments
//////////////////////////////////////////
var dir = Argument{
	Name:        "directory",
	Description: "The starting directory for the command",
	Required:    true,
}

//////////////////////////////////////////
// Commands
//////////////////////////////////////////
var rootCmd = Command{
	Name:      "example(.exe)",
	ShortDesc: "",
	LongDesc:  "An example of the cli-framework for Go! Try out the sub-command.",
	Behavior:  rootBehavior,
}

var childCmd = Command{
	Name:      "run",
	ShortDesc: "Run the example",
	LongDesc:  "Print the working directory while ascending the filesystem 'n' times.",
	Options:   &[]Option{num, verbose, label},
	Argument:  dir,
	Behavior:  childBehavior,
}

//////////////////////////////////////////
// Command Behaviors
//////////////////////////////////////////

// root command will just rerun itself with the 'help' flag set to true
func rootBehavior(ctx Context) {
	// print help string
	fmt.Print(ctx.HelpStr())
}

func childBehavior(ctx Context) {
	// Read cli options
	n := ctx.Args["n"].(int)
	verbose := ctx.Args["verbose"].(bool) // guaranteed to either be true of false (not nil)
	label := ""
	if ctx.Args["label"] != nil {
		label = Blue("[" + ctx.Args["label"].(string) + "] ") // pad with a space at the end
	}

	// Read cli argument 'directory'
	dir := ctx.Args["directory"].(string)

	// build the bash command
	cmd := fmt.Sprintf("cd %s", dir)
	for i := 0; i < n; i++ {
		cmd = fmt.Sprintf(`
		%s
		pwd
		cd ..
		`, cmd)
	}

	// run the bash command
	if verbose {
		fmt.Println(fmt.Sprintf("%sStarting...", label))
		res := BashStreamLabel(cmd, true, true, label)
		if res.Err != nil {
			fmt.Println(res.Err)
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("%sFinished", label))
	} else {
		fmt.Println("Starting...")
		res := Bash(cmd)
		if res.Err != nil {
			fmt.Println(res.Err)
			os.Exit(1)
		}
		fmt.Println("Finished (run again with `--verbose` to see the output)")
	}
}

//////////////////////////////////////////
// main.go
//////////////////////////////////////////
//
// func main() {
// 	cli := NewCli(&rootCmd)

// 	// add Children
// 	cli.AddChild(&rootCmd, &childCmd)

// 	// execute
// 	cli.Exec()
// }
