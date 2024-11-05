package timeLock

import (
	"context"
	"sync"
	"time"
)

type TimeLock struct {
	lock    sync.Mutex
	timeout int
}

func NewTimeLock(t int) *TimeLock {
	//time.Second*t不行是因为他两不是一个类型的，*10可以是因为10是常量整形，可以隐式转化为任何整数类
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//ctx, cancel := context.WithTimeout(context.Background(), t*time.Second)
	return &TimeLock{
		timeout: t,
	}
}

func (t *TimeLock) Lock() bool {
	ctx, cancle := context.WithTimeout(context.Background(), time.Duration(t.timeout)*time.Second)
	defer cancle()
	done := make(chan struct{})
	go func() {
		t.lock.Lock()
		close(done)
	}()
	select {
	case <-ctx.Done():
		defer t.lock.Unlock()
		return false
	case <-done:
		return true
	}
}

func (t *TimeLock) Unlock() {
	t.lock.Unlock()
}
