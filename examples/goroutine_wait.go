package main

import (
	"fmt"
	"sync"
	"time"
)

// 陷阱：未等待 Goroutine 完成
// 问题：主程序在 goroutine 完成前退出，导致 goroutine 被强制终止

func main() {
	fmt.Println("=== 陷阱示例：未等待 Goroutine 完成 ===")
	
	// 错误示例：主程序立即退出
	fmt.Println("\n错误示例：")
	wrongWay()
	
	// 正确示例：使用 WaitGroup 等待
	fmt.Println("\n正确示例：")
	correctWay()
}

// 错误方式：主程序可能在 goroutine 完成前就退出了
func wrongWay() {
	for i := 0; i < 3; i++ {
		go func(id int) {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Goroutine %d 完成\n", id)
		}(i)
	}
	// 主程序立即退出，goroutine 可能还没执行完
	fmt.Println("主程序退出（goroutine 可能未完成）")
}

// 正确方式：使用 sync.WaitGroup
func correctWay() {
	var wg sync.WaitGroup
	
	for i := 0; i < 3; i++ {
		wg.Add(1) // 增加计数
		go func(id int) {
			defer wg.Done() // 完成后减少计数
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Goroutine %d 完成\n", id)
		}(i)
	}
	
	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("所有 goroutine 已完成，主程序退出")
}

