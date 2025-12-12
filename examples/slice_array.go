package main

import "fmt"

// 陷阱：切片和数组的区别
// 问题：混淆切片和数组，导致意外的行为

func main() {
	fmt.Println("=== 陷阱示例：切片和数组的区别 ===")
	
	// 陷阱1：数组是值类型，切片是引用类型
	fmt.Println("\n陷阱1：数组是值类型")
	trap1()
	
	// 陷阱2：切片共享底层数组
	fmt.Println("\n陷阱2：切片共享底层数组")
	trap2()
	
	// 陷阱3：切片的 append 行为
	fmt.Println("\n陷阱3：切片的 append 行为")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：数组是值类型，赋值会复制
func trap1() {
	// 数组：长度是类型的一部分
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1 // 复制整个数组
	arr2[0] = 99
	
	fmt.Printf("arr1: %v\n", arr1) // [1 2 3]
	fmt.Printf("arr2: %v\n", arr2) // [99 2 3]
	
	// 切片：是引用类型
	slice1 := []int{1, 2, 3}
	slice2 := slice1 // 共享底层数组
	slice2[0] = 99
	
	fmt.Printf("slice1: %v\n", slice1) // [99 2 3]
	fmt.Printf("slice2: %v\n", slice2) // [99 2 3]
}

// 陷阱2：切片共享底层数组
func trap2() {
	original := []int{1, 2, 3, 4, 5}
	slice1 := original[1:4] // [2 3 4]
	slice2 := original[2:5] // [3 4 5]
	
	// 修改 slice1 会影响 slice2
	slice1[1] = 99
	
	fmt.Printf("original: %v\n", original) // [1 2 99 4 5]
	fmt.Printf("slice1: %v\n", slice1)    // [2 99 4]
	fmt.Printf("slice2: %v\n", slice2)    // [99 4 5]
}

// 陷阱3：append 可能创建新数组
func trap3() {
	original := []int{1, 2, 3}
	slice1 := original[:2] // [1 2]
	
	// append 可能触发重新分配
	slice2 := append(slice1, 4, 5) // [1 2 4 5]
	
	slice2[0] = 99
	
	fmt.Printf("original: %v\n", original) // [1 2 3] 或 [99 2 3]
	fmt.Printf("slice1: %v\n", slice1)     // [1 2] 或 [99 2]
	fmt.Printf("slice2: %v\n", slice2)     // [99 2 4 5]
	
	// 如果 slice2 的容量足够，会修改 original
	// 如果容量不足，会创建新数组，不会修改 original
}

// 正确方式1：使用 copy 创建独立切片
func correctWay() {
	original := []int{1, 2, 3, 4, 5}
	
	// 创建独立副本
	independent := make([]int, len(original))
	copy(independent, original)
	
	independent[0] = 99
	
	fmt.Printf("original: %v\n", original)    // [1 2 3 4 5]
	fmt.Printf("independent: %v\n", independent) // [99 2 3 4 5]
}

// 正确方式2：使用完整切片表达式
func correctWay2() {
	original := []int{1, 2, 3, 4, 5}
	
	// 完整切片表达式：array[low:high:max]
	// max 限制切片的容量
	slice := original[1:3:3] // 容量为 2，无法扩展
	fmt.Printf("限制容量的切片: %v, 容量: %d\n", slice, cap(slice))
	
	// slice = append(slice, 6) // 会创建新数组，不影响 original
}

// 数组和切片的区别总结
func summary() {
	// 1. 数组：长度固定，是值类型
	var arr [3]int
	arr2 := arr // 复制
	
	// 2. 切片：长度可变，是引用类型
	slice := []int{1, 2, 3}
	slice2 := slice // 共享底层数组
	
	// 3. 数组作为参数会复制
	modifyArray(arr) // 不会修改原数组
	
	// 4. 切片作为参数传递的是引用
	modifySlice(slice) // 会修改原切片
	
	fmt.Println(arr, arr2, slice, slice2)
}

func modifyArray(arr [3]int) {
	arr[0] = 99
}

func modifySlice(slice []int) {
	slice[0] = 99
}

