package flag

const (
	//pdmqd 消息channel 大小
	MSGCHANSIZE = 9999
)

type Config struct {
	TCPAddress  string
	HTTPAddress string

	MsgChanSize int64
}

//配置化
func Construct() *Config {
	return &Config{
		TCPAddress:  "127.0.0.1:9998",
		HTTPAddress: "127.0.0.1:9999",
		MsgChanSize: MSGCHANSIZE,
	}
}
