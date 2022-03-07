Command line interface framework in Go.

# Installing

```bash
go get github.com/aSquidsBody/gocli@v1.0.3
```

# Usage

Commands can be templated using the `Command` struct

```go
// example-commands.go
package main

import "github.com/aSquidsBody/gocli"

var ExampleCommand = gocli.Command{
    Name: "command",
    LongDesc: "An example command",
    ShortDesc: "",
    // ...more configurations
}
```

A behavior can be defined for the command

```go
// example-commands.go
package main

import (
    "fmt"
    "github.com/aSquidsBody/gocli"
)

var ExampleCommand = gocli.Command{
    Name: "command",
    LongDesc: "An example command",
    ShortDesc: "",
    Behavior: ExampleBehavior
}

func ExampleBehavior(ctx gocli.Context) {
    fmt.Println("Hello world!")
}
```

The command can be used as an entrypoint for the CLI

```go
// main.go
package main

import "github.com/aSquidsBody/gocli"

func main() {
    // init the cli
    cli := gocli.NewCli(&ExampleCommand)

    // run the cli
    cli.Exec()
}
```

The build the go executable with either

```bash
go build -o <executable-name>
```

or

```bash
go install
```

See a more in-depth example in [example.go](./example.go).

# API Documentation

## Command

A configuration template for CLI commands

_Command struct Fields_

### Command.Name

_Required_

Type: `String`

Name of the command (as referenced in the CLI)

### Command.Behavior

_Required_

Type: `func(ctx Context) void`

A function which preforms the behavior of the command.

### Command.ShortDesc

_Optional_

Type: `String`

Description that is printed when "--help" is included in the _parent_ command.

### Command.LongDesc

_Optional_

Type: `String`

Description that is printed when "--help" is included in the command.

### Command.Options

_Optional_

Type: `*[]Option`

`Command.Options` is a pointer to a list of `Option` structs for the command

### Command.Argument

_Optional_

Type: `Argument`

An `Argument` struct to define argument for the command

## Context

Context is the object that is passed to `Command.Behavior` when a command is run. It is populated
with CLI arguments.

### Context.Referrer

Type: `string`

The collective name of the commands which built the context. For example, if the user ran `root child1 child2 [OPTIONS]` in the command prompt,
then the Context was built in the `child2` command; however, the _Referrer_ field would equal `root child1 child2`

### Context.Children

Type: `[]*Command`

An array of pointers to Command configurations for each sub-command

### Context.Options

Type: `[]Option`

An array of Option configurations for the command

### Context.StrArgs

Type: `[]string`

An array of unprocessed CLI args

### Context.Args

Type: `map[string]interface{}`

A map of arguments that have been parsed and type-casted to their
respective types. An argument defaults to the value of <nil> if it
is not included in the cli command

### [METHOD] Context.HelpStr()

Parameters: None

Prints the help string for the command

## Option

A configuration template for cli options

### Option.Type

_Required_

Type: `string`

Indicates the expected datatype of the option. Can be "bool", "string", "int", and "float"

### Option.Description

_Optional_

Type: `string`

A description of the option

### Option.Short

_Optional_

Type: `string`

A short name for the option. The "-" should be omitted, e.g. if you want to configure
_-v_ as a short name for the option, then set _Short_ to _v_.

### Option.Long

_Optional_

Type: `string`

A long name for the option. The "--" should be omitted, e.g. if you want to configure
_--verbose_ as a long name for the option, then set _Long_ to _verbose_.

### Option.Required

_Optional_

Type: `bool`

Indicates whether the option is required (true) or optional (false). Defaults to false.

## BashResult

### BashResult.Stdout

Type: `string`

The _stdout_ of a bash command.

### BashResult.Stderr

Type: `string`

The _stderr_ of a bash command.

### BashResult.Err

Type: `error`

An error which may have occurred during the bash command. If the bash command fails (exits with non-zero code), then
_BashResult.Err_ will not be <nil>

## [Function] Bash

Parameters: _cmd_ `string`

Returns: `BashResult`

Runs a bash command.

## [Function] BashStream

Parameters:

- _cmd_ `string`
- _stdout_ `bool`
- _stderr_ `bool`

Returns: `BashResult`

Runs a bash command. If _stdout_ is true, the shell's stdout is copied in real time to the terminal. If _stderr_ is true, the shell's stderr is copied in real time to the terminal.

## [Function] BashStreamLabel

Parameters:

- _cmd_ `string`
- _stdout_ `bool`
- _stderr_ `bool`
- _label_ `string`

Returns: `BashResult`

Runs a bash command in the same manner as _BashStream_ with the addition behavior of prepending _label_ to each line of the streamed stdout/stderr.
