package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	AllowedChars = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

func Validate(inputAsText string) ([]string, error) {
	var (
		stringsToMatch []string
		err            error
	)

	if len(inputAsText) == 0 {
		return nil, errors.New("input is empty")
	}

	for i, stringToMatch := range strings.Split(inputAsText, ",") {
		cleanInput := clearString(stringToMatch)

		err = containsDisallowedChars(cleanInput)
		if err != nil {
			return stringsToMatch, err
		}

		fmt.Printf("Word %d: %s \n", i, cleanInput)

		stringsToMatch = append(stringsToMatch, clearString(stringToMatch))
	}

	return stringsToMatch, nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func containsDisallowedChars(word string) error {
	for _, char := range strings.Split(word, "") {
		if !strings.Contains(AllowedChars, char) {
			return errors.New(fmt.Sprintf("Character %s not allowed in word %s. Try it in leetspeak for example by changing o to 0 (zero)!", char, word))
		}
	}

	return nil
}

func ParseYn(answer string) bool {
	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)

	return answer == "y" || answer == "yes"
}
