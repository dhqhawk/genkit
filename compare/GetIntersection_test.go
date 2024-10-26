package compare

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Intersection(t *testing.T) {
	testCases := []struct {
		Name   string
		SetA   []int
		SetB   []int
		Expect []int
	}{
		{
			Name:   "1",
			SetA:   []int{1, 2, 3, 4, 5},
			SetB:   []int{1, 2, 3, 6, 7},
			Expect: []int{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			c := BasicCherker[int]{}.GetIntersection(tc.SetA, tc.SetB)
			assert.Equal(t, tc.Expect, c)
		})
	}
}
