package pdmqd

import (
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"os"
)

type ArgvConfig struct {
	Help       bool
	Version    bool
	TcpListen  string
	HttpListen string
}

func ParseFlag(config *PDMQDConfig) *PDMQDConfig {
	help := flag.Bool("help", false, "show full pmq function guide")
	h := flag.Bool("h", false, "show full pmq function guide")
	version := flag.Bool("version", false, "show pmq version")
	v := flag.Bool("v", false, "show pmq version")
	tcpListen := flag.String("tcp_listen", config.TCPAddress, "tcp address host and port")
	httpListen := flag.String("http_listen", config.HTTPAddress, "tcp address host and port")

	logPath := flag.String("log_path", config.LogPath, "log path")
	flag.Parse()
	if *help || *h {
		showHelp()
	}
	if *version || *v {
		showVersion()
	}
	config.TCPAddress = *tcpListen
	config.HTTPAddress = *httpListen
	config.LogPath = *logPath

	logger, err := seelog.LoggerFromConfigAsFile(*logPath)
	if err != nil {
		fmt.Println("err parsing config log file", err)
	}
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()
	return config
}
func showHelp() {
	fmt.Println("Usage\n")
	fmt.Println("    --version   show webserver version\n")
	fmt.Println("    --help      show full pmq function guide\n")
	fmt.Println("    --tcp_listen    set tcp address host and port\n")
	fmt.Println("    --http_listen    set http address host and port\n")
	os.Exit(0)
}

func showVersion() {
	fmt.Println("version    1.0.0")
	os.Exit(0)
}
