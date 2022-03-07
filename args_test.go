package gocli

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func alphabet() (a []string) {
	for i := 0; i < 26; i++ {
		a = append(a, strings.Trim(fmt.Sprintf("%q", rune('A'+i)), "'"))
	}
	for i := 32; i < 26+32; i++ {
		a = append(a, strings.Trim(fmt.Sprintf("%q", rune('A'+i)), "'"))
	}
	return
}

func TestIsShortFlag(t *testing.T) {
	for _, l := range alphabet() {
		if !isShortFlag(fmt.Sprintf("-%s", l)) {
			t.Errorf("isShortFlag(\"-%s\") returned false", l)
		}
	}
	for i := 0; i < 10; i++ {
		if isShortFlag("-" + fmt.Sprint(i) + "a") {
			t.Errorf("isShortFlag(\"-%da\") returned true", i)
		}
	}
	if isShortFlag("--o") {
		t.Errorf("isShortFlag(\"--o\") returned true")
	}
	if isShortFlag("o") {
		t.Errorf("isShortFlag(\"o\") returned true")
	}
	if isShortFlag("-example") {
		t.Errorf("isShortFlag(\"-example\") returned true")
	}
	if !isShortFlag("-e=2etnd") {
		t.Errorf("isShortFlag(\"-e=2etnd\") returned false")
	}
}

func TestIsLongFlag(t *testing.T) {
	if !isLongFlag("--etd") {
		t.Errorf("isLongFlag(\"--etd\") returned false")
	}
	if !isLongFlag("--etd=test") {
		t.Errorf("isLongFlag(\"--etd=test\") returned false")
	}
	if isLongFlag("-etd=test") {
		t.Errorf("isLongFlag(\"-etd=test\") returned true")
	}
	if isLongFlag("--2etd=test") {
		t.Errorf("isLongFlag(\"-etd=test\") returned true")
	}
	if isLongFlag("---test") {
		t.Errorf("isLongFlag(\"---test\") returned true")
	}
}

func TestParseFlag(t *testing.T) {
	v1, v2 := parseFlag("test=value")
	if v1 != "test" || v2 != "value" {
		t.Errorf("parseFlag(\"test=value\") did not return \"test\", \"value\"")
	}
	v1, v2 = parseFlag("test==value")
	if v1 != "test" || v2 != "=value" {
		t.Errorf("parseFlag(\"test==value\") did not return \"test\", \"=value\"")
	}
}

func TestMatchShort(t *testing.T) {
	// should match
	options := []Option{{
		Short: "x",
	},
		{Short: "y"}}

	matched, success := matchShort("x", options)
	if matched.Short != "x" || !success {
		t.Errorf("matchShort failed to find the match")
	}

	// should not match
	options = []Option{{
		Short: "x",
	},
		{Short: "y"}}

	matched, success = matchShort("z", options)
	if matched.Short != "" || success {
		t.Errorf("matchShort return success when it shouldn't have")
	}
}

func TestMatchLong(t *testing.T) {
	// should match
	options := []Option{{
		Long: "x",
	},
		{Long: "y"}}

	matched, success := matchLong("x", options)
	if matched.Long != "x" || !success {
		t.Errorf("matchLong failed to find the match")
	}

	// should not match
	options = []Option{{
		Long: "x",
	},
		{Long: "y"}}

	matched, success = matchLong("z", options)
	if matched.Long != "" || success {
		t.Errorf("matchLong return success when it shouldn't have")
	}
}

func TestFirstCastValue(t *testing.T) {
	// Boolean, empty string (res == true)
	option := Option{
		Type: "bool",
	}
	value := ""
	res, err := firstCastValue(option, value)
	if err != nil || res != true {
		t.Errorf("firstCastValue failed: [Boolean, empty string (res == true)]. Result = %v, Error = %s", res, err)
	}

	// Boolean, nonempty string (err != nil)
	value = "example"
	res, err = firstCastValue(option, value)
	if err == nil {
		t.Errorf("firstCastValue failed: [Boolean, nonempty string (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// String, empty string (res == nil)
	option.Type = "string"
	value = ""
	res, err = firstCastValue(option, value)
	if err != nil || res != nil {
		t.Errorf("firstCastValue failed: [String, empty string (res == nil)]. Result = %v, Error = %s", res, err)
	}

	// String, nonempty string (res == value)
	option.Type = "string"
	value = "test"
	res, err = firstCastValue(option, value)
	if err != nil || res != res {
		t.Errorf("firstCastValue failed: [String, nonempty string (res == value)]. Result = %v, Error = %s", res, err)
	}

	// int, empty string (res == nil)
	option.Type = "int"
	value = ""
	res, err = firstCastValue(option, value)
	if err != nil || res != nil {
		t.Errorf("firstCastValue failed: [int, empty string (res == nil)]. Result = %v, Error = %s", res, err)
	}

	// int, non-int string (err != nil)
	option.Type = "int"
	value = "test"
	res, err = firstCastValue(option, value)
	if err == nil {
		t.Errorf("firstCastValue failed: [int, non-int string (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// int, int-convertable string (res == -12)
	option.Type = "int"
	value = "-12"
	res, err = firstCastValue(option, value)
	if err != nil || res != -12 {
		t.Errorf("firstCastValue failed: [int, int-convertable string (res == -12)]. Result = %v, Error = %s", res, err)
	}

	// float, empty string (res == nil)
	option.Type = "float"
	value = ""
	res, err = firstCastValue(option, value)
	if err != nil || res != nil {
		t.Errorf("firstCastValue failed: [float, empty string (res == nil)]. Result = %v, Error = %s", res, err)
	}

	// float, non-float string (err != nil)
	option.Type = "float"
	value = "test"
	res, err = firstCastValue(option, value)
	if err == nil {
		t.Errorf("firstCastValue failed: [float, non-float string (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// float, float-convertable string (res == -12.5)
	option.Type = "float"
	value = "-12.5"
	res, err = firstCastValue(option, value)
	if err != nil || res != -12.5 {
		t.Errorf("firstCastValue failed: [float, float-convertable string (res == -12.5)]. Result = %v, Error = %s", res, err)
	}

	// invalid type, any string (err != nil)
	option.Type = "test"
	value = "test"
	res, err = firstCastValue(option, value)
	if err == nil {
		t.Errorf("firstCastValue failed: [invalid type, any string (err != nil)]. Result = %v, Error = %s", res, err)
	}
}

func TestSecondCastValue(t *testing.T) {
	// bool, nil (res == false)
	option := Option{
		Type: "bool",
	}
	var value interface{}
	res, err := secondCastValue(option, value)
	if err != nil || res != false {
		t.Errorf("secondCastValue failed: [bool, nil (res == false)]. Result = %v, Error = %s", res, err)
	}

	// bool, any value (res == true)
	option.Type = "bool"
	value = "test"
	res, err = secondCastValue(option, value)
	if err != nil || res != true {
		t.Errorf("secondCastValue failed: [bool, any value (res == true)]. Result = %v, Error = %s", res, err)
	}

	// string, nil (err != nil)
	option.Type = "string"
	value = nil
	res, err = secondCastValue(option, value)
	if err == nil {
		t.Errorf("secondCastValue failed: [string, nil (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// string, string (res == value)
	option.Type = "string"
	value = "test"
	res, err = secondCastValue(option, value)
	if err != nil || res != value {
		t.Errorf("secondCastValue failed: [string, string (res == value)]. Result = %v, Error = %s", res, err)
	}

	// int, nil (err != nil)
	option.Type = "int"
	value = nil
	res, err = secondCastValue(option, value)
	if err == nil {
		t.Errorf("secondCastValue failed: [int, nil (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// int, some int (res == value)
	option.Type = "int"
	value = 12
	res, err = secondCastValue(option, value)
	if err != nil || res != value {
		t.Errorf("secondCastValue failed: [int, some int (res == value)]. Result = %v, Error = %s", res, err)
	}

	// float, nil (err != nil)
	option.Type = "float"
	value = nil
	res, err = secondCastValue(option, value)
	if err == nil {
		t.Errorf("secondCastValue failed: [float, nil (err != nil)]. Result = %v, Error = %s", res, err)
	}

	// float, some float (res == value)
	option.Type = "float"
	value = 12.5
	res, err = secondCastValue(option, value)
	if err != nil || res != value {
		t.Errorf("secondCastValue failed: [float, some float (res == value)]. Result = %v, Error = %s", res, err)
	}
}

func TestNoValue(t *testing.T) {
	// bool, no value (noValue() == false)
	o := Option{
		Type: "bool",
	}
	m := matchedOption{
		option: o,
	}
	if noValue(m) {
		t.Errorf("noValue failed: [bool, no value (noValue() == false)].")
	}

	// bool, with value (noValue() == false)
	m.option.Type = "bool"
	m.value = "true"
	if noValue(m) {
		t.Errorf("noValue failed: [bool, with value (noValue() == false)].")
	}

	// non-bool type, no value (noValue() == true)
	m.option.Type = "test"
	m.value = ""
	if !noValue(m) {
		t.Errorf("noValue failed: [non-bool type, no value (noValue() == true)].")
	}

	// non-bool type, no value (noValue() == false)
	m.option.Type = "test"
	m.value = "test"
	if noValue(m) {
		t.Errorf("noValue failed: [non-bool type, no value (noValue() == false)].")
	}
}

func TestShortMatchedOption(t *testing.T) {
	// matched
	short := "-t=test"
	o1 := Option{Short: "x", Type: "string"}
	o2 := Option{Short: "t", Type: "string"}
	options := []Option{o1, o2}
	m, err := shortMatchedOption(short, options)
	expectedM := matchedOption{option: o2, flag: "-t", value: "test", casted: "test"}
	if !reflect.DeepEqual(m, expectedM) || err != nil {
		t.Errorf("shortMatchedOption failed to match. matched: %+v, err: %s\n", m, err)
	}

	// matched no value
	short = "-t"
	o1 = Option{Short: "x", Type: "string"}
	o2 = Option{Short: "t", Type: "string"}
	options = []Option{o1, o2}
	m, err = shortMatchedOption(short, options)
	expectedM = matchedOption{option: o2, flag: "-t", value: "", casted: nil}
	if !reflect.DeepEqual(m, expectedM) || err != nil {
		t.Errorf("shortMatchedOption failed to match no value. matched: %+v, err: %s\n", m, err)
	}

	// not matched
	short = "-t=test"
	o1 = Option{Short: "x", Type: "string"}
	o2 = Option{Short: "z", Type: "string"}
	options = []Option{o1, o2}
	m, err = shortMatchedOption(short, options)
	if err == nil {
		t.Errorf("shortMatchedOption match when it shouldn't have. matched: %+v, err: %s\n", m, err)
	}
}

func TestLongMatchedOption(t *testing.T) {
	// matched
	long := "--t=2"
	o1 := Option{Long: "x", Type: "int"}
	o2 := Option{Long: "t", Type: "int"}
	options := []Option{o1, o2}

	m, err := longMatchedOption(long, options)
	expectedM := matchedOption{option: o2, flag: "--t", value: "2", casted: 2}
	if !reflect.DeepEqual(m, expectedM) || err != nil {
		t.Errorf("longMatchedOption failed to match. matched: %+v, err: %s\n", m, err)
	}

	// matched no value
	long = "--t"
	o1 = Option{Long: "x", Type: "int"}
	o2 = Option{Long: "t", Type: "int"}
	options = []Option{o1, o2}

	m, err = longMatchedOption(long, options)
	expectedM = matchedOption{option: o2, flag: "--t", value: "", casted: nil}
	if !reflect.DeepEqual(m, expectedM) || err != nil {
		t.Errorf("longMatchedOption failed to match no value. matched: %+v, err: %s\n", m, err)
	}

	// not matched
	long = "--t=test"
	o1 = Option{Long: "x", Type: "string"}
	o2 = Option{Long: "z", Type: "string"}
	options = []Option{o1, o2}

	m, err = longMatchedOption(long, options)
	if err == nil {
		t.Errorf("longMatchedOption match when it shouldn't have. matched: %+v, err: %s\n", m, err)
	}
}
