package unsafex

import "unsafe"

// 想要使用切片进行转换，我们得知道他的底层架构
// 切片的底层是 指针-》数据,len, cap
// (*[3]uintptr)指向切片的指针
// string 底层是 unsafe.Pointer 指针+len

// 这里转换是不安全的，因为string的底层是不可修改的
// 所以我们需要不会被cg的
func UnsafeToSring(src []byte) string {
	bp := (*[3]uintptr)(unsafe.Pointer(&src))
	sh := [2]uintptr{bp[0], bp[1]}
	return *(*string)(unsafe.Pointer(&sh))
}

func UnsafeToBytes(src string) []byte {
	sp := (*[2]uintptr)(unsafe.Pointer(&src))
	sh := [3]uintptr{sp[0], sp[1], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&sh))
}
