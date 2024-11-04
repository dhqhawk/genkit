package reentrant

import (
	"github.com/kortschak/goroutine"
	"sync"
	"sync/atomic"
)

//同一协程可重入锁

type ReentrantLock struct {
	lock      sync.Mutex
	gid       int64
	recursion int64
}

func (r *ReentrantLock) Lock() {
	gid := goroutine.ID()
	if gid == atomic.LoadInt64(&r.gid) {
		atomic.AddInt64(&r.recursion, 1)
		return
	}
	//当其他协程不使用defer解锁触发
	if atomic.LoadInt64(&r.gid) != 0 {
		panic("lock by other goroutine")
	}
	//如果前面到这里,就说名之前没有gid，没有加锁
	r.lock.Lock()
	atomic.StoreInt64(&r.gid, gid)
}

func (r *ReentrantLock) UnLock() {
	gid := goroutine.ID()
	//当其他协程使用defer解锁触发
	if atomic.LoadInt64(&r.gid) != gid {
		panic("unlock by other goroutine")
	}
	if atomic.LoadInt64(&r.recursion) > 0 {
		atomic.AddInt64(&r.recursion, -1)
		return
	}
	//释放锁，需要重置gid
	atomic.StoreInt64(&r.gid, 0)
	r.lock.Unlock()
}
