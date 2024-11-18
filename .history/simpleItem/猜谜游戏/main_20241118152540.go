package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano()) // 将程序启动的时间戳来初始化随机数种子
	secretNum := rand.Intn(maxNum)
	fmt.Println("The secret number is: ", secretNum)
}
