package main

import (
	routeChat "evenkey/route"
	"flag"
	"fmt"
)

func main() {

	c := &routeChat.Client{}

	// 定义命令行参数
	flag.IntVar(&c.Port, "port", 50051, "gRPC server port")
	flag.StringVar(&c.Addr, "addr", "localhost", "gRPC server address")
	flag.StringVar(&c.Username, "username", "root", "username")
	flag.StringVar(&c.Password, "password", "root", "password")
	flag.IntVar(&c.Keydelay, "delay", 10, "key press delay")
	// 解析命令行参数
	flag.Parse()

	// 打印解析出来的参数值
	fmt.Printf("port=%d, addr=%s, Username=%s, Password=%s\n",
		c.Port, c.Addr, c.Username, c.Password)

	c.Run()
}
