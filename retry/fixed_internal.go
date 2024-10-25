package retry

import (
	"genkit/errs"
	"time"
)

const defaultmaxcnt int8 = 10
const defaultinterval time.Duration = time.Second * 3

type FixedRetry struct {
	Interval time.Duration
	cnt      int8
	maxcnt   int8
}

// Next 认为只要出现错误，用户就不应该使用我的重试机制了
func (f *FixedRetry) Next() error {
	f.cnt++
	if f.cnt > f.maxcnt {
		return errs.NewErrorBeyondMaxcnt()
	}
	time.Sleep(f.Interval)
	//fmt.Println(f.Interval)
	return nil
}

func NewFixedRetry() *FixedRetry {
	return &FixedRetry{
		Interval: defaultinterval,
		cnt:      0,
		maxcnt:   defaultmaxcnt,
	}
}

// WithMaxCnt 自定义最大重试的次数
func (f *FixedRetry) WithMaxCnt(maxcnt int8) *FixedRetry {
	if maxcnt <= 0 || maxcnt > 99 {
		maxcnt = defaultmaxcnt
	}
	f.maxcnt = maxcnt
	return f
}

// WithInterval 自定义重试的间隔时间
func (f *FixedRetry) WithInterval(interval time.Duration) *FixedRetry {
	if interval <= 0 {
		interval = defaultinterval
	}
	f.Interval = interval
	return f
}
