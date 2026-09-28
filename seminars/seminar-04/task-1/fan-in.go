package fanin

import "sync"

//  Реализовать метод, который сливает данные из нескольких каналов типа int в один канал.

func FanIn[T any](chans []<-chan T) <-chan T {
	res := make(chan T)
	wg := &sync.WaitGroup{}

	go func() {
		for _, ch := range chans {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					res <- v
				}
			}()
		}
		wg.Wait()
		close(res)
	}()

	return res
}
