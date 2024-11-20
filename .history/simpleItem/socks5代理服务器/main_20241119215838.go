package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const socks5Ver = 0x05
const cmdBind = 0x01
const atypIPv4 = 0x01
const atypeHostname = 0x03
const atypeIPv6 = 0x04

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
	err := auth(reader, conn)
	if err != nil {
		log.Panicf("client %v auth failed: %v", conn.RemoteAddr(), err)
		return 
	}
	log.Println("auth success")
}

func auth(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
	// X’00’ NO AUTHENTICATION REQUIRED
	// X’02’ USERNAME/PASSWORD

	ver,  err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read ver failed: %v", err)
	}
	if ver != socks5Ver {
		return fmt.Errorf("not support ver: %v", ver)
	}

	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read methodSize failed: %v", err)
	}
	method := make([]byte, methodSize)
	_, err 
}
