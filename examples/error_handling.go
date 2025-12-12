package main

import (
	"errors"
	"fmt"
	"os"
)

// 陷阱：错误处理
// 问题：忽略错误或错误处理不当，导致程序行为异常

func main() {
	fmt.Println("=== 陷阱示例：错误处理 ===")
	
	// 陷阱1：忽略错误
	fmt.Println("\n陷阱1：忽略错误")
	trap1()
	
	// 陷阱2：错误比较不当
	fmt.Println("\n陷阱2：错误比较不当")
	trap2()
	
	// 陷阱3：错误包装丢失原始错误
	fmt.Println("\n陷阱3：错误包装")
	trap3()
	
	// 正确方式
	fmt.Println("\n正确方式：")
	correctWay()
}

// 陷阱1：忽略错误
func trap1() {
	// 错误：忽略错误
	file, _ := os.Open("不存在的文件.txt")
	defer file.Close() // 如果 file 是 nil，这里会 panic
	
	// 应该检查错误
	if file != nil {
		file.Close()
	}
}

// 陷阱2：使用 == 比较错误
func trap2() {
	err := doSomething()
	
	// 错误：直接比较错误值
	if err == errors.New("something went wrong") {
		// 这永远不会为 true，因为每次 errors.New 都创建新实例
		fmt.Println("错误匹配")
	}
	
	// 正确：使用 errors.Is 或定义错误变量
	var ErrSomething = errors.New("something went wrong")
	if err == ErrSomething {
		fmt.Println("错误匹配")
	}
}

func doSomething() error {
	return errors.New("something went wrong")
}

// 陷阱3：错误包装丢失上下文
func trap3() {
	err := processFile("test.txt")
	if err != nil {
		// 如果只是返回新错误，会丢失原始错误信息
		fmt.Printf("错误: %v\n", err)
	}
}

func processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// 错误：丢失原始错误
		return fmt.Errorf("无法处理文件")
		
		// 正确：包装原始错误
		// return fmt.Errorf("无法处理文件: %w", err)
	}
	defer file.Close()
	return nil
}

// 正确方式1：始终检查错误
func correctWay() {
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()
	
	// 继续处理文件
	fmt.Println("文件打开成功")
}

// 正确方式2：使用 errors.Is 和 errors.As
func correctWay2() {
	err := doSomething()
	
	// 使用 errors.Is 检查错误链
	if errors.Is(err, ErrSomething) {
		fmt.Println("是预期的错误")
	}
	
	// 使用 errors.As 提取特定类型的错误
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		fmt.Printf("路径错误: %s\n", pathErr.Path)
	}
}

var ErrSomething = errors.New("something went wrong")

// 正确方式3：错误包装和展开
func correctWay3() {
	err := processFile2("test.txt")
	if err != nil {
		// 使用 %w 包装错误
		wrapped := fmt.Errorf("处理失败: %w", err)
		fmt.Printf("包装后的错误: %v\n", wrapped)
		
		// 使用 errors.Unwrap 展开错误
		original := errors.Unwrap(wrapped)
		fmt.Printf("原始错误: %v\n", original)
	}
}

func processFile2(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	return nil
}

// 正确方式4：自定义错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validate(data string) error {
	if data == "" {
		return &ValidationError{
			Field:   "data",
			Message: "不能为空",
		}
	}
	return nil
}

func demonstrateCustomError() {
	err := validate("")
	if err != nil {
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("验证错误 - 字段: %s, 消息: %s\n", valErr.Field, valErr.Message)
		}
	}
}

// 最佳实践
func bestPractices() {
	// 1. 永远不要忽略错误
	// _, err := doSomething()
	// if err != nil {
	//     return err
	// }
	
	// 2. 使用 errors.Is 检查错误
	// if errors.Is(err, targetErr) {
	//     // 处理
	// }
	
	// 3. 使用 errors.As 提取错误类型
	// var targetErr *MyError
	// if errors.As(err, &targetErr) {
	//     // 使用 targetErr
	// }
	
	// 4. 使用 %w 包装错误，保留错误链
	// return fmt.Errorf("context: %w", err)
	
	// 5. 提供有意义的错误消息
	// return fmt.Errorf("failed to open %s: %w", filename, err)
}

