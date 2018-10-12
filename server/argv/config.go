package flag

type Config struct {
	TCPAddress  string
	HTTPAddress string
}

//配置化
func Construct() *Config {
	return &Config{
		TCPAddress:  "127.0.0.1:9998",
		HTTPAddress: "127.0.0.1:9999",
	}
}
