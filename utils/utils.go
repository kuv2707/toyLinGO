package utils

import (
	"strconv"
	"toylingo/globals"
)

func DoNothing(args ...interface{}) {}

func IsNumber(temp string) bool {
	for i := 0; i < len(temp); i++ {
		if temp[i] < '0' || temp[i] > '9' {
			if temp[i] != '.' {
				return false
			}
		}
	}
	return true
}

func StringToNumber(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return float64(num)
}

func StringToBoolean(str string) bool {
	if str == globals.Booleans[0] {
		return true
	} else if str == globals.Booleans[1] {
		return false
	} else {
		panic("invalid boolean value " + str)
	}
}


func IsOneOf(temp string, options []string) bool {
	for i := 0; i < len(options); i++ {
		if temp == (options[i]) {
			return true
		}
	}
	return false
}

func IsOneOfArr(str string, options []string) bool {
	for _, word := range options {
		if word == str {

			return true
		}
	}
	return false
}

var QUOTES = []string{"`"}

//todo use regex
func InQuotes(s string) bool {
	return IsOneOf(s[0:1], QUOTES) && IsOneOf(s[len(s)-1:], QUOTES)
}

//todo use regex and add more constraints
func IsValidVariableName(s string) bool {
	return !InQuotes(s)
}

func StringVal(s interface{}) string {
	if s == nil {
		return ("nil")
	}
	switch s.(type) {
	case string:
		return (s.(string))
	case float64:
		a := s.(float64)
		//todo understand why this is needed
		return strconv.FormatFloat(float64(a), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(s.(bool))
	default:
		panic("unknown type")
	}
}

func ParseEscapeSequence(s string) string {
	ret := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			switch s[i+1] {
			case 'n':
				ret += "\n"
			case 't':
				ret += "\t"
			case 'r':
				ret += "\r"
			case 'b':
				ret += "\b"
			case 'f':
				ret += "\f"
			case '\\':
				ret += "\\"
			case '\'':
				ret += "`"
			case '`': //fixme: might be a bug
				ret += "\""
			default:
				ret += string(s[i+1])
			}
			i++
		} else {
			ret += string(s[i])
		}
	}
	return ret
}


func IsLiteral(s string) bool {
	return IsOneOf(s, []string{"NUMBER", "STRING_LITERAL", "BOOLEAN"})
}

func IsOperator(temp string) bool {
	return IsOneOf(temp, globals.Operators)
}

func IsBoolean(temp string) bool {
	return IsOneOf(temp, globals.Booleans)
}

func ClosingBracket(bracket string) string {
	switch bracket {
	case "(":
		return ")"
	case "[":
		return "]"
	case "{":
		return "}"
	default:
		panic("invalid bracket " + bracket)
	}
}

func IsOpenBracket(bracket string) bool {
	return IsOneOf(bracket, []string{"(", "[", "{"})
}
func IsCloseBracket(bracket string) bool {
	return IsOneOf(bracket, []string{")", "]", "}"})
}