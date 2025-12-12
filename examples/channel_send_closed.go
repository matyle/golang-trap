package main

import (
	"fmt"
	"sync"
	"time"
)

// 陷阱：向已关闭通道发送数据
// 问题：向已关闭的通道发送数据会导致 panic

func main() {
	fmt.Println("=== 陷阱示例：向已关闭通道发送数据 ===")
	
	// 错误示例：向已关闭通道发送
	fmt.Println("\n错误示例：")
	// wrongWay() // 取消注释会 panic
	
	// 正确示例：检查通道状态
	fmt.Println("\n正确示例：")
	correctWay()
	
	time.Sleep(100 * time.Millisecond)
}

// 错误方式：向已关闭的通道发送数据
func wrongWay() {
	ch := make(chan int)
	
	go func() {
		close(ch)
	}()
	
	time.Sleep(10 * time.Millisecond)
	
	// panic: send on closed channel
	ch <- 42
}

// 正确方式1：使用 sync.Once 确保只关闭一次
func correctWay() {
	ch := make(chan int)
	var once sync.Once
	
	// 发送方
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
		once.Do(func() {
			close(ch)
			fmt.Println("通道已关闭")
		})
	}()
	
	// 接收方
	go func() {
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
		}
	}()
	
	time.Sleep(50 * time.Millisecond)
}

// 正确方式2：使用 context 控制发送
func correctWay2() {
	ch := make(chan int)
	done := make(chan struct{})
	
	// 发送方
	go func() {
		defer close(ch)
		for i := 0; i < 3; i++ {
			select {
			case ch <- i:
				fmt.Printf("发送: %d\n", i)
			case <-done:
				return
			}
		}
	}()
	
	// 接收方
	go func() {
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
		}
		close(done)
	}()
	
	time.Sleep(50 * time.Millisecond)
}

// 正确方式3：使用 recover 捕获 panic（不推荐，但可以用于防御性编程）
func safeSend(ch chan int, val int) (sent bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到 panic: %v\n", r)
			sent = false
		}
	}()
	
	ch <- val
	return true
}

// 最佳实践：明确通道的所有权
// 1. 谁创建通道，谁负责关闭
// 2. 只在一个 goroutine 中关闭通道
// 3. 使用 done 通道或 context 来通知停止发送

func bestPractice() {
	ch := make(chan int)
	done := make(chan struct{})
	
	// 只有一个 goroutine 负责关闭
	go func() {
		defer close(ch)
		for {
			select {
			case ch <- 1:
				time.Sleep(10 * time.Millisecond)
			case <-done:
				return
			}
		}
	}()
	
	// 主程序控制何时停止
	time.Sleep(100 * time.Millisecond)
	close(done)
	
	// 等待通道关闭
	for range ch {
		// 消费剩余数据
	}
}

