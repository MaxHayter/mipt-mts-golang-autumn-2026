package main

import "fmt"

// Если бы select выбирал первый совпавший case по порядку в коде, как
// if/else, один из двух постоянно готовых источников голодал бы. Внутри
// select проверяет кейсы в псевдослучайном порядке — эта программа считает
// частоту выбора a и b, чтобы увидеть это не на словах, а на числах.
func main() {
	a := make(chan struct{}, 1)
	b := make(chan struct{}, 1)

	counts := map[string]int{}
	const rounds = 100000

	for range rounds {
		a <- struct{}{}
		b <- struct{}{}

		select {
		case <-a:
			counts["a"]++
			<-b
		case <-b:
			counts["b"]++
			<-a
		}
	}

	fmt.Printf("оба канала были готовы одновременно %d раз\n", rounds)
	fmt.Printf("select выбрал a: %d раз (%.1f%%)\n", counts["a"], 100*float64(counts["a"])/rounds)
	fmt.Printf("select выбрал b: %d раз (%.1f%%)\n", counts["b"], 100*float64(counts["b"])/rounds)
}
