package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(in string) []string {
	freqMap := make(map[string]int)
	for _, value := range strings.Fields(in) {
		freqMap[value]++
	}
	uniqueWords := make([]string, 0, len(freqMap))
	for key := range freqMap {
		uniqueWords = append(uniqueWords, key)
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		firstWord := uniqueWords[i]
		secondWord := uniqueWords[j]
		if freqMap[firstWord] > freqMap[secondWord] {
			return true
		}
		if freqMap[firstWord] == freqMap[secondWord] && firstWord < secondWord {
			return true
		}
		return false
	})
	switch {
	case len(uniqueWords) == 0:
		return nil
	case len(uniqueWords) < 10:
		return uniqueWords
	default:
		return uniqueWords[0:10]
	}
}
