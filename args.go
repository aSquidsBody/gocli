package gocli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// return true if a short flag, e.g. "-g" but not "--global"
func isShortFlag(s string) bool {
	matched, err := regexp.Match("^-[a-z || A-Z](=.*)*$", []byte(s))
	if err != nil {
		panic(fmt.Sprintf("isShortFlag failed with error %s", err))
	}
	return matched
}

// return true if a long flag, e.g. "--global" but not "-g"
func isLongFlag(s string) bool {
	b := []byte(s)

	tooLong, err := regexp.Match("^---.*$", b)
	if err != nil {
		panic(fmt.Sprintf("isLongFlag failed with error %s", err))
	}
	if tooLong {
		return false
	}

	numStart, err := regexp.Match("^--[0-9].*$", b)
	if err != nil {
		panic(fmt.Sprintf("isLongFlag failed with error %s", err))
	}
	if numStart {
		return false
	}

	matched, err := regexp.Match("^--[a-z || A-Z || 0-9 || -]+(=.*)*$", []byte(s))
	if err != nil {
		panic(fmt.Sprintf("isLongFlag failed with error %s", err))
	}
	return matched
}

func parseLong(f string) (string, string) {
	return parseFlag(f[2:])
}

func parseShort(f string) (string, string) {
	return parseFlag(f[1:])
}

// return flag name and a value associated with it
func parseFlag(f string) (string, string) {
	idx := strings.Index(f, "=")
	if idx == -1 {
		return f, ""
	}
	return f[:idx], f[idx+1:]
}

// Return the option which matches the flag
func matchShort(flag string, options []Option) (Option, bool) {
	for _, opt := range options {
		if flag == opt.Short {
			return opt, true
		}
	}
	return Option{}, false
}

// return the option which matches the flag
func matchLong(flag string, options []Option) (Option, bool) {
	for _, opt := range options {
		if flag == opt.Long {
			return opt, true
		}
	}
	return Option{}, false
}

func emptyOption(o Option) bool {
	return o.Short == "" && o.Long == "" && o.Description == "" && o.Type == "" && o.Required == false
}

// allows for empty int/float values
func firstCastValue(option Option, value string) (v interface{}, err error) {
	switch option.Type {
	case "bool":
		if value == "" {
			v = true // false booleans are handled after all other args are processed (outside this function)
		} else {
			err = fmt.Errorf("Received invalid value: %s", value)
		}
	case "string":
		v = value
		if value == "" {
			v = nil
		}

	case "int":
		if value == "" {
			v = nil
		} else {
			var i int
			i, err = strconv.Atoi(value)
			if err != nil {
				err = fmt.Errorf("Received invalid int value: %s", value)
			} else {
				v = i
			}
		}

	case "float":
		if value == "" {
			v = nil
		} else {
			var f float64
			f, err = strconv.ParseFloat(value, 64)
			if err != nil {
				err = fmt.Errorf("Received invalid float value: %s", value)
			} else {
				v = f
			}
		}

	default:
		err = fmt.Errorf("Received invalid option configuration: '%s' is configured to receive type '%s'", option.Name(), option.Type)
	}
	return
}

func secondCastValue(option Option, value interface{}) (v interface{}, err error) {
	switch option.Type {
	case "bool":
		// value should be nil or true
		if value == nil {
			v = false
		} else {
			v = true
		}
	case "string":
		// value is a string or nil
		if value == nil {
			err = fmt.Errorf("Received empty string value")
		} else {
			v = value
		}
	case "int":
		// value is an int or nil
		if value == nil {
			err = fmt.Errorf("Received invalid int value")
		} else {
			v = value
		}
	case "float":
		if value == nil {
			err = fmt.Errorf("Received invalid float value")
		} else {
			v = value
		}
	default:
		err = fmt.Errorf("Received invalid option configuration: '%s' is configured to receive type '%s'", option.Name(), option.Type)
	}

	return
}

type matchedOption struct {
	option Option
	flag   string
	value  string
	casted interface{}
}

func noValue(m matchedOption) (b bool) {
	switch m.option.Type {
	case "bool":
		b = false
	default:
		b = m.value == ""
	}

	return
}

// create a matchedOption if the option matches the arg, if not, return an error
func shortMatchedOption(short string, options []Option) (m matchedOption, err error) {
	name, value := parseShort(short)
	flag := "-" + name

	if opt, ok := matchShort(name, options); !ok {
		// if there is no match, return an error
		err = fmt.Errorf("Unexpected option `%s`", flag)
	} else {
		// if there is a match, firstCast the value and return the matchedOption
		var casted interface{}
		casted, err = firstCastValue(opt, value)
		if err != nil {
			err = fmt.Errorf("Error parsing `%s`: %s", flag, err)
			return
		}
		m = matchedOption{
			option: opt,
			flag:   flag,
			value:  value,
			casted: casted,
		}
	}
	return
}

// create a matchedOption if the option matches the arg, if not, return an error
func longMatchedOption(long string, options []Option) (m matchedOption, err error) {
	name, value := parseLong(long)
	flag := "--" + name

	if opt, ok := matchLong(name, options); !ok {
		// if there is no match, return an error
		err = fmt.Errorf("Unexpected option `%s`", flag)
	} else {
		// if there is a match, firstCast the value and return the matchedOption
		var casted interface{}
		casted, err = firstCastValue(opt, value)
		if err != nil {
			err = fmt.Errorf("Error parsing `%s`: %s", flag, err)
			return
		}
		m = matchedOption{
			option: opt,
			flag:   flag,
			value:  value,
			casted: casted,
		}
	}

	return
}

// match options with cli flags and preform the first cast
func firstPass(options []Option, argDef Argument, args []string) (map[string]interface{}, map[Option]matchedOption, error) {
	result := map[string]interface{}{}
	resultOpt := map[Option]matchedOption{}

	for _, opt := range options {
		resultOpt[opt] = matchedOption{casted: nil}
		if opt.Short != "" {
			result[opt.Short] = nil
		}
		if opt.Long != "" {
			result[opt.Long] = nil
		}
	}

	if argDef.Name != "" {
		result[argDef.Name] = nil
	}

	var prev Option
	for _, arg := range args {
		if isShortFlag(arg) {
			var matched matchedOption
			var err error
			matched, err = shortMatchedOption(arg, options)
			if err != nil {
				return result, resultOpt, err
			}
			prev = matched.option
			if resultOpt[matched.option].flag != "" {
				return result, resultOpt, fmt.Errorf("Option entered twice `%s`", matched.option.Name())
			}
			resultOpt[matched.option] = matched

		} else if isLongFlag(arg) {
			var matched matchedOption
			var err error
			matched, err = longMatchedOption(arg, options)
			if err != nil {
				return result, resultOpt, err
			}
			prev = matched.option
			if resultOpt[matched.option].flag != "" {
				return result, resultOpt, fmt.Errorf("Option entered twice `%s`", matched.option.Name())
			}
			resultOpt[matched.option] = matched

		} else if !emptyOption(prev) && noValue(resultOpt[prev]) {
			// this block runs if the prev arg was a flag and the value for the flag is empty
			// assume this arg is the value for the previous flag
			prevMatched := resultOpt[prev]
			flag := prevMatched.flag

			// cast the value
			casted, err := firstCastValue(prev, arg)
			if err != nil {
				return result, resultOpt, fmt.Errorf("Error parsing `%s`: %s", flag, err)
			}

			// save the casted value
			prevMatched.value = arg
			prevMatched.casted = casted

			resultOpt[prev] = prevMatched
			prev = Option{}

		} else {
			// this block runs if the previous arg was not an empty flag
			// assume this arg is an expected arg for the command

			// check if the arg is already defined or if not expected
			if v := result[argDef.Name]; v != nil || argDef.Name == "" {
				return result, resultOpt, fmt.Errorf("Received unknown argument '%s'", arg)
			}

			result[argDef.Name] = arg
			prev = Option{}

		}
	}

	// check that required argument has value
	if result[argDef.Name] == nil && argDef.Required {
		return result, resultOpt, fmt.Errorf("Missing required argument '%s'.", argDef.Name)
	}

	return result, resultOpt, nil
}

// read in string cli args and parse them
func ParseArgs(options []Option, argDef Argument, args []string) (map[string]interface{}, error) {

	result, resultOpt, err := firstPass(options, argDef, args)
	if err != nil {
		return result, err
	}

	// fmt.Println(fmt.Sprintf("Result Opt: %+v", resultOpt))
	// fmt.Println(fmt.Sprintf("Result: %+v", result))

	// check that all the required options have values
	missing := []string{}
	for _, opt := range options {
		matched := resultOpt[opt]
		// noValue always returns false for booleans
		if noValue(matched) && opt.Required {
			// required values that were not matched
			missing = append(missing, opt.Name())
		}
	}

	if len(missing) > 0 {
		return result, fmt.Errorf("The following options are missing or empty: '%s'.", strings.Join(missing, "', '"))
	}

	// second cast the matched options as they are placed into 'result'
	for _, opt := range options {
		matched := resultOpt[opt]
		casted, err := secondCastValue(opt, matched.casted)
		if err != nil && matched.flag != "" {
			return result, fmt.Errorf("Error parsing `%s`: %s", matched.flag, err)
		} else {
			if opt.Long != "" {
				result[opt.Long] = casted
			}
			if opt.Short != "" {
				result[opt.Short] = casted
			}

		}
	}

	return result, nil
}
