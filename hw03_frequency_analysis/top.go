package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(enterString string) []string {
	var sliceStr []string
	uniqueWordsMap := make(map[string]int)
	sliceStr = strings.Fields(enterString)

	for _, word := range sliceStr {
		uniqueWordsMap[word]++
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
