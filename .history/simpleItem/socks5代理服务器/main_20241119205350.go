package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		panic(err)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		// 启动一个子线程处理该连接
		go process(client)
	}
}

func process(conn net.Conn) {
	// 表示在函数退出的时候1一定要把连接关掉，因为该连接的生命周期就是整个函数的生命周期
	defer conn.Close()
	// 创建一个连接，只读的带缓冲的流
	reader := bufio.NewReader(conn)
	for {
		// 每次读一个字节，有缓冲，不会很慢
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		// 写入字节
		_, err = conn.Write([]byte{b})
		if err != nil {
			break
		}
	}
}
