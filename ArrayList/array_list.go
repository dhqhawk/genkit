package ArrayList

import (
	"fmt"
	"genkit/errs"
)

// this is dhqhawk's editor.
var (
	_ ListAny[any] = &ArrayList[any]{}
)

type ArrayList[T any] struct {
	vals []T
}

func (a *ArrayList[T]) Get(index int) (T, error) {
	var t T
	if index > len(a.vals) || index < 0 {
		return t, fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", len(a.vals), index)
	}
	return a.vals[index], nil
}

func (a *ArrayList[T]) Append(ts ...T) error {
	a.vals = append(a.vals, ts...)
	return nil
}

func (a *ArrayList[T]) Add(index int, t T) error {
	if index > len(a.vals) || index < 0 {
		return fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", len(a.vals), index)
	}
	// 创建一个新的切片，长度比原切片多1，容纳新的元素
	//cap为0的话，不能copy
	//valtmp := make([]T, 0, len(a.vals)+1)
	valtmp := make([]T, len(a.vals)+1)
	// 复制 index 之前的元素
	copy(valtmp, a.vals[:index])

	// 插入新元素
	valtmp[index] = t

	// 复制 index 之后的元素
	copy(valtmp[index+1:], a.vals[index:])

	// 更新原切片
	a.vals = valtmp

	return nil
}

func Add[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index > length {
		return nil, errs.NewErrIndexOutOfRange(length, index)
	}

	//先将src扩展一个元素
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}
	src[index] = element
	return src, nil
}

func (a *ArrayList[T]) Set(index int, t T) error {
	if index >= len(a.vals) || index < 0 {
		return fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", len(a.vals), index)
	}
	a.vals[index] = t
	return nil
}

func (a *ArrayList[T]) Delete(index int) (T, error) {
	var t T
	if index >= len(a.vals) || index < 0 {
		return t, fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", len(a.vals), index)
	}
	t = a.vals[index]
	valtmp := make([]T, len(a.vals)-1)
	copy(valtmp, a.vals[:index])
	copy(valtmp[index:], a.vals[index+1:])
	a.vals = valtmp
	a.vals = Shrink(a.vals)
	return t, nil
}

func (a *ArrayList[T]) Len() int {
	return len(a.vals)
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.vals)
}

func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for i, v := range a.vals {
		err := fn(i, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArrayList[T]) AsSlice() []T {
	if len(a.vals) == 0 {
		return make([]T, 0)
	}
	valtmp := make([]T, len(a.vals))
	copy(valtmp, a.vals)
	return valtmp
}

func NewArrayListOf[T any](t []T) *ArrayList[T] {
	a := &ArrayList[T]{
		vals: make([]T, len(t)),
	}
	copy(a.vals, t)
	return a
}

func calCapacity(c, l int) (int, bool) {
	if c <= 64 {
		return c, false
	}
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false
}

func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}
