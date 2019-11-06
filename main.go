package main

import (
	"flag"
	wpsgrabber "github.com/crossi-T2/wpsgrabber/wpsgrabber"
)

func main() {

	configFile := flag.String("c", "/etc/wpsgrabber/config.json", "Configuration file path")
	flag.Parse()

	wpsgrabber.New(*configFile)
}
