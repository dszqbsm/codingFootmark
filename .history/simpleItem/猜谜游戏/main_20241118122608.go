package main

import (
	"fmt"
	"math/rand"
)

func main() {
	maxNum := 100
	secretNum := rand.Intn(maxNum) + 1
	fmt.Println("The secret number is: ", secretNum)
}
