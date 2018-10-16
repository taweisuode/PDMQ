package flag

import (
	"flag"
	"fmt"
	"os"
)

type ArgvConfig struct {
	Help       bool
	Version    bool
	TcpListen  string
	HttpListen string
}

func ParseFlag() *ArgvConfig {
	config := Construct()
	help := flag.Bool("help", false, "show full pmq function guide")
	h := flag.Bool("h", false, "show full pmq function guide")
	version := flag.Bool("version", false, "show pmq version")
	v := flag.Bool("v", false, "show pmq version")
	tcpListen := flag.String("tcp_listen", config.TCPAddress, "tcp address host and port")
	httpListen := flag.String("http_listen", config.HTTPAddress, "tcp address host and port")
	flag.Parse()
	if *help || *h {
		showHelp()
	}
	if *version || *v {
		showVersion()
	}
	return &ArgvConfig{Help: *help, Version: *version, TcpListen: *tcpListen, HttpListen: *httpListen}
}
func showHelp() {
	fmt.Println("Usage\n")
	fmt.Println("    --version   show webserver version\n")
	fmt.Println("    --help      show full pmq function guide\n")
	fmt.Println("    --tpc_listen    set tcp address host and port\n")
	fmt.Println("    --http_listen    set http address host and port\n")
	os.Exit(0)
}

func showVersion() {
	fmt.Println("version    1.0.0")
	os.Exit(0)
}
