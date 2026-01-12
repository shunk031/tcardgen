package text

import (
	"strings"
	
	budoux "github.com/soundkitchen/budoux/golang"
)

var (
	japaneseParser *budoux.Parser
)

func init() {
	japaneseParser = budoux.NewDefaultJapaneseParser()
}

// SegmentationLevel represents the granularity of text segmentation
type SegmentationLevel int

const (
	// DefaultLevel uses BudoX default segmentation
	DefaultLevel SegmentationLevel = iota
	// FineLevel provides finer segmentation by splitting long segments
	FineLevel
	// UltraFineLevel provides the finest segmentation by splitting even more
	UltraFineLevel
)

// SegmentForLineBreaks segments text using BudoX for better line breaking in Japanese text.
// It returns segments that can be used for more natural line breaking.
func SegmentForLineBreaks(text string) []string {
	return SegmentForLineBreaksWithLevel(text, DefaultLevel)
}

// SegmentForLineBreaksWithLevel segments text with specified granularity level
func SegmentForLineBreaksWithLevel(text string, level SegmentationLevel) []string {
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
	
	// Apply additional segmentation based on level
	switch level {
	case FineLevel:
		return applyFineSegmentation(result)
	case UltraFineLevel:
		return applyUltraFineSegmentation(result)
	default:
		return result
	}
}

// applyFineSegmentation splits segments that are longer than a threshold
func applyFineSegmentation(segments []string) []string {
	var result []string
	const maxSegmentLength = 6 // Characters
	
	for _, segment := range segments {
		if len([]rune(segment)) > maxSegmentLength {
			// Split long segments at natural boundaries
			subSegments := splitLongSegment(segment, maxSegmentLength)
			result = append(result, subSegments...)
		} else {
			result = append(result, segment)
		}
	}
	
	return result
}

// applyUltraFineSegmentation provides very fine segmentation
func applyUltraFineSegmentation(segments []string) []string {
	var result []string
	const maxSegmentLength = 3 // Characters
	
	for _, segment := range segments {
		if len([]rune(segment)) > maxSegmentLength {
			subSegments := splitLongSegment(segment, maxSegmentLength)
			result = append(result, subSegments...)
		} else {
			result = append(result, segment)
		}
	}
	
	return result
}

// splitLongSegment splits a long segment into smaller pieces
func splitLongSegment(segment string, maxLength int) []string {
	runes := []rune(segment)
	var result []string
	
	for i := 0; i < len(runes); i += maxLength {
		end := i + maxLength
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}
	
	return result
}

// ParseSegmentationLevel converts string configuration to SegmentationLevel
func ParseSegmentationLevel(levelStr string) SegmentationLevel {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "fine":
		return FineLevel
	case "ultra-fine", "ultrafine":
		return UltraFineLevel
	default:
		return DefaultLevel
	}
}

