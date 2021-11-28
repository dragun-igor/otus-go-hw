package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Количество стейджей
	stageAmount := len(stages)
	// Создаём слайс каналов для хранения промежуточных и конечного варианта
	results := make([]Bi, 0, stageAmount+1)
	for i := 0; i < len(stages)+1; i++ {
		results = append(results, make(Bi))
	}
	// Остановка расчётов по каналу done
	// Не знаю, правильно ли так делать, ибо могут остаться работающие горутины и открытые каналы
	// Можно добавить закрытие всех каналов, а не только последнего
	go func() {
		<-done
		close(results[stageAmount])
	}()
	// Переброска значений в первую ячейку слайса, после окончания закрываем канал
	go func() {
		defer close(results[0])
		for val := range in {
			results[0] <- val
		}
	}()
	// Создание горутин со стейджами
	for i, stage := range stages {
		go func(stage Stage, i int) {
			defer close(results[i+1])
			for val := range stage(results[i]) {
				results[i+1] <- val
			}
		}(stage, i)
	}
	// Возвращаем канал из последней ячейки слайса
	return results[stageAmount]
}
