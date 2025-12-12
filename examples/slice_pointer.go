package main

import "fmt"

// 陷阱：切片中的指针问题
// 问题：切片中存储指针时，容易产生意外的行为

func main() {
	fmt.Println("=== 陷阱示例：切片中的指针问题 ===")
	
	// 陷阱1：切片中存储指针，共享同一个变量
	fmt.Println("\n陷阱1：切片中存储指针，共享同一个变量")
	trap1()
	
	// 陷阱2：切片扩容导致指针失效
	fmt.Println("\n陷阱2：切片扩容导致指针失效")
	trap2()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：在循环中创建指针切片
func trap1() {
	var pointers []*int
	
	// 错误：所有指针都指向同一个变量
	for i := 0; i < 3; i++ {
		pointers = append(pointers, &i) // 所有指针都指向 i
	}
	
	// 打印时，i 已经是循环结束后的值
	for _, p := range pointers {
		fmt.Printf("值: %d\n", *p) // 可能都打印 3
	}
}

// 陷阱2：切片扩容导致指针失效
func trap2() {
	// 创建初始切片
	slice := make([]*int, 0, 2)
	
	val1 := 1
	val2 := 2
	slice = append(slice, &val1, &val2)
	
	// 保存第一个元素的指针
	firstPtr := &slice[0]
	
	// 扩容可能导致底层数组重新分配
	val3 := 3
	val4 := 4
	val5 := 5
	slice = append(slice, &val3, &val4, &val5)
	
	// firstPtr 可能指向旧的底层数组
	fmt.Printf("第一个元素: %d\n", **firstPtr)
	fmt.Printf("切片第一个元素: %d\n", *slice[0])
}

// 正确方式1：在循环中创建新变量
func correctWay() {
	var pointers []*int
	
	// 正确：每次循环创建新变量
	for i := 0; i < 3; i++ {
		val := i // 创建局部变量
		pointers = append(pointers, &val)
	}
	
	for i, p := range pointers {
		fmt.Printf("索引 %d 的值: %d\n", i, *p)
	}
}

// 正确方式2：直接存储值而不是指针
func correctWay2() {
	// 如果不需要指针，直接存储值
	values := []int{1, 2, 3}
	
	for i, v := range values {
		fmt.Printf("索引 %d 的值: %d\n", i, v)
	}
}

// 正确方式3：使用函数创建指针
func correctWay3() {
	var pointers []*int
	
	for i := 0; i < 3; i++ {
		ptr := new(int) // 创建新指针
		*ptr = i
		pointers = append(pointers, ptr)
	}
	
	for i, p := range pointers {
		fmt.Printf("索引 %d 的值: %d\n", i, *p)
	}
}

// 实际应用：结构体切片
type Person struct {
	Name string
	Age  int
}

func demonstrateStructSlice() {
	// 错误：所有指针指向同一个结构体
	var people []*Person
	p := &Person{Name: "Alice", Age: 30}
	for i := 0; i < 3; i++ {
		people = append(people, p) // 所有元素指向同一个 Person
	}
	
	// 修改会影响所有元素
	people[0].Name = "Bob"
	fmt.Println(people[1].Name) // 也是 "Bob"
	
	// 正确：创建新的结构体
	var people2 []*Person
	for i := 0; i < 3; i++ {
		people2 = append(people2, &Person{
			Name: fmt.Sprintf("Person %d", i),
			Age:  20 + i,
		})
	}
}

