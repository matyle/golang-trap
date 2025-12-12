package main

import "fmt"

// 陷阱：Defer 的执行顺序
// 问题：defer 语句的执行顺序和参数求值时机容易混淆

func main() {
	fmt.Println("=== 陷阱示例：Defer 的执行顺序 ===")
	
	// 陷阱1：defer 的参数立即求值
	fmt.Println("\n陷阱1：defer 的参数立即求值")
	trap1()
	
	// 陷阱2：defer 的执行顺序（LIFO）
	fmt.Println("\n陷阱2：defer 的执行顺序（LIFO）")
	trap2()
	
	// 陷阱3：defer 修改返回值
	fmt.Println("\n陷阱3：defer 修改返回值")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：defer 的参数在调用时立即求值
func trap1() {
	i := 0
	defer fmt.Println("defer 1:", i) // i 的值是 0（立即求值）
	
	i++
	defer fmt.Println("defer 2:", i) // i 的值是 1（立即求值）
	
	i++
	fmt.Println("函数结束:", i) // i 的值是 2
	
	// 输出顺序：
	// 函数结束: 2
	// defer 2: 1
	// defer 1: 0
}

// 陷阱2：defer 的执行顺序是 LIFO（后进先出）
func trap2() {
	defer fmt.Println("第一个 defer")
	defer fmt.Println("第二个 defer")
	defer fmt.Println("第三个 defer")
	
	fmt.Println("函数执行")
	
	// 输出顺序：
	// 函数执行
	// 第三个 defer
	// 第二个 defer
	// 第一个 defer
}

// 陷阱3：defer 可以修改命名返回值
func trap3() {
	fmt.Println("返回值:", returnValue1()) // 返回 2
	fmt.Println("返回值:", returnValue2()) // 返回 1
}

// 命名返回值，defer 可以修改
func returnValue1() (result int) {
	defer func() {
		result++ // 修改返回值
	}()
	return 1 // 实际返回 2
}

// 匿名返回值，defer 不能修改
func returnValue2() int {
	result := 1
	defer func() {
		result++ // 修改局部变量，不影响返回值
	}()
	return result // 返回 1
}

// 正确方式1：使用闭包访问最新值
func correctWay() {
	i := 0
	defer func() {
		fmt.Println("defer:", i) // 使用闭包，访问最新的 i
	}()
	
	i++
	fmt.Println("函数结束:", i)
	
	// 输出：
	// 函数结束: 1
	// defer: 1
}

// 正确方式2：理解 defer 的执行时机
func correctWay2() {
	fmt.Println("开始")
	
	defer func() {
		fmt.Println("defer 1")
	}()
	
	defer func() {
		fmt.Println("defer 2")
	}()
	
	fmt.Println("结束")
	
	// 输出：
	// 开始
	// 结束
	// defer 2
	// defer 1
}

// 实际应用：资源清理
func resourceCleanup() {
	fmt.Println("打开资源")
	
	defer func() {
		fmt.Println("清理资源")
	}()
	
	fmt.Println("使用资源")
	
	// 即使发生 panic，defer 也会执行
	// panic("错误")
}

// 实际应用：错误处理
func errorHandling() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	
	// 可能 panic 的代码
	panic("测试错误")
	
	return nil
}

// 注意事项
func notes() {
	// 1. defer 的参数在调用时立即求值
	// 2. defer 在函数返回前执行（LIFO 顺序）
	// 3. defer 可以修改命名返回值
	// 4. defer 中的闭包会捕获最新的变量值
	// 5. defer 即使发生 panic 也会执行
}

