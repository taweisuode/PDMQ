package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":9400", time.Second)
	if err != nil {
		fmt.Println("tcp connect error ", err.Error())
		return
	}
	defer conn.Close()
	// 下面进行读写
	var wg sync.WaitGroup
	wg.Add(1)
	for {
		go handleWrite(conn, &wg)
		time.Sleep(5 * time.Second)
	}
	//go handleRead(conn, &wg)
	wg.Wait()
	select {}
	/*for {
		var buf = make([]byte, 32)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			break
		} else {
			if string(buf[:n]) == "exit" {
				fmt.Println("connect exit")
				break
			}
			fmt.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
		}
	}*/
	/*for {
		for {
			conn, err := topicObject.tcpListener.Accept()
			if err != nil {
				fmt.Println("tcp accept fail", err.Error())
			}
			id++
			if tconn, ok := conn.(*net.TCPConn); ok {
				go HandleConn(tconn, id)
			}
		}
		fmt.Println("connect close")
		defer conn.Close()
	}*/
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := conn.Write([]byte("V1"))
	if err != nil {
		fmt.Println("write data error", err)
	}
	for {
		a := []byte("SUB world ch \n")
		fmt.Println(a)
		conn.Write(a)
		//缓存 conn 中的数据
		buf := make([]byte, 1024)
		cnt, _ := conn.Read(buf)
		//回显服务器端回传的信息
		fmt.Print("服务器端回复" + string(buf[0:cnt]))
	}
	/*fmt.Println("connect pdmqd tcp listener...")
	work := ""
	fmt.Println("print your job")
	for {
		fmt.Scanf("%s", &work)
		switch work {
		case "send_protocol":
			fmt.Println("send your message")
			_, err := conn.Write([]byte("V1"))
			if err != nil {
				fmt.Println("write data error", err)
			}
			continue
		case "exit":
			os.Exit(1)
			break
		default:
			continue
		}
	}*/
	//defer conn.Close()
}
