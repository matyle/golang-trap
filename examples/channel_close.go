package main

import (
	"fmt"
	"time"
)

// 陷阱：未关闭通道导致泄漏
// 问题：通道未正确关闭，导致接收方永远阻塞

func main() {
	fmt.Println("=== 陷阱示例：未关闭通道导致泄漏 ===")
	
	// 错误示例：通道未关闭
	fmt.Println("\n错误示例：")
	wrongWay()
	
	time.Sleep(200 * time.Millisecond)
	
	// 正确示例：正确关闭通道
	fmt.Println("\n正确示例：")
	correctWay()
	
	time.Sleep(200 * time.Millisecond)
}

// 错误方式：通道未关闭，接收方可能永远阻塞
func wrongWay() {
	ch := make(chan int)
	
	// 发送方
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
		// 忘记关闭通道！
	}()
	
	// 接收方会一直等待
	go func() {
		for {
			val, ok := <-ch
			if !ok {
				break
			}
			fmt.Printf("接收: %d\n", val)
		}
		fmt.Println("接收完成")
	}()
	
	time.Sleep(100 * time.Millisecond)
	fmt.Println("主程序退出（接收方可能还在等待）")
}

// 正确方式：发送方关闭通道
func correctWay() {
	ch := make(chan int)
	
	// 发送方
	go func() {
		defer close(ch) // 确保通道被关闭
		for i := 0; i < 3; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
	}()
	
	// 接收方
	go func() {
		for val := range ch { // range 会在通道关闭时自动退出
			fmt.Printf("接收: %d\n", val)
		}
		fmt.Println("接收完成")
	}()
	
	time.Sleep(100 * time.Millisecond)
	fmt.Println("所有操作完成")
}

// 正确方式2：使用 context 控制
func correctWay2() {
	ch := make(chan int)
	done := make(chan bool)
	
	// 发送方
	go func() {
		for i := 0; i < 3; i++ {
			select {
			case ch <- i:
				fmt.Printf("发送: %d\n", i)
			case <-done:
				return
			}
		}
		close(ch)
	}()
	
	// 接收方
	go func() {
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
		}
		done <- true
	}()
	
	time.Sleep(100 * time.Millisecond)
}

// 最佳实践：谁创建通道，谁负责关闭
func bestPractice() {
	// 生产者函数创建并返回通道
	ch := producer()
	
	// 消费者从通道读取
	consumer(ch)
}

func producer() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch) // 创建者负责关闭
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()
	return ch
}

func consumer(ch <-chan int) {
	for val := range ch {
		fmt.Printf("消费: %d\n", val)
	}
}

