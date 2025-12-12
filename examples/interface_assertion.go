package main

import "fmt"

// 陷阱：接口类型断言
// 问题：类型断言失败时未检查 ok 值，导致 panic

type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func main() {
	fmt.Println("=== 陷阱示例：接口类型断言 ===")
	
	// 错误示例：未检查类型断言
	fmt.Println("\n错误示例：")
	// wrongWay() // 取消注释会 panic
	
	// 正确示例：检查类型断言
	fmt.Println("\n正确示例：")
	correctWay()
	
	// 类型断言的两种形式
	fmt.Println("\n类型断言的两种形式：")
	twoForms()
}

// 错误方式：直接使用类型断言，失败会 panic
func wrongWay() {
	var a Animal = Dog{Name: "Buddy"}
	
	// 如果类型断言失败，会 panic
	cat := a.(Cat) // panic: interface conversion: main.Animal is main.Dog, not main.Cat
	fmt.Println(cat.Speak())
}

// 正确方式1：使用 ok 值检查
func correctWay() {
	var a Animal = Dog{Name: "Buddy"}
	
	// 使用两个返回值的形式
	dog, ok := a.(Dog)
	if ok {
		fmt.Printf("是 Dog: %s\n", dog.Speak())
	} else {
		fmt.Println("不是 Dog")
	}
	
	cat, ok := a.(Cat)
	if ok {
		fmt.Printf("是 Cat: %s\n", cat.Speak())
	} else {
		fmt.Println("不是 Cat")
	}
}

// 正确方式2：使用 type switch
func correctWay2(a Animal) {
	switch v := a.(type) {
	case Dog:
		fmt.Printf("是 Dog: %s\n", v.Speak())
	case Cat:
		fmt.Printf("是 Cat: %s\n", v.Speak())
	default:
		fmt.Printf("未知类型: %T\n", v)
	}
}

// 类型断言的两种形式
func twoForms() {
	var a Animal = Dog{Name: "Buddy"}
	
	// 形式1：单值形式（失败会 panic）
	// dog := a.(Dog) // 如果失败会 panic
	
	// 形式2：双值形式（安全）
	dog, ok := a.(Dog)
	if ok {
		fmt.Printf("类型断言成功: %s\n", dog.Speak())
	}
	
	// 形式3：只检查类型，不获取值
	_, ok = a.(Cat)
	if !ok {
		fmt.Println("不是 Cat 类型")
	}
}

// 实际应用：处理多种类型
func processAnimal(a Animal) {
	if dog, ok := a.(Dog); ok {
		fmt.Printf("处理狗: %s\n", dog.Name)
	} else if cat, ok := a.(Cat); ok {
		fmt.Printf("处理猫: %s\n", cat.Name)
	} else {
		fmt.Println("未知动物类型")
	}
}

