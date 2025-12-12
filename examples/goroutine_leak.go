package main

import (
	"fmt"
	"time"
)

// 陷阱：Goroutine 泄漏
// 问题：goroutine 因为通道阻塞而永远无法退出，造成内存泄漏

func main() {
	fmt.Println("=== 陷阱示例：Goroutine 泄漏 ===")
	
	// 错误示例：goroutine 永远阻塞
	fmt.Println("\n错误示例：")
	wrongWay()
	
	time.Sleep(200 * time.Millisecond)
	
	// 正确示例：使用 context 或关闭通道
	fmt.Println("\n正确示例：")
	correctWay()
	
	time.Sleep(200 * time.Millisecond)
}

// 错误方式：goroutine 永远阻塞在通道上
func wrongWay() {
	ch := make(chan int)
	
	// 这个 goroutine 会永远阻塞，因为没有人会向通道发送数据
	go func() {
		val := <-ch // 永远阻塞在这里
		fmt.Printf("收到值: %d\n", val)
	}()
	
	fmt.Println("Goroutine 已启动（但会永远阻塞）")
	// 主程序退出，但 goroutine 仍在运行，造成泄漏
}

// 正确方式1：使用带缓冲的通道或确保有发送者
func correctWay() {
	ch := make(chan int, 1) // 带缓冲的通道
	
	go func() {
		val := <-ch
		fmt.Printf("收到值: %d\n", val)
	}()
	
	ch <- 42 // 发送数据
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Goroutine 正常完成")
}

// 正确方式2：使用 context 控制 goroutine 生命周期
func correctWay2() {
	ch := make(chan int)
	done := make(chan bool)
	
	go func() {
		select {
		case val := <-ch:
			fmt.Printf("收到值: %d\n", val)
		case <-done:
			fmt.Println("收到退出信号")
			return
		}
	}()
	
	// 如果不需要继续运行，发送退出信号
	close(done)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Goroutine 正常退出")
}

