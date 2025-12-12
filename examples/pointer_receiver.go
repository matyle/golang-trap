package main

import "fmt"

// 陷阱：指针接收者 vs 值接收者
// 问题：混淆指针接收者和值接收者的使用场景，导致意外的行为

type Counter struct {
	value int
}

// 值接收者：不会修改原始值
func (c Counter) IncrementByValue() {
	c.value++ // 只修改副本
}

// 指针接收者：会修改原始值
func (c *Counter) IncrementByPointer() {
	c.value++ // 修改原始值
}

// 值接收者：返回新值
func (c Counter) GetValue() int {
	return c.value
}

func main() {
	fmt.Println("=== 陷阱示例：指针接收者 vs 值接收者 ===")
	
	// 示例1：值接收者不会修改原始值
	fmt.Println("\n示例1：值接收者")
	c1 := Counter{value: 0}
	c1.IncrementByValue()
	fmt.Printf("调用值接收者后: %d\n", c1.GetValue()) // 仍然是 0
	
	// 示例2：指针接收者会修改原始值
	fmt.Println("\n示例2：指针接收者")
	c2 := Counter{value: 0}
	c2.IncrementByPointer()
	fmt.Printf("调用指针接收者后: %d\n", c2.GetValue()) // 变为 1
	
	// 示例3：值类型调用指针接收者方法（Go 会自动转换）
	fmt.Println("\n示例3：值类型调用指针接收者方法")
	c3 := Counter{value: 0}
	c3.IncrementByPointer() // Go 会自动转换为 (&c3).IncrementByPointer()
	fmt.Printf("值类型调用指针接收者后: %d\n", c3.GetValue()) // 变为 1
	
	// 示例4：指针类型调用值接收者方法（Go 会自动解引用）
	fmt.Println("\n示例4：指针类型调用值接收者方法")
	c4 := &Counter{value: 0}
	c4.IncrementByValue() // Go 会自动转换为 (*c4).IncrementByValue()
	fmt.Printf("指针类型调用值接收者后: %d\n", c4.GetValue()) // 仍然是 0
	
	// 陷阱：接口实现
	fmt.Println("\n陷阱：接口实现")
	demonstrateInterface()
}

// 接口实现陷阱
type Incrementer interface {
	Increment()
	GetValue() int
}

type ValueCounter struct {
	value int
}

// 值接收者实现接口
func (v ValueCounter) Increment() {
	v.value++
}

func (v ValueCounter) GetValue() int {
	return v.value
}

type PointerCounter struct {
	value int
}

// 指针接收者实现接口
func (p *PointerCounter) Increment() {
	p.value++
}

func (p *PointerCounter) GetValue() int {
	return p.value
}

func demonstrateInterface() {
	// 值类型可以实现接口
	var v1 Incrementer = ValueCounter{value: 0}
	v1.Increment()
	fmt.Printf("值接收者接口: %d\n", v1.GetValue()) // 仍然是 0
	
	// 指针类型也可以实现接口
	var v2 Incrementer = &PointerCounter{value: 0}
	v2.Increment()
	fmt.Printf("指针接收者接口: %d\n", v2.GetValue()) // 变为 1
	
	// 陷阱：值类型不能赋值给需要指针接收者的接口
	// var v3 Incrementer = PointerCounter{value: 0} // 编译错误！
	// 必须使用指针：
	var v3 Incrementer = &PointerCounter{value: 0}
	v3.Increment()
	fmt.Printf("指针接收者接口（正确用法）: %d\n", v3.GetValue())
}

