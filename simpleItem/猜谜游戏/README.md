# 猜谜游戏

该小项目是字节跳动青训营GO入门课程中的实战项目<https://juejin.cn/course/bytetech/7140987981803814919/section/7141265019756675103>

## 需求分析

1. 程序生成0-100随机整数
2. 玩家每次需要输入一个猜测数字
3. 程序需要告诉玩家该猜测数字是高于还是低于谜底随机数
4. 玩家经过反复猜测，若猜对则玩家游戏胜利，退出程序

## 实现思路

1. 生成随机数

```go
    // 生成随机数v1版本
	maxNum := 100
	secretNum := rand.Intn(maxNum)
	fmt.Println("The secret number is: ", secretNum)
```

需要用程序启动的时间戳来初始化随机数种子，不然生成的随机数会相同

```go
	maxNum := 100
	rand.Seed(time.Now().UnixNano()) // 将程序启动的时间戳来初始化随机数种子
	secretNum := rand.Intn(maxNum)
	fmt.Println("The secret number is: ", secretNum)
```

2. 读取用户输入

使用os库来控制输入，需要将输入转化成只读的流，这样才有更多的操作手段，可以从流中读取一行，但是每次读取行的末尾会多出一个换行符，需要单独删除该换行符，最后需要将该流转成一个数字，这样才最终得到用户输入的数字

```go
	maxNum := 100
	rand.Seed(time.Now().UnixNano()) // 用程序启动的时间戳来初始化随机数种子
	secretNum := rand.Intn(maxNum)	// 生成一个随机数
	fmt.Println("The secret number is: ", secretNum)

	fmt.Println("Please input your guess: ")
	reader := bufio.NewReader(os.Stdin)	
	input, err := reader.ReadString('\n')	// 读取用户输入，直到用户输入回车键
	if err != nil {
		fmt.Println("An error occured while reading input, Please try again", err)
		return
	}
	input = strings.TrimSuffix(input, "\r\n")	// 去除换行符，这里注意用户输入的enter一般是\r\n

	guess, err := strconv.Atoi(input)	// 将用户输入的字符串转换为整数
	if err != nil {
		fmt.Println("Invaild input, please enter an integer values", err)
		return
	}
	fmt.Println("You guessed is: ", guess)
```

> 这里需要注意的是，'\r'是回车，'\n'是换行，前者使光标到行首，后者使光标下移一格。通常用的Enter是两个加起来


当然这里也可以直接使用scanf来读取用户输入，这样的话就需要先声明存放用户输入的变量，然后在scanln中将该变量的地址传进去，直到用户输入回车，读取结束

```go
	var guess int
	fmt.Scanln(&guess)
	fmt.Println("You guessed is: ", guess)
```

3. 实现判断逻辑

比较用户输入和随机数的大小，如果用户输入大于随机数，则提示用户猜大了，如果用户输入小于随机数，则提示用户猜小了，如果用户输入等于随机数，则提示用户猜对了，退出程序

```go 
	if guess > secretNum {
		fmt.Println("You guess is bigger than the secret number, please try again")
	} else if guess < secretNum {
		fmt.Println("You guess is smaller than the secret number, please try again")
	} else {
		fmt.Println("Correct, you win!")
	}
```

4. 实现游戏循环

为了让游戏可以一直进行，需要将判断逻辑和读取用户输入的逻辑放在一个循环中，每次循环读取用户输入，然后判断用户输入是否正确，如果正确则退出循环，如果不正确则继续循环

```go
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
```
## 总结

通过这个实战项目，我学习到以下的内容

1. 生成随机数时，需要使用程序启动的时间戳来初始化随机数种子`rand.Seed(time.Now().UnixNano())`，不然生成的随机数会相同

2. `reader := bufio.NewReader(os.Stdin)`这行代码创建了一个bufio.Reader对象，它是一个包装了io.Reader的缓冲区读取器，`input, err := reader.ReadString('\n')`这行代码调用了bufio.Reader的ReadString方法，该方法从缓冲区读取数据直到遇到指定的字符

3. `input = strings.TrimSuffix(input, "\r\n")`使用strings.TrimSuffix函数去除字符串input的后缀，即删除字符串input的最后一个字符，该字符为`\r\n`，即回车和换行

4. `guess, err := strconv.Atoi(input)`将字符串转换为整数



