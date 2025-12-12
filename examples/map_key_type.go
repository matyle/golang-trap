package main

import "fmt"

// 陷阱：Map 键类型限制
// 问题：map 的键类型必须是可比较的类型

func main() {
	fmt.Println("=== 陷阱示例：Map 键类型限制 ===")
	
	// 陷阱1：使用不可比较的类型作为键
	fmt.Println("\n陷阱1：使用不可比较的类型作为键")
	trap1()
	
	// 陷阱2：使用切片作为键
	fmt.Println("\n陷阱2：使用切片作为键")
	trap2()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：使用不可比较的类型作为键
func trap1() {
	// 错误：切片不能作为 map 的键
	// m := make(map[[]int]string) // 编译错误！
	
	// 错误：map 不能作为 map 的键
	// m := make(map[map[string]int]string) // 编译错误！
	
	// 错误：函数不能作为 map 的键
	// m := make(map[func()]string) // 编译错误！
}

// 陷阱2：使用包含不可比较类型的结构体作为键
func trap2() {
	// 结构体包含切片，不能作为键
	type BadKey struct {
		Name  string
		Items []int // 切片不可比较
	}
	
	// m := make(map[BadKey]string) // 编译错误！
	
	// 结构体包含 map，不能作为键
	type BadKey2 struct {
		Name string
		Data map[string]int // map 不可比较
	}
	
	// m := make(map[BadKey2]string) // 编译错误！
}

// 正确方式1：使用可比较的类型作为键
func correctWay() {
	// 基本类型都可以作为键
	m1 := make(map[int]string)
	m1[1] = "one"
	
	m2 := make(map[string]int)
	m2["one"] = 1
	
	m3 := make(map[bool]string)
	m3[true] = "true"
	
	// 数组可以作为键（如果元素类型可比较）
	m4 := make(map[[3]int]string)
	m4[[3]int{1, 2, 3}] = "array"
	
	fmt.Printf("m1: %v\n", m1)
	fmt.Printf("m2: %v\n", m2)
	fmt.Printf("m3: %v\n", m3)
	fmt.Printf("m4: %v\n", m4)
}

// 正确方式2：使用结构体作为键（所有字段都可比较）
func correctWay2() {
	type GoodKey struct {
		Name  string
		ID    int
		Valid bool
	}
	
	m := make(map[GoodKey]string)
	m[GoodKey{Name: "Alice", ID: 1, Valid: true}] = "value"
	
	fmt.Printf("m: %v\n", m)
}

// 正确方式3：将不可比较类型转换为可比较类型
func correctWay3() {
	// 将切片转换为字符串（如果合适）
	slice := []int{1, 2, 3}
	key := fmt.Sprintf("%v", slice) // 转换为字符串
	
	m := make(map[string]string)
	m[key] = "value"
	
	fmt.Printf("m: %v\n", m)
}

// 实际应用：使用指针作为键
func demonstratePointerKey() {
	type Data struct {
		Value int
	}
	
	// 指针是可比较的
	m := make(map[*Data]string)
	
	d1 := &Data{Value: 1}
	d2 := &Data{Value: 2}
	
	m[d1] = "first"
	m[d2] = "second"
	
	fmt.Printf("d1: %v\n", m[d1])
	fmt.Printf("d2: %v\n", m[d2])
}

// 可比较的类型总结
func comparableTypes() {
	// 可比较的类型：
	// - 布尔类型
	// - 数值类型（int, float, complex 等）
	// - 字符串
	// - 指针
	// - 通道
	// - 接口（如果动态类型可比较）
	// - 数组（如果元素类型可比较）
	// - 结构体（如果所有字段可比较）
	
	// 不可比较的类型：
	// - 切片
	// - map
	// - 函数
}

