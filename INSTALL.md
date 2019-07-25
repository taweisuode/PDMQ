## 快速开始

#### 下载与安装

```
1. 下载 https://github.com/taweisuode/PDMQ.git
2. 进入server/app/
3. 执行build命令 go build -o pdmqd main.go
```
#### GO_PDMQ 客户端下载与安装

```
1. 下载https://github.com/taweisuode/GO_PDMQ.git
2. 进入main/文件夹
3. 执行build命令 go build -o consumer main.go
4. 执行consumer  ./consumer
```
#### API 使用
> ping 探活接口
```
curl  http://127.0.0.1:9401/ping
```

> API 发布消息请求

```
 curl -d "you are message" http://127.0.0.1:9401/pub?topic=name
```
#### 启动步骤


```
1. 启动pdmqd 服务server/app/pdmqd   
```

![](http://ww1.sinaimg.cn/large/006tNc79ly1g5cdthhpfxj30uw044ab0.jpg)

可以看出pdmqd 监听这tcp 以及http端口，http 有2个方法 一个ping探活接口，一个是pub 消息发布接口
```
2.启动consumer main/consumer 
```
![](http://ww2.sinaimg.cn/large/006tNc79ly1g5ce90oxewj30o5027aac.jpg)
这里显示连接成功，可以接收生产的消息数据了

同时pdmqd接收到了consumer的连接，已经有2次请求
![](http://ww1.sinaimg.cn/large/006tNc79ly1g5cedkssmij30nk02i0t7.jpg)

sub是注册请求，RDY是客户端用来流量控制的请求
```
3.通过pub api 进行投递消息

curl -d "you are message" http://127.0.0.1:9401/pub?topic=name
```
```
func (this *PDMQHandler) HandleMessage(message *pdmq.Message) error {
	fmt.Printf("[PDMQ CONSUMER] [%+v] get handler msg is [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), string(message.Body))
	return nil
}
```
consumer客户端 在HandleMessage中接收到了投递的消息

![](http://ww4.sinaimg.cn/large/006tNc79ly1g5ceizjxa5j30t8031mxr.jpg)


#### 相关代码
##### consumer 完整连接代码
```
package main

import (
	pdmq "GO_PDMQ"
	"fmt"
	"time"
)

type PDMQHandler struct {
}

func (this *PDMQHandler) HandleMessage(message *pdmq.Message) error {
	fmt.Printf("[PDMQ CONSUMER] [%+v] get handler msg is [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), string(message.Body))
	return nil
}
func main() {
	config := pdmq.NewConfig()
	consumer, err := pdmq.NewConsumer("name", "hello", config)
	if err != nil {
		fmt.Println(err)
	}
	consumer.AddHandler(&PDMQHandler{})

	err = consumer.ConnectToPDMQD("127.0.0.1:9400")
	if err != nil {
		fmt.Println(err)
	}

	select {}
}

```

可以看出往pdmqd注册了topic，channel 以及连接pdmqd 地址是 9400端口
