package main

import "fmt"

// 陷阱：Interface 接收者问题
// 问题：接口方法接收者的选择影响接口实现

type Writer interface {
	Write([]byte) (int, error)
}

type MyWriter struct {
	data []byte
}

// 值接收者实现接口
func (w MyWriter) Write(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

// 指针接收者实现接口
func (w *MyWriter) WritePointer(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

type PointerWriter interface {
	WritePointer([]byte) (int, error)
}

func main() {
	fmt.Println("=== 陷阱示例：Interface 接收者问题 ===")

	// 陷阱1：值接收者 vs 指针接收者
	fmt.Println("\n陷阱1：值接收者 vs 指针接收者")
	trap1()

	// 陷阱2：接口赋值问题
	fmt.Println("\n陷阱2：接口赋值问题")
	trap2()

	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：值接收者实现接口
func trap1() {
	var w Writer

	// 值类型可以实现接口
	mw1 := MyWriter{}
	w = mw1
	w.Write([]byte("test"))
	fmt.Printf("值接收者: %v\n", mw1.data) // 空，因为修改的是副本

	// 指针类型也可以实现接口（Go 自动转换）
	mw2 := &MyWriter{}
	w = mw2
	w.Write([]byte("test"))
	fmt.Printf("指针类型调用值接收者: %v\n", mw2.data) // 仍然是空
}

// 陷阱2：指针接收者实现接口
func trap2() {
	var pw PointerWriter

	// 值类型不能赋值给需要指针接收者的接口
	// mw1 := MyWriter{}
	// pw = mw1 // 编译错误！

	// 必须使用指针
	mw2 := &MyWriter{}
	pw = mw2
	pw.WritePointer([]byte("test"))
	fmt.Printf("指针接收者: %v\n", mw2.data) // 有数据
}

// Counter 类型定义
type Counter struct {
	count int
}

// 指针接收者：可以修改
func (c *Counter) Increment() {
	c.count++
}

// 值接收者：不能修改
func (c Counter) GetCount() int {
	return c.count
}

// 正确方式：根据需求选择接收者类型
func correctWay() {
	// 如果方法需要修改接收者，使用指针接收者
	c := Counter{}
	c.Increment()
	fmt.Printf("计数: %d\n", c.GetCount())
}

// 实际应用：接口设计原则
func designPrinciples() {
	// 1. 如果方法需要修改接收者，使用指针接收者
	// 2. 如果接收者是大结构体，使用指针接收者（避免复制）
	// 3. 如果接收者是值类型（如 int, string），使用值接收者
	// 4. 保持一致性：一个类型的所有方法应该使用相同的接收者类型
}

// 陷阱：接口方法集
type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type File struct {
	data []byte
}

// 值接收者实现 Read
func (f File) Read(p []byte) (int, error) {
	copy(p, f.data)
	return len(f.data), nil
}

// 指针接收者实现 Write
func (f *File) Write(p []byte) (int, error) {
	f.data = append(f.data, p...)
	return len(p), nil
}

func demonstrateMethodSet() {
	// 值类型的方法集只包含值接收者的方法
	// var f1 File
	// var rw1 ReadWriter = f1 // 编译错误！f1 没有实现 Write

	// 指针类型的方法集包含值接收者和指针接收者的方法
	var f2 *File = &File{}
	var rw2 ReadWriter = f2 // 可以，因为 *File 实现了所有方法

	_ = rw2
}
