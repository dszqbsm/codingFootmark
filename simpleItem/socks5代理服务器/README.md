# SOCKS5代理服务器

该小项目是字节跳动青训营GO入门课程中的实战项目<https://juejin.cn/course/bytetech/7140987981803814919/section/7141265019756675103>

## SOCKS5介绍

SOCKS5协议都是明文协议，无法用来翻墙

若企业为了确保内网安全性，配置了很严格的防火墙策略，副作用就是访问内网中的资源会变得很麻烦，而SOCKS5协议相当于在防火墙上开了个口子，让授权的用户可以通过单个端口访问内部的所有资源

![SOCKS5原理](../images/socks5yuanli.jpg)

浏览器首先要跟SOCKS5代理服务器建立连接，再由代理服务器去和真正的服务器建立TCP连接

第一个阶段：协商阶段（协议版本号等信息）

第二个阶段：认证阶段（本项目不涉及，因为实现的是一个不加密的代理服务器）

第三个阶段：请求阶段

第四个阶段：relay阶段，代理服务器简单的将响应转发给浏览器，不关心流量的细节，因此流量可以是http、tcp等流量

## 实现思路

1. 构建一个简单的TCP echo server，用来测试编写的代理服务是否正确

该代理服务器功能简单，即发送啥就回复啥，利用goroutine开子线程处理，开销比操作系统子线程子进程少很多，可以轻松的处理上万的并发，这也是go的优势之一

```go
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
```

2. 认证阶段

认证流程，首先浏览器会给服务器发送一个报文，第一个字段是version协议版本号，固定是5

```go
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
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method failed: %v", err)
	}
	log.Println("ver", ver, "method", method)
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	_, err = conn.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write falied: %v", err)
	}
	return nil
}
```

3. 请求阶段

因为前四个字段长度相同，所以一次性读取，用4个字节的缓冲区，用`ReadFull`一下子填满，从而可以读取到这4个字节，然后逐个验证合法性

也就是按协议的字段定义规则，把字段都读取，然后进行验证分析，最后能够得到对应的IP和端口字段

```go
func connect(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05
	// CMD 0x01表示CONNECT请求
	// RSV 保留字段，值为0x00
	// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值
	// DST.PORT 目标端口，固定2个字节

	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed: %v", err)
	}
	ver, cmd, atyp := buf[0], buf[1], buf[3]
	if ver != socks5Ver {
		return fmt.Errorf("not support ver: %v", ver)
	}
	if cmd != cmdBind {
		return fmt.Errorf("not support cmd: %v", cmd)
	}
	addr := ""
	switch atyp {
	case atypIPv4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("read atyp failed: %v", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case atypeHOST:
		// 读取一个字节，即域名的长度
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostSize failed: %v", err)
		}
		// 按域名长度读取
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed: %v", err)
		}
		addr = string(host)
	case atypeIPv6:
		return errors.New("IPv6: no supported yet")
	default:
		return errors.New("invaild atyp")
	}
	// 复用前面四字节的缓冲区buf，用切片语法裁剪成2字节的缓冲区
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed: %v", err)
	}
	// 按照大端字节序解析出整型数值端口号
	port := binary.BigEndian.Uint16(buf[:2])

	log.Println("dial", addr, port)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址，四个字节，需要四个0
	// BND.PORT 服务绑定的端口DST.PORT，两个字节，需要两个0

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed: %v", err)
	}
	return nil
}
```
4. relay阶段

本阶段代理服务器会与真正的服务器建立tcp连接，需要建立浏览器和下游服务器的双向数据转换，io库中的Copy函数可以实现单向数据转化`func Copy(dst Write, src Reader) (written int64, err error)`会将src只读流中的数据用一个死循环逐步的拷贝到dst这个可写流中

这里需要启动两个协程，一个负责读取浏览器发来的数据，另一个负责读取下游服务器发来的数据，然后通过io库的Copy函数实现双向数据转换

需要等待任何一个方向的copy失败，即某一方关闭连接，才能终止整个连接

```go
	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return fmt.Errorf("dial failed: %v", err)
	}
	defer dest.Close()

	ctx, concel := context.WithCancel(context.Background())
	defer concel()

	go func() {
		_, _ = io.Copy(dest, reader)
		concel()
	}()

	go func() {
		_, _ = io.Copy(conn, dest)
		concel()
	}()

	<-ctx.Done()
```

## 总结

1. 学习了socks5协议的基本原理，了解了socks5协议的报文格式，以及各个字段的含义

2. 学习了go语言的net包，了解了net包的基本用法，以及net包的基本原理

3. 进一步熟悉bufio库的用法
