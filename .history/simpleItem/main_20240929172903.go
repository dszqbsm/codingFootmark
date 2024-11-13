package main

import (
	"errors"
	"fmt"
)

func divide(number, divisor float32) (float32, error) {
	if divisor == 0 {
		return 0.0, errors.New("除数不能为0")
	}
	return number / divisor, nil
}

func main() {
	var number, divisor float32
	fmt.Println("请输入一个数字")
	fmt.Scanln(&number)
	fmt.Println("请输入一个除数")
	fmt.Scanln(&divisor)

	result, err := divide(number, divisor)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
