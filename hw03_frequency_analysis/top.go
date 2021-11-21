package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type freqStruct struct {
	word   string
	amount int
}

func Top10(in string) []string {
	// Инициализируем карту
	freqMap := make(map[string]int)
	// Разделяем строку по пробелам и подсчитываем количество слов => map[слово] = количество повторений
	for _, value := range strings.Fields(in) {
		freqMap[value]++
	}
	// Создаём слайс, в котором будут храниться структуры freqstruct. Это необходимо для того, чтобы произвести сортировку
	freqSlice := make([]freqStruct, 0, len(freqMap))
	for key, value := range freqMap {
		freqSlice = append(freqSlice, freqStruct{key, value})
	}
	// Сортируем структуры по количеству повторений
	sort.Slice(freqSlice, func(i, j int) bool { return freqSlice[i].amount > freqSlice[j].amount })
	// Создаём карту, в которую размещаем слайсы структур, сгруппированных по количеству повторений
	repNum := make(map[int][]freqStruct)
	for _, value := range freqSlice {
		repNum[value.amount] = append(repNum[value.amount], value)
	}
	// Создаём слайс в котором будут хранится отсортированные значения повторений и лексиграфически сортируем слайсы
	keySlice := make([]int, 0, len(repNum))
	for amount := range repNum {
		keySlice = append(keySlice, amount)
		sort.Slice(repNum[amount], func(i, j int) bool { return repNum[amount][i].word < repNum[amount][j].word })
	}
	sort.Slice(keySlice, func(i, j int) bool { return keySlice[i] > keySlice[j] })
	// Создаём слайс, в котором будет храниться результат
	result := make([]string, 0, 10)
	num := 0
	for _, key := range keySlice {
		for _, value := range repNum[key] {
			result = append(result, value.word)
			num++
			if num >= 10 {
				break
			}
		}
		if num >= 10 {
			break
		}
	}
	return result
}
