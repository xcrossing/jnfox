package util

import (
	"sync"
)

type pool struct {
	wg sync.WaitGroup
	ch chan string
}

type funcPool struct {
	wg sync.WaitGroup
	ch chan func()
}

func MakePool(threadCount int, fn func(str string)) *pool {
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

func (p *pool) Add(str string) {
	p.ch <- str
}

func (p *pool) Wait() {
	close(p.ch)
	p.wg.Wait()
}

func MakeFuncPool(threadCount int) *funcPool {
	p := new(funcPool)
	p.wg.Add(threadCount)
	p.ch = make(chan func())

	for thread := 0; thread < threadCount; thread++ {
		go func() {
			for {
				fn, ok := <-p.ch
				if !ok {
					break
				}
				fn()
			}
			p.wg.Done()
		}()
	}

	return p
}

func (p *funcPool) Add(fn func()) {
	p.ch <- fn
}

func (p *funcPool) Wait() {
	close(p.ch)
	p.wg.Wait()
}
