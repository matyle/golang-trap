package main

import (
	"fmt"
	"reflect"
)

// 陷阱：空接口的使用
// 问题：过度使用空接口 interface{}，失去类型安全

func main() {
	fmt.Println("=== 陷阱示例：空接口的使用 ===")
	
	// 陷阱：失去类型安全
	fmt.Println("\n陷阱：失去类型安全")
	trap1()
	
	// 正确方式：使用泛型（Go 1.18+）或具体类型
	fmt.Println("\n正确方式：使用具体类型或泛型")
	correctWay()
	
	// 实际应用：JSON 处理
	fmt.Println("\n实际应用：JSON 处理")
	jsonExample()
}

// 陷阱：使用空接口失去类型安全
func trap1() {
	// 可以存储任何类型
	var data interface{}
	
	data = 42
	fmt.Printf("整数: %v, 类型: %T\n", data, data)
	
	data = "hello"
	fmt.Printf("字符串: %v, 类型: %T\n", data, data)
	
	data = []int{1, 2, 3}
	fmt.Printf("切片: %v, 类型: %T\n", data, data)
	
	// 问题：使用时需要类型断言，容易出错
	// str := data.(string) // 如果 data 不是 string，会 panic
}

// 正确方式1：使用具体类型
func correctWay() {
	// 使用具体类型，编译时检查
	var data string
	data = "hello"
	// data = 42 // 编译错误！
	fmt.Printf("字符串: %s\n", data)
}

// 正确方式2：使用泛型（Go 1.18+）
func correctWay2[T any](data T) T {
	return data
}

// 正确方式3：使用类型断言时检查
func safeTypeAssertion(data interface{}) {
	if str, ok := data.(string); ok {
		fmt.Printf("是字符串: %s\n", str)
	} else if num, ok := data.(int); ok {
		fmt.Printf("是整数: %d\n", num)
	} else {
		fmt.Printf("未知类型: %T\n", data)
	}
}

// 实际应用：JSON 处理
func jsonExample() {
	// JSON 解析时经常使用 map[string]interface{}
	jsonData := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	
	// 安全访问
	if name, ok := jsonData["name"].(string); ok {
		fmt.Printf("姓名: %s\n", name)
	}
	
	if age, ok := jsonData["age"].(float64); ok {
		fmt.Printf("年龄: %.0f\n", age)
	}
	
	// 更好的方式：定义结构体
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}
	
	// 使用结构体解析 JSON，类型安全
	fmt.Println("使用结构体更安全")
}

// 使用反射处理空接口（复杂但灵活）
func reflectExample(data interface{}) {
	v := reflect.ValueOf(data)
	fmt.Printf("类型: %v, 种类: %v\n", v.Type(), v.Kind())
	
	switch v.Kind() {
	case reflect.Int:
		fmt.Printf("整数值: %d\n", v.Int())
	case reflect.String:
		fmt.Printf("字符串值: %s\n", v.String())
	case reflect.Slice:
		fmt.Printf("切片长度: %d\n", v.Len())
	default:
		fmt.Println("其他类型")
	}
}

