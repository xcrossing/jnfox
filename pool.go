package main

import "sync"

type pool struct {
	wg sync.WaitGroup
	ch chan string
}

func makePool(threadCount int, fn func(str string)) *pool {
	p := new(pool)
	p.wg.Add(threadCount)
	p.ch = make(chan string)

	for thread := 0; thread < threadCount; thread++ {
		go func() {
			for {
				str, ok := <-p.ch
				if !ok {
					break
				}
				fn(str)
			}
			p.wg.Done()
		}()
	}

	return p
}

func (p *pool) add(str string) {
	p.ch <- str
}

func (p *pool) wait() {
	close(p.ch)
	p.wg.Wait()
}
