package p2p

import (
	"fmt"
	"log"
	"net"
	"time"
)

func StartSvr() {
	// 1. 服务器启动监听
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9527})
	if err != nil {
		log.Fatal("start upd server err", err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()
	fmt.Println("begin server at ", listener.LocalAddr().String())

	// 定义切片存放两个 udp 地址
	peers := make([]*net.UDPAddr, 2)
	buf := make([]byte, 256)

	// 2. 接下来从2个UPD消息中获得链接的地址A、B
	n, addr, _ := listener.ReadFromUDP(buf)
	fmt.Printf("read from< %s >:%s\n", addr.String(), buf[:n])
	peers[0] = addr
	n, addr, _ = listener.ReadFromUDP(buf)
	fmt.Printf("read from< %s >:%s\n", addr.String(), buf[:n])
	peers[1] = addr
	fmt.Printf("begin nat \n")

	// 3. 将A和B分别介绍给彼此
	_, _ = listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	_, _ = listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	// 4. 睡眠10s确保消息发送完成
	time.Sleep(10 * time.Second)
}
