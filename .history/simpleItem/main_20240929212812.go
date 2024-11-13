package main

import (
	"fmt"
)

type Person struct {
	name  string
	age   int
	email string
}

func main() {
	// 定义一个结构体 Person 包含姓名、年龄和邮箱，编写一个程序创建并打印多个 Person 实例
	var zhangsan Person = Person{name: "张三", age: 20, email: "123456@qq.com"}
	var lisi Person
	lisi.name = "李四"
	lisi.age = 22
	lisi.email = "123@qq.com"
	fmt.Println(zhangsan)
	fmt.Println(lisi)

}
