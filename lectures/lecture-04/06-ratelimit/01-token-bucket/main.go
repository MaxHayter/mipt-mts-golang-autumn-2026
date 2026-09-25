package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 2 запроса в секунду, разрешаем всплеск до 3 сразу
	limiter := rate.NewLimiter(rate.Limit(2), 3)

	start := time.Now()
	for i := range 15 {
		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Println("ошибка:", err)
			return
		}
		fmt.Printf("запрос %d прошёл в %v\n", i, time.Since(start).Round(time.Millisecond))

		if i%5 == 0 {
			time.Sleep(2*time.Second)
		}
	}
}
