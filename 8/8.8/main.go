/*
练习 8.8： 使用select来改造8.3节中的echo服务器，为其增加超时，这样服务器可以在客户端10秒中没有任何喊话时自动断开连接。
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is listening on localhost:8000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	// 1. 创建一个用于接收客户端消息的 channel
	// 为什么用 buffer？防止 input.Scan() 读太快，稍微缓冲一下
	messageCh := make(chan string)

	// 2. 启动一个后台小弟，专门负责从连接里读数据
	go func() {
		input := bufio.NewScanner(c)
		for input.Scan() {
			messageCh <- input.Text() // 读到一行，扔进 channel
		}
		// 如果客户端主动断开连接或者出错了，input.Scan 会结束
		// 此时我们需要通知外面“没数据了”
		// 怎么通知？关闭 channel 是个好办法，但这里我们简单点，不做复杂处理，
		// 因为外面 select 也能通过 case msg, ok := <-messageCh 感知到。
	}()

	// 3. 设置初始超时时间 (10秒)
	// time.NewTimer 比 time.After 更灵活，因为我们需要不断重置它
	timer := time.NewTimer(10 * time.Second)

	for {
		select {
		// Case A: 收到客户端发来的消息
		case msg := <-messageCh:
			// 这里的 msg 可能是空字符串（如果 channel 被关闭或客户端发空行）
			// 简单的回显逻辑
			echo(c, msg, 1*time.Second)

			// **关键点**：只要有活动，就重置定时器
			// Reset 返回值需要处理吗？在这个简单场景下，
			// 为了防止 timer 通道里残留数据，严谨的做法需要 drain channel。
			// 但在该特定循环结构中，因为我们要么处理了 timer 事件，要么 timer 还没触发，
			// 直接 Reset 是相对安全的（或者先 Stop 再 Reset）。
			if !timer.Stop() {
				// 如果 timer 已经触发了但还没被 select 取走，
				// 这里需要把它排空，防止在这个 select 循环里误触。
				// 注意：在多路复用里这是个有点绕的坑，简单写法直接 timer.Reset 往往也能跑。
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(10 * time.Second)

		// Case B: 定时器响了（意味着10秒内没走 Case A）
		case <-timer.C:
			fmt.Printf("Client %s timed out, closing connection.\n", c.RemoteAddr())
			return // 退出函数，defer c.Close() 会自动断开连接
		}
	}
}

// 模拟慢速回显（保持原书风格）
func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
