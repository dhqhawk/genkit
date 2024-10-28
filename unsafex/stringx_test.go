package unsafex

import (
	"log"
	"testing"
)

// b.(type)只能用在switch中,%T可以直接输出类型
// 基准测试标准的 string([]byte) 转换
// 注意这里只有数据足够大的时候才有意义,如果切片太小会被优化，不会出现复制
func Benchmark_StringConversion(b *testing.B) {
	b.Run("direct string", func(b *testing.B) {
		data := make([]byte, 1024*1024) // 1 MB 数据
		for i := range data {
			data[i] = 'a'
		}
		b.ReportAllocs() // 启用内存分配报告
		_ = string(data)
		for i := 0; i < b.N; i++ {
			_ = string(data) // 标准转换
		}
	})
	b.Run("unsafe string", func(b *testing.B) {
		data := make([]byte, 1024*1024) // 1 MB 数据
		for i := range data {
			data[i] = 'a'
		}
		b.ReportAllocs() // 启用内存分配报告

		for i := 0; i < b.N; i++ {
			_ = UnsafeToSring(data) // 非安全转换
		}
	})
}

func Benchmark_ByteConversion(b *testing.B) {
	b.Run("direct byte", func(b *testing.B) {
		data := "hello"
		for i := 0; i < 1024*100; i++ {
			data += "a" // 不断添加字符
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = []byte(data)
		}
	})
	b.Run("unsafe byte", func(b *testing.B) {
		data := "hello"
		for i := 0; i < 1024*100; i++ {
			data += "a" // 不断添加字符
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = UnsafeToBytes(data)
		}
	})
}

func Test_ByteConversion(t *testing.T) {
	g := UnsafeToSring([]byte("hello world"))
	log.Println(g)
}
