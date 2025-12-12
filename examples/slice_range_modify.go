package main

import "fmt"

// 陷阱：切片遍历时修改
// 问题：在遍历切片时修改切片，导致意外的行为

func main() {
	fmt.Println("=== 陷阱示例：切片遍历时修改 ===")
	
	// 陷阱1：遍历时修改元素（值类型）
	fmt.Println("\n陷阱1：遍历时修改元素（值类型）")
	trap1()
	
	// 陷阱2：遍历时添加/删除元素
	fmt.Println("\n陷阱2：遍历时添加/删除元素")
	trap2()
	
	// 陷阱3：遍历时修改底层数组
	fmt.Println("\n陷阱3：遍历时修改底层数组")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：遍历时修改元素（值类型）
func trap1() {
	slice := []int{1, 2, 3, 4, 5}
	
	// 错误：修改的是副本，不会影响原切片
	for _, v := range slice {
		v *= 2 // 只修改副本
	}
	fmt.Printf("修改后: %v\n", slice) // [1 2 3 4 5]，没有变化
}

// 陷阱2：遍历时添加/删除元素
func trap2() {
	slice := []int{1, 2, 3, 4, 5}
	
	// 危险：在遍历时修改切片长度
	for i, v := range slice {
		if v%2 == 0 {
			// 删除元素（错误的方式）
			slice = append(slice[:i], slice[i+1:]...)
			// 这会导致索引错乱和未遍历的元素
		}
	}
	fmt.Printf("修改后: %v\n", slice) // 结果不确定
}

// 陷阱3：遍历时修改底层数组
func trap3() {
	original := []int{1, 2, 3, 4, 5}
	slice := original[1:4] // [2 3 4]
	
	// 遍历 slice，但修改会影响 original
	for i := range slice {
		slice[i] *= 10
	}
	
	fmt.Printf("original: %v\n", original) // [1 20 30 40 5]
	fmt.Printf("slice: %v\n", slice)      // [20 30 40]
}

// 正确方式1：使用索引修改元素
func correctWay() {
	slice := []int{1, 2, 3, 4, 5}
	
	// 正确：使用索引修改
	for i := range slice {
		slice[i] *= 2
	}
	fmt.Printf("修改后: %v\n", slice) // [2 4 6 8 10]
}

// 正确方式2：使用指针遍历
func correctWay2() {
	slice := []*int{}
	for i := 0; i < 5; i++ {
		val := i
		slice = append(slice, &val)
	}
	
	// 通过指针修改
	for _, p := range slice {
		*p *= 2
	}
	
	for _, p := range slice {
		fmt.Printf("%d ", *p)
	}
	fmt.Println()
}

// 正确方式3：先收集要删除的索引，再删除
func correctWay3() {
	slice := []int{1, 2, 3, 4, 5}
	
	// 先收集要删除的索引
	var toDelete []int
	for i, v := range slice {
		if v%2 == 0 {
			toDelete = append(toDelete, i)
		}
	}
	
	// 从后往前删除，避免索引错乱
	for i := len(toDelete) - 1; i >= 0; i-- {
		idx := toDelete[i]
		slice = append(slice[:idx], slice[idx+1:]...)
	}
	
	fmt.Printf("删除偶数后: %v\n", slice) // [1 3 5]
}

// 正确方式4：创建新切片
func correctWay4() {
	slice := []int{1, 2, 3, 4, 5}
	
	// 创建新切片，不修改原切片
	newSlice := make([]int, 0, len(slice))
	for _, v := range slice {
		if v%2 != 0 {
			newSlice = append(newSlice, v)
		}
	}
	
	fmt.Printf("原切片: %v\n", slice)     // [1 2 3 4 5]
	fmt.Printf("新切片: %v\n", newSlice)   // [1 3 5]
}

// 注意事项
func notes() {
	// 1. range 遍历时，v 是元素的副本
	// 2. 修改 v 不会影响原切片
	// 3. 使用 slice[i] 可以修改原切片
	// 4. 在遍历时添加/删除元素是危险的
	// 5. 如果需要修改，先收集索引，再统一处理
}

