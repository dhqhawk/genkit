package timeLock

import (
	"fmt"
	"testing"
	"time"
)

func Test_lockfunc(t *testing.T) {
	lock := NewTimeLock(2)
	var a int
	go func() {
		lock.Lock()
		a++
		time.Sleep(6 * time.Second)
		lock.Unlock()
	}()
	//time.Sleep(time.Millisecond)
	//因为下面这个比上面更快获得锁了，所以才会两次输出加锁成功
	lock.Lock()
	a++
	time.Sleep(3 * time.Second)
	lock.Unlock()
	time.Sleep(5 * time.Second)
	fmt.Println(a)
}

// 测试如果超时之后，是否加锁还在
func Test_lockagain(t *testing.T) {
	lock := NewTimeLock(1)
	var a int
	go func() {
		lock.Lock()
		a++
		time.Sleep(3 * time.Second)
		lock.Unlock()
	}()
	time.Sleep(time.Millisecond)
	//就是它还是会等锁，导致第3个加不上锁
	if lock.Lock() {
		defer lock.Unlock()
	} else {
		fmt.Println("没拿到锁啊")
	}
	time.Sleep(4 * time.Second)
	if lock.Lock() {
		defer lock.Unlock()
	} else {
		fmt.Println("没拿到锁啊")
	}
}
