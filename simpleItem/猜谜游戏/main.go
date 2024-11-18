package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano()) // 用程序启动的时间戳来初始化随机数种子
	secretNum := rand.Intn(maxNum)   // 生成一个随机数
	// fmt.Println("The secret number is: ", secretNum)

	for {
		fmt.Println("Please input your guess: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n') // 读取用户输入，直到用户输入回车键
		if err != nil {
			fmt.Println("An error occured while reading input, Please try again", err)
			continue
		}
		input = strings.TrimSuffix(input, "\r\n") // 去除换行符，这里注意用户输入的enter一般是\r\n

		guess, err := strconv.Atoi(input) // 将用户输入的字符串转换为整数
		if err != nil {
			fmt.Println("Invaild input, please enter an integer values", err)
			continue
		}
		fmt.Println("You guessed is: ", guess)

		if guess > secretNum {
			fmt.Println("You guess is bigger than the secret number, please try again")
		} else if guess < secretNum {
			fmt.Println("You guess is smaller than the secret number, please try again")
		} else {
			fmt.Println("Correct, you win!")
			break
		}
	}
}
