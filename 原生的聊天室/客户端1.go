package main

import (
	"fmt"
    "net"
    "os"
    "strings"
)

func main() {
	fmt.Println("正在连接服务器……")
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("连接服务器成功")

	buf2 := make([]byte, 4096)
	_, err = conn.Read(buf2)
	if err != nil {
		fmt.Println("conn.Read error:", err)
		return
	}

	//        fmt.Println(string(buf2[:n]))
	fmt.Println("&#9888;提示:长时间没有发送消息会被系统强制踢出")
	// 获取用户 网络地址 IP+port
	netAddr := conn.LocalAddr().String()
	//客户端发送消息到服务器
	go func() {
		for {
			buffer1 := make([]byte, 4096)

			//使用Stdin标准输入，因为scanf无法识别空格
			n, err := os.Stdin.Read(buffer1)
			if err != nil {
				fmt.Println("os.Stdin.Read error:", err)
				continue
			}
			_, _ = conn.Write(buffer1[:n]) //写操作出现error的概率比较低，这里省去判断
		}
	}()

	//接收服务器发来的数据
	for {
		buffer2 := make([]byte, 4096)
		n, err := conn.Read(buffer2)
		//fmt.Println("n2:",n)
		if n == 0 {
			fmt.Println("服务器已关闭当前连接，正在退出……")
			return
		}
		if err != nil {
			fmt.Println("conn.Read error:", err)
			return
		}
		fmt.Print(strings.ReplaceAll(string(buffer2[:n]), netAddr, "我"))

	}

}