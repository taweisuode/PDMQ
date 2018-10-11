package flag

import (
	"flag"
	"fmt"
	"os"
)

type FlagConfig struct {
	help    string
	version string
}

func ParseFlag() {
	help := flag.Bool("help", false, "show full pmq function guide")
	h := flag.Bool("h", false, "show full pmq function guide")
	version := flag.Bool("version", false, "show pmq version")
	v := flag.Bool("v", false, "show pmq version")
	flag.Parse()
	if *help || *h {
		showHelp()
	}
	if *version || *v {
		showVersion()
	}
}
func showHelp() {
	fmt.Println("Usage\n")
	fmt.Println("    --version   show webserver version\n")
	fmt.Println("    --help      show full pmq function guide\n")
	os.Exit(0)
}

func showVersion() {
	fmt.Println("version    1.0.0")
	os.Exit(0)
}
