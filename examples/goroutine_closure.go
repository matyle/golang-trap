package main

import (
	"fmt"
	"time"
)

// 陷阱：闭包变量捕获问题
// 问题：在循环中使用 goroutine 时，所有 goroutine 可能共享同一个变量

func main() {
	fmt.Println("=== 陷阱示例：闭包变量捕获 ===")
	
	// 错误示例：所有 goroutine 共享变量 i
	fmt.Println("\n错误示例：")
	wrongWay()
	
	time.Sleep(100 * time.Millisecond)
	
	// 正确示例：通过参数传递或创建局部变量
	fmt.Println("\n正确示例：")
	correctWay()
	
	time.Sleep(100 * time.Millisecond)
}

// 错误方式：所有 goroutine 都读取到循环结束后的 i 值
func wrongWay() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Printf("错误: i = %d\n", i) // 所有 goroutine 可能都打印 5
		}()
	}
	time.Sleep(50 * time.Millisecond)
}

// 正确方式1：通过参数传递
func correctWay() {
	for i := 0; i < 5; i++ {
		go func(val int) {
			fmt.Printf("正确: val = %d\n", val)
		}(i) // 将 i 作为参数传递
	}
	time.Sleep(50 * time.Millisecond)
}

// 正确方式2：在循环内创建局部变量
func correctWay2() {
	for i := 0; i < 5; i++ {
		i := i // 创建局部变量
		go func() {
			fmt.Printf("正确: i = %d\n", i)
		}()
	}
	time.Sleep(50 * time.Millisecond)
}

