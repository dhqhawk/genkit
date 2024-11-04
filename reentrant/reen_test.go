package reentrant

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

type mutexT struct {
	lock ReentrantLock
}

func (t *mutexT) foo() {
	fmt.Println("in foo")
}

func (t *mutexT) bar() {
	fmt.Println("in bar")
}

func TestReentrant(t *testing.T) {
	mt := &mutexT{}
	testCases := []struct {
		name        string
		testfunc    func(*mutexT)
		wantrecover string
	}{
		{
			name: "ok",
			testfunc: func(mt *mutexT) {
				mt.lock.Lock()
				defer mt.lock.UnLock()
				mt.foo()
				mt.lock.Lock()
				defer mt.lock.UnLock()
				mt.bar()
			},
		},
		{
			name: "别的协程抢锁",
			testfunc: func(mt *mutexT) {
				mt.lock.Lock()
				defer mt.lock.UnLock()
				mt.foo()
				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					//这里即使下面panic，也会运行
					//defer mt.lock.UnLock()
					defer wg.Done()
					mt.lock.Lock()
					mt.bar()
					mt.lock.UnLock()
				}()
				wg.Wait()
			},
			wantrecover: "lock is already locked",
		},
		{
			name: "别的协程解锁",
			testfunc: func(mt *mutexT) {
				mt.lock.Lock()
				defer mt.lock.UnLock()
				mt.foo()
				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer mt.lock.UnLock()
					defer wg.Done()
					mt.lock.Lock()
					mt.bar()
				}()
				wg.Wait()
			},
			wantrecover: "unlock by other goroutine",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				//没有捕获到，因为是其他协程panic
				if r := recover(); r != nil {
					require.Equal(t, tc.wantrecover, r)
					t.Errorf("Test case %s panicked: %v", tc.name, r)
				}
			}()
			tc.testfunc(mt)
		})
	}
}
