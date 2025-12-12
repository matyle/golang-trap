# Go 语言常见陷阱总结

本文档总结了 Go 语言开发中容易踩的坑，包括协程、指针、接口、通道等各个方面。每个陷阱都配有可运行的代码示例。

## 目录

1. [协程（Goroutines）陷阱](#1-协程goroutines陷阱)
2. [指针（Pointers）陷阱](#2-指针pointers陷阱)
3. [接口（Interfaces）陷阱](#3-接口interfaces陷阱)
4. [通道（Channels）陷阱](#4-通道channels陷阱)
5. [其他常见陷阱](#5-其他常见陷阱)

---

## 1. 协程（Goroutines）陷阱

### 1.1 闭包变量捕获问题

**问题**：在循环中使用 goroutine 时，所有 goroutine 可能共享同一个变量，导致意外的行为。

**错误示例**：
```go
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // 所有 goroutine 可能都打印 5
    }()
}
```

**正确示例**：
```go
for i := 0; i < 5; i++ {
    go func(val int) {
        fmt.Println(val) // 通过参数传递
    }(i)
}
```

**示例代码**：`examples/goroutine_closure.go`

### 1.2 未等待 Goroutine 完成

**问题**：主程序在 goroutine 完成前退出，导致 goroutine 被强制终止。

**错误示例**：
```go
for i := 0; i < 3; i++ {
    go func(id int) {
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("Goroutine %d 完成\n", id)
    }(i)
}
// 主程序立即退出，goroutine 可能还没执行完
```

**正确示例**：
```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("Goroutine %d 完成\n", id)
    }(i)
}
wg.Wait() // 等待所有 goroutine 完成
```

**示例代码**：`examples/goroutine_wait.go`

### 1.3 Goroutine 泄漏

**问题**：goroutine 因为通道阻塞而永远无法退出，造成内存泄漏。

**错误示例**：
```go
ch := make(chan int)
go func() {
    val := <-ch // 永远阻塞在这里，造成泄漏
    fmt.Println(val)
}()
// 没有人向通道发送数据，goroutine 永远无法退出
```

**正确示例**：
```go
ch := make(chan int, 1) // 带缓冲的通道
go func() {
    val := <-ch
    fmt.Println(val)
}()
ch <- 42 // 发送数据，goroutine 可以正常退出
```

**示例代码**：`examples/goroutine_leak.go`

---

## 2. 指针（Pointers）陷阱

### 2.1 Nil 指针解引用

**问题**：在使用指针前未检查是否为 nil，导致程序 panic。

**错误示例**：
```go
var p *int
fmt.Println(*p) // panic: nil pointer dereference
```

**正确示例**：
```go
var p *int
if p != nil {
    fmt.Println(*p)
} else {
    fmt.Println("指针为 nil")
}
```

**示例代码**：`examples/pointer_nil.go`

### 2.2 返回局部变量指针

**问题**：返回函数内部局部变量的指针，该变量在函数返回后可能被回收。

**注意**：在 Go 中，编译器会进行逃逸分析，通常会自动将变量分配到堆上，所以返回局部变量指针通常是安全的。但理解这个概念很重要。

**示例**：
```go
// Go 编译器会将 val 分配到堆上，这是安全的
func getPointer() *int {
    val := 42
    return &val
}

// 更安全的做法：返回值而不是指针
func getValue() int {
    val := 42
    return val
}
```

**示例代码**：`examples/pointer_local.go`

### 2.3 指针接收者 vs 值接收者

**问题**：混淆指针接收者和值接收者的使用场景，导致意外的行为。

**错误示例**：
```go
type Counter struct {
    value int
}

func (c Counter) Increment() {
    c.value++ // 只修改副本，不会修改原始值
}

c := Counter{value: 0}
c.Increment()
fmt.Println(c.value) // 仍然是 0
```

**正确示例**：
```go
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++ // 修改原始值
}

c := Counter{value: 0}
c.Increment()
fmt.Println(c.value) // 变为 1
```

**示例代码**：`examples/pointer_receiver.go`

---

## 3. 接口（Interfaces）陷阱

### 3.1 Nil 接口值

**问题**：接口值为 nil 但接口类型不为 nil，导致判断错误。

**错误示例**：
```go
var w Writer
var mw *MyWriter = nil

mw == nil  // true
w = mw
w == nil   // false! 因为接口包含类型信息
```

**正确示例**：
```go
var w Writer
var mw *MyWriter = nil

if mw != nil {
    w = mw
}

// 或者使用类型断言检查
if w != nil {
    if mw, ok := w.(*MyWriter); ok && mw != nil {
        // 安全使用
    }
}
```

**示例代码**：`examples/interface_nil.go`

### 3.2 接口类型断言

**问题**：类型断言失败时未检查 ok 值，导致 panic。

**错误示例**：
```go
var a Animal = Dog{}
cat := a.(Cat) // panic: interface conversion
```

**正确示例**：
```go
var a Animal = Dog{}

// 方式1：使用 ok 值检查
dog, ok := a.(Dog)
if ok {
    fmt.Println("是 Dog")
}

// 方式2：使用 type switch
switch v := a.(type) {
case Dog:
    fmt.Println("是 Dog")
case Cat:
    fmt.Println("是 Cat")
default:
    fmt.Println("未知类型")
}
```

**示例代码**：`examples/interface_assertion.go`

### 3.3 空接口的使用

**问题**：过度使用空接口 `interface{}`，失去类型安全。

**错误示例**：
```go
var data interface{}
data = 42
str := data.(string) // 如果 data 不是 string，会 panic
```

**正确示例**：
```go
// 使用具体类型
var data string
data = "hello"

// 或者使用类型断言时检查
if str, ok := data.(string); ok {
    fmt.Println(str)
}

// Go 1.18+ 使用泛型
func process[T any](data T) T {
    return data
}
```

**示例代码**：`examples/interface_empty.go`

---

## 4. 通道（Channels）陷阱

### 4.1 未关闭通道导致泄漏

**问题**：通道未正确关闭，导致接收方永远阻塞。

**错误示例**：
```go
ch := make(chan int)
go func() {
    for i := 0; i < 3; i++ {
        ch <- i
    }
    // 忘记关闭通道！
}()
// 接收方会一直等待
```

**正确示例**：
```go
ch := make(chan int)
go func() {
    defer close(ch) // 确保通道被关闭
    for i := 0; i < 3; i++ {
        ch <- i
    }
}()

for val := range ch { // range 会在通道关闭时自动退出
    fmt.Println(val)
}
```

**示例代码**：`examples/channel_close.go`

### 4.2 向已关闭通道发送数据

**问题**：向已关闭的通道发送数据会导致 panic。

**错误示例**：
```go
ch := make(chan int)
close(ch)
ch <- 42 // panic: send on closed channel
```

**正确示例**：
```go
ch := make(chan int)
var once sync.Once

go func() {
    for i := 0; i < 3; i++ {
        ch <- i
    }
    once.Do(func() {
        close(ch) // 只关闭一次
    })
}()
```

**示例代码**：`examples/channel_send_closed.go`

### 4.3 从已关闭通道读取

**问题**：从已关闭的通道读取会立即返回零值，需要检查通道状态。

**错误示例**：
```go
ch := make(chan int)
close(ch)
val := <-ch
if val == 0 {
    // 无法区分是零值还是通道关闭
}
```

**正确示例**：
```go
ch := make(chan int)
close(ch)

// 使用两个返回值检查
val, ok := <-ch
if !ok {
    fmt.Println("通道已关闭")
}

// 或使用 range
for val := range ch {
    fmt.Println(val)
}
```

**示例代码**：`examples/channel_receive_closed.go`

### 4.4 Select 的 Default Case

**问题**：select 语句中的 default case 可能导致非阻塞行为，影响程序逻辑。

**错误示例**：
```go
ch := make(chan int)
go func() {
    time.Sleep(50 * time.Millisecond)
    ch <- 42
}()

select {
case val := <-ch:
    fmt.Println(val)
default:
    fmt.Println("没有数据，立即返回") // 可能错过数据
}
```

**正确示例**：
```go
// 不使用 default，等待数据
select {
case val := <-ch:
    fmt.Println(val)
}

// 或使用 default 实现超时
select {
case val := <-ch:
    fmt.Println(val)
case <-time.After(100 * time.Millisecond):
    fmt.Println("超时")
}
```

**示例代码**：`examples/channel_select_default.go`

---

## 5. 其他常见陷阱

### 5.1 切片和数组的区别

**问题**：混淆切片和数组，导致意外的行为。

**错误示例**：
```go
// 数组是值类型，赋值会复制
arr1 := [3]int{1, 2, 3}
arr2 := arr1
arr2[0] = 99
// arr1 仍然是 [1 2 3]

// 切片是引用类型，共享底层数组
slice1 := []int{1, 2, 3}
slice2 := slice1
slice2[0] = 99
// slice1 也变成了 [99 2 3]
```

**正确示例**：
```go
// 创建独立切片副本
original := []int{1, 2, 3}
independent := make([]int, len(original))
copy(independent, original)
independent[0] = 99
// original 仍然是 [1 2 3]
```

**示例代码**：`examples/slice_array.go`

### 5.2 Map 的并发读写

**问题**：多个 goroutine 同时读写 map 会导致 panic。

**错误示例**：
```go
m := make(map[string]int)
go func() {
    for i := 0; i < 1000; i++ {
        m["key"] = i // 并发写入
    }
}()
go func() {
    for i := 0; i < 1000; i++ {
        _ = m["key"] // 并发读取，会 panic
    }
}()
```

**正确示例**：
```go
// 方式1：使用 Mutex 保护
m := make(map[string]int)
var mu sync.RWMutex

go func() {
    mu.Lock()
    m["key"] = 1
    mu.Unlock()
}()

go func() {
    mu.RLock()
    val := m["key"]
    mu.RUnlock()
}()

// 方式2：使用 sync.Map
var m sync.Map
m.Store("key", 1)
val, _ := m.Load("key")
```

**示例代码**：`examples/map_concurrent.go`

### 5.3 Defer 的执行顺序

**问题**：defer 语句的执行顺序和参数求值时机容易混淆。

**错误示例**：
```go
i := 0
defer fmt.Println("defer:", i) // i 的值是 0（立即求值）
i++
fmt.Println("函数结束:", i) // i 的值是 1
// 输出：函数结束: 1
//      defer: 0
```

**正确示例**：
```go
i := 0
defer func() {
    fmt.Println("defer:", i) // 使用闭包，访问最新的 i
}()
i++
fmt.Println("函数结束:", i)
// 输出：函数结束: 1
//      defer: 1
```

**注意**：defer 的执行顺序是 LIFO（后进先出），defer 可以修改命名返回值。

**示例代码**：`examples/defer_order.go`

### 5.4 错误处理

**问题**：忽略错误或错误处理不当，导致程序行为异常。

**错误示例**：
```go
// 忽略错误
file, _ := os.Open("file.txt")
defer file.Close() // 如果 file 是 nil，会 panic

// 错误比较不当
if err == errors.New("error") {
    // 永远不会为 true
}
```

**正确示例**：
```go
// 始终检查错误
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("打开文件失败: %w", err)
}
defer file.Close()

// 使用 errors.Is 检查错误
if errors.Is(err, ErrNotFound) {
    // 处理
}

// 使用 errors.As 提取错误类型
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Path)
}
```

**示例代码**：`examples/error_handling.go`

---

## 运行示例

每个示例文件都可以独立运行：

```bash
go run examples/goroutine_closure.go
go run examples/pointer_nil.go
# ... 等等
```

或者运行所有示例：

```bash
./run_all_examples.sh
```

---

## 总结

Go 语言虽然简洁，但在并发、指针、接口等方面有很多细节需要注意。理解这些陷阱可以帮助写出更安全、更可靠的代码。

**关键要点**：
- 始终检查指针是否为 nil
- 在循环中使用 goroutine 时注意变量捕获
- 正确关闭通道，避免泄漏
- 理解接口的 nil 值行为
- 使用 sync 包保护共享资源
- 正确处理错误，不要忽略

