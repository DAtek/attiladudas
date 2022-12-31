package helpers

import "sync"

type waitGroup struct {
	wg         *sync.WaitGroup
	goroutines []func()
}

func NewWaitGroup() *waitGroup {
	return &waitGroup{
		wg:         &sync.WaitGroup{},
		goroutines: []func(){},
	}
}

func (wg *waitGroup) Add(f func()) {
	wg.wg.Add(1)
	wg.goroutines = append(wg.goroutines, f)
}

func (wg *waitGroup) Wait() {
	for _, goro := range wg.goroutines {
		f := goro

		go func() {
			defer wg.wg.Done()
			f()
		}()
	}

	wg.wg.Wait()
	wg.goroutines = []func(){}
}
