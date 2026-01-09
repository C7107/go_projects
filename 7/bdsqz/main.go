package main

import (
	"fmt"
	"math"
	"strings"
)

// --- 1. 定义接口和基本类型 ---

// Expr 接口：代表任意一个表达式积木
type Expr interface {
	// Eval: 给定环境变量 env，计算该表达式的值
	Eval(env Env) float64
	// Check: 静态检查表达式是否有错误，并将用到的变量名加入 vars 集合
	Check(vars map[Var]bool) error
}

// Var: 代表变量，如 "x", "pi"
type Var string

// literal: 代表数字常量，如 3.14
type literal float64

// unary: 代表一元操作，如 -x
type unary struct {
	op rune // '+' 或 '-'
	x  Expr
}

// binary: 代表二元操作，如 x + y
type binary struct {
	op   rune // '+', '-', '*', '/'
	x, y Expr
}

// call: 代表函数调用，如 sin(x)
type call struct {
	fn   string // "pow", "sin", "sqrt"
	args []Expr
}

// Env: 环境变量，用于存储变量的值 (例如: x=10, y=5)
type Env map[Var]float64

// --- 2. 实现 Eval 方法 (计算逻辑) ---

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

// --- 3. 实现 Check 方法 (静态检查逻辑) ---

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true // 记录这个变量被用到了
	return nil
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

// 定义允许的函数和参数个数
var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

// --- 4. 主程序运行演示 ---

func main() {
	// === 场景 1: 计算华氏度转摄氏度 ===
	// 表达式: 5 / 9 * (F - 32)
	// 因为我们没有 Parse 函数，所以要手动把这棵树“搭”出来

	// 1. 构建表达式树
	// 结构: binary('*', binary('/', 5, 9), binary('-', F, 32))
	expr1 := binary{
		op: '*',
		x: binary{
			op: '/',
			x:  literal(5),
			y:  literal(9),
		},
		y: binary{
			op: '-',
			x:  Var("F"),
			y:  literal(32),
		},
	}

	// 2. 静态检查
	vars := make(map[Var]bool)
	if err := expr1.Check(vars); err != nil {
		fmt.Printf("check error: %v\n", err)
		return
	}
	fmt.Printf("表达式1 需要的变量: %v\n", vars)

	// 3. 运行求值 (F = 212)
	env1 := Env{"F": 212}
	result1 := expr1.Eval(env1)
	fmt.Printf("当 F=212 时, 5/9*(F-32) = %.2f (预期 100)\n", result1)

	fmt.Println("------------------------------------------------")

	// === 场景 2: 计算 sqrt(A / pi) ===

	// 1. 构建表达式树
	// 结构: call("sqrt", [binary('/', A, pi)])
	expr2 := call{
		fn: "sqrt",
		args: []Expr{
			binary{
				op: '/',
				x:  Var("A"),
				y:  Var("pi"),
			},
		},
	}

	// 2. 静态检查
	vars2 := make(map[Var]bool)
	if err := expr2.Check(vars2); err != nil {
		fmt.Printf("check error: %v\n", err)
		return
	}
	fmt.Printf("表达式2 需要的变量: %v\n", vars2)

	// 3. 运行求值 (A = 87616, pi = 3.14...)
	env2 := Env{"A": 87616, "pi": math.Pi}
	result2 := expr2.Eval(env2)
	fmt.Printf("当 A=87616 时, sqrt(A/pi) = %.2f (预期 167)\n", result2)

	fmt.Println("------------------------------------------------")

	// === 场景 3: 演示错误检查 (Check 的作用) ===

	// 构造一个错误的表达式: unknown(x) —— 这是一个不存在的函数
	exprError := call{
		fn:   "unknown",
		args: []Expr{Var("x")},
	}

	fmt.Println("正在检查一个错误的表达式: unknown(x)...")
	if err := exprError.Check(make(map[Var]bool)); err != nil {
		fmt.Printf("检查成功拦截错误: %v\n", err)
	} else {
		// 如果没检查直接 Eval，这里就会 panic
		exprError.Eval(Env{"x": 1})
	}
}
