/*
第一关：鸭子类型（Duck Typing）与多态
知识点：定义接口、实现接口、多态函数。
题目描述：
定义一个接口 USB，它要求有一个方法 Connect()。
定义两个结构体：Phone（手机）和 Camera（相机）。
让 Phone 连接时打印 "手机已连接"，让 Camera 连接时打印 "相机已连接"。
写一个函数 Boot(device USB)，在这个函数里调用 Connect()。
在 main 函数里，分别把手机和相机传给 Boot 函数运行。
*/
package main

import (
	"fmt"
)

type USB interface {
	Connect()
}

type Phone struct {
	Name string
}

type Camera struct {
	Name string
}

func (p Phone) Connect() {
	fmt.Println("手机已连接")
}

func (c Camera) Connect() {
	fmt.Println("相机已连接")
}

func Boot(device USB) {
	device.Connect()
}

func main() {
	phone := Phone{Name: "华为"}
	camera := Camera{Name: "尼康"}
	Boot(phone)
	Boot(camera)
}
