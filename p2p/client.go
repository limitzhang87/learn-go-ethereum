package p2p

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func StartCli() {
	// 1. 设定参数
	if len(os.Args) < 5 {
		fmt.Println("./client tag remoteIP remotePort port")
		return
	}

	// 本地要绑定的端口
	port, _ := strconv.Atoi(os.Args[4])
	// 客户端标识
	tag := os.Args[1]
	// 服务器IP
	remoteIP := os.Args[2]
	// 服务器端口
	remotePort, _ := strconv.Atoi(os.Args[3])

	// 为了绑定本地端口
	localAddr := net.UDPAddr{Port: port}

	// 2. 与服务器建立联系
	conn, err := net.DialUDP(
		"udp",
		&localAddr,
		&net.UDPAddr{
			IP:   net.ParseIP(remoteIP),
			Port: remotePort,
		},
	)

	if err != nil {
		log.Panic("Failed to DialUDP", err)
	}

	// 2.1 自我介绍
	_, _ = conn.Write([]byte("我是 ：" + tag))

	// 3. 从服务器获取目标
	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Failed to ReadFromUDP", err)
	}

	_ = conn.Close() // 读取后可以放弃服务器了
	toAddr := ParseAddr(string(buf[:n]))
	// 4. 两个人建立P2P通信
	p2p(&localAddr, &toAddr)
}

func ParseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

// p2p 实现P2P链接
func p2p(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) {
	// 1. 请求和对方建立连接
	conn, _ := net.DialUDP("udp", srcAddr, dstAddr)
	// 2. 发送打洞消息
	_, _ = conn.Write([]byte("打洞消息\n"))

	// 3. 启动一个goroutine监控标准输入
	go func() {
		buf := make([]byte, 256)
		for {
			// 接受UDP消息并打印
			n, _, _ := conn.ReadFromUDP(buf)
			if n > 0 {
				fmt.Printf("收到消息: %s p2p > ", buf[:n])
			}
		}
	}()

	// 4. 监控标准输入,发送给对方
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("p2p> ")
		// 读取标砖输入，已换行为读取标志
		data, _ := reader.ReadString('\n')
		_, _ = conn.Write([]byte(data))
	}
}
