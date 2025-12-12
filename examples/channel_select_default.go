package main

import (
	"fmt"
	"time"
)

// 陷阱：Select 的 Default Case
// 问题：select 语句中的 default case 可能导致非阻塞行为，影响程序逻辑

func main() {
	fmt.Println("=== 陷阱示例：Select 的 Default Case ===")
	
	// 陷阱：default case 导致非阻塞
	fmt.Println("\n陷阱：default case 导致非阻塞")
	trap1()
	
	time.Sleep(100 * time.Millisecond)
	
	// 正确方式：理解 default 的用途
	fmt.Println("\n正确方式：使用 default 实现超时")
	correctWay()
	
	time.Sleep(200 * time.Millisecond)
}

// 陷阱：default case 导致立即返回，可能错过数据
func trap1() {
	ch := make(chan int)
	
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- 42
	}()
	
	// 问题：default case 会立即执行，不会等待通道数据
	select {
	case val := <-ch:
		fmt.Printf("接收到: %d\n", val)
	default:
		fmt.Println("没有数据，立即返回（可能错过数据）")
	}
	
	time.Sleep(100 * time.Millisecond)
	// 此时数据才到达，但已经错过了
}

// 正确方式1：不使用 default，等待数据
func correctWay1() {
	ch := make(chan int)
	
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- 42
	}()
	
	// 没有 default，会阻塞等待
	select {
	case val := <-ch:
		fmt.Printf("接收到: %d\n", val)
	}
}

// 正确方式2：使用 default 实现非阻塞读取
func correctWay() {
	ch := make(chan int, 1) // 带缓冲
	
	// 非阻塞发送
	select {
	case ch <- 42:
		fmt.Println("发送成功")
	default:
		fmt.Println("通道已满，无法发送")
	}
	
	// 非阻塞接收
	select {
	case val := <-ch:
		fmt.Printf("接收到: %d\n", val)
	default:
		fmt.Println("没有数据可读")
	}
}

// 正确方式3：使用 default 实现超时
func correctWay3() {
	ch := make(chan int)
	
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- 42
	}()
	
	// 使用 default 和 time.After 实现超时
	select {
	case val := <-ch:
		fmt.Printf("接收到: %d\n", val)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时：没有在指定时间内收到数据")
	}
}

// 实际应用：超时模式
func timeoutPattern() {
	ch := make(chan string)
	
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "结果"
	}()
	
	select {
	case result := <-ch:
		fmt.Printf("成功: %s\n", result)
	case <-time.After(1 * time.Second):
		fmt.Println("操作超时")
	}
}

// 实际应用：非阻塞操作
func nonBlockingPattern() {
	ch := make(chan int, 1)
	
	// 尝试发送，不阻塞
	select {
	case ch <- 1:
		fmt.Println("发送成功")
	default:
		fmt.Println("通道已满，跳过")
	}
	
	// 尝试接收，不阻塞
	select {
	case val := <-ch:
		fmt.Printf("接收成功: %d\n", val)
	default:
		fmt.Println("没有数据，跳过")
	}
}

// 实际应用：多路复用
func multiplexPattern() {
	ch1 := make(chan int)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- 42
	}()
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- "hello"
	}()
	
	// 等待任意一个通道有数据
	select {
	case val := <-ch1:
		fmt.Printf("从 ch1 接收到: %d\n", val)
	case val := <-ch2:
		fmt.Printf("从 ch2 接收到: %s\n", val)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("超时")
	}
}

// 注意事项
func notes() {
	ch := make(chan int)
	
	// 1. 没有 default 的 select 会阻塞
	// select {
	// case <-ch:
	// }
	// 上面的代码会永远阻塞
	
	// 2. 有 default 的 select 不会阻塞
	select {
	case <-ch:
		fmt.Println("有数据")
	default:
		fmt.Println("没有数据，立即返回")
	}
	
	// 3. 多个 case 都准备好时，随机选择一个
	// 4. 空的 select {} 会永远阻塞
}

