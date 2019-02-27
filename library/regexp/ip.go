package regexp

import (
	"regexp"
)

func IP(text string) []string {
	IPRegExp := regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	return IPRegExp.FindAllString(text, len(text))
}
