package main

import "fmt"

// ----------------------
// 1. 定义接口 (Interface)
// ----------------------
// 这就是“职位描述”。
// 我们规定：任何东西，只要拥有 Jiao() 这个方法，它就是一个 "Jiaoer"（会叫者）。
type Jiaoer interface {
	Jiao()
}

// ----------------------
// 2. 定义具体类型 (Concrete Types)
// ----------------------
// ---- 第一种动物：狗 ----
type Dog struct {
	Name string
}

// 狗实现了 Jiao() 方法
// Go 语言潜规则：只要实现了接口里的方法，就自动算作实现了该接口。
// 不需要写 "Dog implements Jiaoer"。
func (d Dog) Jiao() {
	fmt.Printf("【%s】被打了：汪汪汪！\n", d.Name)
}

// ---- 第二种动物：猫 ----
type Cat struct {
	Name string
}

// 猫也实现了 Jiao() 方法
func (c Cat) Jiao() {
	fmt.Printf("【%s】被打了：喵喵喵！\n", c.Name)
}

// ---- 第三种动物：鸭子 (再加一个为了演示扩展性) ----
type Duck struct {
	Name string
}

func (d Duck) Jiao() {
	fmt.Printf("【%s】被打了：嘎嘎嘎！\n", d.Name)
}

// ----------------------
// 3. 定义通用函数 (Polymorphism)
// ----------------------
// 这个函数是“通用打击工具”。
// 注意参数类型：它接收的是 Jiaoer 接口，而不是具体的 Dog 或 Cat。
// 这意味着：只要是“会叫的”，不管是啥，都能传进来。
func Hit(animal Jiaoer) {
	fmt.Println("--> 准备发起攻击...")

	// 这里发生了神奇的事情（多态）：
	// 如果传进来的是狗，这里就调狗的 Jiao()
	// 如果传进来的是猫，这里就调猫的 Jiao()
	animal.Jiao()

	fmt.Println("--> 攻击结束。")
	fmt.Println() // 打个空行
}

// ----------------------
// 4. 主函数 (Main)
// ----------------------
func main() {
	// 1. 创建几只具体的动物
	dog := Dog{Name: "旺财"}
	cat := Cat{Name: "咪咪"}
	duck := Duck{Name: "唐老鸭"}
	// 2. 调用同一个 Hit 函数，传入不同的动物
	// 编译器会自动检查：dog 有 Jiao() 方法吗？有！那就允许传入。
	Hit(dog)
	Hit(cat)
	Hit(duck)
	// 3. 甚至可以搞一个“动物园”切片
	// 这个切片里装的都是 Jiaoer 接口类型
	zoo := []Jiaoer{dog, cat, duck}

	fmt.Println("--- 动物园集体大合唱 ---")
	for _, animal := range zoo {
		animal.Jiao()
	}
}
