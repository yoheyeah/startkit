package gins

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	validator "gopkg.in/go-playground/validator.v8"
)

func UpperCaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LowerCaseFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	for index, word := range entries {
		if index == 0 {
			entries[index] = UpperCaseFirst(word)
		} else {
			entries[index] = LowerCaseFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}

func ValidationErrorForFieldToText(e *validator.FieldError) string {
	word := Split(e.Field)
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param)
	}
	return fmt.Sprintf("%s is not valid", word)
}

func ValidationErrorToText(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return fmt.Sprintf("Field %s Is Required", e.Field)
	case "max":
		return fmt.Sprintf("Field %s Cannot Be Longer Than %s", e.Field, e.Param)
	case "min":
		return fmt.Sprintf("Field %s Must Be Longer Than %s", e.Field, e.Param)
	case "email":
		return fmt.Sprintf("Invalid Email Format")
	case "len":
		return fmt.Sprintf("Field %s Must Be %s Characters Long", e.Field, e.Param)
	}
	return fmt.Sprintf("Field %s Is Not Valid", e.Field)
}
