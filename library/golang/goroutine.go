package golang

import (
	"sync"
)

/*
	wg = sync.WaitGroup{}
	golang.Goroutine(&wg, func() {
		fmt.Println(1)
	})
	wg.Wait()
*/
func Goroutine(w *sync.WaitGroup, f func()) {
	(*w).Add(1)
	go func() {
		defer (*w).Done()
		f()
	}()
}
