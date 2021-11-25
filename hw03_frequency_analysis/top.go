package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(in string) []string {
	result := make([]string, 0, 10)
	// Инициализируем карту
	freqMap := make(map[string]int)
	// Разделяем строку по пробелам и подсчитываем количество слов => map[слово] = количество повторений
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
	if len(uniqueWords) >= 10 {
		result = uniqueWords[0:10]
	}
	if len(uniqueWords) < 10 {
		result = uniqueWords
	}
	if len(uniqueWords) == 0 {
		result = nil
	}
	return result
}
