package main

import (
	"fmt"
	"strings"
)

// 陷阱：性能问题
// 问题：常见的性能陷阱导致程序运行缓慢

func main() {
	fmt.Println("=== 陷阱示例：性能问题 ===")
	
	// 陷阱1：字符串拼接
	fmt.Println("\n陷阱1：字符串拼接")
	trap1()
	
	// 陷阱2：切片预分配
	fmt.Println("\n陷阱2：切片预分配")
	trap2()
	
	// 陷阱3：不必要的内存分配
	fmt.Println("\n陷阱3：不必要的内存分配")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：使用 + 拼接字符串
func trap1() {
	// 错误：每次拼接都创建新字符串
	var result string
	for i := 0; i < 1000; i++ {
		result += fmt.Sprintf("%d ", i) // 低效
	}
	_ = result
}

// 陷阱2：切片未预分配容量
func trap2() {
	// 错误：频繁扩容
	var slice []int
	for i := 0; i < 1000; i++ {
		slice = append(slice, i) // 可能多次扩容
	}
	_ = slice
}

// 陷阱3：不必要的内存分配
func trap3() {
	// 错误：在循环中创建大对象
	for i := 0; i < 1000; i++ {
		data := make([]byte, 1024*1024) // 每次循环都分配
		_ = data
	}
}

// 正确方式1：使用 strings.Builder
func correctWay() {
	// 正确：使用 strings.Builder
	var builder strings.Builder
	builder.Grow(10000) // 预分配容量
	for i := 0; i < 1000; i++ {
		builder.WriteString(fmt.Sprintf("%d ", i))
	}
	result := builder.String()
	_ = result
}

// 正确方式2：预分配切片容量
func correctWay2() {
	// 正确：预分配容量
	slice := make([]int, 0, 1000) // 预分配容量
	for i := 0; i < 1000; i++ {
		slice = append(slice, i) // 不会扩容
	}
	_ = slice
}

// 正确方式3：复用对象
func correctWay3() {
	// 正确：在循环外分配
	data := make([]byte, 1024*1024)
	for i := 0; i < 1000; i++ {
		// 复用 data
		_ = data
	}
}

// 其他性能陷阱
func otherPitfalls() {
	// 1. 频繁的 map 查找
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		val := m["key"] // 每次都查找
		_ = val
	}
	
	// 正确：缓存查找结果
	val, ok := m["key"]
	if ok {
		// 使用 val
		_ = val
	}
	
	// 2. 不必要的类型转换
	var i interface{} = 42
	for j := 0; j < 1000; j++ {
		_ = i.(int) // 每次都转换
	}
	
	// 正确：转换一次
	intVal := i.(int)
	for j := 0; j < 1000; j++ {
		_ = intVal
	}
	
	// 3. 大结构体按值传递
	type LargeStruct struct {
		data [1000]int
	}
	
	// 错误：按值传递大结构体
	processLarge := func(s LargeStruct) {
		_ = s
	}
	
	// 正确：按指针传递
	processLargePtr := func(s *LargeStruct) {
		_ = s
	}
	
	_ = processLarge
	_ = processLargePtr
}

// 性能优化建议
func optimizationTips() {
	// 1. 使用 strings.Builder 而不是 + 拼接字符串
	// 2. 预分配切片和 map 的容量
	// 3. 避免不必要的内存分配
	// 4. 复用对象而不是创建新对象
	// 5. 大结构体使用指针传递
	// 6. 使用 sync.Pool 复用临时对象
	// 7. 避免频繁的类型断言和转换
	// 8. 使用 pprof 分析性能瓶颈
}

