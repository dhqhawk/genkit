package compare

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicCherker_IsSubset(t *testing.T) {
	testCases := []struct {
		name     string
		superSet []int
		subSet   []int
		wantbool bool
		wantErr  error
	}{
		{
			name:     "int子集只传入一个,是子集",
			superSet: []int{1, 2, 3, 4, 5},
			subSet:   []int{2},
			wantbool: true,
			wantErr:  nil,
		},
		{
			name:     "int子集只传入一个,不是子集",
			superSet: []int{1, 2, 3, 4, 5},
			subSet:   []int{14},
			wantbool: false,
			wantErr:  fmt.Errorf("%v 不在父集中", 14),
		},
		{
			name:     "int子集传入多个,是子集",
			superSet: []int{1, 2, 3, 4, 5},
			subSet:   []int{1, 2},
			wantbool: true,
			wantErr:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := BasicCherker[int]{}.IsSubset(tc.superSet, tc.subSet)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantbool, ok)
		})
	}
	testCase1 := []struct {
		name     string
		superSet []string
		subSet   []string
		wantbool bool
		wantErr  error
	}{
		{
			name:     "string子集只传入一个,是子集",
			superSet: []string{"a", "b", "c", "d", "e"},
			subSet:   []string{"a"},
			wantbool: true,
			wantErr:  nil,
		},
		{
			name:     "string子集只传入一个,不是子集",
			superSet: []string{"a", "b", "c", "d", "e"},
			subSet:   []string{"t"},
			wantbool: false,
			wantErr:  fmt.Errorf("%v 不在父集中", "t"),
		},
		{
			name:     "string子集传入多个,是子集",
			superSet: []string{"a", "b", "c", "d", "e"},
			subSet:   []string{"a", "b", "c", "d", "e"},
			wantbool: true,
			wantErr:  nil,
		},

		{
			name:     "string子集传入多个,不是子集",
			superSet: []string{"a", "b", "c", "d", "e"},
			subSet:   []string{"a", "b", "c", "d", "t"},
			wantbool: false,
			wantErr:  fmt.Errorf("%v 不在父集中", "t"),
		},
	}
	for _, tc := range testCase1 {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := BasicCherker[string]{}.IsSubset(tc.superSet, tc.subSet)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantbool, ok)
		})
	}
}

type User struct {
	Id int
}

func TestBasicCherker_structIsSubset(t *testing.T) {
	testCases := []struct {
		name     string
		superSet []User
		subSet   []User
		wantbool bool
		wantErr  error
	}{
		{
			name:     "struct子集只传入一个,是子集",
			superSet: []User{User{1}, User{2}, User{3}, User{4}, User{5}},
			subSet:   []User{User{1}},
			wantbool: true,
			wantErr:  nil,
		},
		{
			name:     "struct子集只传入一个,不是子集",
			superSet: []User{User{1}, User{2}, User{3}, User{4}, User{5}},
			subSet:   []User{User{11}},
			wantbool: false,
			wantErr:  fmt.Errorf("%v 不在父集中", User{11}),
		},
		{
			name:     "struct子集传入多个,是子集",
			superSet: []User{User{1}, User{2}, User{3}, User{4}, User{5}},
			subSet:   []User{User{1}, User{2}},
			wantbool: true,
			wantErr:  nil,
		},
		{
			name:     "struct子集只传入多个,不是子集",
			superSet: []User{User{1}, User{2}, User{3}, User{4}, User{5}},
			subSet:   []User{User{11}, {13}},
			wantbool: false,
			wantErr:  fmt.Errorf("%v 不在父集中", User{11}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := BasicCherker[User]{}.IsSubset(tc.superSet, tc.subSet)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantbool, ok)
		})
	}
}
