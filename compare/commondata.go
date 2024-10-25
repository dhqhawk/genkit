package compare

import (
	"errors"
	"fmt"
)

type BasicCherker[T comparable] struct {
}

// map是无法使用any的
// 定义map[KeyType]ValueType
// tmpmap := make(map[T any]struct{},len(SuperSet))
// 这里使用反射操作，声明也是不行的，因为make需要在编译时确定类型
// 反射是在运行的时候确认类型的
// tmpmap := make(map[reflect.TypeOf(SubSet[0])]struct{},len(SuperSet))
// 所以我这里只能any改为compareable
// 注意这里只需要将该struct改为compareable即可,接口保持any
func (b BasicCherker[T]) IsSubset(SuperSet []T, SubSet []T) (bool, error) {
	if len(SuperSet) < len(SubSet) {
		return false, errors.New("父集长度过短")
	}
	tmpmap := make(map[T]struct{}, len(SuperSet))
	for _, v := range SuperSet {
		tmpmap[v] = struct{}{}
	}
	for _, v := range SubSet {
		if _, ok := tmpmap[v]; !ok {
			return false, fmt.Errorf("%v 不在父集中", v)
		}
	}
	return true, nil
}

// 这里SetB使用切片，而不是...
// 如果使用...传入一个参数的时候就陷入了 IsSubset的问题
func (b BasicCherker[T]) GetIntersection(SetA []T, SetB []T) []T {
	var result []T
	tmpmap := make(map[T]struct{}, len(SetA))
	for _, v := range SetA {
		tmpmap[v] = struct{}{}
	}
	for _, v := range SetB {
		if _, ok := tmpmap[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
