package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	res, err := http.Get("127.0.0.1:9999?topic=test&message=helloworld")
	if err != nil {
		return
	}
	//执行close之前一定要判断错误，如没有body会崩溃
	defer res.Body.Close()
	//重用连接一定要执行上和下两步
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	fmt.Println(res.Status)
	for k, v := range res.Header {
		fmt.Println(k, strings.Join(v, ""))
	}
}
