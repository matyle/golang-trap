package main

import "fmt"

// 陷阱：nil map 写入
// 问题：向 nil map 写入数据会导致 panic

func main() {
	fmt.Println("=== 陷阱示例：nil map 写入 ===")
	
	// 陷阱1：向 nil map 写入
	fmt.Println("\n陷阱1：向 nil map 写入")
	// trap1() // 取消注释会 panic
	
	// 陷阱2：nil map 读取
	fmt.Println("\n陷阱2：nil map 读取")
	trap2()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：向 nil map 写入
func trap1() {
	var m map[string]int
	
	// panic: assignment to entry in nil map
	m["key"] = 1
}

// 陷阱2：nil map 读取
func trap2() {
	var m map[string]int
	
	// nil map 可以读取，返回零值
	val := m["key"]
	fmt.Printf("读取 nil map: %d\n", val) // 0
	
	// 检查键是否存在
	val, ok := m["key"]
	fmt.Printf("值: %d, 存在: %v\n", val, ok) // 0, false
}

// 正确方式1：初始化 map
func correctWay() {
	// 方式1：使用 make
	m1 := make(map[string]int)
	m1["key"] = 1
	fmt.Printf("m1: %v\n", m1)
	
	// 方式2：使用字面量
	m2 := map[string]int{
		"key": 1,
	}
	fmt.Printf("m2: %v\n", m2)
	
	// 方式3：声明时初始化
	var m3 map[string]int = make(map[string]int)
	m3["key"] = 1
	fmt.Printf("m3: %v\n", m3)
}

// 正确方式2：检查 map 是否为 nil
func correctWay2() {
	var m map[string]int
	
	// 在使用前检查并初始化
	if m == nil {
		m = make(map[string]int)
	}
	
	m["key"] = 1
	fmt.Printf("m: %v\n", m)
}

// 实际应用：map 作为函数参数
func processMap(m map[string]int) {
	// 如果传入 nil map，需要检查
	if m == nil {
		m = make(map[string]int)
	}
	m["processed"] = 1
}

// 注意事项
func notes() {
	// 1. nil map 可以读取，返回零值
	// 2. nil map 不能写入，会 panic
	// 3. 使用 make 或字面量初始化 map
	// 4. 检查 map 是否为 nil 后再写入
	// 5. 删除 nil map 的元素不会 panic（但也没有效果）
	
	var m map[string]int
	delete(m, "key") // 不会 panic，但也没有效果
}

