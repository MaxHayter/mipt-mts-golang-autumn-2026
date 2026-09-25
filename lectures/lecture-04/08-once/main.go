package main

import (
	"fmt"
	"sync"
)

type config struct {
	value string
}

func main() {
	var once sync.Once
	var cfg *config
	initCount := 0

	loadConfig := func() {
		initCount++
		cfg = &config{value: "загружено из файла"}
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			once.Do(loadConfig)
			_ = cfg.value
		})
	}
	wg.Wait()

	fmt.Println("инициализация вызвана раз:", initCount)
	fmt.Println("итоговое значение:", cfg.value)
}
