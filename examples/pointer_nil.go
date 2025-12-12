package main

import "fmt"

// 陷阱：Nil 指针解引用
// 问题：在使用指针前未检查是否为 nil，导致程序 panic

func main() {
	fmt.Println("=== 陷阱示例：Nil 指针解引用 ===")
	
	// 错误示例：直接使用 nil 指针
	fmt.Println("\n错误示例：")
	// wrongWay() // 取消注释会 panic
	
	// 正确示例：检查 nil
	fmt.Println("\n正确示例：")
	correctWay()
}

// 错误方式：直接解引用可能为 nil 的指针
func wrongWay() {
	var p *int
	fmt.Println(*p) // panic: runtime error: invalid memory address or nil pointer dereference
}

// 正确方式：在使用前检查 nil
func correctWay() {
	var p *int
	
	// 方式1：检查 nil
	if p != nil {
		fmt.Println(*p)
	} else {
		fmt.Println("指针为 nil，不能解引用")
	}
	
	// 方式2：使用函数返回指针
	p = getPointer()
	if p != nil {
		fmt.Printf("指针值: %d\n", *p)
	}
}

func getPointer() *int {
	val := 42
	return &val
}

// 常见场景：结构体指针方法
type Person struct {
	Name string
	Age  int
}

func (p *Person) GetName() string {
	// 如果 p 为 nil，这里会 panic
	// 应该先检查 p != nil
	return p.Name
}

func safeGetName(p *Person) string {
	if p == nil {
		return "未知"
	}
	return p.Name
}

