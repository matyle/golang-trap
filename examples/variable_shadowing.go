package main

import (
	"fmt"
	"os"
)

// 陷阱：变量遮蔽（Variable Shadowing）
// 问题：内部作用域的变量遮蔽外部作用域的变量，导致意外的行为

var global = "global"

func main() {
	fmt.Println("=== 陷阱示例：变量遮蔽 ===")
	
	// 陷阱1：短变量声明遮蔽外部变量
	fmt.Println("\n陷阱1：短变量声明遮蔽外部变量")
	trap1()
	
	// 陷阱2：if 语句中的变量遮蔽
	fmt.Println("\n陷阱2：if 语句中的变量遮蔽")
	trap2()
	
	// 陷阱3：错误处理中的变量遮蔽
	fmt.Println("\n陷阱3：错误处理中的变量遮蔽")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：短变量声明遮蔽外部变量
func trap1() {
	x := 1
	
	if true {
		x := 2 // 创建新变量，遮蔽外部的 x
		fmt.Printf("内部 x: %d\n", x) // 2
	}
	
	fmt.Printf("外部 x: %d\n", x) // 1，没有被修改
}

// 陷阱2：if 语句中的变量遮蔽
func trap2() {
	var err error // 声明为 error 类型
	
	// 错误：创建了新变量 err，遮蔽了外部的 err
	if err := doSomething(); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	// 外部的 err 仍然是 nil
	fmt.Printf("外部 err: %v\n", err)
}

func doSomething() error {
	return fmt.Errorf("something went wrong")
}

// 陷阱3：错误处理中的变量遮蔽
func trap3() {
	file, err := os.Open("test.txt")
	if err != nil {
		return
	}
	defer file.Close()
	
	// 错误：创建了新变量 file 和 err
	file, err := os.Open("another.txt")
	if err != nil {
		return // 外部的 file 没有被关闭！
	}
	defer file.Close()
}

// 正确方式1：使用赋值而不是短变量声明
func correctWay() {
	x := 1
	
	if true {
		x = 2 // 赋值，修改外部的 x
		fmt.Printf("内部 x: %d\n", x) // 2
	}
	
	fmt.Printf("外部 x: %d\n", x) // 2，被修改了
}

// 正确方式2：使用不同的变量名
func correctWay2() {
	var err error // 声明为 error 类型
	
	// 使用不同的变量名
	if err2 := doSomething(); err2 != nil {
		fmt.Printf("错误: %v\n", err2)
		return
	}
	
	fmt.Printf("外部 err: %v\n", err)
}

// 正确方式3：在 if 语句外声明变量
func correctWay3() {
	var file *os.File
	var err error
	
	file, err = os.Open("test.txt")
	if err != nil {
		return
	}
	defer file.Close()
	
	// 使用赋值，不创建新变量
	file, err = os.Open("another.txt")
	if err != nil {
		return
	}
	defer file.Close()
}

// 实际应用：循环中的变量遮蔽
func demonstrateLoop() {
	s := []int{1, 2, 3}
	
	// 错误：所有闭包共享同一个 i
	var funcs []func()
	for i := range s {
		funcs = append(funcs, func() {
			fmt.Println(i) // 可能都打印 2
		})
	}
	
	// 正确：创建局部变量
	var funcs2 []func()
	for i := range s {
		i := i // 创建局部变量
		funcs2 = append(funcs2, func() {
			fmt.Println(i)
		})
	}
	
	for _, f := range funcs2 {
		f()
	}
}

// 注意事项
func notes() {
	// 1. := 会创建新变量，即使变量名相同
	// 2. = 会修改现有变量
	// 3. 在 if/for 等语句中使用 := 要小心
	// 4. 使用 go vet 可以检测变量遮蔽
	// 5. 使用不同的变量名可以避免遮蔽
}

