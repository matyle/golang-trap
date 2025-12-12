package main

import (
	"fmt"
	"sync"
	"time"
)

// 陷阱：Map 的并发读写
// 问题：多个 goroutine 同时读写 map 会导致 panic

func main() {
	fmt.Println("=== 陷阱示例：Map 的并发读写 ===")
	
	// 错误示例：并发读写 map
	fmt.Println("\n错误示例：")
	// wrongWay() // 取消注释会 panic
	
	// 正确示例：使用 sync.Mutex 保护
	fmt.Println("\n正确示例1：使用 Mutex")
	correctWay1()
	
	time.Sleep(100 * time.Millisecond)
	
	// 正确示例：使用 sync.Map
	fmt.Println("\n正确示例2：使用 sync.Map")
	correctWay2()
	
	time.Sleep(100 * time.Millisecond)
}

// 错误方式：并发读写 map
func wrongWay() {
	m := make(map[string]int)
	
	// 并发写入
	go func() {
		for i := 0; i < 1000; i++ {
			m["key"] = i
		}
	}()
	
	// 并发读取
	go func() {
		for i := 0; i < 1000; i++ {
			_ = m["key"] // panic: concurrent map read and map write
		}
	}()
	
	time.Sleep(100 * time.Millisecond)
}

// 正确方式1：使用 sync.Mutex 保护
func correctWay1() {
	m := make(map[string]int)
	var mu sync.RWMutex // 读写锁，支持多个并发读
	
	// 写入
	go func() {
		for i := 0; i < 10; i++ {
			mu.Lock()
			m["key"] = i
			mu.Unlock()
			time.Sleep(1 * time.Millisecond)
		}
	}()
	
	// 读取
	go func() {
		for i := 0; i < 10; i++ {
			mu.RLock() // 读锁
			val := m["key"]
			mu.RUnlock()
			fmt.Printf("读取: %d\n", val)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	
	time.Sleep(50 * time.Millisecond)
}

// 正确方式2：使用 sync.Map（适合读多写少的场景）
func correctWay2() {
	var m sync.Map
	
	// 写入
	go func() {
		for i := 0; i < 10; i++ {
			m.Store("key", i)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	
	// 读取
	go func() {
		for i := 0; i < 10; i++ {
			if val, ok := m.Load("key"); ok {
				fmt.Printf("读取: %v\n", val)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
	
	time.Sleep(50 * time.Millisecond)
}

// 正确方式3：使用 channel 串行化访问
func correctWay3() {
	m := make(map[string]int)
	ops := make(chan func(), 100)
	
	// 单 goroutine 处理所有操作
	go func() {
		for op := range ops {
			op()
		}
	}()
	
	// 通过 channel 发送操作
	set := func(key string, val int) {
		ops <- func() {
			m[key] = val
		}
	}
	
	get := func(key string) int {
		result := make(chan int, 1)
		ops <- func() {
			result <- m[key]
		}
		return <-result
	}
	
	set("key", 42)
	fmt.Printf("读取: %d\n", get("key"))
	close(ops)
}

// 实际应用：线程安全的计数器
type SafeCounter struct {
	mu    sync.RWMutex
	count map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		count: make(map[string]int),
	}
}

func (c *SafeCounter) Increment(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count[key]++
}

func (c *SafeCounter) Get(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count[key]
}

func demonstrateCounter() {
	counter := NewSafeCounter()
	
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment("test")
		}()
	}
	
	wg.Wait()
	fmt.Printf("计数: %d\n", counter.Get("test"))
}

// 注意事项
func notes() {
	// 1. map 的并发读写会导致 panic，不是数据竞争
	// 2. 即使只是并发读取，如果有写入也会 panic
	// 3. sync.Map 适合读多写少的场景
	// 4. 对于写多读少的场景，使用 Mutex 保护普通 map 可能更高效
}

