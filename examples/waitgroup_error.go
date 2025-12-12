package main

import (
	"fmt"
	"sync"
	"time"
)

// 陷阱：WaitGroup 使用错误
// 问题：WaitGroup 使用不当导致死锁或 goroutine 泄漏

func main() {
	fmt.Println("=== 陷阱示例：WaitGroup 使用错误 ===")
	
	// 陷阱1：Add 和 Done 不匹配
	fmt.Println("\n陷阱1：Add 和 Done 不匹配")
	// trap1() // 取消注释会 panic
	
	// 陷阱2：在 goroutine 外调用 Done
	fmt.Println("\n陷阱2：在 goroutine 外调用 Done")
	trap2()
	
	// 陷阱3：Add 调用时机错误
	fmt.Println("\n陷阱3：Add 调用时机错误")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：Add 和 Done 不匹配
func trap1() {
	var wg sync.WaitGroup
	
	wg.Add(2) // 添加 2 个计数
	go func() {
		defer wg.Done() // 只完成 1 个
		fmt.Println("Goroutine 1")
	}()
	
	wg.Wait() // 永远等待，因为计数不匹配
	// 或者 Done 调用次数超过 Add，会 panic
}

// 陷阱2：在 goroutine 外调用 Done
func trap2() {
	var wg sync.WaitGroup
	
	wg.Add(1)
	wg.Done() // 在 goroutine 外调用，可能导致计数错误
	
	go func() {
		fmt.Println("Goroutine 执行")
		// 忘记调用 wg.Done()
	}()
	
	// 可能立即返回，也可能永远等待
	wg.Wait()
	fmt.Println("完成")
}

// 陷阱3：Add 调用时机错误
func trap3() {
	var wg sync.WaitGroup
	
	// 错误：在 goroutine 启动后才 Add
	go func() {
		wg.Add(1) // 可能太晚了
		defer wg.Done()
		fmt.Println("Goroutine 执行")
	}()
	
	// 主程序可能在 Add 之前就 Wait 了
	time.Sleep(10 * time.Millisecond)
	wg.Wait()
}

// 正确方式1：确保 Add 和 Done 匹配
func correctWay() {
	var wg sync.WaitGroup
	
	// 在启动 goroutine 之前 Add
	wg.Add(3)
	
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done() // 确保 Done 被调用
			fmt.Printf("Goroutine %d 执行\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Println("所有 goroutine 完成")
}

// 正确方式2：使用 defer 确保 Done 被调用
func correctWay2() {
	var wg sync.WaitGroup
	
	for i := 0; i < 3; i++ {
		wg.Add(1) // 在循环中每次 Add
		go func(id int) {
			defer wg.Done() // 使用 defer 确保 Done 被调用
			fmt.Printf("Goroutine %d 执行\n", id)
		}(i)
	}
	
	wg.Wait()
}

// 正确方式3：使用函数封装
func runWithWaitGroup(fn func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
	wg.Wait()
}

// 注意事项
func notes() {
	// 1. Add 必须在 Wait 之前调用
	// 2. Done 必须在对应的 goroutine 中调用
	// 3. Add 和 Done 的次数必须匹配
	// 4. 使用 defer wg.Done() 确保即使发生 panic 也会调用
	// 5. WaitGroup 不能复制，必须传递指针
}

