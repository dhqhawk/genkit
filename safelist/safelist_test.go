package safelist

import (
	"fmt"
	"sync"
	"testing"
)

func TestSafeList_PushBackAll(t *testing.T) {
	list := NewSafeList()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			list.PushBackAll(1, 2, 3)
		}()
		go func() {
			defer wg.Done()
			list.PopBack()
		}()
	}
	wg.Wait()
	fmt.Println(list.Len())
}
