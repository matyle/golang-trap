package main

import (
	"fmt"
	"time"
)

// 陷阱：从已关闭通道读取
// 问题：从已关闭的通道读取会立即返回零值，需要检查通道状态

func main() {
	fmt.Println("=== 陷阱示例：从已关闭通道读取 ===")
	
	// 陷阱：无法区分零值和通道关闭
	fmt.Println("\n陷阱：无法区分零值和通道关闭")
	trap1()
	
	time.Sleep(100 * time.Millisecond)
	
	// 正确方式：检查通道状态
	fmt.Println("\n正确方式：检查通道状态")
	correctWay()
	
	time.Sleep(100 * time.Millisecond)
}

// 陷阱：无法区分零值和通道关闭
func trap1() {
	ch := make(chan int)
	
	go func() {
		ch <- 0  // 发送零值
		ch <- 1
		close(ch)
	}()
	
	time.Sleep(10 * time.Millisecond)
	
	// 问题：无法区分接收到的 0 是实际值还是通道关闭后的零值
	for {
		val := <-ch
		fmt.Printf("接收到: %d\n", val)
		if val == 0 {
			// 错误：无法判断是零值还是通道关闭
			break
		}
	}
}

// 正确方式1：使用两个返回值检查通道状态
func correctWay() {
	ch := make(chan int)
	
	go func() {
		ch <- 0  // 发送零值
		ch <- 1
		ch <- 2
		close(ch)
	}()
	
	time.Sleep(10 * time.Millisecond)
	
	// 正确：使用 ok 值检查通道是否关闭
	for {
		val, ok := <-ch
		if !ok {
			fmt.Println("通道已关闭")
			break
		}
		fmt.Printf("接收到: %d\n", val)
	}
}

// 正确方式2：使用 range 循环
func correctWay2() {
	ch := make(chan int)
	
	go func() {
		ch <- 0
		ch <- 1
		ch <- 2
		close(ch)
	}()
	
	time.Sleep(10 * time.Millisecond)
	
	// range 会在通道关闭时自动退出
	for val := range ch {
		fmt.Printf("接收到: %d\n", val)
	}
	fmt.Println("通道已关闭，循环退出")
}

// 实际应用：工作池模式
func workerPool() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	
	// 启动 3 个 worker
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for job := range jobs { // 使用 range，通道关闭时自动退出
				fmt.Printf("Worker %d 处理任务 %d\n", id, job)
				results <- job * 2
			}
			fmt.Printf("Worker %d 退出\n", id)
		}(w)
	}
	
	// 发送任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // 关闭通道，通知 worker 没有更多任务
	
	// 收集结果
	for i := 1; i <= 5; i++ {
		result := <-results
		fmt.Printf("结果: %d\n", result)
	}
}

// 注意事项
func notes() {
	ch := make(chan int)
	close(ch)
	
	// 1. 从已关闭通道读取会立即返回零值
	val := <-ch
	fmt.Printf("从已关闭通道读取: %d\n", val) // 0
	
	// 2. 可以多次从已关闭通道读取
	val2 := <-ch
	fmt.Printf("再次读取: %d\n", val2) // 0
	
	// 3. 使用两个返回值检查
	val3, ok := <-ch
	fmt.Printf("值: %d, 通道打开: %v\n", val3, ok) // 0, false
	
	// 4. 从已关闭通道读取不会阻塞
	fmt.Println("不会阻塞")
}

