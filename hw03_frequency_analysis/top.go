package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type freqstruct struct {
	word   string
	amount int
}

func Top10(in string) []string {
	// Инициализируем карту
	freqmap := make(map[string]int)
	// Разделяем строку по пробелам и подсчитываем количество слов => map[слово] = количество повторений
	for _, value := range strings.Fields(in) {
		freqmap[value]++
	}
	// Создаём слайс, в котором будут храниться структуры freqstruct. Это необходимо для того, чтобы произвести сортировку
	freqslice := make([]freqstruct, 0, len(freqmap))
	for key, value := range freqmap {
		freqslice = append(freqslice, freqstruct{key, value})
	}
	// Сортируем структуры по количеству повторений
	sort.Slice(freqslice, func(i, j int) bool { return freqslice[i].amount > freqslice[j].amount })
	// Создаём карту, в которую размещаем слайсы структур, сгруппированных по количеству повторений
	repnum := make(map[int][]freqstruct)
	for _, value := range freqslice {
		repnum[value.amount] = append(repnum[value.amount], value)
	}
	// Создаём слайс в котором будут хранится отсортированные значения повторений и лексиграфически сортируем слайсы
	keyslice := make([]int, 0, len(repnum))
	for amount := range repnum {
		keyslice = append(keyslice, amount)
		sort.Slice(repnum[amount], func(i, j int) bool { return repnum[amount][i].word < repnum[amount][j].word })
	}
	sort.Slice(keyslice, func(i, j int) bool { return keyslice[i] > keyslice[j] })
	// Создаём слайс, в котором будет храниться результат
	result := make([]string, 0, 10)
	num := 0
	for _, key := range keyslice {
		for _, value := range repnum[key] {
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
