# PDMQ

    由于最近在学习go语言，并且nsq是go语言编写的一个分布式消息队列中间件，用来处理每天数十亿级别的消息推送。
    NSQ 具有分布式和去中心化拓扑结构，该结构具有无单点故障、故障容错、高可用性以及能够保证消息的可靠传递的特征，
    是一个成熟的、已在大规模生成环境下应用的产品。


#### 第一步：认识NSQ
    
    
![image](https://f.cloud.github.com/assets/187441/1700696/f1434dc8-6029-11e3-8a66-18ca4ea10aca.gif)


#### 第二部：PDMQD的基本功能实现
    1.PDMQD 跟nsq 一致，是采用go语言进行编写，这样能够最大化契合nsq的代码开发模式
    2.实现最基本的server，client 交互模式
      - 具体是：server 实现 tcp ，http 两种协议的监听，
      - 支持多种配置加载 终端配置>默认配置文件配置
      - client端拥有简单的交互模式，一问一答形式进行test
    3.消息分发类似nsq，分为topic->channel->consumer形式
    

      
