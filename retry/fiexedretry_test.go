package retry

import (
	"genkit/errs"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_FixedInterval(t *testing.T) {
	testCases := []struct {
		Name          string
		retrytimes    int
		maxcnt        int8
		retryInterval time.Duration
		wantError     error
	}{
		{
			Name:          "cnt < maxcnt,defaultmaxcnt",
			retrytimes:    3,
			retryInterval: time.Second,
			wantError:     nil,
		},
		{
			Name:          "超过最大重试",
			retrytimes:    7,
			maxcnt:        5,
			retryInterval: time.Second,
			wantError:     errs.NewErrorBeyondMaxcnt(),
		},
		{
			Name:          "错误的最大次数,但是正常运行",
			retrytimes:    7,
			maxcnt:        -1,
			retryInterval: time.Second,
			wantError:     nil,
		},
		{
			Name:          "错误的最大次数,超过重试次数",
			retrytimes:    11,
			maxcnt:        -1,
			retryInterval: time.Second,
			wantError:     errs.NewErrorBeyondMaxcnt(),
		},
		{
			Name:          "错误的时间",
			retrytimes:    11,
			maxcnt:        -1,
			retryInterval: time.Duration(0),
			wantError:     errs.NewErrorBeyondMaxcnt(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			retryhandler := NewFixedRetry()
			retryhandler.WithMaxCnt(tc.maxcnt).WithInterval(tc.retryInterval)
			var err error
			for i := 0; i < tc.retrytimes; i++ {
				err = retryhandler.Next()
				if err != nil {
					break
				}
			}
			require.Equal(t, tc.wantError, err)
		})
	}
}
