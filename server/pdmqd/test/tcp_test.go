/**
 * @Time : 2019-07-11 21:38
 * @Author : zhuangjingpeng
 * @File : tcp_test.go
 * @Desc : file function description
 */
package test

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/**
 *  @desc:  需要将包改成main包 再执行，可以模拟tcp client方法
 *  @input:
 *  @resp:
 *
**/
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9998")
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	inputReader := bufio.NewReader(os.Stdin)

	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err: %v\n", err)
			break
		}
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "Q" {
			break
		}
		_, err = conn.Write([]byte(trimmedInput))

		if err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}
	}
}
