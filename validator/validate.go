package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func Validate(inputAsText string) ([]string, error) {
	var stringsToMatch []string

	if len(inputAsText) == 0 {
		return nil, errors.New("input is empty")
	}

	for i, stringToMatch := range strings.Split(inputAsText, ",") {
		cleanInput := clearString(stringToMatch)

		fmt.Printf("Word %d: %s \n", i, cleanInput)

		stringsToMatch = append(stringsToMatch, clearString(stringToMatch))
	}

	return stringsToMatch, nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
