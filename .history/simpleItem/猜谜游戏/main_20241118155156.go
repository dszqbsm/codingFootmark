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
	rand.Seed(time.Now().UnixNano()) // 将程序启动的时间戳来初始化随机数种子
	secretNum := rand.Intn(maxNum)
	fmt.Println("The secret number is: ", secretNum)

	fmt.Println("Please input your guess: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input, Please try again", err)
		return
	}
	input = strings.TrimSuffix(input, "\r")

	guess, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invaild input, please enter an integer values", err)
		return
	}
	fmt.Println("You guessed is: ", guess)
}
