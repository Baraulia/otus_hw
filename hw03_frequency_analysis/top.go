package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var taskWithAsterisk bool
var wordPatternAnySymbol = regexp.MustCompile(`\S+`)

func Top10(enterString string) []string {
	taskWithAsterisk = true // change to true for the task with an asterisk
	var sliceStr []string
	uniqueWordsMap := make(map[string]int)
	switch taskWithAsterisk {
	case true:
		sliceStr = wordPatternAnySymbol.FindAllString(enterString, -1)
		for _, word := range sliceStr {
			if word == "-" {
				continue
			}

			word = cleanWord(word)
			word = strings.ToLower(word)
			uniqueWordsMap[word]++
		}
	default:
		sliceStr = strings.Fields(enterString)
		_ = wordPatternAnySymbol.FindAllString(enterString, -1) //can be used instead of the  strings.Fields()

		for _, word := range sliceStr {
			uniqueWordsMap[word]++
		}
	}

	uniqueSliceStr := make([]string, 0, len(uniqueWordsMap))
	for word := range uniqueWordsMap {
		uniqueSliceStr = append(uniqueSliceStr, word)
	}

	sort.Slice(uniqueSliceStr, func(i, j int) bool {
		var result bool
		countI := uniqueWordsMap[uniqueSliceStr[i]]
		countJ := uniqueWordsMap[uniqueSliceStr[j]]

		if countI == countJ {
			result = uniqueSliceStr[i] < uniqueSliceStr[j]
		} else {
			result = countI > countJ
		}

		return result
	})

	if len(uniqueSliceStr) > 10 {
		uniqueSliceStr = uniqueSliceStr[:10]
	}

	return uniqueSliceStr
}

func cleanWord(word string) string {
	cleanedWord := strings.TrimFunc(word, func(r rune) bool {
		return strings.ContainsRune(".,!;", r)
	})

	return cleanedWord
}
