package main

import "fmt"

// 陷阱：Nil 接口值
// 问题：接口值为 nil 但接口类型不为 nil，导致判断错误

type Writer interface {
	Write([]byte) (int, error)
}

type MyWriter struct {
	data []byte
}

func (w *MyWriter) Write(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func main() {
	fmt.Println("=== 陷阱示例：Nil 接口值 ===")
	
	// 陷阱1：接口值为 nil，但接口本身不为 nil
	fmt.Println("\n陷阱1：接口值为 nil，但接口本身不为 nil")
	trap1()
	
	// 陷阱2：nil 指针实现接口
	fmt.Println("\n陷阱2：nil 指针实现接口")
	trap2()
	
	// 正确方式：检查接口值和类型
	fmt.Println("\n正确方式：检查接口值和类型")
	correctWay()
}

func trap1() {
	var w Writer
	var mw *MyWriter = nil
	
	// mw 是 nil 指针
	fmt.Printf("mw == nil: %v\n", mw == nil) // true
	
	// 但是将 nil 指针赋值给接口后，接口不为 nil
	w = mw
	fmt.Printf("w == nil: %v\n", w == nil) // false!
	
	// 因为接口包含类型信息 (*MyWriter) 和值 (nil)
	// 所以接口本身不为 nil
}

func trap2() {
	var w Writer = (*MyWriter)(nil)
	
	// 接口不为 nil
	if w != nil {
		fmt.Println("接口不为 nil，可以调用方法")
		// 但是调用方法会 panic，因为底层值是 nil
		// w.Write([]byte("test")) // panic: runtime error: invalid memory address
	}
}

func correctWay() {
	var w Writer
	var mw *MyWriter = nil
	
	// 方式1：在赋值前检查
	if mw != nil {
		w = mw
	}
	
	// 方式2：使用类型断言检查
	w = mw
	if w != nil {
		if mw, ok := w.(*MyWriter); ok && mw != nil {
			mw.Write([]byte("safe"))
			fmt.Println("安全调用")
		} else {
			fmt.Println("接口值或类型为 nil，不能调用")
		}
	}
	
	// 方式3：使用反射检查（更复杂但更准确）
	// import "reflect"
	// if w != nil && reflect.ValueOf(w).IsNil() {
	//     // 处理 nil 情况
	// }
}

// 实际应用：错误处理
type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	if e == nil {
		return "nil error"
	}
	return e.msg
}

func returnError() error {
	var err *MyError = nil
	return err // 返回的 error 接口不为 nil！
}

func demonstrateError() {
	err := returnError()
	if err != nil {
		fmt.Println("错误不为 nil") // 会执行这里
		fmt.Println(err.Error())    // 输出 "nil error"
	}
}

