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
		cond1 := freqMap[uniqueWords[i]] > freqMap[uniqueWords[j]]
		cond2 := freqMap[uniqueWords[i]] == freqMap[uniqueWords[j]] && uniqueWords[i] < uniqueWords[j]
		if cond1 || cond2 {
			return true
		}
		return false
	})
	switch {
	case len(uniqueWords) >= 10:
		return uniqueWords[0:10]
	case len(uniqueWords) > 0:
		return uniqueWords
	default:
		return nil
	}
}
