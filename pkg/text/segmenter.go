package text

import (
	budoux "github.com/soundkitchen/budoux/golang"
)

var (
	japaneseParser *budoux.Parser
)

func init() {
	japaneseParser = budoux.NewDefaultJapaneseParser()
}

// SegmentForLineBreaks segments text using BudoX for better line breaking in Japanese text.
// It returns segments that can be used for more natural line breaking.
func SegmentForLineBreaks(text string) []string {
	if text == "" {
		return []string{}
	}

	// Use BudoX to segment Japanese text
	segments := japaneseParser.Parse(text)
	
	// Filter out completely empty segments, but preserve internal spaces
	var result []string
	for _, segment := range segments {
		if segment != "" {
			result = append(result, segment)
		}
	}
	
	return result
}

