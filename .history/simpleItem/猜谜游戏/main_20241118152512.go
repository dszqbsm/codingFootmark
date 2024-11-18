package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	secretNum := rand.Intn(maxNum)
	fmt.Println("The secret number is: ", secretNum)
}
