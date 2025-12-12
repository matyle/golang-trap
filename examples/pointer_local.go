package main

import "fmt"

// 陷阱：返回局部变量指针
// 问题：返回函数内部局部变量的指针，该变量在函数返回后可能被回收
// 注意：Go 编译器会进行逃逸分析，通常会自动将变量分配到堆上，但理解这个概念很重要

func main() {
	fmt.Println("=== 陷阱示例：返回局部变量指针 ===")

	// 在 Go 中，返回局部变量指针通常是安全的（编译器会处理）
	// 但理解内存管理很重要

	fmt.Println("\n示例：返回局部变量指针（Go 中通常是安全的）")
	safeExample()

	fmt.Println("\n示例：返回局部变量的值（更安全）")
	saferExample()
}

// 在 Go 中，返回局部变量指针是安全的
// 编译器会进行逃逸分析，将变量分配到堆上
func safeExample() *int {
	val := 42   // 局部变量
	return &val // Go 编译器会将 val 分配到堆上
}

// 更安全的做法：返回值而不是指针
func saferExample() int {
	val := 42
	return val // 返回值的副本
}

// 陷阱场景：返回局部数组/切片的指针
func getArrayPointer() *[3]int {
	arr := [3]int{1, 2, 3} // 局部数组
	return &arr            // 在 Go 中这是安全的，编译器会处理
}

// 更好的做法：返回切片（切片本身包含指针）
func getSlice() []int {
	arr := [3]int{1, 2, 3}
	return arr[:] // 返回切片
}

// 实际使用示例
func demonstrate() {
	p := safeExample()
	fmt.Printf("指针值: %d\n", *p)

	val := saferExample()
	fmt.Printf("值: %d\n", val)

	arrPtr := getArrayPointer()
	fmt.Printf("数组: %v\n", *arrPtr)

	slice := getSlice()
	fmt.Printf("切片: %v\n", slice)
}
